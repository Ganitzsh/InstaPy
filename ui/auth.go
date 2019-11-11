package main

type token struct {
	accessToken  string
	refreshToken string
}

func newToken() *token {
	return &token{}
}

func authenticateUser(username, password string) (*token, error) {
	return nil, nil
}

func verifyToken(t *token) error {
	return nil
}
