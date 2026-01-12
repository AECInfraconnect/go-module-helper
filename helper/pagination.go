package helper

import (
	"math"

	"gorm.io/gorm"
)

const (
	// PSQL_TOTAL_ROW_KEY is the column name used to store the total row count in paginated queries.
	// Use this in SQL queries like: SELECT *, COUNT(*) OVER() as total FROM table
	PSQL_TOTAL_ROW_KEY = "total"
)

// Paginator handles pagination data including current page, items per page, and total counts.
// It provides methods to calculate pagination info and integrate with GORM queries.
type Paginator struct {
	Page            int `json:"page"`       // Current page number (1-indexed)
	Limit           int `json:"limit"`      // Number of items per page
	TotalPages      int `json:"total_page"` // Total number of pages
	TotalEntrySizes int `json:"total_rows"` // Total number of items across all pages
}

// NewPaginator creates a new Paginator with default values.
// Default: Page = 1, Limit = 20
func NewPaginator() Paginator {
	return Paginator{Page: 1, Limit: 20}
}

// NewPaginatorWithParams creates a new Paginator with custom page and limit values.
func NewPaginatorWithParams(page int, limit int) Paginator {
	return Paginator{Page: page, Limit: limit}
}

// SetPaginatorByAllRows sets the total number of rows and calculates total pages.
// This is useful when you already know the total count from a separate query.
func (p *Paginator) SetPaginatorByAllRows(allRows int) {
	p.setTotalEntrySizes(allRows)
	p.setTotalPages()
}

// setTotalEntrySizes sets the total number of entries.
func (p *Paginator) setTotalEntrySizes(allRows int) {
	p.TotalEntrySizes = allRows
}

// setTotalPages calculates and sets the total number of pages based on total rows and per page.
func (p *Paginator) setTotalPages() {
	totalRows := p.TotalEntrySizes
	limit := p.Limit
	totalPage := math.Ceil(float64(totalRows) / float64(limit))
	p.TotalPages = int(totalPage)
}

// SetTotalFromMap sets pagination info from a map that contains a 'total' field.
// Use this when scanning GORM rows into a map[string]any.
//
// Example:
//
//	var result map[string]any
//	db.Raw("SELECT *, COUNT(*) OVER() as total FROM users LIMIT 1").Scan(&result)
//	paginator.SetTotalFromMap(result)
func (p *Paginator) SetTotalFromMap(data map[string]any) {
	if totalVal, ok := data[PSQL_TOTAL_ROW_KEY]; ok {
		var total int
		switch v := totalVal.(type) {
		case int64:
			total = int(v)
		case int32:
			total = int(v)
		case int:
			total = v
		case float64:
			total = int(v)
		}
		if total > 0 {
			p.setTotalEntrySizes(total)
			p.setTotalPages()
		}
	}
}

// SetPaginatorByGORM sets pagination info by counting total rows from GORM query.
// This method performs a COUNT query to get the total number of records.
//
// Example:
//
//	paginator := helper.NewPaginatorWithParams(1, 20)
//	db := gormDB.Where("active = ?", true)
//	if err := paginator.SetPaginatorByGORM(db); err != nil {
//	    return err
//	}
func (p *Paginator) SetPaginatorByGORM(db *gorm.DB) error {
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return err
	}
	p.SetPaginatorByAllRows(int(total))
	return nil
}

// ApplyToGORM applies offset and limit to GORM query based on current page and per_page.
// Returns a new GORM DB instance with pagination applied.
//
// Example:
//
//	paginator := helper.NewPaginatorWithParams(2, 20)
//	var users []User
//	paginator.ApplyToGORM(db.Where("active = ?", true)).Find(&users)
func (p *Paginator) ApplyToGORM(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit)
}

// PaginateGORM performs both count and paginated query in one call.
// It counts total records, applies pagination, and executes the query.
// This is the most convenient method for simple pagination needs.
//
// Example:
//
//	paginator := helper.NewPaginatorWithParams(1, 20)
//	var users []User
//	db := gormDB.Where("active = ?", true)
//	if err := paginator.PaginateGORM(db, &users); err != nil {
//	    return err
//	}
//	* paginator now contains total_pages and total_rows
//	* users contains the paginated results
func (p *Paginator) PaginateGORM(db *gorm.DB, dest any) error {
	// Count total records
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return err
	}
	p.SetPaginatorByAllRows(int(total))

	// Apply pagination and query
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit).Find(dest).Error
}
