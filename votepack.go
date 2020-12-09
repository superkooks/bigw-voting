package main

import (
	"bigw-voting/ui"
	"bigw-voting/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Votepack represents a JSON-encoded pack of the candidates + trustee votes
type Votepack struct {
	Candidates   []string
	TrusteeVotes []string
}

// NewVotepackFromFile opens a file and loads the Votepack
func NewVotepackFromFile(filename string) *Votepack {
	f, err := os.Open(filename)
	if err != nil {
		ui.Stop()
		panic(fmt.Errorf("unable to open votepack: %v", err))
	}

	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		ui.Stop()
		panic(fmt.Errorf("unable to read votepack: %v", err))
	}

	return NewVotepackFromJSON(b)
}

// NewVotepackFromJSON unmarshalls JSON into a Votepack
func NewVotepackFromJSON(marshalled []byte) *Votepack {
	v := new(Votepack)
	err := json.Unmarshal(marshalled, v)
	if err != nil {
		ui.Stop()
		panic(fmt.Errorf("unable to parse votepack: %v", err))
	}

	return v
}

// Export returns the JSON representation of the Votepack
func (v *Votepack) Export() []byte {
	b, err := json.Marshal(v)
	if err != nil {
		util.Errorf("error while exporting votepack: %v", err)
	}

	return b
}

// ExportToFile exports the Votepack to a file
func (v *Votepack) ExportToFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		ui.Stop()
		panic(err)
	}

	defer f.Close()
	_, err = f.Write(v.Export())
	if err != nil {
		util.Errorf("error while exporting votepack to file: %v", err)
	}
}
