package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/IzePhanthakarn/go-basic-shop/config"
	"github.com/IzePhanthakarn/go-basic-shop/modules/entities"
	"github.com/IzePhanthakarn/go-basic-shop/modules/files"
	"github.com/IzePhanthakarn/go-basic-shop/modules/files/filesUsecases"
	"github.com/IzePhanthakarn/go-basic-shop/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type filesHandlersErrCode string

const (
	uploadErr filesHandlersErrCode = "files-001"
	deleteErr filesHandlersErrCode = "files-002"
)

type IFilesHandler interface {
	UploadFile(c fiber.Ctx) error
	DeleteFile(c fiber.Ctx) error
}

type filesHandler struct {
	cfg           config.IConfig
	filesUsecases filesUsecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, filesUsecases filesUsecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:           cfg,
		filesUsecases: filesUsecases,
	}
}

// @Summary Upload File
// @Description Upload File
// @Tags Files
// @Accept multipart/form-data
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param files formData file true "Files to upload"
// @Param destination formData string true "Destination path"
// @Success 200 {array} files.FileRes
// @Router /files/upload [post]
func (h *filesHandler) UploadFile(c fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// File ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadErr),
				"invalid file extension",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.StatusBadRequest,
				string(uploadErr),
				fmt.Sprintf("file size must be less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			Extension:   ext,
			FileName:    filename,
		})
	}

	res, err := h.filesUsecases.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

// @Summary Delete File
// @Description Delete File
// @Tags Files
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param request body files.DeleteFileReq true "Files to delete"
// @Success 200 {array} nil
// @Router /files/delete [delete]
func (h *filesHandler) DeleteFile(c fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.Bind().Body(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.StatusBadRequest,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.filesUsecases.DeleteFile(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
