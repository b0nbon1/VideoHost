'use client';

import { useRef } from "react";
import React from "react";
import videojs from "video.js";
import VideoPlayer from "./VideoPlayer";

function Player() {

  const videoSrc = "http://127.0.0.1:4500/api/v1/videos/stream/bZOl1D";

  const playerRef = useRef(null);

  const videoJsOptions = {
    autoplay: true,
    controls: true,
    responsive: true,
    fluid: true,
    sources: [
      {
        src: videoSrc,
        type: "application/x-mpegURL",
      },
    ],
  };

  const handlePlayerReady = (player) => {
    playerRef.current = player;

    // You can handle player events here, for example:
    player.on("waiting", () => {
      videojs.log("player is waiting");
    });

    player.on("dispose", () => {
      videojs.log("player will dispose");
    });
  };

  return (
    <>
      <h1>Video streaming app</h1>
      <VideoPlayer options={videoJsOptions} onReady={handlePlayerReady} />
    </>
  );
}

export default Player;
