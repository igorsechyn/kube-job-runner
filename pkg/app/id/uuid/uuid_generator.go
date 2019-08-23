package uuid

import "github.com/satori/go.uuid"

type Generator struct{}

func (generator Generator) Generate() string {
	return uuid.NewV4().String()
}
