import React from "react"
import { IClient, IPacket, IRoom, ISpyGame } from "../types"
import { socket } from "./socket.ts"
import RemoteImage from "./RemoteImage.tsx"

interface ISpyGameProps {
    gameData: ISpyGame
    roomData: IRoom
    clientData: IClient
}

const SpyGame = ({gameData, roomData, clientData}: ISpyGameProps) => {
    const ReadyUp = () => {
        socket.send(JSON.stringify({
            from: "",
            to: "",
            type: "toGame",
            data: "ready"
        } as IPacket))
    }

    const NextPlayer = () => {
        socket.send(JSON.stringify({
            from: "",
            to: "",
            type: "toGame",
            data: "nextPlayer"
        } as IPacket))
    }

    if (!gameData.isReady) {
        return(
            <div className="SpyGame">
                {
                    gameData.isSpy
                    ? <div>
                        <svg fill="#ffffff" version="1.1" id="Capa_1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 390.704 390.704">
                    <g>
                        <g>
                            <path d="M379.711,326.556L265.343,212.188c30.826-54.189,23.166-124.495-23.001-170.663c-55.367-55.366-145.453-55.366-200.818,0
                                c-55.365,55.366-55.366,145.452,0,200.818c46.167,46.167,116.474,53.827,170.663,23.001l114.367,114.369
                                c14.655,14.655,38.503,14.654,53.157,0C394.367,365.059,394.368,341.212,379.711,326.556z M214.057,214.059
                                c-39.77,39.771-104.479,39.771-144.25,0c-39.77-39.77-39.77-104.48,0-144.25c39.771-39.77,104.48-39.77,144.25,0
                                C253.828,109.579,253.827,174.29,214.057,214.059z"/>
                        </g>
                    </g>
                    </svg>
                    <h1>You're the SPY</h1>
                    <p>That means you have to guess what everyone is talking about</p>
                    <p>don't feel left out (loser)</p>
                    </div>
                    : <div>
                        <svg fill="#ffffff" version="1.1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                        <g>
                            <path d="M396.138,85.295c-13.172-25.037-33.795-45.898-59.342-61.03C311.26,9.2,280.435,0.001,246.98,0.001
                                c-41.238-0.102-75.5,10.642-101.359,25.521c-25.962,14.826-37.156,32.088-37.156,32.088c-4.363,3.786-6.824,9.294-6.721,15.056
                                c0.118,5.77,2.775,11.186,7.273,14.784l35.933,28.78c7.324,5.864,17.806,5.644,24.875-0.518c0,0,4.414-7.978,18.247-15.88
                                c13.91-7.85,31.945-14.173,58.908-14.258c23.517-0.051,44.022,8.725,58.016,20.717c6.952,5.941,12.145,12.594,15.328,18.68
                                c3.208,6.136,4.379,11.5,4.363,15.574c-0.068,13.766-2.742,22.77-6.603,30.442c-2.945,5.729-6.789,10.813-11.738,15.744
                                c-7.384,7.384-17.398,14.207-28.634,20.479c-11.245,6.348-23.365,11.932-35.612,18.68c-13.978,7.74-28.77,18.858-39.701,35.544
                                c-5.449,8.249-9.71,17.686-12.416,27.641c-2.742,9.964-3.98,20.412-3.98,31.071c0,11.372,0,20.708,0,20.708
                                c0,10.719,8.69,19.41,19.41,19.41h46.762c10.719,0,19.41-8.691,19.41-19.41c0,0,0-9.336,0-20.708c0-4.107,0.467-6.755,0.917-8.436
                                c0.773-2.512,1.206-3.14,2.47-4.668c1.29-1.452,3.895-3.674,8.698-6.331c7.019-3.946,18.298-9.276,31.07-16.176
                                c19.121-10.456,42.367-24.646,61.972-48.062c9.752-11.686,18.374-25.758,24.323-41.968c6.001-16.21,9.242-34.431,9.226-53.96
                                C410.243,120.761,404.879,101.971,396.138,85.295z"/>
                            <path d="M228.809,406.44c-29.152,0-52.788,23.644-52.788,52.788c0,29.136,23.637,52.772,52.788,52.772
                                c29.136,0,52.763-23.636,52.763-52.772C281.572,430.084,257.945,406.44,228.809,406.44z"/>
                        </g>
                        </svg>
                        <h2>The Prompt Is:</h2>
                        <h1>{gameData.prompt}</h1>
                        <p>Try not to be too obvious with your questions</p>
                    </div>
                }
                <input type="button" className="bigButton" value="      Ready ?      " onClick={ReadyUp}/>
                <h2>{gameData.readyString} Ready</h2>
            </div>
        )
    }

    if (gameData.questionClient.imguuid) {
        return(
            <div>
                {
                    gameData.questionClient.id === clientData.id
                    ? <div>
                        <p>{gameData.questionClient.id}--{clientData.id}</p>
                        <h1>It's Your Turn To Ask A Question</h1>
                        <p>Ask your question then press the button</p>
                        <input className="bigButton" type="button" value="      Next Player      " onClick={NextPlayer}/>
                    </div>
                    : <div>
                        <h1>{gameData.questionClient.name}s Time To Ask A Question</h1>
                        <RemoteImage uuid={gameData.questionClient.imguuid} />
                    </div>
                }
            </div>
        )
    } else {
        return (
            <div>
                <h1>{gameData.readyString}</h1>
            </div>
        )
    }
}

export default SpyGame