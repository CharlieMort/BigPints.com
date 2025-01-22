import React, { useEffect, useState } from 'react';
import './App.css';
import { socket, SOCKURL } from './Comps/socket.ts';
import { IClient, IPacket, IRoom, ISpyGame } from './types.ts';
import RoomJoin from './Comps/RoomJoin.tsx';
import Lobby from './Comps/Lobby.tsx';
import SpyGame from './Comps/SpyGame.tsx';

let c = false
let currPacket = undefined

function App() {
  // const [packet, setPacket] = useState<IPacket>()
  const [clientData, setClientData] = useState<IClient>()
  const [roomData, setRoomData] = useState<IRoom>()
  const [gameData, setGameData] = useState<ISpyGame>()
  const [connected, setConnected] = useState(false)

  const clientSetup = () => {
    setConnected(true)
    c = true
    if (sessionStorage.tabID == undefined) {
        sessionStorage.tabID = crypto.randomUUID()
    }
    socket.send(JSON.stringify({
      from: "0",
      to: "0",
      type: "setup",
      data: `clientconnect ${sessionStorage.tabID}`
    } as IPacket))
    socket.send("ping")
  }

  useEffect(() => {
    socket.onopen = () => {
      console.log("wagwan")
      clientSetup()
    }

    socket.onclose = () => {
      setConnected(false)
      c = false
      console.log("closed bruv")
    }

    socket.onmessage = (event) => {
      if (event.data == "pong") {
        socket.send("ping")
      } else {
        console.log("New Message "+JSON.parse(event.data).type)
        if (!c) {
          console.log("wagwan pt2")
          clientSetup()
        }
        let packet = JSON.parse(event.data)
        switch (packet.type) {
          case "clientData":
            setClientData(JSON.parse(packet.data))
            break;
          case "roomData":
            setRoomData(JSON.parse(packet.data))
            break;
          case "gameData":
            setGameData(JSON.parse(packet.data))
            break
        }
      }
    }
  }, [])

  return (
    <div className="App">
      <h1>🍺BigPints.com</h1>
      <div className='content'>
        {
          clientData === undefined || !c
          ? <div>
              <h2>Connection Status: {c?"Online":"Offline"}</h2>
              <button className='bigButton' onClick={() => window.location.reload()}>Reconnect</button>
            </div>
          : roomData === undefined
            ? <RoomJoin client={clientData} />
            : gameData === undefined
              ? <Lobby client={clientData} room={roomData} />
              : <SpyGame gameData={gameData} roomData={roomData} clientData={clientData} />
        }
      </div>
    </div>
  );
}

export default App;
