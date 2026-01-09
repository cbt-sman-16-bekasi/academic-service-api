package pagination

// SearchKey for search queries
type SearchKey struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

// SortOrder for ordering queries
type SortOrder struct {
	Field string `json:"field" form:"field"` // Column name to sort by
	Dir   string `json:"dir" form:"dir"`     // "asc" or "desc"
}

// Request represents pagination request parameters
type Request[T any] struct {
	Page   int       `json:"page" form:"page"`
	Size   int       `json:"size" form:"size"`
	Sort   SortOrder `json:"sort" form:"sort"`
	Filter *T        `json:"filter" form:"filter"`
	Search SearchKey `json:"search" form:"search"`
}

// DefaultSort returns the sort order, defaulting to "id desc" if not specified
func (r *Request[T]) DefaultSort(defaultField string) SortOrder {
	if r.Sort.Field == "" {
		return SortOrder{Field: defaultField, Dir: "desc"}
	}
	// Normalize direction
	dir := r.Sort.Dir
	if dir != "asc" && dir != "desc" {
		dir = "desc"
	}
	return SortOrder{Field: r.Sort.Field, Dir: dir}
}
