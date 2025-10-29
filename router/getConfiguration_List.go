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

//this is for list and batch configurations

// IsPaginationEnabled determines whether pagination should be enabled for the current request based on configuration and query parameters.
func (r *ApiRouter[T]) IsPaginationEnabled(c *gin.Context) (bool, error) {
	paginationConf, err := r.GetConfiguration(configuration.PaginationType, route.GetList)
	if err != nil {
		return false, err
	}
	forcedPaginationConf, err := r.GetConfiguration(configuration.PaginationClientControlType, route.GetList)
	if err != nil {
		return false, err
	}
	clientCanForcePaginationUsingParameter, err := strconv.ParseBool(forcedPaginationConf.Values[0])
	if err != nil {
		return false, err
	}
	basepaginationBool, err := strconv.ParseBool(paginationConf.Values[0])
	if err != nil {
		return false, err
	}
	if !clientCanForcePaginationUsingParameter {
		return basepaginationBool, nil
	}
	forcedParameterConf, err := r.GetConfiguration(configuration.PaginationParameterNameType, route.GetList)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(c.DefaultQuery(forcedParameterConf.Values[0], paginationConf.Values[0]))
}

// GetPage extracts the page number from the request query parameters and returns it as a zero-indexed value.
func (r *ApiRouter[T]) GetPage(c *gin.Context) (int, error) {
	pageParameter, err := r.GetConfiguration(configuration.PageParameterNameType, route.GetList)
	if err != nil {
		return 0, err
	}
	page, err := strconv.Atoi(c.DefaultQuery(pageParameter.Values[0], "1"))
	if err != nil {
		return 0, err
	}
	return page - 1, nil
}

// GetItemPerPage extracts the items per page value from request query parameters, enforcing the configured maximum limit.
func (r *ApiRouter[T]) GetItemPerPage(c *gin.Context) (int, error) {
	defaultItemPerPage, err := r.GetConfiguration(configuration.ItemPerPageType, route.GetList)
	if err != nil {
		return 0, err
	}
	maxItemPerPage, err := r.GetConfiguration(configuration.MaxItemPerPageType, route.GetList)
	if err != nil {
		return 0, err
	}
	itemPerPageParameter, err := r.GetConfiguration(configuration.ItemPerPageParameterNameType, route.GetList)
	if err != nil {
		return 0, err
	}
	itemPerPage, err := strconv.Atoi(c.DefaultQuery(itemPerPageParameter.Values[0], defaultItemPerPage.Values[0]))
	if err != nil {
		return 0, err
	}
	maxItemPerPageValue, err := strconv.Atoi(maxItemPerPage.Values[0])
	if itemPerPage > maxItemPerPageValue {
		itemPerPage = maxItemPerPageValue
	}
	return itemPerPage, err
}

// GetSortOrder extracts sorting parameters from the request, validating against allowed fields and returning a map of field names to sort directions.
func (r *ApiRouter[T]) GetSortOrder(c *gin.Context) (map[string]string, error) {
	sortParams := make(map[string]string)

	sortEnabled, err := r.GetConfiguration(configuration.SortingClientControlType, route.GetList)
	if err != nil {
		return nil, err
	}
	enabled, parseErr := strconv.ParseBool(sortEnabled.Values[0])
	if parseErr != nil {
		return nil, parseErr
	}

	// I should replace this by a default map[string]string
	defaultSortOrder, err := r.GetConfiguration(configuration.SortingType, route.GetList)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(defaultSortOrder.Values); i += 2 {
		if i+1 < len(defaultSortOrder.Values) {
			sortParams[defaultSortOrder.Values[i]] = defaultSortOrder.Values[i+1]
		}
	}

	if !enabled {
		return sortParams, nil
	}

	// get the sort paramter name and allowed fields for sorting
	sortParam, err := r.GetConfiguration(configuration.SortingParameterNameType, route.GetList)
	if err != nil {
		return nil, err
	}
	SortableFields, err := r.GetConfiguration(configuration.SortableFieldsType, route.GetList)
	if err != nil {
		return nil, err
	}

	queryParams := c.QueryMap(sortParam.Values[0])
	for field, order := range queryParams {
		order = strings.ToUpper(order)
		if (order != "ASC" && order != "DESC") || !slices.Contains(SortableFields.Values, field) {
			return nil, errors.ErrBadRequest
		}
		sortParams[field] = order
	}
	if len(sortParams) == 0 {
		sortParams["id"] = defaultSortOrder.Values[0]
	}

	return sortParams, nil
}

// IsBatchGetOrGetList determines whether the request is a BatchGet or GetList operation.
// BatchGet returns a list of objects using given ids, while GetList returns paginated results.
func (r *ApiRouter[T]) IsBatchGetOrGetList(c *gin.Context) route.RouteType {
	//first if BatchGet or GetList is not allowed, it is not a BatchGet or GetList
	if _, ok := r.Routes[route.BatchGet]; !ok {
		return route.GetList
	}
	if _, ok := r.Routes[route.GetList]; !ok {
		return route.BatchGet
	}

	ids, _ := r.GetConfiguration(configuration.BatchIdsParameterNameType, route.BatchGet)
	idsParameter := ids.Values[0]
	exists := false
	if _, exists = c.GetQuery(idsParameter); !exists {
		_, exists = c.GetQuery(idsParameter + "[]")
	}
	if exists {
		return route.BatchGet
	}
	return route.GetList
}

// GetListOrBatchGet routes the request to either BatchGet or GetList based on query parameters.
func (r *ApiRouter[T]) GetListOrBatchGet(c *gin.Context) {
	rr := r.IsBatchGetOrGetList(c)
	if rr == route.BatchGet {
		r.BatchGet(c)
	} else {
		r.GetList(c)
	}
}

// GetIds extracts the list of IDs from query parameters for batch operations.
// Supports both array notation (ids[]=1&ids[]=2) and comma-separated (ids=1,2,3).
func (r *ApiRouter[T]) GetIds(c *gin.Context) []string {
	ids, _ := r.GetConfiguration(configuration.BatchIdsParameterNameType, route.BatchGet)
	idsParameter := ids.Values[0]
	exists := false
	if _, exists = c.GetQuery(idsParameter + "[]"); exists {
		return c.QueryArray(idsParameter + "[]")
	}

	idsValues := c.QueryArray(ids.Values[0])
	if len(idsValues) == 1 && len(strings.Split(idsValues[0], ",")) > 1 {
		return strings.Split(idsValues[0], ",")
	}
	return idsValues
}
