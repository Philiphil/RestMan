package mongorepository

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Specification defines the interface for MongoDB query specifications.
type Specification interface {
	GetFilter() bson.M
	GetSort() bson.D
}

type filterSpecification struct {
	filter bson.M
}

func (s filterSpecification) GetFilter() bson.M {
	return s.filter
}

func (s filterSpecification) GetSort() bson.D {
	return nil
}

type sortSpecification struct {
	sort bson.D
}

func (s sortSpecification) GetFilter() bson.M {
	return nil
}

func (s sortSpecification) GetSort() bson.D {
	return s.sort
}

type joinSpecification struct {
	specifications []Specification
	operator       string
}

func (s joinSpecification) GetFilter() bson.M {
	filters := make([]bson.M, 0, len(s.specifications))

	for _, spec := range s.specifications {
		if filter := spec.GetFilter(); filter != nil {
			filters = append(filters, filter)
		}
	}

	if len(filters) == 0 {
		return nil
	}

	if len(filters) == 1 {
		return filters[0]
	}

	return bson.M{s.operator: filters}
}

func (s joinSpecification) GetSort() bson.D {
	return nil
}

// And combines multiple specifications with the MongoDB $and operator.
func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		operator:       "$and",
	}
}

// Or combines multiple specifications with the MongoDB $or operator.
func Or(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		operator:       "$or",
	}
}

type notSpecification struct {
	Specification
}

func (s notSpecification) GetFilter() bson.M {
	if filter := s.Specification.GetFilter(); filter != nil {
		return bson.M{"$nor": []bson.M{filter}}
	}
	return nil
}

func (s notSpecification) GetSort() bson.D {
	return nil
}

// Not negates a specification using the MongoDB $nor operator.
func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

// Equal creates a specification for equality comparison in MongoDB.
func Equal[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: value},
	}
}

// GreaterThan creates a specification for greater-than comparison using MongoDB $gt operator.
func GreaterThan[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$gt": value}},
	}
}

// GreaterOrEqual creates a specification for greater-than-or-equal comparison using MongoDB $gte operator.
func GreaterOrEqual[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$gte": value}},
	}
}

// LessThan creates a specification for less-than comparison using MongoDB $lt operator.
func LessThan[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$lt": value}},
	}
}

// LessOrEqual creates a specification for less-than-or-equal comparison using MongoDB $lte operator.
func LessOrEqual[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$lte": value}},
	}
}

// In creates a specification for checking if a field value is in a list using MongoDB $in operator.
func In[T any](field string, values []T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$in": values}},
	}
}

// Like creates a specification for case-sensitive pattern matching using MongoDB $regex operator.
func Like(field string, pattern string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$regex": pattern, "$options": ""}},
	}
}

// Ilike creates a specification for case-insensitive pattern matching using MongoDB $regex operator with i option.
func Ilike(field string, pattern string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$regex": pattern, "$options": "i"}},
	}
}

// IsNull creates a specification for checking if a field is null in MongoDB.
func IsNull(field string) Specification {
	return filterSpecification{
		filter: bson.M{field: nil},
	}
}

// IsNotNull creates a specification for checking if a field is not null using MongoDB $ne operator.
func IsNotNull(field string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$ne": nil}},
	}
}

// OrderBy creates a specification for sorting results by a field and direction in MongoDB.
func OrderBy(field string, direction string) Specification {
	order := 1
	if direction == "desc" || direction == "DESC" {
		order = -1
	}

	return sortSpecification{
		sort: bson.D{{Key: field, Value: order}},
	}
}

// Custom creates a specification from a custom MongoDB filter.
func Custom(filter bson.M) Specification {
	return filterSpecification{
		filter: filter,
	}
}
