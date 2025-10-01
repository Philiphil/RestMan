package router_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
)

type Resource struct {
	entity.BaseEntity
	SubResources []SubResource
}

type SubResource struct {
	entity.BaseEntity
	ResourceID uint
}

func (e Resource) GetId() entity.ID {
	return e.Id
}
func (e Resource) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}
func (e SubResource) GetId() entity.ID {
	return e.Id
}
func (e SubResource) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}
func (r Resource) ToEntity() Resource {
	return r
}
func (r Resource) FromEntity(entity Resource) any {
	return entity
}
func (s SubResource) ToEntity() SubResource {
	return s
}
func (s SubResource) FromEntity(entity SubResource) any {
	return entity
}

func TestApiRouter_SubResources(t *testing.T) {
	// Create main resource router
	resourceRouter := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Resource](getDB())),
		route.DefaultApiRoutes(),
	)

	// Create subresource router
	subResourceRouter := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[SubResource](getDB())),
		route.DefaultApiRoutes(),
	)

	// Add subresource to main resource
	resourceRouter.AddSubresource(subResourceRouter)

	// Register routes
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resourceRouter.AllowRoutes(router)

	// Verify routes are registered
	routes := router.Routes()

	// Check that main resource routes exist
	foundMainGet := false
	foundMainList := false

	// Check that subresource routes exist
	foundSubGet := false
	foundSubList := false

	for _, r := range routes {
		if r.Path == "/api/resource/:id" && r.Method == "GET" {
			foundMainGet = true
		}
		if r.Path == "/api/resource" && r.Method == "GET" {
			foundMainList = true
		}
		if r.Path == "/api/resource/:id/sub_resource/:id" && r.Method == "GET" {
			foundSubGet = true
		}
		if r.Path == "/api/resource/:id/sub_resource" && r.Method == "GET" {
			foundSubList = true
		}
	}

	if !foundMainGet {
		t.Error("Expected main resource GET route to be registered")
	}
	if !foundMainList {
		t.Error("Expected main resource LIST route to be registered")
	}
	if !foundSubGet {
		t.Error("Expected subresource GET route to be registered")
	}
	if !foundSubList {
		t.Error("Expected subresource LIST route to be registered")
	}
}

func TestApiRouter_NestedSubResources(t *testing.T) {
	// Create main resource router
	resourceRouter := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Resource](getDB())),
		route.DefaultApiRoutes(),
	)

	// Create subresource router with nested subresource
	subResourceRouter := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[SubResource](getDB())),
		route.DefaultApiRoutes(),
	)

	// Create nested subresource
	nestedSubRouter := NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
	)

	// Build hierarchy
	subResourceRouter.AddSubresource(nestedSubRouter)
	resourceRouter.AddSubresource(subResourceRouter)

	// Register routes
	gin.SetMode(gin.TestMode)
	router := gin.New()
	resourceRouter.AllowRoutes(router)

	// Verify nested route exists
	routes := router.Routes()

	foundNested := false

	for _, r := range routes {
		// The nested route should be: /api/resource/:id/sub_resource/:id/test/:id
		if r.Path == "/api/resource/:id/sub_resource/:id/test/:id" && r.Method == "GET" {
			foundNested = true
			break
		}
	}

	if !foundNested {
		t.Error("Expected nested subresource route to be registered")
	}
}
