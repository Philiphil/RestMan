package gormrepository_reflection

import (
	"context"
	"reflect"
	"sync"

	"github.com/philiphil/restman/orm/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewRepository(db *gorm.DB, modelType reflect.Type, entityType reflect.Type) *GormRepository {
	return &GormRepository{
		db:         db,
		modelType:  modelType,
		entityType: entityType,
	}
}

type GormRepository struct {
	db                 *gorm.DB
	modelType          reflect.Type // Type of the GORM model (implements entity.DatabaseModel)
	entityType         reflect.Type // Type of the domain entity (implements entity.Entity)
	assocationsLoaded  bool
	preloadAssocations bool
	associations       []string
}

func (r *GormRepository) EnablePreloadAssociations() *GormRepository {
	r.preloadAssocations = true
	return r
}

func (r *GormRepository) DisablePreloadAssociations() *GormRepository {
	r.preloadAssocations = false
	return r
}

func (r *GormRepository) SetPreloadAssociations(association bool) *GormRepository {
	if association {
		r.EnablePreloadAssociations()
	} else {
		r.DisablePreloadAssociations()
	}
	return r
}

func (r *GormRepository) setAssociations(modelInstance any) *GormRepository {
	sc, _ := schema.Parse(modelInstance, &sync.Map{}, r.db.NamingStrategy)
	for _, i := range sc.Relationships.Many2Many {
		r.associations = append(r.associations, i.Name)
	}
	r.assocationsLoaded = true
	return r
}

func (r *GormRepository) Insert(ctx context.Context, entityValue any) error {
	model := reflect.New(r.modelType).Interface()
	dbModel, ok := model.(entity.DatabaseModel[entity.Entity])
	if !ok {
		panic("model does not implement entity.DatabaseModel")
	}

	dbModel.FromEntity(entityValue.(entity.Entity))

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}

	// Update the original entity with the new ID and other fields
	reflect.ValueOf(entityValue).Elem().Set(reflect.ValueOf(dbModel.ToEntity()))
	return nil
}

func (r *GormRepository) DeleteByID(ctx context.Context, id entity.ID) error {
	model := reflect.New(r.modelType).Interface()
	err := r.db.WithContext(ctx).Delete(model, id).Error // GORM handles ID correctly
	if err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) Upsert(ctx context.Context, entityValue any) error {
	model := reflect.New(r.modelType).Interface()
	dbModel, ok := model.(entity.DatabaseModel[entity.Entity])
	if !ok {
		panic("model does not implement entity.DatabaseModel")
	}
	dbModel.FromEntity(entityValue.(entity.Entity))

	err := r.db.WithContext(ctx).Save(model).Error
	if err != nil {
		return err
	}

	reflect.ValueOf(entityValue).Elem().Set(reflect.ValueOf(dbModel.ToEntity()))
	return nil
}

func (r *GormRepository) FindByID(ctx context.Context, id entity.ID) (any, error) {
	model := reflect.New(r.modelType).Interface()

	err := r.getPreWarmDbForSelect(ctx).First(model, id).Error
	if err != nil {
		return nil, err
	}

	dbModel, ok := model.(entity.DatabaseModel[entity.Entity])
	if !ok {
		panic("model does not implement entity.DatabaseModel")
	}
	return dbModel.ToEntity(), nil
}

func (r *GormRepository) Find(ctx context.Context, specifications ...Specification) ([]any, error) {
	return r.FindWithLimit(ctx, -1, -1, specifications...)
}

func (r *GormRepository) getPreWarmDbForSelect(ctx context.Context, specification ...Specification) *gorm.DB {
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
			// Create a new instance of the model type for schema parsing
			modelInstance := reflect.New(r.modelType).Interface()
			r.setAssociations(modelInstance)
		}
		for _, association := range r.associations {
			dbPrewarm = dbPrewarm.Preload(association)
		}
	}
	return dbPrewarm
}

