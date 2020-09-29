package message

import (
	"fmt"
	"github.com/gocql/gocql"
)

type UUID gocql.UUID

func (i UUID) String() string {
	return i.String()
}

func NewID() gocql.UUID {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		fmt.Println(err)
	}
	return uuid
}

type Number struct {
	MagicNumber int `json:"magic_number"`
}

type Message struct {
	Id          gocql.UUID `json:"id"`
	Email       string     `json:"email" validate:"regexp=^[0-9a-z]+(\\.[0-9a-z]+)+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	MagicNumber int        `json:"magic_number" validate:"nonzero"`
	Title       string     `json:"title" validate:"nonzero"`
	Content     string     `json:"content" validate:"nonzero"`
}
