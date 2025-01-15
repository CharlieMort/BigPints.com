package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
)

type Packet struct {
	From string `json:"from"` //ClientID who sent msg - 0 if from server
	To   string `json:"to"`   //Recipent
	Type string `json:"type"` //Type of packet
	Data string `json:"data"` //The actual msg of the data
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]*Room
	broadcast  chan Packet
	register   chan *Client
	unregister chan *Client
}

type Room struct {
	RoomCode string    `json:"roomCode"`
	Host     *Client   `json:"host"`
	Game     Game      `json:"game"`
	Clients  []*Client `json:"clients"`
}

type RoomJSON struct {
	RoomCode string       `json:"roomCode"`
	Host     ClientJSON   `json:"host"`
	Clients  []ClientJSON `json:"clients"`
	GameType string       `json:"gameType"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Packet),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
	}
}

func (h *Hub) GetClientFromID(ID string) *Client {
	for client := range h.clients {
		if client.Id == ID {
			return client
		}
	}
	return nil
}

func (h *Hub) TabOpen(id string) *Client {
	for client, alive := range h.clients {
		if client.TabID == id && !alive {
			return client
		}
	}
	return nil
}

func (h *Hub) SendRoomUpdate(roomCode string) {
	cJSON := make([]ClientJSON, 0)
	for _, client := range h.rooms[roomCode].Clients {
		cJSON = append(cJSON, client.ClientJSON)
	}
	gType := ""
	if h.rooms[roomCode].Game != nil {
		gType = h.rooms[roomCode].Game.GetType()
	}
	fmt.Println(gType)
	dat, err := json.Marshal(RoomJSON{
		RoomCode: roomCode,
		Host:     h.rooms[roomCode].Host.ClientJSON,
		Clients:  cJSON,
		GameType: gType,
	})
	if err != nil {
		log.Println("Error Creating RoomJoinDataPacket")
	}

	for _, client := range h.rooms[roomCode].Clients {
		client.SendPacket(Packet{
			From: "0",
			To:   client.Id,
			Type: "roomData",
			Data: string(dat),
		})
	}
}

func (h *Hub) CreateRoom() string {
	roomCode := GetRandomRoomCode()
	h.rooms[roomCode] = &Room{
		RoomCode: roomCode,
		Host:     nil,
		Game:     nil,
		Clients:  make([]*Client, 0),
	}
	log.Printf("Created The Room:%s\n", roomCode)
	return roomCode
}

func (h *Hub) JoinRoom(client *Client, roomCode string) {
	roomCode = strings.ToLower(roomCode)
	if _, ok := h.rooms[roomCode]; ok {
		if h.rooms[roomCode].Host == nil {
			h.rooms[roomCode].Host = client
		}
		h.rooms[roomCode].Clients = append(h.rooms[roomCode].Clients, client)
		client.RoomCode = roomCode
		client.SendClientJSON()
		h.SendRoomUpdate(roomCode)
		log.Printf("Client:%s Joined the Room:%s", client.Id, roomCode)
	} else {
		log.Printf("Client:%s Failed to join the Room:%s", client.Id, roomCode)
	}
}

func (h *Hub) LeaveRoom(client *Client, roomCode string) {
	roomCode = strings.ToLower(roomCode)
	if _, ok := h.rooms[roomCode]; ok {
		room := h.rooms[roomCode]
		if room.Host == client {
			if len(room.Clients) == 1 {
				delete(h.rooms, roomCode)
			} else {
				room.Clients = RemoveClient(room.Clients, slices.Index(room.Clients, client))
				room.Host = room.Clients[0]
				h.SendRoomUpdate(roomCode)
			}
		}
	} else {
		log.Printf("Client:%s Failed to join the Room:%s", client.Id, roomCode)
	}
}

func (h *Hub) SystemPacket(packet Packet) {
	sysCmd := strings.Split(packet.Data, " ")
	switch sysCmd[0] {
	case "createroom":
		client := h.GetClientFromID(packet.From)
		if client == nil {
			log.Printf("Client:" + packet.From + " Couldnt Be Found")
			return
		}

	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Client Connect ------------------------------------------------")
			client.Print()
			fmt.Println("---------------------------------------------------------------")
			client.SendClientJSON()
		case client := <-h.unregister:
			fmt.Println("Client Disconnect ---------------------------------------------")
			client.Print()
			if client.RoomCode != "" {
				h.LeaveRoom(client, client.RoomCode)
			}
			fmt.Println("---------------------------------------------------------------")
			h.clients[client] = false
		case packet := <-h.broadcast:
			fmt.Println(packet)
			//h.SendPacket(packet)
		}
	}
}
