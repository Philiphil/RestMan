package router

import (
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
	routerValue, found := r.Configuration[configurationType]
	if len(routeType) == 1 {
		for _, route_ := range r.Routes {
			if route_.RouteType == routeType[0] {
				routeValue, exists := route_.Configuration[configurationType]
				if exists {
					return routeValue, nil
				}
			}
		}
	}
	if !found {
		return routerValue, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
	}

	return routerValue, nil
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
