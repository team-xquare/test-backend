package main

import (
	"log"

	"github.com/team-xquare/deployment-platform/internal/app/auth"
	"github.com/team-xquare/deployment-platform/internal/app/github"
	"github.com/team-xquare/deployment-platform/internal/app/project"
	"github.com/team-xquare/deployment-platform/internal/app/user"
	"github.com/team-xquare/deployment-platform/internal/pkg/config"
	"github.com/team-xquare/deployment-platform/internal/pkg/db/mysql"
	"github.com/team-xquare/deployment-platform/internal/pkg/db/redis"
	"github.com/team-xquare/deployment-platform/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	redisClient, err := redis.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	mysqlDB, err := mysql.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer mysqlDB.Close()

	authRepo := redis.NewAuthRepository(redisClient)
	userRepo := mysql.NewUserRepository(mysqlDB)
	projectRepo := mysql.NewProjectRepository(mysqlDB)
	githubRepo := mysql.NewGitHubRepository(mysqlDB)

	authService := auth.NewService(authRepo, userRepo)
	userService := user.NewService(userRepo)
	projectService := project.NewService(projectRepo, githubRepo)
	githubService := github.NewService(githubRepo)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	projectHandler := project.NewHandler(projectService)
	githubHandler := github.NewHandler(githubService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())

	api := router.Group("/api/v1")
	{
		authHandler.RegisterRoutes(api)
		userHandler.RegisterRoutes(api)
		projectHandler.RegisterRoutes(api)
		githubHandler.RegisterRoutes(api)
	}

	log.Printf("Starting server on port %s", config.AppConfig.AppPort)
	if err := router.Run(":" + config.AppConfig.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
