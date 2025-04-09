package controllers

import (
	"github.com/IamMaheshGurung/NewPortfolio/backend/models"
	"github.com/IamMaheshGurung/NewPortfolio/backend/services"
	"github.com/gofiber/fiber/v2"
)

type AdminSkillController struct {
	skillService *services.SkillService
}

func NewAdminSkillController(ss *services.SkillService) *AdminSkillController {
	return &AdminSkillController{skillService: ss}
}

func (ac *AdminSkillController) Create(c *fiber.Ctx) error {
	var skill models.Skill
	if err := c.BodyParser(&skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ac.skillService.Create(&skill); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create skill"})
	}
	return c.Status(201).JSON(skill)
}

func (ac *AdminSkillController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	var skill models.Skill
	if err := c.BodyParser(&skill); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	skill.ID = uint(id)
	if err := ac.skillService.Update(&skill); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update skill"})
	}
	return c.JSON(skill)
}

func (ac *AdminSkillController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	if err := ac.skillService.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete skill"})
	}
	return c.SendStatus(204)
}

func (ac *AdminSkillController) UpdateProficiency(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid skill ID"})
	}

	var body struct {
		Proficiency int `json:"proficiency"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := ac.skillService.UpdateProficiency(uint(id), body.Proficiency); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update proficiency"})
	}
	return c.SendStatus(200)
}
