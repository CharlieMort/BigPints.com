package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const PORT = ":80"
const DEBUG = true

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func GetRandomRoomCode() string {
	rnd := "0123456789"
	lng := 5
	out := ""
	for i := 0; i < lng; i++ {
		out = out + string(rnd[rand.Intn(len(rnd))])
	}
	return out
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("File Uploading")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	fx := strings.ReplaceAll(string(body), "=", "")
	dec, err := base64.RawStdEncoding.DecodeString(fx[strings.IndexByte(fx, ',')+1:])
	if err != nil {
		log.Printf("E0 %s", err.Error())
		return
	}
	uid := uuid.New().String()
	var fileName string
	if DEBUG {
		fileName = "./images/" + uid + ".jpeg"
	} else {
		fileName = "/go/bin/images/" + uid + ".jpeg"
	}
	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("E1 %s", err.Error())
		return
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		log.Printf("E2 %s", err.Error())
		return
	}
	if err := f.Sync(); err != nil {
		log.Printf("E3 %s", err.Error())
		return
	}

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "%s", uid)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	vars := mux.Vars(r)

	var fileName string
	if DEBUG {
		fileName = "./images/" + vars["uuid"] + ".jpeg"
	} else {
		fileName = "/go/bin/images/" + vars["uuid"] + ".jpeg"
	}
	f, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("Fialed Reading FIle")
		log.Println(err)
	}

	bs64 := base64.RawStdEncoding.EncodeToString(f)
	fmt.Fprintf(w, "data:image/jpeg;base64,%s", bs64)
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	indexPath := filepath.Join(h.staticPath, h.indexPath)
	log.Printf("URL: %s outPATH: %s indexPath:%s\n", r.URL.Path, path, indexPath)

	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		log.Println("FILE DONT EXIST")
		http.ServeFile(w, r, indexPath)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		ClientJSON{
			Id:       strconv.Itoa(nextClientID),
			RoomCode: "",
			Name:     "",
			Imguuid:  "",
		},
		ClientData{
			TabID: "",
			Hub:   hub,
			Conn:  conn,
			Send:  make(chan Packet),
		},
	}
	nextClientID += 1
	client.Hub.register <- client

	go client.WritePackets()
	go client.ReadPackets()
}

func Chk(r *http.Request) bool {
	return true
}

var nextClientID int
var upgrader = websocket.Upgrader{
	CheckOrigin: Chk,
}

func main() {
	fmt.Println("Drunk Games Server Running...")

	nextClientID = 1
	hub := NewHub()
	go hub.run()

	r := mux.NewRouter()
	r.HandleFunc("/api/upload", uploadImage)
	r.HandleFunc("/api/image/get/{uuid}", getImage)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		log.Println("Someone Connected to ws")
		serveWs(hub, w, r)
	})
	spa := spaHandler{staticPath: "/go/bin/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	http.ListenAndServe(PORT, r)
}
