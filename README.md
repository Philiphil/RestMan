# RestMan

> A declarative REST API framework for Go that automatically generates routes from structs

RestMan takes your Go structs and creates fully functional REST APIs with minimal boilerplate. Inspired by Symfony's API Platform, built on top of Gin, and designed with Go generics for type safety.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Features

- 🚀 **Zero Boilerplate** - Full REST API from a single struct
- 🎯 **Type-Safe Generics** - Compile-time type checking with Go 1.23+ generics
- 🔄 **Multi-Format Support** - JSON, JSON-LD (Hydra), XML, CSV, MessagePack
- 🔒 **Security First** - Built-in firewall and fine-grained authorization
- 📦 **Multiple ORMs** - GORM and MongoDB out of the box, extensible for others
- 🎭 **Serialization Groups** - Control field visibility per context
- 🌳 **Nested Resources** - Unlimited subresource nesting
- ⚡ **Batch Operations** - Efficient bulk create/update/delete
- 📄 **Pagination** - Configurable pagination with Hydra metadata
- 🔍 **Sorting** - Multi-field sorting with client control
- 💾 **HTTP & Redis Caching** - Cache-Control headers and Redis cache library

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Examples](#examples)
- [Configuration](#configuration)
- [Security](#security)
- [Advanced Usage](#advanced-usage)
- [Contributing](#contributing)
- [Roadmap](#roadmap)

## Installation

```bash
go get github.com/philiphil/restman
```

**Requirements:**
- Go 1.23 or higher
- A database (SQLite, PostgreSQL, MySQL via GORM, or MongoDB)

## Quick Start

### 1. Define Your Entity

```go
package main

import (
    "github.com/philiphil/restman/orm/entity"
)

type Book struct {
    entity.BaseEntity
    Title       string `json:"title" groups:"read,write"`
    Author      string `json:"author" groups:"read,write"`
    ISBN        string `json:"isbn" groups:"read"`
    PublishedAt string `json:"published_at" groups:"read"`
}

func (b Book) GetId() entity.ID { return b.Id }
func (b Book) SetId(id any) entity.Entity {
    b.Id = entity.CastId(id)
    return b
}
```

### 2. Create the API

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/philiphil/restman/orm"
    "github.com/philiphil/restman/orm/gormrepository"
    "github.com/philiphil/restman/route"
    "github.com/philiphil/restman/router"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    db, _ := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
    db.AutoMigrate(&Book{})

    r := gin.Default()

    bookRouter := router.NewApiRouter(
        *orm.NewORM(gormrepository.NewRepository[Book](db)),
        route.DefaultApiRoutes(),
    )

    bookRouter.AllowRoutes(r)

    r.Run(":8080")
}
```

### 3. Use Your API

```bash
# Create a book
curl -X POST http://localhost:8080/api/book \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Programming Language","author":"Alan Donovan"}'

# Get all books
curl http://localhost:8080/api/book

# Get a specific book
curl http://localhost:8080/api/book/1

# Update a book
curl -X PUT http://localhost:8080/api/book/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","author":"Alan Donovan"}'

# Delete a book
curl -X DELETE http://localhost:8080/api/book/1
```

**That's it!** You now have a full REST API with:
- GET `/api/book` - List all books (paginated)
- GET `/api/book/:id` - Get a specific book
- POST `/api/book` - Create a book
- PUT `/api/book/:id` - Full update
- PATCH `/api/book/:id` - Partial update
- DELETE `/api/book/:id` - Delete a book
- HEAD `/api/book/:id` - Check existence
- OPTIONS `/api/book` - Available methods

## Core Concepts

### Entity Interface

Every entity must implement the `entity.Entity` interface:

```go
type Entity interface {
    GetId() ID
    SetId(any) Entity
}
```

Use `entity.BaseEntity` to get this for free, along with `CreatedAt`, `UpdatedAt`, and `DeletedAt`.

### Serialization Groups

Control field visibility using the `groups` tag:

```go
type User struct {
    entity.BaseEntity
    Email    string `json:"email" groups:"read,write"`
    Password string `json:"password" groups:"write"` // Only for input
    Token    string `json:"token" groups:"admin"`    // Only for admins
}
```

Groups are applied automatically:
- **POST/PUT/PATCH**: Uses `write` group
- **GET**: Uses `read` group
- Custom groups can be configured per route

### Repository Pattern

RestMan uses a repository abstraction, allowing you to swap databases easily:

```go
// GORM (SQL databases)
gormRepo := gormrepository.NewRepository[Book](db)

// MongoDB
mongoRepo := mongorepository.NewRepository[Book](collection)

// Custom implementation
type MyRepo struct{}
func (r MyRepo) FindAll(ctx context.Context) ([]Book, error) { ... }
```

### Multi-Format Support

RestMan automatically negotiates content type based on the `Accept` header:

```bash
# JSON (default)
curl http://localhost:8080/api/book

# XML
curl -H "Accept: text/xml" http://localhost:8080/api/book

# CSV
curl -H "Accept: application/csv" http://localhost:8080/api/book

# MessagePack
curl -H "Accept: application/msgpack" http://localhost:8080/api/book

# JSON-LD with Hydra pagination
curl -H "Accept: application/ld+json" http://localhost:8080/api/book
```

## Examples

See the [example/](example/) directory for complete working examples:

- **[basic_router_test.go](example/basic_router_test.go)** - Minimal setup
- **[basic_firewall_test.go](example/basic_firewall_test.go)** - Authentication and authorization
- **[model_entity_separation_test.go](example/model_entity_separation_test.go)** - Separating database models from API entities
- **[router_conf_test.go](example/router_conf_test.go)** - Advanced configuration
- **[subresources_test.go](example/subresources_test.go)** - Nested resources
- **[batch_operations_test.go](example/batch_operations_test.go)** - Bulk operations
- **[custom_serialization_test.go](example/custom_serialization_test.go)** - Group-based serialization
- **[pagination_sorting_test.go](example/pagination_sorting_test.go)** - Pagination and sorting configuration

## Configuration

### Router-Level Configuration

Set defaults for all routes:

```go
import "github.com/philiphil/restman/configuration"

bookRouter := router.NewApiRouter(
    *orm.NewORM(gormrepository.NewRepository[Book](db)),
    route.DefaultApiRoutes(),
)

config := configuration.DefaultRouterConfiguration().
    ItemPerPage(50).
    MaxItemPerPage(500).
    RoutePrefix("v1/books").
    RouteName("library").
    AllowClientPagination(true).
    AllowClientSorting(true).
    DefaultSortOrder(configuration.SortAsc, "title").
    NetworkCachingPolicy(configuration.NetworkCachingPolicy{MaxAge: 3600})

bookRouter.Configure(config)
```

### Route-Level Configuration

Override settings for specific operations:

```go
routes := route.DefaultApiRoutes()

getConfig := configuration.DefaultRouteConfiguration().
    InputSerializationGroups("read", "public").
    ItemPerPage(100)

routes.Get.Configure(getConfig)

postConfig := configuration.DefaultRouteConfiguration().
   InputSerializationGroups("write")

routes.Post.Configure(postConfig)

bookRouter := router.NewApiRouter(orm, routes)
```

### Pagination

```bash
# Default pagination
GET /api/book?page=2

# Custom items per page (if allowed)
GET /api/book?page=1&itemsPerPage=50
```

### Sorting

```bash
# Sort by title ascending
GET /api/book?order[title]=asc

# Multiple field sorting
GET /api/book?order[publishedAt]=desc&order[title]=asc
```

## Security

### Authentication with Firewalls

```go
import (
    "github.com/philiphil/restman/security"
)

type MyFirewall struct{}

func (f MyFirewall) ExtractUser(c *gin.Context) (any, *errors.ApiError) {
    token := c.GetHeader("Authorization")
    if token == "" {
        return nil, errors.NewBlockingError(errors.ErrUnauthorized, "Missing token")
    }

    user := validateToken(token) // Your validation logic
    if user == nil {
        return nil, errors.NewBlockingError(errors.ErrUnauthorized, "Invalid token")
    }

    return user, nil
}

bookRouter.SetFirewall(MyFirewall{})
```

### Authorization

```go
// Control read access
bookRouter.SetReadingRights(func(c *gin.Context, book Book, user any) bool {
    if book.Private && book.AuthorID != user.(User).ID {
        return false // User can't read this private book
    }
    return true
})

// Control write access
bookRouter.SetWritingRights(func(c *gin.Context, book Book, user any) bool {
    return book.AuthorID == user.(User).ID // Only author can modify
})
```

## Advanced Usage

### Subresources

Create nested resource routes:

```go
// Creates routes like: /api/author/:id/books/:id
authorRouter := router.NewApiRouter(
    *orm.NewORM(gormrepository.NewRepository[Author](db)),
    route.DefaultApiRoutes(),
)

bookRouter := router.NewApiRouter(
    *orm.NewORM(gormrepository.NewRepository[Book](db)),
    route.DefaultApiRoutes(),
)

authorRouter.AddSubresource(bookRouter)
authorRouter.AllowRoutes(r)
```

### Batch Operations

```bash
# Batch create
POST /api/book/batch
[
  {"title": "Book 1", "author": "Author 1"},
  {"title": "Book 2", "author": "Author 2"}
]

# Batch get by IDs
GET /api/book/batch?ids=1,2,3

# Batch update
PUT /api/book/batch
[
  {"id": 1, "title": "Updated Book 1"},
  {"id": 2, "title": "Updated Book 2"}
]

# Batch delete
DELETE /api/book/batch?ids=1,2,3
```

### Caching

**HTTP Caching (Headers)**

RestMan supports HTTP caching via Cache-Control headers:

```go
import "github.com/philiphil/restman/configuration"

bookRouter := router.NewApiRouter(
    *orm.NewORM(gormrepository.NewRepository[Book](db)),
    route.DefaultApiRoutes(),
    configuration.NetworkCachingPolicy(3600), // Cache for 1 hour
)
```

This automatically sets `Cache-Control: public, max-age=3600` headers on GET requests.

### Model/Entity Separation

Keep your database models separate from API representations:

```go
// Database model (internal)
type BookModel struct {
    ID          uint
    Title       string
    AuthorID    uint
    InternalRef string // Not exposed in API
}

// API entity (external)
type Book struct {
    entity.BaseEntity
    Title  string `json:"title" groups:"read,write"`
    Author Author `json:"author" groups:"read"`
}

func (b BookModel) ToEntity() Book {
    return Book{
        BaseEntity: entity.BaseEntity{Id: b.ID},
        Title:      b.Title,
        Author:     fetchAuthor(b.AuthorID),
    }
}

func (b BookModel) FromEntity(book Book) any {
    return BookModel{
        ID:       book.Id,
        Title:    book.Title,
        AuthorID: book.Author.Id,
    }
}
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suite
go test ./test/router/...
```

## Roadmap

### TODO/ IDEAS
- [ ] Add random configuration to clarify behavior, suggest best practices and allow flexibility  (right now clarifying backup configuration)
- [ ] Filtering implementation
- [ ] UUID compatibility for entity.ID
- [ ] Force lowercase option for JSON keys
- [ ] Automatic Redis caching integration in router
- [ ] GraphQL support
- [ ] Hooks system for lifecycle events
- [ ] Built-in `requireOwnership` for firewall or something
- [ ] Rate limiting middleware (Ai suggestion)
- [ ] Audit login middleware (Ai suggestion)
- [ ] Validation/constraints (Ai suggestion)
- [ ] Finishing redis implementation
- [ ] OpenAPI/Swagger documentation generation
- [ ] Some UI backoffice ?
- [ ] Graphql like PageInfo object after, before, first, last, pageof 


## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

Inspired by:
- [API Platform](https://api-platform.com/) (PHP/Symfony)
