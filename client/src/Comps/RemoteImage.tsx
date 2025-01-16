import React, {useEffect, useState} from "react"
import { API_URL } from "./socket.ts"

interface IRemoteImageProps {
    uuid: string
}

const RemoteImage = ({uuid}: IRemoteImageProps) => {
    const [imageData, setImageData] = useState<string>("")

    useEffect(() => {
        console.log("Getting Image for ", uuid)
        if (uuid === undefined || uuid === "") {
            return
        }
        fetch(`${API_URL}/api/image/get/${uuid}`, {
            method: "GET"
        }).then((r) => {
            r.text().then((dat) => setImageData(dat))
        })
    }, [uuid])

    if (imageData == "") {
        return <></>
    }

    return(
        <img className="Photo" src={imageData} alt='Taken photo'/>
    )
}

export default RemoteImage