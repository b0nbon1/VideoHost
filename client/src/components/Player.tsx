"use client";

import { useRef, useState } from "react";
import React from "react";
import videojs from "video.js";
import VideoPlayer from "./VideoPlayer";

function Player() {
  const [res, setRes] = useState(720);

  const videoSrc = "http://127.0.0.1:4500/static/videos/6vqlLY/master.m3u8";

  const playerRef = useRef(null);

  const videoJsOptions = {
    autoplay: true,
    controls: true,
    responsive: true,
    fluid: true,
    html5: {
      hls: {
        overrideNative: true,
        limitRenditionByPlayerDimensions: true,
        useDevicePixelRatio: true
        // bandwidth: 16777216,
      },
      nativeAudioTracks: false,
      nativeVideoTracks: false,
      useBandwidthFromLocalStorage: true
    },
    sources: [
      {
        src: videoSrc,
        type: 'application/x-mpegURL',
      },
    ],
    controlBar: {
      pictureInPictureToggle: false
    }
  };

  const handlePlayerReady = (player) => {
    playerRef.current = player;

    // You can handle player events here, for example:
    player.on('waiting', () => {
      videojs.log('player is waiting');
    });

    player.on('dispose', () => {
      videojs.log('player will dispose');
    });
    // player.on('ready', () => {
    //   player.httpSourceSelector();
    // })
  };

  return (
    <>
      <h1>Video streaming app</h1>
      <VideoPlayer src={videoSrc} />
    </>
  );
}

export default Player;
