package main

import (
	"log"

	"berita-acara-api/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS — allow Next.js dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Disposition", "Content-Length"},
		AllowCredentials: false,
	}))

	api := r.Group("/api")
	{
		generate := api.Group("/generate")
		{
			generate.POST("/docx", handlers.GenerateDOCX)
			generate.POST("/pdf", handlers.GeneratePDF)
		}
	}

	log.Println("🚀 Berita Acara API listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
