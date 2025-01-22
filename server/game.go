package main

type Game interface {
	SetupGame()
	StartGame()
	GetType() string
	SendUpdateToClient(client *Client)
	HandleClientSwap(oldClient *Client, newClient *Client)
	ReadyUp(client *Client)
	SendGameData()
	GameTick(client *Client)
}
