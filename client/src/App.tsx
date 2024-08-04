import { useRef, useState } from "react";
import React from 'react';
import videojs from "video.js";
import VideoPlayer from "./VideoPlayer";

function App() {
  const [count, setCount] = useState(0);

  const videoSrc = 'http://127.0.0.1:4500/static/videos/7yjNCBU8nr/index.m3u8';

  const playerRef = useRef(null);

  const videoJsOptions = {
    autoplay: true,
    controls: true,
    responsive: true,
    fluid: true,
    sources: [{
      src: videoSrc,
      type: 'application/x-mpegURL'
    }],
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
  };

  return (
    <>
      <h1>Video streaming app</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
      <video id="streamVideo" src="http://127.0.0.1:4500/api/v1/videos/stream" controls>
        Your browser does not support the video tag.
      </video>

      <h1>HLS player down here</h1>
      {/* <ShakaPlayer autoPlay src="http://127.0.0.1:4500/static/videos/ZJy8cTMcMW/index.m3u8" /> */}
      <VideoPlayer options={videoJsOptions} onReady={handlePlayerReady} />
    </>
  );
}

export default App;
