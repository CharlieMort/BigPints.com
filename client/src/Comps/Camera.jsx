import React, { useState, useRef, useEffect } from "react";
import {Camera} from "react-camera-pro";
import { API_URL } from "./socket.ts";

const Cam = ({setImgUUID}) => {
  const camera = useRef(null);

  const takePhoto = () => {
    fetch(`${API_URL}/api/upload`, {
      method: "POST",
      body: camera.current.takePhoto(),
      headers: {
        "Content-type": "text/plain"
      }
    }).then((r) => {
      r.text().then((uid) => setImgUUID(uid))
    });
  }

  return (
    <div className="Camera">
      <h1>Take A Photo !</h1>
      <Camera ref={camera} aspectRatio={1} />
      <button className="Camera_Button" onClick={takePhoto}>📷</button>
    </div>
  );
}

export default Cam;