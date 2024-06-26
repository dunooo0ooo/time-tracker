package repositories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"time-tracker/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// GetUsers retrieves all users from the database
func (ur *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	result := ur.DB.Find(&users)
	return users, result.Error
}

// GetUserEfforts retrieves tasks of a user within a given period
func (ur *UserRepository) GetUserEfforts(userID, startDate, endDate string) ([]models.Task, error) {
	var tasks []models.Task
	result := ur.DB.Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).
		Order("end_time - start_time DESC").
		Find(&tasks)
	return tasks, result.Error
}

// StartTask starts a task for a user
func (ur *UserRepository) StartTask(task *models.Task) error {
	task.StartTime = time.Now()
	return ur.DB.Create(task).Error
}

// StopTask stops a task for a user
func (ur *UserRepository) StopTask(task *models.Task) error {
	task.EndTime = time.Now()
	return ur.DB.Model(task).Updates(map[string]interface{}{"end_time": task.EndTime}).Error
}

// DeleteUser deletes a user from the database
func (ur *UserRepository) DeleteUser(userID string) error {
	return ur.DB.Delete(&models.User{}, userID).Error
}

// UpdateUser updates a user's information in the database
func (ur *UserRepository) UpdateUser(user *models.User) error {
	return ur.DB.Save(user).Error
}

// AddUser adds a new user to the database
func (ur *UserRepository) AddUser(user *models.User) error {
	enrichedUser, err := ur.EnrichUserInfo(user.PassportNumber)
	if err != nil {
		return err
	}
	user.Surname = enrichedUser.Surname
	user.Name = enrichedUser.Name
	user.Patronymic = enrichedUser.Patronymic
	user.Address = enrichedUser.Address
	return ur.DB.Create(user).Error
}

// EnrichUserInfo enriches user information by calling an external API
func (ur *UserRepository) EnrichUserInfo(passportNumber string) (*models.User, error) {
	passportParts := strings.Split(passportNumber, " ")
	if len(passportParts) != 2 {
		return nil, fmt.Errorf("invalid passport number format")
	}
	passportSerie := passportParts[0]
	passportNum := passportParts[1]

	apiURL := fmt.Sprintf("http://external-api-url/info?passportSerie=%s&passportNumber=%s", passportSerie, passportNum)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call external API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned non-200 status code: %d", resp.StatusCode)
	}

	var externalData struct {
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
		Address    string `json:"address"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&externalData); err != nil {
		return nil, fmt.Errorf("failed to decode external API response: %v", err)
	}

	return &models.User{
		Surname:    externalData.Surname,
		Name:       externalData.Name,
		Patronymic: externalData.Patronymic,
		Address:    externalData.Address,
	}, nil
}
