package gormrepository

import (
	"context"
	"slices"
	"strings"
	"sync"
	"unicode"

	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"gorm.io/gorm/schema"
)

// Create implements RestRepository.Create by inserting entities.
func (r *GormRepository[M, E]) Create(ctx context.Context, entities []*E) error {
	return r.BatchInsert(ctx, entities)
}

// Update implements RestRepository.Update by updating entities.
func (r *GormRepository[M, E]) Update(ctx context.Context, entities []*E) error {
	return r.BatchUpdate(ctx, entities)
}

// Read implements RestRepository.Read by finding entities by IDs.
func (r *GormRepository[M, E]) Read(ctx context.Context, ids []entity.ID) ([]*E, error) {
	return r.FindByIDs(ctx, ids)
}

// Delete implements RestRepository.Delete by deleting entities.
func (r *GormRepository[M, E]) Delete(ctx context.Context, entities []*E) error {
	return r.BatchDelete(ctx, entities)
}

// List implements RestRepository.List by finding entities with pagination and sorting.
func (r *GormRepository[M, E]) List(ctx context.Context, limit int, offset int, order map[string]string) ([]E, error) {
	orderSpecification := make([]Specification, 0, len(order))
	for k, v := range order {
		orderSpecification = append(orderSpecification, OrderBy(k, v))
	}
	return r.FindWithLimit(ctx, limit, offset, orderSpecification...)
}

// New implements RestRepository.New by creating a new entity instance.
func (r *GormRepository[M, E]) New() E {
	return r.NewEntity()
}

// Count implements RestRepository.Count by returning the total number of entities.
func (r *GormRepository[M, E]) Count(ctx context.Context) (i int64, err error) {
	model := new(M)
	err = r.getPreWarmDbForSelect(ctx).Model(model).Count(&i).Error
	return
}

func ConvertToSnakeCase(input string) string {
	if input == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(input) + 5)

	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// LoadAssociations implements RestRepository.LoadAssociations by returning a copy that preloads the named associations (snake case).
// The receiver is left untouched, so it stays safe to share between requests.
func (r *GormRepository[M, E]) LoadAssociations(associations []string) (orm.RestRepository[entity.DatabaseModel[E], E], error) {
	if len(associations) == 0 {
		return r, nil
	}
	var model M
	sc, err := schema.Parse(&model, &sync.Map{}, r.db.NamingStrategy)
	if err != nil {
		return nil, err
	}
	scoped := *r
	scoped.preloadAssociations = true
	scoped.assocationsLoaded = true
	scoped.associations = nil
	for _, relation := range sc.Relationships.Relations {
		if slices.Contains(associations, ConvertToSnakeCase(relation.Name)) {
			scoped.associations = append(scoped.associations, relation.Name)
		}
	}
	return &scoped, nil
}
