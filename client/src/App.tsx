import React, { useEffect, useState } from 'react';
import './App.css';
import { socket } from './Comps/socket.ts';
import { IClient, IPacket, IRoom } from './types.ts';
import RoomJoin from './Comps/RoomJoin.tsx';
import Lobby from './Comps/Lobby.tsx';

function App() {
  const [packet, setPacket] = useState<IPacket>()
  const [clientData, setClientData] = useState<IClient>()
  const [roomData, setRoomData] = useState<IRoom>()

  useEffect(() => {
    socket.onopen = () => {
      console.log("wagwan")
      if (sessionStorage.tabID == undefined) {
          sessionStorage.tabID = crypto.randomUUID()
      }
      socket.send(JSON.stringify({
        from: "0",
        to: "0",
        type: "setup",
        data: `clientconnect ${sessionStorage.tabID}`
      } as IPacket))
    }

    socket.onclose = () => {
      console.log("closed bruv")
    }

    socket.onmessage = (event) => {
      setPacket(JSON.parse(event.data))
    }
  }, [])

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
      <h1>🍺BigPint.com</h1>
      {
        clientData === undefined
        ? <div>
            <h2>Waiting To Connect...</h2>
          </div>
        : roomData === undefined
          ? <RoomJoin client={clientData} />
          : <Lobby client={clientData} room={roomData} />
      }
    </div>
  );
}

export default App;
