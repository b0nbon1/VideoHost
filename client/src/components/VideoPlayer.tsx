import React, { useEffect, useRef } from "react";
import videojs from "video.js";
import "videojs-hls-quality-selector";
import "video.js/dist/video-js.css";
import httpSourceSelector from "videojs-http-source-selector";

const VideoPlayer = ({ src }) => {
  const videoRef = useRef(null);
  const playerRef = useRef(null);

  useEffect(() => {
    if (!playerRef.current) {
      playerRef.current = videojs(videoRef.current, {
        controls: true,
        autoplay: true,
        responsive: true,
        fluid: true,
        html5: {
          vhs: {
            enableLowInitialPlaylist: true,
          },
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
        controlBar: {
          pictureInPictureToggle: false
        },
        // height: "480px"
      });

      // Register the quality selector plugin
      videojs.registerPlugin("httpSourceSelector", httpSourceSelector);
      playerRef.current.httpSourceSelector(); // Enable the plugin

      // Load HLS source
      playerRef.current.src({
        src: src,
        type: "application/x-mpegURL",
      });

      playerRef.current.ready(() => {
        playerRef.current.hlsQualitySelector({
          displayCurrentQuality: true, // Shows selected quality
        });
      });
    }

    return () => {
      if (playerRef.current) {
        playerRef.current.dispose();
        playerRef.current = null;
      }
    };
  }, [src]);

  return (
    <div style={{ width: "640px", height: "480px" }}>
    <div data-vjs-player>
      <video ref={videoRef} className="video-js vjs-default-skin" />
    </div>
    </div>
  );
};

export default VideoPlayer;
