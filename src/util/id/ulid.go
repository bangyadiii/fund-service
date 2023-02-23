package id

import (
	"strings"

	u "github.com/oklog/ulid/v2"
)

type ulid struct {
}

func NewUlid() IDGenerator {
	return &ulid{}
}

func (*ulid) Generate() string {
	unique := u.Make().String()
	return strings.ToLower(unique)
}
