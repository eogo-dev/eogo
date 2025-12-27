package resource

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// Resource defines the interface for API resources (transformers)
// Resource provides a way to format API responses consistently.
type Resource interface {
	// ToMap converts the resource to a map
	ToMap() map[string]interface{}
}

// ResourceWithMeta can include additional meta data
type ResourceWithMeta interface {
	Resource
	Meta() map[string]interface{}
}

// ResourceWithLinks can include HATEOAS links
type ResourceWithLinks interface {
	Resource
	Links() map[string]string
}

// BaseResource provides a base implementation
type BaseResource struct {
	data interface{}
}

// NewResource creates a new base resource
func NewResource(data interface{}) *BaseResource {
	return &BaseResource{data: data}
}

// ToMap returns the underlying data as a map
func (r *BaseResource) ToMap() map[string]interface{} {
	result := make(map[string]interface{})

	// Try to convert data to map
	if m, ok := r.data.(map[string]interface{}); ok {
		return m
	}

	// Use JSON encoding/decoding to convert struct to map
	data, err := json.Marshal(r.data)
	if err != nil {
		result["data"] = r.data
		return result
	}

	if err := json.Unmarshal(data, &result); err != nil {
		result["data"] = r.data
	}

	return result
}

// Collection represents a collection of resources
type Collection[T Resource] struct {
	items []T
	meta  map[string]interface{}
	links map[string]string
}

// NewCollection creates a new resource collection
func NewCollection[T Resource](items []T) *Collection[T] {
	return &Collection[T]{
		items: items,
		meta:  make(map[string]interface{}),
		links: make(map[string]string),
	}
}

// WithMeta adds metadata to the collection
func (c *Collection[T]) WithMeta(meta map[string]interface{}) *Collection[T] {
	c.meta = meta
	return c
}

// WithLinks adds links to the collection
func (c *Collection[T]) WithLinks(links map[string]string) *Collection[T] {
	c.links = links
	return c
}

// ToSlice converts the collection to a slice of maps
func (c *Collection[T]) ToSlice() []map[string]interface{} {
	result := make([]map[string]interface{}, len(c.items))
	for i, item := range c.items {
		result[i] = item.ToMap()
	}
	return result
}

// ToResponse converts to a full response with meta and links
func (c *Collection[T]) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"data": c.ToSlice(),
	}

	if len(c.meta) > 0 {
		response["meta"] = c.meta
	}

	if len(c.links) > 0 {
		response["links"] = c.links
	}

	return response
}

// PaginatedCollection wraps paginated data
type PaginatedCollection[T Resource] struct {
	*Collection[T]
	page       int
	perPage    int
	total      int64
	totalPages int
}

// NewPaginatedCollection creates a paginated collection
func NewPaginatedCollection[T Resource](items []T, page, perPage int, total int64) *PaginatedCollection[T] {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &PaginatedCollection[T]{
		Collection: NewCollection(items),
		page:       page,
		perPage:    perPage,
		total:      total,
		totalPages: totalPages,
	}
}

// ToResponse converts to paginated response
func (c *PaginatedCollection[T]) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": c.ToSlice(),
		"meta": map[string]interface{}{
			"current_page": c.page,
			"per_page":     c.perPage,
			"total":        c.total,
			"total_pages":  c.totalPages,
		},
	}
}

// --- Response Helpers ---

// Respond sends a resource as JSON response
func Respond(c *gin.Context, status int, resource Resource) {
	response := map[string]interface{}{
		"data": resource.ToMap(),
	}

	// Add meta if available
	if rm, ok := resource.(ResourceWithMeta); ok {
		if meta := rm.Meta(); len(meta) > 0 {
			response["meta"] = meta
		}
	}

	// Add links if available
	if rl, ok := resource.(ResourceWithLinks); ok {
		if links := rl.Links(); len(links) > 0 {
			response["links"] = links
		}
	}

	c.JSON(status, response)
}

// RespondCollection sends a collection as JSON response
func RespondCollection[T Resource](c *gin.Context, status int, collection *Collection[T]) {
	c.JSON(status, collection.ToResponse())
}

// RespondPaginated sends a paginated collection as JSON response
func RespondPaginated[T Resource](c *gin.Context, status int, collection *PaginatedCollection[T]) {
	c.JSON(status, collection.ToResponse())
}

// --- Anonymous Resource ---

// AnonymousResource allows creating resources inline
type AnonymousResource struct {
	toMap func() map[string]interface{}
	meta  map[string]interface{}
	links map[string]string
}

// NewAnonymousResource creates an anonymous resource
func NewAnonymousResource(toMap func() map[string]interface{}) *AnonymousResource {
	return &AnonymousResource{
		toMap: toMap,
		meta:  make(map[string]interface{}),
		links: make(map[string]string),
	}
}

// ToMap implements Resource
func (r *AnonymousResource) ToMap() map[string]interface{} {
	if r.toMap != nil {
		return r.toMap()
	}
	return make(map[string]interface{})
}

// SetMeta sets metadata
func (r *AnonymousResource) SetMeta(meta map[string]interface{}) *AnonymousResource {
	r.meta = meta
	return r
}

// Meta returns metadata
func (r *AnonymousResource) Meta() map[string]interface{} {
	return r.meta
}

// SetLinks sets links
func (r *AnonymousResource) SetLinks(links map[string]string) *AnonymousResource {
	r.links = links
	return r
}

// Links returns links
func (r *AnonymousResource) Links() map[string]string {
	return r.links
}
