package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type api struct {
}

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
				MinFollowers: 500,
				MaxFollowers: 1000,
				MinFollowing: 50,
				MaxFollowing: 100,
			},
		})
	}
	router := gin.Default()

	router.POST("/auth", authenticate)

	auth := router.Group("/api", validateTokenMiddleware)
	auth.GET("/tickets/:ticketID", getTicketStatus)
	auth.GET("/tickets", myTickets)
	auth.GET("/me", me)
	auth.POST("/jobs", runJob)

	return http.ListenAndServe(":8080", router)
}
