package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
)

type SpyGame struct {
	Room          *Room
	Type          string
	Spies         []*Client
	Prompt        string
	Ready         []string
	QuestionStack []*Client
}

type SpyGameData struct {
	Prompt         string     `json:"prompt"`
	IsSpy          bool       `json:"isSpy"`
	IsReady        bool       `json:"isReady"`
	ReadyString    string     `json:"readyString"`
	QuestionClient ClientJSON `json:"questionClient"`
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

	if slices.Contains(game.Ready, client.Id) {
		sgd.IsReady = true
	}
	sgd.ReadyString = fmt.Sprintf(`%d/%d`, len(game.Ready), len(game.Room.Clients))

	if len(game.QuestionStack) > 0 {
		sgd.QuestionClient = game.QuestionStack[0].ClientJSON
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
		fmt.Println("Swapped ", oldClient.Id, newClient.Id)
		game.Spies[slices.Index(game.Spies, oldClient)] = newClient
	}
}

func (game *SpyGame) ReadyUp(client *Client) {
	if !slices.Contains(game.Ready, client.Id) {
		game.Ready = append(game.Ready, client.Id)
		if len(game.Ready) == len(game.Room.Clients) {
			game.QuestionStack = make([]*Client, len(game.Room.Clients))
			if copy(game.QuestionStack, game.Room.Clients) == 0 {
				fmt.Println("It just didnt copy")
			}
			for i := 0; i < 10; i++ {
				idx1 := rand.IntN(len(game.QuestionStack))
				idx2 := rand.IntN(len(game.QuestionStack))
				tmp := game.QuestionStack[idx1]
				game.QuestionStack[idx1] = game.QuestionStack[idx2]
				game.QuestionStack[idx2] = tmp
			}
		}
		game.SendGameData()
	}
}

func (game *SpyGame) GameTick(client *Client) {
	if len(game.QuestionStack) > 0 {
		game.QuestionStack = RemoveClient(game.QuestionStack, 0)
		game.SendGameData()
	}
}
