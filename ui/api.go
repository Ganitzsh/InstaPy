package main

type api struct{}

func newAPI() *api {
	return &api{}
}

func (a *api) Start() error {
	return nil
}
