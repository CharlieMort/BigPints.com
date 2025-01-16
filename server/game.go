package main

type Game interface {
	SetupGame()
	StartGame()
	GetType() string
}
