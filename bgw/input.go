package bgw

import (
	"bigw-voting/shamir"
	"sort"
)

// Input represents either an input or a gate
type Input interface {
	GetOutput() int
}

// ShamirInput is an input that takes a single number and returns the points
// from a Shamir polynomial, one point is retained from the set, the rest should
// be distributed to peers. ShamirInput also requires peer's points to be entered
type ShamirInput struct {
	Point int
}

// NewLocalShamirInput returns a new ShamirInput with a local input and shares for peers
func NewLocalShamirInput(in int, numOfParties int, externalIP string, peerIPs []string) (*ShamirInput, map[string]int) {
	s := &ShamirInput{}

	sortedIPs := append(peerIPs, externalIP)
	sort.Strings(sortedIPs)

	points := shamir.ConstructPoints(in, numOfParties, numOfParties)
	out := make(map[string]int)
	for k, v := range points {
		if sortedIPs[k] == externalIP {
			s.Point = v[1]
		} else {
			out[sortedIPs[k]] = v[1]
		}
	}

	return s, out
}

// NewRemoteShamirInput returns a new ShamirInput to input peer's shares in
func NewRemoteShamirInput() *ShamirInput {
	return &ShamirInput{}
}

// IsReady returns whether a ShamirInput is ready. Useful for remote
func (s *ShamirInput) IsReady() bool {
	if s.Point == 0 {
		return false
	}

	return true
}

// AddPeerShare adds a peer's share to the input
func (s *ShamirInput) AddPeerShare(share int) {
	s.Point = share
}

// GetOutput returns the points from a Shamir polynomial.
func (s *ShamirInput) GetOutput() int {
	return s.Point
}
