package commands

import (
	"bigw-voting/p2p"
	"bigw-voting/util"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
)

// TrusteeVote is an IRV map and a name
type TrusteeVote struct {
	Name  string
	Votes map[string]int
}

// LocalTrusteeVotes is a list of all trustee votes loaded by us
var LocalTrusteeVotes []TrusteeVote

// AllTrusteeVotes is a list of all trustee votes loaded by all peers.
// Used for preventing overlaps.
var AllTrusteeVotes []string

// VotepackTrustees is the list of trustees mentioned in the votepack
var VotepackTrustees []string

// CommandLoadTrustee loads a new trustee vote
func CommandLoadTrustee(args []string) {
	f, err := os.Open(args[0])
	if err != nil {
		util.Errorf("unable to open trustee vote: %v\n", err)
		return
	}

	defer f.Close()

	// Deobfuscate votes. This is not encryption, it simply prevents people from reading
	// the plaintext.
	archiver, err := gzip.NewReader(f)
	if err != nil {
		util.Errorf("unable to deobfuscate layer-2 trustee vote: %v\n", err)
		return
	}

	layer0Reader := base64.NewDecoder(base64.StdEncoding, archiver)
	layer0, err := ioutil.ReadAll(layer0Reader)
	if err != nil {
		util.Errorf("unable to deobfuscate layer-1 trustee vote: %v\n", err)
		return
	}

	// Unmarhsal trustee vote JSON
	unmarshaled := &TrusteeVote{}
	err = json.Unmarshal(layer0, unmarshaled)
	if err != nil {
		util.Errorf("unable to unmarshal trustee vote: %v\n", err)
		return
	}

	// Inform peers of our new trustee vote
	err = p2p.BroadcastMessage([]byte("TrusteeVote "+unmarshaled.Name), 0)
	if err != nil {
		util.Errorln(err)
		return
	}

	// Store our trustee vote
	LocalTrusteeVotes = append(LocalTrusteeVotes, *unmarshaled)

	var found bool
	for _, v := range VotepackTrustees {
		if v == unmarshaled.Name {
			found = true
			break
		}
	}

	if !found {
		util.Errorln("tried to add trustee vote not in votepack.\n")
		panic("trying to cheat the election? you've been caught.")
	}

	found = false
	for _, v := range AllTrusteeVotes {
		if v == unmarshaled.Name {
			found = true
			break
		}
	}

	if found {
		util.Errorln("tried to add trustee vote for person who already has trustee.\n")
		panic("trying to cheat the election? you've been caught.")
	}

	AllTrusteeVotes = append(AllTrusteeVotes, unmarshaled.Name)
}
