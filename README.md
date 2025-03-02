# gofiber-template

My backend template using Golang + Fiber

### Prerequisites

-   Golang 1.24.0 or newer
-   golangci-lint
-   pre-commit

### Rename Package

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

### Pre-commit

Install pre-commit and set up hooks:

```bash
brew install pre-commit
pre-commit install
```

### Commit Lint

Install and initialize commitlint to enforce commit message conventions:

```bash
go install github.com/conventionalcommit/commitlint@latest
commitlint init
```

Example commit message:

```bash
feat: add user authentication
```

### Post-Rename Dependency Cleanup

After renaming the module, ensure dependencies are updated:

```bash
go mod tidy
```
