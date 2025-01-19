import React, { useEffect, useState } from 'react';
import './App.css';
import { socket, SOCKURL } from './Comps/socket.ts';
import { IClient, IPacket, IRoom, ISpyGame } from './types.ts';
import RoomJoin from './Comps/RoomJoin.tsx';
import Lobby from './Comps/Lobby.tsx';

function App() {
  const [packet, setPacket] = useState<IPacket>()
  const [clientData, setClientData] = useState<IClient>()
  const [roomData, setRoomData] = useState<IRoom>()
  const [gameData, setGameData] = useState<ISpyGame>()
  const [connected, setConnected] = useState(false)

  useEffect(() => {
    socket.onopen = () => {
      console.log("wagwan")
      setConnected(true)
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

    socket.onclose = () => {
      setConnected(false)
      console.log("closed bruv")
    }

    socket.onmessage = (event) => {
      if (event.data == "pong") {
        socket.send("ping")
      } else {
        setPacket(JSON.parse(event.data))
      }
    }
  }, [])

  useEffect(() => {
    if (packet === undefined) {
      return
    }
    console.log("New Packet")
    console.log(packet)
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
  }, [packet])

  return (
    <div className="App">
      <h1>🍺BigPints.com</h1>
      <div className='content'>
        {
          clientData === undefined || !connected
          ? <div>
              <h2>Connection Status: {connected?"Online":"Offline"}</h2>
              <button className='bigButton' onClick={() => window.location.reload()}>Reconnect</button>
            </div>
          : roomData === undefined
            ? <RoomJoin client={clientData} />
            : gameData === undefined
              ? <Lobby client={clientData} room={roomData} />
              : <div>
                  <h2>{gameData.isSpy?"SPY":gameData.prompt}</h2>
                </div>
        }
      </div>
    </div>
  );
}

export default App;
