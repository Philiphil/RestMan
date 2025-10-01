package gormrepository_reflection

import (
	"fmt"
	"strings"
)

// Specification interface remains the same as it's not generic
type Specification interface {
	GetQuery() string
	GetValues() []any
}

type joinSpecification struct {
	specifications []Specification
	separator      string
}

func (s joinSpecification) GetQuery() string {
	queries := make([]string, 0, len(s.specifications))

	for _, spec := range s.specifications {
		queries = append(queries, spec.GetQuery())
	}

	return strings.Join(queries, fmt.Sprintf(" %s ", s.separator))
}

func (s joinSpecification) GetValues() []any {
	values := make([]any, 0)

	for _, spec := range s.specifications {
		values = append(values, spec.GetValues()...)
	}

	return values
}

func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

func Or(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "OR",
	}
}

type notSpecification struct {
	Specification
}

func (s notSpecification) GetQuery() string {
	return fmt.Sprintf(" NOT (%s)", s.Specification.GetQuery())
}

func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

// binaryOperatorSpecification now uses `any` for value type as generics are removed.
// The actual type checking or conversion might need to happen at runtime if specific GORM features depend on it,
// but for query building with '?', `any` is generally fine.
type binaryOperatorSpecification struct {
	field    string
	operator string
	value    any // Changed from generic T to any
}

func (s binaryOperatorSpecification) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binaryOperatorSpecification) GetValues() []any {
	return []any{s.value}
}

func Equal(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "=",
		value:    value,
	}
}

func GreaterThan(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: ">",
		value:    value,
	}
}

func GreaterOrEqual(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

func LessThan(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "<",
		value:    value,
	}
}

func LessOrEqual(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "<=",
		value:    value,
	}
}

func In(field string, value []any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "IN",
		value:    value, // value is already []any, matching the struct field
	}
}

func Like(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "LIKE",
		value:    value,
	}
}

func Ilike(field string, value any) Specification {
	return binaryOperatorSpecification{
		field:    field,
		operator: "ILIKE",
		value:    value,
	}
}

type stringSpecification string

func (s stringSpecification) GetQuery() string {
	return string(s)
}

func (s stringSpecification) GetValues() []any {
	return nil
}

func IsNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NULL", field))
}
func IsNotNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NOT NULL", field))
}

type orderSpecification struct {
	field     string
	direction string
}

func (s orderSpecification) GetQuery() string {
	return s.field + " " + s.direction
}

func (s orderSpecification) GetValues() []any {
	return nil
}

func OrderBy(field string, direction string) Specification {
	return orderSpecification{
		field:     field,
		direction: direction,
	}
}
