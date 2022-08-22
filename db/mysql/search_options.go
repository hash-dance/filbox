package mysql

const (
	// Query 条件查询
	Query = "query"
	// Search 模糊搜索
	Search = "search"
)

// Pagination 翻页
type Pagination struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"page_size,omitempty"`
	OrderBy  string `json:"order"`
	Total    int    `json:"total"`
}

// Condition 条件查询
type Condition struct {
	Name  string      `json:"name"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

// QueryCondition query condition
type QueryCondition struct {
	Condition []*Condition `json:"conditions,omitempty"`
}

// PageQueryInfo pagination & QueryCondition
type PageQueryInfo struct {
	*Pagination
	*QueryCondition
	ConditionKey string
}

// FuzzyQuery 模糊搜索条件
type FuzzyQuery struct {
	Key    string   `json:"key"`
	Fields []string `json:"fields"`
}
