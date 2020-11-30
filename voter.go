package main

import "bigw-voting/p2p"

// Voter represents a voter voting in the election
type Voter struct {
	Peer   *p2p.Peer
	Status string
}

// NewVoter creates a new Voter
func NewVoter(peer *p2p.Peer) *Voter {
	return &Voter{
		Peer:   peer,
		Status: "Voting InProgress",
	}
}
