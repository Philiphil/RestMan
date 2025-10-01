package orm

import (
	"github.com/philiphil/restman/orm/entity"
	"gorm.io/gorm"
)

// RepositoryProvider defines an interface for objects that can provide a RestRepository.
type RepositoryProvider[T entity.Entity] interface {
	Provide() RestRepository[entity.DatabaseModel[T], T]
}

// DefaultRepositoryProvider is a default implementation of RepositoryProvider.
// It uses a dbProvider function to get a database connection and a
// repositoryProvider function to create the RestRepository instance.
type DefaultRepositoryProvider[T entity.Entity] struct {
	dbProvider         func() *gorm.DB
	repositoryProvider func(db *gorm.DB) RestRepository[entity.DatabaseModel[T], T]
}

// NewDefaultRepositoryProvider creates a new DefaultRepositoryProvider.
// It requires:
//   - dbProvider: A function that returns a *gorm.DB database connection.
//   - repositoryProvider: A function that takes a *gorm.DB and returns an instance of
//     RestRepository[entity.DatabaseModel[T], T] for the entity type T.
func NewDefaultRepositoryProvider[T entity.Entity](
	dbProvider func() *gorm.DB,
	repositoryProvider func(db *gorm.DB) RestRepository[entity.DatabaseModel[T], T],
) *DefaultRepositoryProvider[T] {
	if dbProvider == nil {
		panic("DefaultRepositoryProvider: dbProvider cannot be nil")
	}
	if repositoryProvider == nil {
		panic("DefaultRepositoryProvider: repositoryProvider cannot be nil")
	}
	return &DefaultRepositoryProvider[T]{
		dbProvider:         dbProvider,
		repositoryProvider: repositoryProvider,
	}
}

// Provide creates and returns a RestRepository instance.
// It uses the configured dbProvider to get a database connection and
// the repositoryProvider to create the RestRepository instance.
func (p *DefaultRepositoryProvider[T]) Provide() RestRepository[entity.DatabaseModel[T], T] {
	db := p.dbProvider()
	if db == nil {
		// This indicates an issue with the dbProvider.
		panic("DefaultRepositoryProvider: dbProvider returned a nil DB connection")
	}

	// Create the repository using the provided constructor function.
	repo := p.repositoryProvider(db)
	if repo == nil {
		// This indicates an issue with the repositoryProvider.
		panic("DefaultRepositoryProvider: repositoryProvider returned a nil repository instance")
	}
	return repo
}

