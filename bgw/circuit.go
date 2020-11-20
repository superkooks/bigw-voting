package bgw

// Circuit is the overall structure of the gates
type Circuit struct {
	Inputs []*ShamirInput

	Head Input // The head of the circuit (highest topological order)
}

// Execute runs the circuit
func (c *Circuit) Execute() int {
	return c.Head.GetOutput()
}
