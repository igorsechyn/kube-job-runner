package uuid

import uuid "github.com/satori/go.uuid"

type Generator struct{}

func (generator Generator) Generate() string {
	return uuid.NewV4().String()
}
