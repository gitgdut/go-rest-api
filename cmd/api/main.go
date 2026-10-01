package main

import (
	"go-api-rest/internal/config"
	"go-api-rest/internal/database"
	"go-api-rest/internal/handlers"
	"go-api-rest/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error
	cfg, err = config.Load()

	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	var pool *pgxpool.Pool
	pool, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer pool.Close()

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(c *gin.Context) {
		//map[string]interface{}
		//map[string]any{}
		c.JSON(200, gin.H{
			"message":  "Todo API is running",
			"status":   "success",
			"database": "connected",
		})
	})

	router.POST("/auth/register", handlers.CreateUserHnadler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	protected := router.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg))

	{
		protected.POST("", handlers.CreateTodoHandler(pool))
		protected.GET("", handlers.GetAllTodosHandler(pool))
		protected.GET("/:id", handlers.GetToDoByIDHandler(pool))
		protected.PUT("/:id", handlers.UpdateToDoHandler(pool))
		protected.DELETE("/:id", handlers.DeleteToDoHandler(pool))
	}
	//Middleware Test Route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())
	router.Run(":" + cfg.Port)

}
