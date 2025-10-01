# Configuration Package

The Configuration package provides the default behavior for ApiRouter by defining:
- `ConfigurationStruct` - Main configuration structure
- `ConfigurationType` - Different configuration types

## Purpose

When processing requests (e.g., `GET /api/item`), ApiRouter determines how to respond (e.g., with `GetList`) based on configuration settings for:
- Sorting
- Pagination
- Filtering
- And more...

## Naming Convention

Configuration types follow the naming pattern: `[Feature][OptionalModifier]Type`

### Optional Modifiers
There are primarily two types of optional modifiers:
- `ClientControl`: Determines if clients can override settings via query parameters
- `ParameterName`: Defines the query parameter name used when client control is enabled

### Examples:
- `PaginationType`: Controls if pagination is enabled
- `PaginationClientControlType`: Determines if clients can override pagination settings via query parameters
- `PaginationClientParameterName`: Defines the query parameter name if client control is enabled

## Configuration Hierarchy

- Router: Contains a complete list of configurations (provided by `DefaultConfiguration()`)
- Routes: Can override specific router configurations as needed