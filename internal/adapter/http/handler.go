package http

import (
	"MemeCraft/internal/adapter/http/dto"
	"MemeCraft/internal/port"
	"MemeCraft/internal/service/meme"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"strings"
)

type Handler struct {
	memeGenerator   *meme.Generator
	storageProvider port.StorageProvider
}

func (h *Handler) GenerateMeme(c *fiber.Ctx) error {
	presetId := c.Params("preset_id")
	if presetId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "preset not found",
		})
	}

	_, err := h.memeGenerator.GetPresetById(presetId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var payload dto.CreateMemeRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid payload",
		})
	}

	imageOverlay, err := DownloadImageAsBytes(payload.Overlay)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	result, err := h.memeGenerator.Generate(&meme.Config{
		PresetId:   presetId,
		ResizeMode: payload.ResizeMode,
		Overlay:    imageOverlay,
		Text:       payload.Text,
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *Handler) GetAllPreset(c *fiber.Ctx) error {
	presets := h.memeGenerator.GetAllPreset()
	return c.JSON(presets)
}

func (h *Handler) GetPresetById(c *fiber.Ctx) error {
	presetId := c.Params("preset_id")
	p, err := h.memeGenerator.GetPresetById(presetId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(p)
}

const maxUploadSize = 3 * 1024 * 1024 // 3MB
func (h *Handler) UploadFile(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file not found",
		})
	}

	if fileHeader.Size > maxUploadSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file too large (max 3MB)",
		})
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "only JPG and PNG are allowed",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to open file",
		})
	}
	defer file.Close()

	ctx := context.Background()
	result, err := h.storageProvider.Upload(ctx, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": fmt.Sprintf("upload failed: %v", err),
		})
	}

	return c.JSON(result)
}

func NewHandler(memeGenerator *meme.Generator, storageProvider port.StorageProvider) *Handler {
	return &Handler{
		memeGenerator:   memeGenerator,
		storageProvider: storageProvider,
	}
}
