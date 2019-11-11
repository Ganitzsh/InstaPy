package main

type potencyMode string

const (
	potencyModePositive potencyMode = "positive"
	potencyModeNegative potencyMode = "negative"
)

type botSettings struct {
	hashtags     []string
	comments     []string
	totalLikes   int
	potency      potencyMode
	perUser      int
	maxFollowers int
	minFollowers int
	maxFollowing int
	minFollowing int
}
