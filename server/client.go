package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type ClientJSON struct {
	Id       string `json:"id"`
	RoomCode string `json:"roomCode"`
	Name     string `json:"name"`
	Imguuid  string `json:"imguuid"`
}

type ClientData struct {
	TabID string
	Hub   *Hub
	Conn  *websocket.Conn
	Send  chan Packet
}

type Client struct {
	ClientJSON
	ClientData
}

func (client *Client) BroadcastPacket(packet Packet) {
	client.Hub.broadcast <- packet
}

func (client *Client) SendPacket(packet Packet) {
	log.Println("Sending Packet To Client client.SendPacket")
	client.Send <- packet
}

func (client *Client) GetJSON() string {
	dat, err := json.Marshal(client.ClientJSON)
	if err != nil {
		log.Println("Couldnt Parse Client JSON")
		return ""
	}
	return string(dat)
}

func (client *Client) SendClientJSON() {
	fmt.Println("Sending Client JSON")
	dat := client.GetJSON()
	client.SendPacket(Packet{
		From: "0",
		To:   client.Id,
		Type: "clientData",
		Data: dat,
	})
}

func (client *Client) Print() {
	fmt.Printf("ID:%s\nName:%s\nRoomCode:%s\n", client.Id, client.Name, client.RoomCode)
}

func (client *Client) ReadPackets() {
	defer func() {
		client.Hub.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, packetJson, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		if string(packetJson) == "ping" {
			fmt.Printf("#")
			client.Send <- Packet{Data: "pong"}
			continue
		}
		var packet Packet
		err = json.Unmarshal(packetJson, &packet)
		fmt.Println("Received Packet From Client ---------------------------------------------")
		fmt.Printf("From:%s\nTo:%s\nType:%s\nData:%s\n", packet.From, packet.To, packet.Type, packet.Data)
		fmt.Println("-------------------------------------------------------------------------")
		if err != nil {
			log.Printf("Error ReadPackets(1) %v", err)
		}
		switch packet.Type {
		case "setup":
			client.TabID = packet.Data
			oldClient := client.Hub.TabOpen(client.TabID)
			if oldClient != nil {
				fmt.Println("Client Exists ------------------------------------")
				client.ClientJSON = oldClient.ClientJSON
				if client.RoomCode != "" {
					client.Hub.LeaveRoom(oldClient, client.RoomCode)
					if !client.Hub.JoinRoom(client, client.RoomCode) {
						client.SendClientJSON()
					} else {
						if client.Hub.rooms[client.RoomCode].Game != nil {
							client.Hub.rooms[client.RoomCode].Game.HandleClientSwap(oldClient, client)
							client.Hub.rooms[client.RoomCode].Game.SendUpdateToClient(client)
							client.Hub.SendRoomUpdate(client.RoomCode)
						}
					}
				} else {
					client.SendClientJSON()
				}
				delete(client.Hub.clients, oldClient)
			} else {
				fmt.Println("Client Doesn't Exist -----------------------------")
			}
		case "toSystem":
			sysCmd := strings.SplitN(packet.Data, " ", 2)
			switch sysCmd[0] {
			case "setclientname":
				client.Name = sysCmd[1]
				client.SendClientJSON()
			case "setclientimage":
				client.Imguuid = sysCmd[1]
				client.SendClientJSON()
			case "joinroom":
				roomCode := sysCmd[1]
				client.Hub.JoinRoom(client, roomCode)
			case "createroom":
				roomCode := client.Hub.CreateRoom()
				client.Hub.JoinRoom(client, roomCode)
			case "startgame":
				switch sysCmd[1] {
				case "spygame":
					game := &SpyGame{Room: client.Hub.rooms[client.RoomCode]}
					game.SetupGame()

					client.Hub.rooms[client.RoomCode].Game = game
					client.Hub.SendRoomUpdate(client.RoomCode)

					game.StartGame()
				}
			}
		}
	}
}

func (client *Client) WritePackets() {
	defer func() {
		log.Panicln("Write Packet Close")
		client.Conn.Close()
	}()

	for {
		select {
		case packet, ok := <-client.Send:
			if !ok {
				log.Println("Error WritePackets (0)")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			var err error
			if packet.Data == "pong" {
				go func() {
					time.Sleep(time.Second * 10)
					err = client.Conn.WriteMessage(websocket.TextMessage, []byte("pong"))
				}()
			} else {
				fmt.Println("Sending Packet To ", client.Id)
				err = client.Conn.WriteJSON(packet)
			}
			if err != nil {
				log.Println("Error WritePackets (1)")
				log.Println(err)
			}
		}
	}
}
