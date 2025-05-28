package dto

type QueryUserLogRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (q *QueryUserLogRequest) SetDefaultPagination() {
	if q.Page < 1 {
		q.Page = 1
	}

	if q.PageSize < 1 {
		q.PageSize = 10
	}

	if q.PageSize > 100 {
		q.PageSize = 100
	}
}
