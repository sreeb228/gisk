package gisk

type ValueInterface interface {
	Parse(gisk *Gisk) (Value, error)
}
