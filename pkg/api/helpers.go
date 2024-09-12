package api

import (
	"digital-wallet/pkg/errs"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetPageAndLimit(c *fiber.Ctx) (int, int, error) {
	page, limit := 1, 10
	var err error
	if c.Query("page") != "" {
		if page, err = strconv.Atoi(c.Query("page")); err != nil {
			return 0, 0, errs.NewBadRequestError("invalid page", "INVALID_PAGE_PARAM", err)
		}
	}
	if c.Query("limit") != "" {
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			return 0, 0, errs.NewBadRequestError("invalid limit", "INVALID_LIMIT_PARAM", err)
		}
	}
	return page, limit, nil
}
