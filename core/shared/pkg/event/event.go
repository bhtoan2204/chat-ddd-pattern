package event

type Event interface {
	GetName() string
	GetData() interface{}
}
