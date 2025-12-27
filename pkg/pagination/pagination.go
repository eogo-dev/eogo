package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Request represents a standard pagination request
type Request struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// GetPage returns the current page, default 1
func (r *Request) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

// GetPageSize returns the page size, default 10, max 100
func (r *Request) GetPageSize() int {
	if r.PageSize <= 0 {
		return 10
	}
	if r.PageSize > 100 {
		return 100
	}
	return r.PageSize
}

// GetOffset calculates the offset for SQL queries
func (r *Request) GetOffset() int {
	return (r.GetPage() - 1) * r.GetPageSize()
}

// FromQuery creates pagination request from query string
func FromQuery(query map[string][]string) *Request {
	req := &Request{}

	if pages, ok := query["page"]; ok && len(pages) > 0 {
		if page, err := strconv.Atoi(pages[0]); err == nil {
			req.Page = page
		}
	}

	if pageSizes, ok := query["page_size"]; ok && len(pageSizes) > 0 {
		if pageSize, err := strconv.Atoi(pageSizes[0]); err == nil {
			req.PageSize = pageSize
		}
	}

	if keywords, ok := query["keyword"]; ok && len(keywords) > 0 {
		req.Keyword = keywords[0]
	}

	return req
}

// FromContext extracts pagination from Gin context
func FromContext(c *gin.Context) *Request {
	req := &Request{}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		req.Page = page
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil {
		req.PageSize = pageSize
	}

	req.Keyword = c.Query("keyword")

	return req
}

// Result represents pagination result metadata
type Result struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	LastPage int   `json:"last_page"`
	From     int   `json:"from"`
	To       int   `json:"to"`
}

// BuildResult creates pagination result metadata
func BuildResult(total int64, page, pageSize int) *Result {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	lastPage := int(math.Ceil(float64(total) / float64(pageSize)))
	if lastPage < 1 {
		lastPage = 1
	}

	from := (page-1)*pageSize + 1
	to := from + pageSize - 1

	if total == 0 {
		from = 0
		to = 0
	}

	return &Result{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		LastPage: lastPage,
		From:     from,
		To:       to,
	}
}

// Paginate executes paginated query on GORM
func Paginate[T any](db *gorm.DB, req *Request) ([]T, *Result, error) {
	var items []T
	var total int64

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Get items
	offset := req.GetOffset()
	if err := db.Offset(offset).Limit(req.GetPageSize()).Find(&items).Error; err != nil {
		return nil, nil, err
	}

	result := BuildResult(total, req.GetPage(), req.GetPageSize())
	return items, result, nil
}

// PaginateFromContext paginates using Gin context
func PaginateFromContext[T any](c *gin.Context, db *gorm.DB) ([]T, *Result, error) {
	req := FromContext(c)
	return Paginate[T](db, req)
}
