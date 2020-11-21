package bgw

import (
	"bigw-voting/util"
	"fmt"
	"sort"
)

// NewVotingCircuit creates a standard voting circuit, returning the head of the circuit and
// the peer shares to distribute for each gate
func NewVotingCircuit(in int, externalIP string, peerIPs []string) (Input, map[string][]int) {
	allParties := append(peerIPs, externalIP)
	sort.Strings(allParties)

	var inputs []*ShamirInput
	var allPeerShares []map[string]int

	// Create Shamir inputs for all parties
	for _, p := range allParties {
		var newIn *ShamirInput
		if p != externalIP {
			// Create a remote Shamir input for party
			newIn = NewRemoteShamirInput()
		} else {
			// Create a local Shamir input so we can vote
			var peerShares map[string]int
			newIn, peerShares = NewLocalShamirInput(in, len(allParties), externalIP, peerIPs)
			allPeerShares = append(allPeerShares, peerShares)
		}

		inputs = append(inputs, newIn)
	}

	fmt.Printf("Num of inputs: %v\n", len(inputs))

	// Insert addition gates for all inputs
	var nextIn Input
	for _, v := range inputs {
		if nextIn == nil {
			nextIn = v
			continue
		}

		fmt.Println("Adding addition gate")
		nextIn = &AdditionGate{
			Inputs: []Input{nextIn, v},
		}
	}

	// Change string of maps to map of strings
	m := make(map[string][]int)
	for _, v := range allPeerShares {
		for key, point := range v {
			m[key] = append(m[key], point)
		}
	}

	// Head of the circuit is the last gate
	return nextIn, m
}

// DescendCircuit passes through the circuit filling in shares from
// lowest to highest topological order
func DescendCircuit(head Input, shares []int) {
	descending := descendNode(head)

	// Iterate over the slice in reverse (ascending order)
	for i := len(descending) - 1; i >= 0; i-- {
		if shamirInput, ok := descending[i].(*ShamirInput); ok {
			shamirInput.AddPeerShare(shares[len(shares)-i-1])
		} else {
			util.Errorln("cannot add input to non-shamir input or local shamir input")
			return
		}
	}
}

func descendNode(node Input) []Input {
	fmt.Printf("%t\n", node)

	switch n := node.(type) {
	case *ShamirInput:
		if n.Point == 0 {
			return []Input{node}
		}

		return []Input{}

	case Gate:
		ins := n.GetInputs()

		var out []Input
		for _, v := range ins {
			out = append(out, descendNode(v)...)
		}

		return out

	default:
		util.Errorln("cannot descend through unknown gate")
		return []Input{}
	}
}
