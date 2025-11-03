package router

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/route"
)

// This function return either the router wide configuration or the route specific configuration
// If the routeType is not provided, it will return the router wide configuration
// If the routeType is provided, it will return the route specific configuration
// error is returned if the configuration is not found
// by default error should always be nil if you use NewApiRouter
func (r *ApiRouter[T]) GetConfiguration(configurationType configuration.ConfigurationType, routeType ...route.RouteType) (configuration.Configuration, error) {
	if len(routeType) == 1 {
		if routeValue, err := r.GetRouteWideConfiguration(configurationType, routeType[0]); err == nil {
			return routeValue, nil
		}
	}
	return r.GetRouterWideConfiguration(configurationType)
}

func (r *ApiRouter[T]) GetRouterWideConfiguration(configurationType configuration.ConfigurationType) (configuration.Configuration, error) {
	routerValue, found := r.Configuration[configurationType]
	if !found {
		// this should not happen, there should always be a router wide configuration
		// it would mean that the router was not properly initialized
		return configuration.Configuration{}, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
	}
	return routerValue, nil
}

func (r *ApiRouter[T]) GetRouteWideConfiguration(configurationType configuration.ConfigurationType, routeType route.RouteType) (configuration.Configuration, error) {
	configuration, err := r._GetRouteWideConfiguration(configurationType, routeType)
	if err == nil {
		return configuration, nil
	}

	// check if we can fallback to another route configuration
	return r.GetRouteWideConfigurationFallback(configurationType, routeType)
}

func (r *ApiRouter[T]) _GetRouteWideConfiguration(configurationType configuration.ConfigurationType, routeType route.RouteType) (configuration.Configuration, error) {
	routeConfig, exists := r.Routes[routeType]
	if !exists {
		//we trying to get configuration for a route that does not exist
		//should not happen
		return configuration.Configuration{}, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
	}
	routeValue, exists := routeConfig.Configuration[configurationType]
	if exists {
		return routeValue, nil
	}
	return configuration.Configuration{}, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
}

// if a given route-wide configuration is not found, there's case we want to fallback to another route-wide configuration
// if the router's configuration allows it, batch route will fallback to there single entity counterpart
// BatchGet -> Get
// router's configuration might also allow CREATE/UPDATE routes to fallback to READ route configuration
// POST, PUT, PATCH -> GET
// In case of BATCH + WRITE route, both fallback mechanisms can be applied
// does this makes sense?
// BatchPost -> BatchGet -> Post -> Get
func (r *ApiRouter[T]) GetRouteWideConfigurationFallback(configurationType configuration.ConfigurationType, routeType route.RouteType) (configuration.Configuration, error) {
	//not implemented yet

	listOfFallbackableConfigurationType := []configuration.ConfigurationType{
		configuration.OutputSerializationGroupsType,
	}

	if !slices.Contains(listOfFallbackableConfigurationType, configurationType) {
		return configuration.Configuration{}, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
	}

	return configuration.Configuration{}, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
}

func (r *ApiRouter[T]) IsOutputSerializationGroupOverwriteEnabled(c *gin.Context) (bool, error) {
	groupOverwriteConf, err := r.GetConfiguration(configuration.OutputSerializationGroupOverwriteClientControlType, route.GetList)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(c.DefaultQuery(groupOverwriteConf.Values[0], "false"))
}

func (r *ApiRouter[T]) GetOverwriteGroups(c *gin.Context, routeType route.RouteType) ([]string, error) {
	groupsParamConf, err := r.GetConfiguration(configuration.OutputSerializationGroupOverwriteParameterNameType, routeType)
	if err != nil {
		return nil, err
	}
	groupsParam := c.Query(groupsParamConf.Values[0])
	if groupsParam == "" {
		return []string{}, nil
	}
	groups := strings.Split(groupsParam, ",")
	return groups, nil
}

func (r *ApiRouter[T]) GetEffectiveOutputSerializationGroups(c *gin.Context, routeType route.RouteType) ([]string, error) {
	groups, err := r.GetConfiguration(configuration.OutputSerializationGroupsType, routeType)
	if err != nil {
		return nil, err
	}
	effectiveGroups := groups.Values

	enabled, err := r.IsOutputSerializationGroupOverwriteEnabled(c)
	if err != nil {
		return nil, err
	}
	if enabled {
		overwriteGroups, err := r.GetOverwriteGroups(c, routeType)
		if err != nil {
			return nil, err
		}
		if len(overwriteGroups) > 0 {
			effectiveGroups = overwriteGroups
		}
	}
	return effectiveGroups, nil
}
