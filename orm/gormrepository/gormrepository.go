package gormrepository

//this package is the layer that interacts with the database
//its a gorm wraper essentially
import (
	"context"
	"sync"

	"github.com/philiphil/restman/orm/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// NewRepository creates a new GormRepository instance with the provided database connection.
func NewRepository[M entity.DatabaseModel[E], E entity.Entity](db *gorm.DB) *GormRepository[M, E] {
	return &GormRepository[M, E]{
		db: db,
	}
}

// GormRepository is a GORM-based implementation of the RestRepository interface.
type GormRepository[M entity.DatabaseModel[E], E entity.Entity] struct {
	db                 *gorm.DB
	assocationsLoaded  bool
	preloadAssocations bool
	associations       []string
}

// EnablePreloadAssociations enables automatic preloading of entity associations.
func (r *GormRepository[M, E]) EnablePreloadAssociations() *GormRepository[M, E] {
	r.preloadAssocations = true
	return r
}
// DisablePreloadAssociations disables automatic preloading of entity associations.
func (r *GormRepository[M, E]) DisablePreloadAssociations() *GormRepository[M, E] {
	r.preloadAssocations = false
	return r
}

// SetPreloadAssociations sets whether associations should be preloaded.
func (r *GormRepository[M, E]) SetPreloadAssociations(association bool) *GormRepository[M, E] {
	if association {
		r.EnablePreloadAssociations()
	} else {
		r.DisablePreloadAssociations()
	}
	return r
}

func (r *GormRepository[M, E]) setAssociations(model *M) *GormRepository[M, E] {
	sc, _ := schema.Parse(model, &sync.Map{}, r.db.NamingStrategy)
	for _, i := range sc.Relationships.Many2Many {
		r.associations = append(r.associations, i.Name)
	}
	r.assocationsLoaded = true
	return r
}

// Insert creates a new entity in the database.
func (r *GormRepository[M, E]) Insert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

// DeleteByID removes an entity from the database by its ID.
func (r *GormRepository[M, E]) DeleteByID(ctx context.Context, id entity.ID) error {
	var start M
	err := r.db.WithContext(ctx).Delete(&start, &id).Error
	if err != nil {
		return err
	}

	return nil
}

// Upsert creates or updates an entity in the database.
func (r *GormRepository[M, E]) Upsert(ctx context.Context, entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.WithContext(ctx).Save(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()
	return nil
}

// FindByID retrieves an entity by its ID.
func (r *GormRepository[M, E]) FindByID(ctx context.Context, id entity.ID) (E, error) {
	var model M

	err := r.getPreWarmDbForSelect(ctx).First(&model, id).Error
	if err != nil {
		return *new(E), err
	}

	return model.ToEntity(), nil
}

// Find retrieves entities matching the provided specifications.
func (r *GormRepository[M, E]) Find(ctx context.Context, specifications ...Specification) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1, specifications...)
}

func (r *GormRepository[M, E]) getPreWarmDbForSelect(ctx context.Context, specification ...Specification) *gorm.DB {
	dbPrewarm := r.db.WithContext(ctx)

	for _, s := range specification {
		switch spec := s.(type) {
		case orderSpecification:
			dbPrewarm = dbPrewarm.Order(spec.GetQuery())
		default:
			dbPrewarm = dbPrewarm.Where(s.GetQuery(), s.GetValues()...)
		}
	}

	if r.preloadAssocations {
		if !r.assocationsLoaded {
			var start M
			r.setAssociations(&start)
		}
		for _, association := range r.associations {
			dbPrewarm = dbPrewarm.Preload(association)
		}
	}
	return dbPrewarm
}

// FindWithLimit retrieves entities matching the provided specifications with pagination.
func (r *GormRepository[M, E]) FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]E, error) {
	var models []M

	dbPrewarm := r.getPreWarmDbForSelect(ctx, specifications...)

	err := dbPrewarm.Limit(limit).Offset(offset).Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make([]E, 0, len(models))
	for _, row := range models {
		result = append(result, row.ToEntity())
	}

	return result, nil
}

// FindAll retrieves all entities matching the provided specifications.
func (r *GormRepository[M, E]) FindAll(ctx context.Context, specification ...Specification) ([]E, error) {
	return r.FindWithLimit(ctx, -1, -1, specification...)
}

// FindByIDs retrieves multiple entities by their IDs.
func (r *GormRepository[M, E]) FindByIDs(ctx context.Context, ids []entity.ID) ([]*E, error) {
	var models []M
	err := r.db.WithContext(ctx).Find(&models, ids).Error
	if err != nil {
		return nil, err
	}
	result := make([]*E, 0, len(models))
	for _, row := range models {
		entity := row.ToEntity()
		result = append(result, &entity)
	}

	return result, nil
}

// DeleteByIDs removes multiple entities by their IDs.
func (r *GormRepository[M, E]) DeleteByIDs(ctx context.Context, ids []entity.ID) error {
	var start M
	err := r.db.WithContext(ctx).Delete(&start, ids).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchDelete removes multiple entities in a single operation.
func (r *GormRepository[M, E]) BatchDelete(ctx context.Context, entities []*E) error {
	ids := make([]entity.ID, 0, len(entities))
	for _, entity := range entities {
		ids = append(ids, (*entity).GetId())
	}
	return r.DeleteByIDs(ctx, ids)
}

// BatchUpdate updates multiple entities in a transaction.
func (r *GormRepository[M, E]) BatchUpdate(ctx context.Context, entities []*E) error {
	var models []M
	for _, entity := range entities {
		var start M
		model := start.FromEntity(*entity).(M)
		models = append(models, model)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&models).Error; err != nil {
			return err
		}
		for i := range entities {
			*entities[i] = models[i].ToEntity()
		}
		return nil
	})
}

// BatchInsert creates multiple entities in a transaction.
func (r *GormRepository[M, E]) BatchInsert(ctx context.Context, entities []*E) error {
	var models []M
	for _, entity := range entities {
		var start M
		model := start.FromEntity(*entity).(M)
		models = append(models, model)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models).Error; err != nil {
			return err
		}
		for i := range entities {
			*entities[i] = models[i].ToEntity()
		}
		return nil
	})
}

// GetDB returns the underlying GORM database connection.
func (r GormRepository[M, E]) GetDB() *gorm.DB {
	return r.db
}

// NewEntity creates a new empty entity instance.
func (r GormRepository[M, E]) NewEntity() E {
	var entity E
	return entity
}
