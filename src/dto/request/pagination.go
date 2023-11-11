package request

type PaginationParam struct {
	Limit int64 `json:"limit" form:"limit" gorm:"-"`
	Page  int64 `json:"page" form:"page" gorm:"-"`
}
