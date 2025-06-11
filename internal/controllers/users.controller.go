package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chand-magar/SolidBaseGoStructure/internal/dto"
	"github.com/chand-magar/SolidBaseGoStructure/internal/interfaces"
	"github.com/chand-magar/SolidBaseGoStructure/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserController struct {
	Service interfaces.UserService
}

func NewUserController(service interfaces.UserService) *UserController {
	return &UserController{Service: service}
}

func (ctrl *UserController) Create(c *gin.Context) {

	var request dto.RequestDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if err := validator.New().Struct(request); err != nil {
		validationErrs := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, utils.ValidationError(validationErrs))
		return
	}

	now := time.Now().UTC()
	request.CreatedAt = now
	request.UpdatedAt = now
	request.CreatedBy = 1
	request.UpdatedBy = 1

	id, err := ctrl.Service.Create(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": true, "profile_id": id})
}

func (ctrl *UserController) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	search := c.DefaultQuery("search", "")
	status := c.DefaultQuery("status", "")
	sortBy := c.DefaultQuery("sort_by", "user_full_name")
	order := strings.ToUpper(c.DefaultQuery("order", "ASC"))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	params := dto.PaginationParams{
		Page:   page,
		Size:   size,
		Search: search,
		Status: status,
		SortBy: sortBy,
		Order:  order,
	}

	users, totalRecords, totalPages, err := ctrl.Service.GetAll(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        true,
		"data":          users,
		"total_records": totalRecords,
		"total_pages":   totalPages,
	})
}

func (ctrl *UserController) FindOne(c *gin.Context) {

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	user, err := ctrl.Service.FindOne(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   user,
	})
}

func (ctrl *UserController) Update(c *gin.Context) {

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var data dto.RequestDTO
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %s", err.Error())})
		return
	}

	if err := ctrl.Service.Update(id, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "message": "User updated successfully"})
}
