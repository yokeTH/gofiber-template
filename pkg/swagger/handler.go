package swagger

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/yokeTH/gofiber-template/docs"
)

type swagger struct {
}

func New() *swagger {
	return &swagger{}
}

func (s *swagger) Handler() fiber.Handler {
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: docs.SwaggerInfo.ReadDoc(),
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})
}

var (
	instance = New()
	Handler  = instance.Handler()
)
