package Payload

type Payload interface {
	FromString(payload string)
	ToString() string
}
