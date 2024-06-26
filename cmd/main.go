package main

import (
	"fmt"
	"log"
	"time-tracker/internal/config"
	"time-tracker/internal/controllers"
	"time-tracker/internal/repositories"

	"github.com/gin-gonic/gin"
	migrate "github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(gormPostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run migrations
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("could not get db instance: %v", err)
	}
	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:/internal/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}

	userRepo := repositories.UserRepository{DB: db}
	userController := controllers.UserController{UserRepository: userRepo}

	r := gin.Default()

	r.GET("/users", userController.GetUsers)
	r.GET("/users/:id/efforts", userController.GetUserEfforts)
	r.POST("/users/:id/start", userController.StartTask)
	r.POST("/users/:id/stop", userController.StopTask)
	r.DELETE("/users/:id", userController.DeleteUser)
	r.PUT("/users/:id", userController.UpdateUser)
	r.POST("/users", userController.AddUser)

	// Swagger documentation route
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()
}
