package types

type SourceType string

const (
	SourceTypeStripe SourceType = "stripe"
)

func (s SourceType) String() string {
	return string(s)
}
