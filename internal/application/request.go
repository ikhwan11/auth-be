package application

import "github.com/google/uuid"

type CreateApplicationRequest struct {
	Name      string     `json:"name" binding:"required"`
	Code      string     `json:"code" binding:"required"`
	URL       string     `json:"url" binding:"required,url"`
	IconID    *uuid.UUID `json:"icon_id"`
	IsDefault bool       `json:"is_default"`
}

type UpdateApplicationRequest struct {
	Name      string     `json:"name" binding:"required"`
	Code      string     `json:"code" binding:"required"`
	URL       string     `json:"url" binding:"required,url"`
	IconID    *uuid.UUID `json:"icon_id"`
	IsDefault bool       `json:"is_default"`
	IsActive  bool       `json:"is_active"`
}
