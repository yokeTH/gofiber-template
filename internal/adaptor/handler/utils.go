package handler

import "github.com/gofiber/fiber/v2"

func extractPaginationControl(c *fiber.Ctx) (int, int) {
	page := c.QueryInt("page", 1)
	if page <= 0 {
		page = 1
	}

	limit := c.QueryInt("limit", 10)
	if limit <= 0 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	return page, limit
}
