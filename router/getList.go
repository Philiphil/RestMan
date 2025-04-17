package router

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/route"
)

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
		if order != "ASC" && order != "DESC" || !slices.Contains(SortableFields.Values, field) {
			return nil, errors.ErrBadRequest
		}
		sortParams[field] = order
	}
	if len(sortParams) == 0 {
		sortParams["id"] = defaultSortOrder.Values[0]
	}

	return sortParams, nil
}

func (r *ApiRouter[T]) GetList(c *gin.Context) {
	paginate, err := r.IsPaginationEnabled(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	itemPerPage, err := r.GetItemPerPage(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	sortOrder, err := r.GetSortOrder(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	page, err := r.GetPage(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetConfiguration(configuration.SerializationGroupsType, route.GetList)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}

	var objects []T
	if paginate {
		objects, err = r.Orm.GetPaginatedList(itemPerPage, page, sortOrder)
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
			return
		}
		count, err := r.Orm.Count()
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
			return
		}
		params := map[string]string{}
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		if responseFormat == format.JSONLD {
			c.Render(
				200,
				SerializerRenderer{
					Data:   JsonldCollection(objects, c.Request.URL.String(), page+1, params, int((count+int64(itemPerPage)-1)/int64(itemPerPage))),
					Format: responseFormat,
					Groups: groups.Values,
				},
			)
			return
		}
	} else {
		objects, err = r.Orm.GetAll(sortOrder)
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		}
	}

	c.Render(200,
		SerializerRenderer{
			Data:   objects,
			Format: responseFormat,
			Groups: groups.Values,
		})

}
