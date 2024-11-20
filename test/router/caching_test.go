package router_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	. "github.com/philiphil/restman/router"
	"github.com/philiphil/restman/security"
)

type securedTest struct {
	entity.BaseEntity
}

func (e securedTest) GetId() entity.ID {
	return e.Id
}

func (e securedTest) SetId(id any) entity.Entity {
	e.Id = entity.CastId(id)
	return e
}
func (t securedTest) ToEntity() securedTest {
	return t
}

func (t securedTest) FromEntity(entity securedTest) any {
	return entity
}

func (e securedTest) GetReadingRights() security.AuthorizationFunction {
	return security.AuthenticationRequired
}

func TestApiRouter_HandleCaching(t *testing.T) {
	repo := orm.NewORM[Test](gormrepository.NewRepository[Test, Test](getDB()))
	test_ := NewApiRouter[Test](
		*repo,
		route.DefaultApiRoutes(),
		configuration.NetworkCachingPolicy(5),
	)
	test_.AllowRoutes(SetupRouter())

	// Créer la requête AVANT d'appeler HandleCaching
	recorder := httptest.NewRecorder()
	testContext, _ := gin.CreateTestContext(recorder)
	testContext.Request = httptest.NewRequest("GET", "/api/test/1", nil)

	test_.HandleCaching(route.Get, testContext)

	cacheControl := recorder.Header().Get("Cache-Control")
	expected := "public, max-age=5"

	if cacheControl != expected {
		t.Errorf("Cache-Control: %s", cacheControl)
	}

	test_.Configuration[configuration.NetworkCachingPolicyType] = configuration.NetworkCachingPolicy(0)

	test_.HandleCaching(route.Get, testContext)
	recorder = httptest.NewRecorder()
	cacheControl = recorder.Header().Get("Cache-Control")
	if cacheControl != "" {
		t.Errorf("Cache-Control: %s", cacheControl)
	}

	//now private
	repo2 := orm.NewORM(gormrepository.NewRepository[securedTest](getDB()))
	test2_ := NewApiRouter(
		*repo2,
		route.DefaultApiRoutes(),
		configuration.NetworkCachingPolicy(15),
	)
	test2_.AllowRoutes(SetupRouter())

	// Créer la requête AVANT d'appeler HandleCaching
	recorder = httptest.NewRecorder()
	testContext, _ = gin.CreateTestContext(recorder)
	testContext.Request = httptest.NewRequest("GET", "/api/test/1", nil)

	test2_.HandleCaching(route.Get, testContext)

	cacheControl = recorder.Header().Get("Cache-Control")
	expected = "private, max-age=15"

	if cacheControl != expected {
		t.Errorf("Cache-Control: %s", cacheControl)
	}

}
