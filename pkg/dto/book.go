package dto

type CreateBookRequest struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
}

type BookResponse struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UpdateBookRequest struct {
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}
