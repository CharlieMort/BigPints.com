import React, { useEffect, useState } from 'react';
import './App.css';
import { ReconnectSocket, socket, SOCKURL } from './Comps/socket.ts';
import { IClient, IPacket, IRoom } from './types.ts';
import RoomJoin from './Comps/RoomJoin.tsx';
import Lobby from './Comps/Lobby.tsx';

function App() {
  const [packet, setPacket] = useState<IPacket>()
  const [clientData, setClientData] = useState<IClient>()
  const [roomData, setRoomData] = useState<IRoom>()
  const [connected, setConnected] = useState(false)
  const [retry, setRetry] = useState(false)

  useEffect(() => {
    socket.onopen = () => {
      console.log("wagwan")
      setConnected(true)
      if (sessionStorage.tabID == undefined) {
          sessionStorage.tabID = crypto.randomUUID()
      }
      try {
        socket.send(JSON.stringify({
          from: "0",
          to: "0",
          type: "setup",
          data: `clientconnect ${sessionStorage.tabID}`
        } as IPacket))
        socket.send("ping")
      } catch {
        console.log("Failed Send")
      }
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
  }, [socket])

  useEffect(() => {
      let ReConnectTimer = setInterval(() => {
        setRetry(!retry)
        if (connected == false){
          ReconnectSocket()
        }
      }, 2000)
      return () => {
        clearInterval(ReConnectTimer)
      }
    }, [retry])

  useEffect(() => {
    if (packet === undefined) {
      return
    }
    console.log("New Packet")
    switch (packet.type) {
      case "clientData":
        setClientData(JSON.parse(packet.data))
        break;
      case "roomData":
        setRoomData(JSON.parse(packet.data))
        break;
    }
  }, [packet])

  return (
    <div className="App">
      <h1>🍺BigPint.com {retry?"1":"0"}</h1>
      {
        clientData === undefined || !connected
        ? <div>
            <h2>Connection Status: {connected?"Online":"Offline"}</h2>
          </div>
        : roomData === undefined
          ? <RoomJoin client={clientData} />
          : <Lobby client={clientData} room={roomData} />
      }
    </div>
  );
}

export default App;
