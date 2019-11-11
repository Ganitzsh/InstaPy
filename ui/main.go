package main

import "github.com/sirupsen/logrus"

func main() {
	api := newAPI()

	if err := api.Start(); err != nil {
		logrus.Fatalf("Could not start the API: %v", err)
	}
}
