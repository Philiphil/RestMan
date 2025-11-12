package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/route"
)

// SparseFieldRequest represents parsed sparse fields from query parameters.
//
//	{
//	 "id": 123,
//
// "name": "John",
//
//		 "posts": [...]
//		} three syntaxes:
//	  - Dot notation: ?fields=id,name,posts.title
//	  - Separate params: ?fields=id,name&fields.posts=title
//	  - JSON:API style: ?fields[users]=id,name&fields[posts]=title
//
// ResourceFields maps resource name to its field list:
//   - "" (empty key) = root resource fields
//   - "posts" = nested posts fields
//   - "posts.category" = deeply nested fields
type SparseFieldRequest struct {
	ResourceFields map[string][]string
}

// GetSparseFields parses sparse field query parameters and returns a structured request.
// Returns nil if sparse fields are disabled or no fields parameter provided.
//
// Behavior:
//   - If disabled via configuration: returns nil
//   - If ?fields absent or empty: returns nil (no filtering)
//   - If fields present: parses and returns SparseFieldRequest
//
// Parsing supports:
//  1. Dot notation: fields=id,name,posts.title,posts.category.name
//  2. Separate params: fields=id&fields.posts=title&fields.category=name
//  3. JSON:API: fields[users]=id&fields[posts]=title
func (r *ApiRouter[T]) GetSparseFields(c *gin.Context, routeType route.RouteType) (*SparseFieldRequest, error) {
	// Check if sparse fields enabled
	fieldsEnabledConf, err := r.GetConfiguration(configuration.SparseFieldsEnabledType, routeType)
	if err != nil {
		return nil, err
	}
	enabled, err := fieldsEnabledConf.GetValueAsBool()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, nil
	}

	// Get parameter name (default: "fields")
	fieldsParamConf, err := r.GetConfiguration(configuration.SparseFieldsParameterNameType, routeType)
	if err != nil {
		return nil, err
	}
	paramName := fieldsParamConf.Values[0]

	request := &SparseFieldRequest{
		ResourceFields: make(map[string][]string),
	}

	// Parse main ?fields=... parameter (dot notation or simple)
	fieldsParam := c.Query(paramName)
	if fieldsParam != "" {
		r.ParseDotNotationFields(fieldsParam, request)
	}

	// Parse ?fields.resource=... parameters (separate params syntax)
	r.ParseSeparateParamFields(c, paramName, request)

	// Parse ?fields[resource]=... parameters (JSON:API syntax)
	r.ParseJsonApiFields(c, request)

	// If no fields parsed, return nil (treat as no filtering)
	if len(request.ResourceFields) == 0 {
		return nil, nil
	}

	return request, nil
}

// ParseDotNotationFields parses dot notation: fields=id,name,posts.title,posts.category.name
func (r *ApiRouter[T]) ParseDotNotationFields(fieldsParam string, request *SparseFieldRequest) {
	fields := strings.SplitSeq(fieldsParam, ",")

	for field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		if strings.Contains(field, ".") {
			// Nested field: posts.title or posts.category.name
			parts := strings.SplitN(field, ".", 2)
			resourcePath := parts[0]
			fieldName := parts[1]

			// Handle multi-level: posts.category.name
			// Resource path = "posts.category", field = "name"
			for strings.Contains(fieldName, ".") {
				subParts := strings.SplitN(fieldName, ".", 2)
				resourcePath = resourcePath + "." + subParts[0]
				fieldName = subParts[1]
			}

			request.ResourceFields[resourcePath] = append(
				request.ResourceFields[resourcePath],
				fieldName,
			)
		} else {
			// Root field
			request.ResourceFields[""] = append(
				request.ResourceFields[""],
				field,
			)
		}
	}
}

// ParseSeparateParamFields parses: ?fields=id&fields.posts=title&fields.category=name
func (r *ApiRouter[T]) ParseSeparateParamFields(c *gin.Context, paramName string, request *SparseFieldRequest) {
	// Iterate over all query params
	for key, values := range c.Request.URL.Query() {
		// Check if it's a fields.xxx parameter
		if !strings.HasPrefix(key, paramName+".") {
			continue
		}

		// Extract resource name: "fields.posts" -> "posts"
		resourceName := strings.TrimPrefix(key, paramName+".")
		if resourceName == "" {
			continue
		}

		// Parse comma-separated fields
		for _, value := range values {
			fields := strings.Split(value, ",")
			for _, field := range fields {
				field = strings.TrimSpace(field)
				if field != "" {
					request.ResourceFields[resourceName] = append(
						request.ResourceFields[resourceName],
						field,
					)
				}
			}
		}
	}
}

// ParseJsonApiFields parses: ?fields[users]=id,name&fields[posts]=title
func (r *ApiRouter[T]) ParseJsonApiFields(c *gin.Context, request *SparseFieldRequest) {
	// Iterate over all query params
	for key, values := range c.Request.URL.Query() {
		// Check if it's a fields[xxx] parameter
		if !strings.HasPrefix(key, "fields[") || !strings.HasSuffix(key, "]") {
			continue
		}

		// Extract resource name: "fields[posts]" -> "posts"
		resourceName := key[7 : len(key)-1] // Remove "fields[" and "]"
		if resourceName == "" {
			continue
		}

		// Parse comma-separated fields
		for _, value := range values {
			fields := strings.Split(value, ",")
			for _, field := range fields {
				field = strings.TrimSpace(field)
				if field != "" {
					// Map to root resource if it's the main entity
					// For now, treat all non-empty resource names as-is
					request.ResourceFields[resourceName] = append(
						request.ResourceFields[resourceName],
						field,
					)
				}
			}
		}
	}
}
