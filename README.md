# gofiber-template

My backend template using Golang + Fiber
## Table of content
- [How to use](#how-to-use)
- [Already Implemented Feature Guide](#already-implemented-feature-guide)

### How to use
#### Prerequisites

-   Golang 1.24.0 or newer
-   golangci-lint
-   pre-commit

#### Rename Package

1. Rename Go module name:

    ```bash
    go mod edit -module YOUR_MODULE_NAME
    ```

    Example:

    ```bash
    go mod edit -module github.com/yourusername/yourprojectname
    ```

2. Find all occurrences of `github.com/yokeTH/gofiber-template/internal` and replace them with `YOUR_MODULE_NAME`:

    ```bash
    find . -type f -name '*.go' -exec sed -i '' 's|github.com/yokeTH/gofiber-template|YOUR_MODULE_NAME|g' {} +
    ```

#### Pre-commit

Install pre-commit and set up hooks:

```bash
brew install pre-commit
pre-commit install
```

#### Commit Lint

Install and initialize commitlint to enforce commit message conventions:

```bash
go install github.com/conventionalcommit/commitlint@latest
commitlint init
```

Example commit message:

```bash
feat: add user authentication
```

#### Post-Rename Dependency Cleanup

After renaming the module, ensure dependencies are updated:

```bash
go mod tidy
```

### Already Implemented Feature Guide
- [App Error](#app-error)
- [DTO](#dto)
- [Pagination](#pagination)

#### App Error
The `apperror` package ensures consistent and structured error handling across repositories, services, and handlers.

##### Example: Using `apperror` in the `CreateBook` Flow

###### **1. Repository Layer (`BookRepository`)**
In the repository layer, if a database operation fails, an `InternalServerError` is returned:

```go
func (r *BookRepository) CreateBook(book *domain.Book) error {
	if err := r.db.Create(book).Error; err != nil {
		return apperror.InternalServerError(err, "failed to create book")
	}
	return nil
}
```

###### **2. Service Layer (BookService)**
The service layer simply propagates the error from the repository:
```go
func (s *BookService) CreateBook(book *domain.Book) error {
	return s.BookRepository.CreateBook(book)
}
```

###### **3. Handler Layer (BookHandler)**
In the handler, errors are processed and returned appropriately:
```go
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	body := new(dto.CreateBookRequest)
	if err := c.BodyParser(body); err != nil {
		return apperror.BadRequestError(err, err.Error())
	}

	book := &domain.Book{
		Author: body.Author,
		Title:  body.Title,
	}

	if err := h.BookService.CreateBook(book); err != nil {
		if apperror.IsAppError(err) {
			return err // Return the structured AppError as-is
		}
		return apperror.InternalServerError(err, "create book service failed")
	}

	res := dto.BookResponse{
		ID:     book.ID,
		Author: book.Author,
		Title:  book.Title,
	}

	return c.Status(fiber.StatusCreated).JSON(dto.Success(res))
}
```

#### DTO
benefit to using with swagger it wrapped with `SuccessResponse` or `ErrorResponse` (if return as apperror)

```go
// for data with pagination
dto.SuccessPagination(data []T, currentPage int, lastPage int, limit int, total int)
// for data without pagination
dto.Success(data T)
```
swagger comment docs example:
```go
// @response 201 {object} dto.SuccessResponse[dto.BookResponse] "Created"
// @response 400 {object} dto.ErrorResponse "Bad Request"
// @response 500 {object} dto.ErrorResponse "Internal Server Error"
```

#### Pagination

Here is an example of how to use `database.Pagination` in your repository:

```go
func (r *BookRepository) GetBooks(limit, page, total, last *int) ([]domain.Book, error) {
    var books []domain.Book

    if err := r.db.Scopes(database.Paginate(domain.Book{}, limit, page, total, last)).Find(&books).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, apperror.NotFoundError(err, "books not found")
        }
        return nil, apperror.InternalServerError(err, "failed to get books")
    }
    return books, nil
}
```