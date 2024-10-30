package routes

import (
	"social/util"

	"github.com/gofiber/fiber/v2"
)

func CreateImages(c *fiber.Ctx) error {
	bf := util.BlackForest{}
	bf.Init()

	requestId := bf.Request()
	bf.Poll(requestId)

	result := map[string]interface{}{
		"requestId": requestId,
		"status":    200,
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
