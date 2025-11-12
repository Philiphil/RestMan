package configuration_test

import (
	"testing"

	. "github.com/philiphil/restman/configuration"
)

func TestConfiguration_Builders(t *testing.T) {
	cfgs := []Configuration{
		NetworkCachingPolicy(60),
		RoutePrefix("api", "v2"),
		RouteName("books"),
		InputSerializationGroups("read", "public"),
		OutputSerializationGroups("read"),
		MaxItemPerPage(500),
		ItemPerPage(50),
		Pagination(true),
		PaginationClientControl(true),
		PaginationParameterName("paginate"),
		PageParameterName("p"),
		ItemPerPageParameterName("limit"),
		BatchIdsName("ids"),
		Sorting(map[string]string{"id": "asc", "createdAt": "desc"}),
		SortingParameterName("order"),
		SortingClientControl(true),
		SortableFields("id", "createdAt"),
		OutputSerializationGroupOverwriteClientControl(true),
		OutputSerializationGroupOverwriteParameterName("overwrite"),
		WriteRouteOutputShouldDefaultToReadOutput(true),
		BatchRouteConfigurationDefaultToSingleRouteConfiguration(true),
	}

	if len(cfgs) == 0 { // impossible, but keeps linter calm if altered
		t.Fatal("no configurations built")
	}

	// Spot check a few specific expectations
	if cfgs[0].Type != NetworkCachingPolicyType || cfgs[0].Values[0] != "60" {
		t.Errorf("unexpected NetworkCachingPolicy: %#v", cfgs[0])
	}
	// Sorting returns slice of alternating key/value pairs – ensure even length >= 4
	foundSorting := false
	for _, c := range cfgs {
		if c.Type == SortingType {
			foundSorting = true
			if len(c.Values)%2 != 0 || len(c.Values) < 4 {
				t.Errorf("sorting values malformed: %#v", c.Values)
			}
		}
	}
	if !foundSorting {
		t.Error("sorting configuration not constructed")
	}
}

// TestDefaultConfiguration_Sanity ensures the default configuration map contains
// a few critical expected entries.
func TestDefaultConfiguration_Sanity(t *testing.T) {
	def := DefaultConfiguration()
	required := []ConfigurationType{
		RoutePrefixType,
		NetworkCachingPolicyType,
		InputSerializationGroupsType,
		OutputSerializationGroupsType,
		PaginationType,
		ItemPerPageType,
		MaxItemPerPageType,
		SortingType,
	}
	for _, r := range required {
		if _, ok := def[r]; !ok {
			t.Errorf("missing default configuration for %v", r)
		}
	}
	if def[RoutePrefixType].Values[0] != "api" {
		t.Errorf("expected default route prefix 'api', got %v", def[RoutePrefixType].Values)
	}
}
