package main

type potencyMode string

const (
	potencyModePositive potencyMode = "positive"
	potencyModeNegative potencyMode = "negative"
)

type botSettings struct {
	Hashtags     []string
	Comments     []string
	TotalLikes   int
	Potency      potencyMode
	PerUser      int
	MaxFollowers int
	MinFollowers int
	MaxFollowing int
	MinFollowing int
}
