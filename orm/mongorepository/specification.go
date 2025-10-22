package mongorepository

import (
	"go.mongodb.org/mongo-driver/bson"
)

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

func And(specifications ...Specification) Specification {
	return joinSpecification{
		specifications: specifications,
		operator:       "$and",
	}
}

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

func Not(specification Specification) Specification {
	return notSpecification{
		specification,
	}
}

func Equal[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: value},
	}
}

func GreaterThan[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$gt": value}},
	}
}

func GreaterOrEqual[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$gte": value}},
	}
}

func LessThan[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$lt": value}},
	}
}

func LessOrEqual[T any](field string, value T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$lte": value}},
	}
}

func In[T any](field string, values []T) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$in": values}},
	}
}

func Like(field string, pattern string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$regex": pattern, "$options": ""}},
	}
}

func Ilike(field string, pattern string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$regex": pattern, "$options": "i"}},
	}
}

func IsNull(field string) Specification {
	return filterSpecification{
		filter: bson.M{field: nil},
	}
}

func IsNotNull(field string) Specification {
	return filterSpecification{
		filter: bson.M{field: bson.M{"$ne": nil}},
	}
}

func OrderBy(field string, direction string) Specification {
	order := 1
	if direction == "desc" || direction == "DESC" {
		order = -1
	}

	return sortSpecification{
		sort: bson.D{{Key: field, Value: order}},
	}
}

func Custom(filter bson.M) Specification {
	return filterSpecification{
		filter: filter,
	}
}
