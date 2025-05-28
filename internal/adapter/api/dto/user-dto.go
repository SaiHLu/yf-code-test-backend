package dto

type CreateUserRequest struct {
	Name            string `json:"name" binding:"required,max=250"`
	Email           string `json:"email" binding:"required,email,max=250"`
	Password        string `json:"password" binding:"required,min=6,max=250"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=250"`
}

type QueryUserRequest struct {
	Search   string `form:"search" binding:"omitempty,max=50"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,max=50"`
	Email    string `json:"email" binding:"omitempty,max=50,email"`
	Password string `json:"password" binding:"omitempty,min=6,max=250"`
}

type UserIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (q *QueryUserRequest) SetDefaultPagination() {
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
