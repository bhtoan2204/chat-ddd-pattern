package utils

type QueryOptions struct {
	Conditions     []Condition
	Limit          *int
	Offset         *int
	OrderBy        string
	OrderDirection string
}

type Condition struct {
	Field    string
	Value    interface{}
	Operator Operator
}

type Operator string

const (
	Equal              Operator = "="
	NotEqual           Operator = "!="
	GreaterThan        Operator = ">"
	GreaterThanOrEqual Operator = ">="
	LessThan           Operator = "<"
	LessThanOrEqual    Operator = "<="
	In                 Operator = "in"
	NotIn              Operator = "not in"
	Like               Operator = "like"
	NotLike            Operator = "not like"
	IsNull             Operator = "is null"
	IsNotNull          Operator = "is not null"
)

func (c Condition) BuildCondition() string {
	switch c.Operator {
	case IsNull, IsNotNull:
		return c.Field + " " + string(c.Operator)
	default:
		return c.Field + " " + string(c.Operator) + " " + "?"
	}
}
