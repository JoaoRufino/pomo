package models

type Protocol struct {
	Cid     string
	Payload Payload
}

type Payload interface {
}
