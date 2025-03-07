package dto

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

// Success returns a SuccessResponse containing the provided data.
// It creates a new instance of SuccessResponse[T] with Data set to the value of data, enabling consistent handling of successful responses in JSON.
func Success[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{Data: data}
}

// SuccessPagination constructs a PaginationResponse[T] that encapsulates a data slice along with pagination metadata.
// The pagination metadata comprises the current page, the last available page, the maximum number of items per page (limit),
// and the total number of items, facilitating a standardized structure for paginated API responses.
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
