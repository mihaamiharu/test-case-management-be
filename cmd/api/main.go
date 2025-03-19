package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mihaamiharu/test-case-management-be/internal/api"
	"github.com/mihaamiharu/test-case-management-be/internal/config"
	"github.com/mihaamiharu/test-case-management-be/internal/db"
	"github.com/mihaamiharu/test-case-management-be/internal/repository"
	"github.com/mihaamiharu/test-case-management-be/internal/service"
	"github.com/mihaamiharu/test-case-management-be/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database)
	projectRepo := repository.NewProjectRepository(database)
	projectAccessRepo := repository.NewProjectAccessRepository(database)
	testSuiteRepo := repository.NewTestSuiteRepository(database)
	testCaseRepo := repository.NewTestCaseRepository(database)
	tagRepo := repository.NewTagRepository(database)

	// Initialize services
	authService := service.NewAuthService(userRepo, cfg)
	projectService := services.NewProjectService(projectRepo)
	projectAccessService := services.NewProjectAccessService(projectAccessRepo, projectRepo, userRepo)
	testCaseService := service.NewTestCaseService(testCaseRepo, tagRepo)
	testSuiteService := service.NewTestSuiteService(testSuiteRepo)
	tagService := service.NewTagService(tagRepo)

	// Initialize handlers
	authHandler := api.NewAuthHandler(authService)
	projectHandler := api.NewProjectHandler(projectService, projectAccessService)
	testSuiteHandler := api.NewTestSuiteHandler(testSuiteService)
	testCaseHandler := api.NewTestCaseHandler(testCaseService)
	tagHandler := api.NewTagHandler(tagService)

	// Initialize router
	router := gin.Default()
	api.SetupRouter(router, authHandler, projectHandler, testSuiteHandler, testCaseHandler, tagHandler)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
