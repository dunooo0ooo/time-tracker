package controllers

// @Summary Get all users
// @Description Get all users with pagination and filtering
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Router /users [get]

import (
	"net/http"
	"time"
	"time-tracker/internal/models"
	"time-tracker/internal/repositories"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserRepository repositories.UserRepository
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.UserRepository.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUserEfforts(c *gin.Context) {
	userID := c.Param("id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	efforts, err := uc.UserRepository.GetUserEfforts(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, efforts)
}

func (uc *UserController) StartTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.StartTime = time.Now()

	if err := uc.UserRepository.StartTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (uc *UserController) StopTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.EndTime = time.Now()

	if err := uc.UserRepository.StopTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := uc.UserRepository.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.UserRepository.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) AddUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	enrichedInfo, err := uc.UserRepository.EnrichUserInfo(user.PassportNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Surname = enrichedInfo.Surname
	user.Name = enrichedInfo.Name
	user.Patronymic = enrichedInfo.Patronymic
	user.Address = enrichedInfo.Address

	if err := uc.UserRepository.AddUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
