package hsm

import (
	"reflect"
	"strings"
)

// Printer public definition of HSM printers.
type Printer[C any] interface {
	// Print returns a HSM representation in whatever format the printer decides to (e.g. PNG, plain-text, etc).
	Print(hsm *HSM[C]) []byte
}

func fnSignatureString(f interface{}) string {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return "<not a function>"
	}

	buf := strings.Builder{}

	buf.WriteString("func (")

	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}

		buf.WriteString(t.In(i).String())
	}

	buf.WriteString(")")

	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			buf.WriteString(" (")
		} else {
			buf.WriteString(" ")
		}

		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				buf.WriteString(", ")
			}

			buf.WriteString(t.Out(i).String())
		}

		if numOut > 1 {
			buf.WriteString(")")
		}
	}

	return buf.String()
}
