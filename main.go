package main

import (
	"fmt"
	"net/http"
	"todo-list/configs"
	"todo-list/models"
	"todo-list/routes"

	"context"
	"log"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "Welcome to app todo",
		})
	})

	models.InitValidator()
	//run database
	configs.ConnectDB()

	// All route
	routes.GetAllRoute(router)

	// get PORT config
	port := os.Getenv("PORT")
	if port == "" {
		port = "6001"
	}
	router.Run("0.0.0.0:" + port)

	// Graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting xxx")
}
