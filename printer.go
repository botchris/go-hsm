package hsm

// Printer public definition of HSM printers
type Printer interface {
	// Print returns a HSM representation in whatever format the printer decides to (e.g. PNG, plain-text, etc).
	Print(hsm *HSM) []byte
}
