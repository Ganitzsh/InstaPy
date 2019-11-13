package main

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type instagramAccount struct {
	id       string
	label    string
	username string
	password string
}

func newInstagramAccount() *instagramAccount {
	return &instagramAccount{}
}

type user struct {
	id                string
	username          string
	password          string
	instagramAccounts []*instagramAccount
	settings          *botSettings
}

func newUser() *user {
	return &user{}
}

type userRepository struct {
	*mongo.Collection
}

func newUserRepository(col *mongo.Collection) *userRepository {
	return &userRepository{col}
}

func (r *userRepository) save(u *user) (*user, error) {
	if u.id == "" {
		u.id = uuid.New().String()
	}
	if u.instagramAccounts != nil {
		for _, account := range u.instagramAccounts {
			if account.id == "" {
				account.id = uuid.New().String()
			}
		}
	}
	if _, err := r.InsertOne(context.Background(), u); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) findUserByUsername(username string) (*user, error) {
	ret := newUser()

	res := r.FindOne(context.Background(), bson.M{"username": username})
	if err := res.Err(); err != nil {
		return nil, err
	}
	if err := res.Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}
