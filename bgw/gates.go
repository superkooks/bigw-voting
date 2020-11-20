package bgw

import "bigw-voting/shamir"

//
// AdditionGate has two inputs and adds them together
type AdditionGate struct {
	Inputs []Input
}

// GetInputs returns the required inputs
func (a *AdditionGate) GetInputs() []Input {
	return a.Inputs
}

// GetOutput returns the output from the gate, evaluting previous gates
func (a *AdditionGate) GetOutput() int {
	in0 := a.Inputs[0].GetOutput()
	in1 := a.Inputs[1].GetOutput()

	return (in0 + in1) % shamir.FieldSize
}

//
// ConstMultiplicationGate has one input and multiplies it by a constant
type ConstMultiplicationGate struct {
	Input    Input
	Constant int
}

// GetInputs returns the required inputs
func (c *ConstMultiplicationGate) GetInputs() []Input {
	return []Input{c.Input}
}

// GetOutput returns the output from the gate, evaluting previous gates
func (c *ConstMultiplicationGate) GetOutput() int {
	in := c.Input.GetOutput()

	return (in * c.Constant) % shamir.FieldSize
}
