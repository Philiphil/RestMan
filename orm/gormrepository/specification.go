package gormrepository

import (
	"fmt"
	"strings"
)

// Specification defines the interface for query specifications used to filter database queries.
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

// And combines multiple specifications with the AND logical operator.
func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		separator:      "AND",
	}
}

// Or combines multiple specifications with the OR logical operator.
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

// Not negates a specification with the NOT logical operator.
func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

type binaryOperatorSpecification[T any] struct {
	field    string
	operator string
	value    T
}

func (s binaryOperatorSpecification[T]) GetQuery() string {
	return fmt.Sprintf("%s %s ?", s.field, s.operator)
}

func (s binaryOperatorSpecification[T]) GetValues() []any {
	return []any{s.value}
}

// Equal creates a specification for equality comparison.
func Equal[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "=",
		value:    value,
	}
}

// GreaterThan creates a specification for greater-than comparison.
func GreaterThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">",
		value:    value,
	}
}

// GreaterOrEqual creates a specification for greater-than-or-equal comparison.
func GreaterOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: ">=",
		value:    value,
	}
}

// LessThan creates a specification for less-than comparison.
func LessThan[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<",
		value:    value,
	}
}

// LessOrEqual creates a specification for less-than-or-equal comparison.
func LessOrEqual[T comparable](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "<=",
		value:    value,
	}
}

// In creates a specification for checking if a field value is in a list.
func In[T any](field string, value []T) Specification {
	return binaryOperatorSpecification[[]T]{
		field:    field,
		operator: "IN",
		value:    value,
	}
}

// Like creates a specification for case-sensitive pattern matching.
func Like[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
		field:    field,
		operator: "LIKE",
		value:    value,
	}
}

// Ilike creates a specification for case-insensitive pattern matching.
func Ilike[T any](field string, value T) Specification {
	return binaryOperatorSpecification[T]{
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

// IsNull creates a specification for checking if a field is NULL.
func IsNull(field string) Specification {
	return stringSpecification(fmt.Sprintf("%s IS NULL", field))
}
// IsNotNull creates a specification for checking if a field is NOT NULL.
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

// OrderBy creates a specification for sorting results by a field and direction.
func OrderBy(field string, direction string) Specification {
	return orderSpecification{
		field:     field,
		direction: direction,
	}
}
