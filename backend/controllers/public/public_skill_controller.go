package controllers

import (
	"github.com/IamMaheshGurung/NewPortfolio/backend/services"
	"github.com/gofiber/fiber/v2"
)

type PublicSkillController struct {
	skillService *services.SkillService
}

func NewPublicSkillController(ss *services.SkillService) *PublicSkillController {
	return &PublicSkillController{skillService: ss}
}

func (pc *PublicSkillController) ListAll(c *fiber.Ctx) error {
	skills, err := pc.skillService.ListAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch skills"})
	}
	return c.JSON(skills)
}

func (pc *PublicSkillController) ListFeatured(c *fiber.Ctx) error {
	skills, err := pc.skillService.ListFeatured()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch featured skills"})
	}
	return c.JSON(skills)
}

func (pc *PublicSkillController) ListByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	skills, err := pc.skillService.ListByCategory(category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch skills by category"})
	}
	return c.JSON(skills)
}
