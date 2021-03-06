package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type api struct{}

func newAPI() *api {
	usersCollection := mongoDBClient.Database("instabot").Collection("users")
	uRepo = newUserRepository(usersCollection)

	return &api{}
}

func (a *api) Start() error {
	usersCount, err := uRepo.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return err
	}

	if usersCount == 0 {
		logrus.Println("No existing users found, creating it")
		uRepo.save(&user{
			Username: "instabot",
			Password: "instabot!",
			InstagramAccounts: []*instagramAccount{
				newInstagramAccount().setLabel("Main account"),
			},
			Settings: &botSettings{
				Hashtags: []string{
					"hello",
					"world",
				},
				Comments: []string{
					"nice pic!",
					"I love it!",
				},
				TotalLikes:   10,
				Potency:      potencyModePositive,
				PerUser:      2,
				MinPosts:     10,
				MinFollowers: 500,
				MaxFollowers: 1000,
				MinFollowing: 50,
				MaxFollowing: 100,
			},
		})
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.StaticFile("/index", "./assets/index.html")
	router.StaticFile("/assets", "./assets/log-modal.js")
	router.NoRoute(noRoute)
	router.GET("/script.js", script)

	router.POST("/auth", authenticate)

	auth := router.Group("/api", validateTokenMiddleware)
	auth.GET("/tickets/:ticketID", getTicketStatus)
	auth.GET("/tickets/:ticketID/logs", getTicketLogs)
	auth.GET("/tickets", myTickets)
	auth.GET("/me", me)
	auth.POST("/jobs", runJob)
	auth.POST("/settings", saveSettings)

	return http.ListenAndServe(":8080", router)
}
