package p2p

import (
	"encoding/binary"
	"time"
)

// Message has information about the message such as sequence number and data.
// Of a message received, the first 10 bytes are the header.
type Message struct {
	SequenceNumber int
	Ack            bool
	Broadcast      bool
	MaxBounces     int8
	Data           []byte

	sentAt time.Time
}

// NewMessage creates a new message with header
func (p *Peer) NewMessage(data []byte, ack bool, broadcast bool) *Message {
	p.latestSeqNumber++
	return &Message{
		SequenceNumber: p.latestSeqNumber,
		Ack:            ack,
		Broadcast:      broadcast,
		Data:           data,
		sentAt:         time.Now(),
	}
}

// Serialize returns the []byte representation of a message
func (m *Message) Serialize() []byte {
	headeredMsg := make([]byte, 11+len(m.Data))
	binary.LittleEndian.PutUint64(headeredMsg, uint64(m.SequenceNumber))

	if m.Ack {
		headeredMsg[8] = byte(1)
	} else {
		headeredMsg[8] = byte(0)
	}

	if m.Broadcast {
		headeredMsg[9] = byte(1)
	} else {
		headeredMsg[9] = byte(0)
	}

	headeredMsg[10] = byte(m.MaxBounces)

	copy(headeredMsg[11:], m.Data)

	return headeredMsg
}

// Deserialize turns a []byte into a message structure
func (m *Message) Deserialize(msg []byte) {
	seqNum := binary.LittleEndian.Uint64(msg[:8])

	if msg[8] == byte(1) {
		m.Ack = true
	} else {
		m.Ack = false
	}

	if msg[9] == byte(1) {
		m.Broadcast = true
	} else {
		m.Broadcast = false
	}

	m.MaxBounces = int8(msg[10])
	m.SequenceNumber = int(seqNum)

	m.Data = msg[11:]
}
