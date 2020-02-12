package main

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type instagramAccount struct {
	ID       string
	Label    string
	Username string
	Password string
}

func (igacc *instagramAccount) setUsername(value string) *instagramAccount {
	igacc.Username = value
	return igacc
}

func (igacc *instagramAccount) setLabel(value string) *instagramAccount {
	igacc.Label = value
	return igacc
}

func newInstagramAccount() *instagramAccount {
	return &instagramAccount{}
}

type user struct {
	ID                string
	Username          string
	Password          string
	InstagramAccounts []*instagramAccount
	Settings          *botSettings
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
	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	if u.InstagramAccounts != nil {
		for _, account := range u.InstagramAccounts {
			if account.ID == "" {
				account.ID = uuid.New().String()
			}
		}
	}

	opts := options.FindOneAndReplace().SetUpsert(true)
	r.FindOneAndReplace(context.Background(), bson.M{"_id": u.ID}, u, opts)

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
