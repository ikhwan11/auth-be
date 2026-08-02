package application_icon

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Create godoc
// @Tags application_icons
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Icon name"
// @Param file formData file true "SVG icon"
// @Success 201 {object} ApplicationIcon
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /application-icons [post]
func (h *Handler) Create(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Icon name is required",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "SVG file is required",
		})
		return
	}

	if !strings.HasSuffix(
		strings.ToLower(file.Filename),
		".svg",
	) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Only SVG files are allowed",
		})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to open SVG file",
		})
		return
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read SVG file",
		})
		return
	}

	icon, err := h.service.Create(
		c.Request.Context(),
		name,
		fileData,
		"image/svg+xml",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create application icon",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    icon,
	})
}

// List godoc
// @Tags application_icons
// @Security BearerAuth
// @Produce json
// @Success 200 {object} []ApplicationIcon
// @Failure 500 {object} map[string]interface{}
// @Router /application-icons [get]
func (h *Handler) List(c *gin.Context) {
	icons, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get application icons",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    icons,
	})
}

// Get godoc
// @Tags application_icons
// @Security BearerAuth
// @Produce json
// @Param id path string true "Application Icon ID"
// @Success 200 {object} ApplicationIcon
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /application-icons/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application icon ID",
		})
		return
	}

	icon, err := h.service.FindByID(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrApplicationIconNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Application icon not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get application icon",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    icon,
	})
}

// File godoc
// @Tags application_icons
// @Security BearerAuth
// @Produce image/svg+xml
// @Param id path string true "Application Icon ID"
// @Success 200 {file} binary
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /application-icons/{id}/file [get]
func (h *Handler) File(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid application icon ID",
		})
		return
	}

	icon, err := h.service.GetFile(
		c.Request.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, ErrApplicationIconNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Application icon not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get application icon",
			"error":   err.Error(),
		})
		return
	}

	c.Data(
		http.StatusOK,
		icon.MimeType,
		icon.FileData,
	)
}
