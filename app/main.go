package main

import (
	"log"
	"os"

	"github.com/frisk038/wordcraft/adapters/repository"
	"github.com/frisk038/wordcraft/app/handlers"
	"github.com/frisk038/wordcraft/business"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	repo, err := repository.NewClient(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("db init fail %s", err)
	}
	bPick := business.NewBPicks(repo)
	bUser := business.NewBUsers(repo)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	rPick := router.Group("/pick")
	rPick.GET("/getLetters", handlers.GetDailyLetters(&bPick))
	rPick.POST("/checkWord", handlers.CheckWordExists(&bPick))

	rUser := router.Group("/user")
	rUser.POST("/register", handlers.InsertUser(&bUser))
	rUser.POST("/score", handlers.InsertScore(&bUser))
	rUser.GET("/leaderboard", handlers.GetLeaderScore(&bUser))

	router.Run(":" + port)
}
