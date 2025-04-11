package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yokeTH/gofiber-template/internal/adaptor/dto"
	"github.com/yokeTH/gofiber-template/internal/usecase/file"
	"github.com/yokeTH/gofiber-template/pkg/apperror"
	"github.com/yokeTH/gofiber-template/pkg/storage"
)

type fileHandler struct {
	fileUseCase file.FileUseCase
	dto         dto.FileDto
}

func NewFileHandler(uc file.FileUseCase, private storage.Storage, public storage.Storage) *fileHandler {
	return &fileHandler{
		fileUseCase: uc,
		dto:         dto.NewFileDto(private, public),
	}
}

func (h *fileHandler) CreatePrivateFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "invalid file")
	}

	fileData, err := h.fileUseCase.CreatePrivateFile(c.Context(), file)
	if err != nil {
		return err
	}

	response := h.dto.ToResponse(*fileData)

	return c.Status(201).JSON(dto.Success(response))
}

func (h *fileHandler) CreatePublicFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequestError(err, "invalid file")
	}

	fileData, err := h.fileUseCase.CreatePublicFile(c.Context(), file)
	if err != nil {
		return err
	}

	response := h.dto.ToResponse(*fileData)

	return c.Status(201).JSON(dto.Success(response))
}

func (h *fileHandler) List(c *fiber.Ctx) error {
	page, limit := extractPaginationControl(c)
	files, last, total, err := h.fileUseCase.List(limit, page)
	if err != nil {
		return err
	}

	response := h.dto.ToResponseList(files)

	return c.Status(200).JSON(dto.SuccessPagination(response, page, last, limit, total))
}

func (h *fileHandler) GetInfo(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return apperror.BadRequestError(err, "invalid id")
	}

	file, err := h.fileUseCase.GetByID(id)
	if err != nil {
		return err
	}

	response := h.dto.ToResponse(*file)

	return c.Status(200).JSON(dto.Success(response))
}