func (r *GormRepository) FindWithLimit(ctx context.Context, limit int, offset int, specifications ...Specification) ([]any, error) {
	// Create a slice of the model type
	modelsSlicePtr := reflect.New(reflect.SliceOf(r.modelType))

	dbPrewarm := r.getPreWarmDbForSelect(ctx, specifications...)

	err := dbPrewarm.Limit(limit).Offset(offset).Find(modelsSlicePtr.Interface()).Error
	if err != nil {
		return nil, err
	}

	modelsSlice := modelsSlicePtr.Elem()
	result := make([]any, modelsSlice.Len())
	for i := 0; i < modelsSlice.Len(); i++ {
		modelInstance := modelsSlice.Index(i).Addr().Interface() // Get pointer to the element
		dbModel, ok := modelInstance.(entity.DatabaseModel[entity.Entity])
		if !ok {
			panic("model does not implement entity.DatabaseModel")
		}
		result[i] = dbModel.ToEntity()
	}

	return result, nil
}

func (r *GormRepository) FindAll(ctx context.Context, specification ...Specification) ([]any, error) {
	return r.FindWithLimit(ctx, -1, -1, specification...)
}

func (r *GormRepository) FindByIDs(ctx context.Context, ids []entity.ID) ([]any, error) {
	modelsSlicePtr := reflect.New(reflect.SliceOf(r.modelType))
	err := r.db.WithContext(ctx).Find(modelsSlicePtr.Interface(), ids).Error // GORM handles slice of IDs
	if err != nil {
		return nil, err
	}

	modelsSlice := modelsSlicePtr.Elem()
	result := make([]any, modelsSlice.Len())
	for i := 0; i < modelsSlice.Len(); i++ {
		modelInstance := modelsSlice.Index(i).Addr().Interface()
		dbModel, ok := modelInstance.(entity.DatabaseModel[entity.Entity])
		if !ok {
			panic("model does not implement entity.DatabaseModel")
		}
		result[i] = dbModel.ToEntity()
	}
	return result, nil
}

func (r *GormRepository) DeleteByIDs(ctx context.Context, ids []entity.ID) error {
	model := reflect.New(r.modelType).Interface()
	err := r.db.WithContext(ctx).Delete(model, ids).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormRepository) BatchDelete(ctx context.Context, entities []any) error {
	ids := make([]entity.ID, len(entities))
	for i, e := range entities {
		ent, ok := e.(entity.Entity)
		if !ok {
			panic("item in entities slice does not implement entity.Entity")
		}
		ids[i] = ent.GetId()
	}
	return r.DeleteByIDs(ctx, ids)
}

func (r *GormRepository) BatchUpdate(ctx context.Context, entities []any) error {
	modelsSlice := reflect.MakeSlice(reflect.SliceOf(r.modelType), len(entities), len(entities))

	for i, e := range entities {
		modelVal := reflect.New(r.modelType).Elem() // Get a Value of the model type
		dbModel, ok := modelVal.Addr().Interface().(entity.DatabaseModel[entity.Entity])
		if !ok {
			panic("model does not implement entity.DatabaseModel")
		}
		dbModel.FromEntity(e.(entity.Entity))
		modelsSlice.Index(i).Set(modelVal)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(modelsSlice.Interface()).Error; err != nil {
			return err
		}
		for i := range entities {
			modelInstance := modelsSlice.Index(i).Addr().Interface()
			dbModel := modelInstance.(entity.DatabaseModel[entity.Entity])
			// Update original entity pointer
			reflect.ValueOf(entities[i]).Elem().Set(reflect.ValueOf(dbModel.ToEntity()))
		}
		return nil
	})
}

func (r *GormRepository) BatchInsert(ctx context.Context, entities []any) error {
	modelsSlice := reflect.MakeSlice(reflect.SliceOf(r.modelType), len(entities), len(entities))

	for i, e := range entities {
		modelVal := reflect.New(r.modelType).Elem()
		dbModel, ok := modelVal.Addr().Interface().(entity.DatabaseModel[entity.Entity])
		if !ok {
			panic("model does not implement entity.DatabaseModel")
		}
		dbModel.FromEntity(e.(entity.Entity))
		modelsSlice.Index(i).Set(modelVal)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(modelsSlice.Interface()).Error; err != nil {
			return err
		}
		for i := range entities {
			modelInstance := modelsSlice.Index(i).Addr().Interface()
			dbModel := modelInstance.(entity.DatabaseModel[entity.Entity])
			reflect.ValueOf(entities[i]).Elem().Set(reflect.ValueOf(dbModel.ToEntity()))
		}
		return nil
	})
}

func (r *GormRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *GormRepository) NewEntity() any {
	// Create a new instance of the entity type (pointer to struct)
	return reflect.New(r.entityType).Interface()
}
