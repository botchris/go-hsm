package hsm

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strings"
)

// PlantUMLPrinter provides a simple PlantUML plain-text notation printer.
//
// Usage:
//
// ```go
// printer := hsm.NewPlantUMLPrinter()
// out := printer.Print(MyMachine)
// println(string(out))
// ```.
type PlantUMLPrinter[C any] struct {
	machine   *HSM[C]
	ids       map[string]uint32
	allStates []*Vertex[C]
}

// NewPlantUMLPrinter returns a new printer.
func NewPlantUMLPrinter[C any]() Printer[C] {
	return &PlantUMLPrinter[C]{
		ids: make(map[string]uint32),
	}
}

// Print prints the given HSM.
func (p *PlantUMLPrinter[C]) Print(hsm *HSM[C]) []byte {
	p.init(hsm)

	return []byte(p.print())
}

func (p *PlantUMLPrinter[C]) init(hsm *HSM[C]) *PlantUMLPrinter[C] {
	p.machine = hsm
	for _, s := range p.machine.states {
		p.ids[s.id] = p.fNV32a(s.id)
	}

	var (
		merge  = []*Vertex[C]{p.machine.errorState}
		choice []*Vertex[C]
		entry  []*Vertex[C]
		start  []*Vertex[C]
		final  []*Vertex[C]
		state  []*Vertex[C]
	)

	for _, v := range p.machine.states {
		switch v.kind {
		case vertexKindChoice:
			choice = append(choice, v)
		case vertexKindStart:
			start = append(start, v)
		case vertexKindFinal:
			final = append(final, v)
		case vertexKindState:
			state = append(state, v)
		}

		if v.entryState != nil {
			entry = append(entry, v.entryState)
		}
	}

	merge = append(merge, choice...)
	merge = append(merge, entry...)
	merge = append(merge, start...)
	merge = append(merge, final...)
	merge = append(merge, state...)

	p.allStates = merge

	return p
}

func (p *PlantUMLPrinter[C]) print() string {
	var (
		out     = ""
		caption = fmt.Sprintf("caption HSM %s@%s\n", p.machine.name, p.machine.Current().id)
		roots   []*Vertex[C]
	)

	for _, v := range p.allStates {
		if v.parent == nil {
			roots = append(roots, v)
		}
	}

	for _, v := range roots {
		out += p.renderVertex(v)
	}

	// ensure THE final state is at the very end
	reg := regexp.MustCompile(`(.+) --\> \[\*\].*`)
	matches := reg.FindAllString(out, 10)
	final := strings.Join(matches, "\n")
	out = reg.ReplaceAllLiteralString(out, "")

	template := `@startuml
%s
%s
%s
@enduml`

	return fmt.Sprintf(template, caption, out, final)
}

func (p *PlantUMLPrinter[C]) renderVertex(v *Vertex[C]) string {
	content := ""
	children := p.children(v)
	alias := p.alias(v)
	template := ""

	switch v.kind {
	case vertexKindError:
		template = fmt.Sprintf("state %q as %s #Red\n", v.id, alias)
		template += "%s"
	case vertexKindChoice:
		template = fmt.Sprintf("state %s <<choice>>\n", alias)
		template += "%s\n"
	case vertexKindEntry:
		template += "%s\n"
	case vertexKindStart:
		template = "%s\n"
	case vertexKindFinal:
		template = "%s\n"
	case vertexKindState:
		template = fmt.Sprintf("state %q as %s {\n", v.id, alias)
		template += "%s\n"
		template += "}\n"
	default:
		template = "%s\n"
	}

	if v.onEntry != nil {
		content += fmt.Sprintf("%s : entry / %s\n", alias, v.onEntry)
	}

	if v.onExit != nil {
		content += fmt.Sprintf("%s : exit / %s\n", alias, v.onExit)
	}

	for _, t := range v.edges.list() {
		content += p.renderTransitionFor(v, t)
	}

	if len(children) > 0 {
		for _, c := range children {
			content += p.renderVertex(c)
		}
	}

	return fmt.Sprintf(template, content)
}

func (p *PlantUMLPrinter[C]) renderTransitionFor(v *Vertex[C], t *Transition[C]) string {
	var (
		out          = ""
		currentState = p.machine.Current()
		green        bool
		from         = p.alias(v)
		to           = p.alias(t.nextStatePtr)
	)

	if from == "" && to == "" {
		return out
	}

	label := p.renderTransitionLabelFor(v, t)

	for _, possible := range currentState.edges.list() {
		if possible == t {
			if possible.guard == nil || possible.guard.method(p.machine.context) {
				green = true

				break
			}
		}
	}

	switch t.kind {
	case transitionKindInternal:
		if label != "" {
			if green {
				label = "<color:green>" + label
			}

			label = " : " + label
		}

		out = fmt.Sprintf("%s %s\n", from, label)
	case transitionKindNormal:
		if label != "" {
			label = " : " + label
		}

		arrow := "-->"
		if green {
			arrow = "-[#green]->"
		}

		out += fmt.Sprintf("%s %s %s%s\n", from, arrow, to, label)
	}

	return out
}

func (p *PlantUMLPrinter[C]) renderTransitionLabelFor(from *Vertex[C], t *Transition[C]) string {
	trigger := strings.Replace(p.machine.kind(t.signal), "*", "", 1)
	guard := ""
	effect := ""

	if t.signal == nil {
		trigger = ""
	}

	if t.guard != nil {
		guard = fmt.Sprintf(`[%s]`, t.guard.label)
	}

	if t.effect != nil {
		effect = fmt.Sprintf(`/ %s`, t.effect.label)
	}

	if from.kind == vertexKindChoice && trigger == "" && guard == "" && effect == "" {
		return "[else]"
	}

	return strings.TrimSpace(strings.Join([]string{trigger, guard, effect}, " "))
}

func (p *PlantUMLPrinter[C]) alias(v *Vertex[C]) string {
	switch v.kind {
	case vertexKindEntry:
		return "[*]"
	case vertexKindStart:
		return "[*]"
	case vertexKindFinal:
		return "[*]"
	case vertexKindChoice:
		return fmt.Sprintf("choice_%d", p.ids[v.id])
	case vertexKindError:
		return fmt.Sprintf("error_%d", p.ids[v.id])
	}

	return fmt.Sprintf("state_%d", p.ids[v.id])
}

func (p *PlantUMLPrinter[C]) children(v *Vertex[C]) []*Vertex[C] {
	var children []*Vertex[C]

	for _, s := range p.allStates {
		if s.parent != nil && s.parent == v {
			children = append(children, s)
		}
	}

	return children
}

func (p *PlantUMLPrinter[C]) fNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	_, _ = algorithm.Write([]byte(text))

	return algorithm.Sum32()
}
