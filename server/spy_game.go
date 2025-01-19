package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"slices"
)

type SpyGame struct {
	Room   *Room
	Type   string
	Spies  []*Client
	Prompt string
}

type SpyGameData struct {
	Prompt string `json:"prompt"`
	IsSpy  bool   `json:"isSpy"`
}

func GetRandomPrompt() string {
	prompts := []string{
		"Beach",
		"VirginMedia",
		"Reading",
		"London",
		"Essex",
		"Weatherspoons",
		"England",
		"Germany",
		"Syria",
		"America",
		"Pryzm",
		"Pop World",
		"Cruise Ship",
		"Disney World",
		"Greenland",
		"Brighton",
		"Wales",
		"The Artic Monkeys",
		"Paramore",
		"VirginMediaO2 Reading Office Cafe Brownie",
		"Pepsi",
		"Lana Del Ray",
		"The Smiths",
		"Stella Artois",
		"Bud Light",
		"Strongbow",
		"Inches",
		"The Thames",
		"Mexico",
		"India",
		"Pakistan",
		"Japan",
		"Korea",
		"China",
		"Russia",
		"Donald Trump",
		"Vladmir Putin",
		"Jeremy Corbin",
		"Nigel Farage",
		"Water",
		"Air",
		"Fire",
		"Ice",
		"Jack Robb",
	}
	return prompts[rand.IntN(len(prompts))]
}

func (game *SpyGame) SetupGame() {
	game.Prompt = GetRandomPrompt()
	game.Spies = make([]*Client, 0)
	spyAmt := 1
	for i := 0; i < spyAmt; i++ {
		game.Spies = append(game.Spies, game.Room.Clients[rand.IntN(len(game.Room.Clients))])
	}
}

func (game *SpyGame) StartGame() {
	game.SendGameData()
}

func (game *SpyGame) SendUpdateToClient(client *Client) {
	var sgd SpyGameData
	if slices.Contains(game.Spies, client) {
		sgd.IsSpy = true
	} else {
		sgd.Prompt = game.Prompt
	}

	dat, err := json.Marshal(sgd)
	if err != nil {
		log.Println("Error Making SpyGame Data")
		return
	}
	client.SendPacket(Packet{
		From: "0",
		To:   client.Id,
		Type: "gameData",
		Data: string(dat),
	})
}

func (game *SpyGame) SendGameData() {
	for _, client := range game.Room.Clients {
		game.SendUpdateToClient(client)
	}
}

func (game *SpyGame) GetType() string {
	return "spygame"
}

func (game *SpyGame) HandleClientSwap(oldClient *Client, newClient *Client) {
	if slices.Contains(game.Spies, oldClient) {
		game.Spies[slices.Index(game.Spies, oldClient)] = newClient
	}
}
