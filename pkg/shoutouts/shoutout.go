package shoutouts

import uuid "github.com/satori/go.uuid"

// ShoutoutKind is the kind out shoutout.
type ShoutoutKind struct {
	Name     string
	Abbrev   string
	Desc     string
	ThumbURL string
}

// Shoutout is the shoutout a user sends or announces to another.
type Shoutout struct {
	ID            string
	Kind          *ShoutoutKind
	SenderID      string
	SenderName    string
	RecipientID   string
	RecipientName string
	TeamID        string
	Comment       string
	ResponseURL   string
}

// New initializes a new Shoutout.
func New() *Shoutout {
	return &Shoutout{
		ID:   uuid.NewV4().String(),
		Kind: &ShoutoutKind{},
	}
}
