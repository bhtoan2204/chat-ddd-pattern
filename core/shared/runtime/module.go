package modruntime

type Module interface {
	Start() error
	Stop() error
}
