package valueobject

type PlainPassword struct {
	value string
}

func NewPlainPassword(value string) (PlainPassword, error) {
	normalized, err := normalizePasswordValue(value)
	if err != nil {
		return PlainPassword{}, err
	}
	return PlainPassword{value: normalized}, nil
}

func (p PlainPassword) Value() string {
	return p.value
}
