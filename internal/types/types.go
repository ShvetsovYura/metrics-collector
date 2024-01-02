package types

type Gauge float64
type Counter int64

type Sender interface {
	Send(string)
}

type Stringer interface {
	ToString() string
}

type Stored interface {
}
