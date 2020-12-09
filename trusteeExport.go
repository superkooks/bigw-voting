package main

import (
	"bigw-voting/commands"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ExportTrusteeVote creates a new trustee vote to be sent to a trustee
func ExportTrusteeVote() {
	// Input name
	var name string
	fmt.Print("Enter name (this must match the votepack): ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		panic(err)
	}

	// Input IRV Votes
	irv := make(map[string]int)
	for _, cand := range votepack.Candidates {
		var v string
		fmt.Print(cand + ": ")
		_, err = fmt.Scanln(&v)
		if err != nil {
			panic(err)
		}

		irv[cand], err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}

	vote := &commands.TrusteeVote{Name: name, Votes: irv}

	// Validate votes so that IRV works
	for k, v := range irv {
		if v > len(irv) {
			panic("votes are invalid")
		}

		for l, w := range irv {
			if v == w && k != l {
				panic("votes are not transferrable")
			}
		}
	}

	fmt.Println("\nValidated votes")

	// Marshal vote into JSON
	marshaled, err := json.Marshal(vote)
	if err != nil {
		panic(err)
	}

	var filename string
	fmt.Print("File to output to: ")
	_, err = fmt.Scanln(&filename)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// Obfuscate votes. This is not encryption, it simply prevents people from reading
	// the plaintext.
	archiver := gzip.NewWriter(f)

	layer1Writer := base64.NewEncoder(base64.StdEncoding, archiver)
	layer0 := bytes.NewBuffer(marshaled)
	_, err = io.Copy(layer1Writer, layer0)
	if err != nil {
		panic(err)
	}

	layer1Writer.Close()
	archiver.Close()
	f.Close()

	os.Exit(0)
}
