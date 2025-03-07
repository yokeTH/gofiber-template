package response

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type PaginationResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
}

func Success[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{Data: data}
}

func SuccessPagination[T any](data []T, currentPage int, lastPage int, limit int, total int) PaginationResponse[T] {
	return PaginationResponse[T]{
		Data: data,
		Pagination: Pagination{
			CurrentPage: currentPage,
			LastPage:    lastPage,
			Limit:       limit,
			Total:       total,
		}}
}
