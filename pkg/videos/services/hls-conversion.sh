#!/bin/sh
params=$1
param2=$2

# write the ffmpeg command to split the video into segments with different bitrates and resolutions
ffmpeg -i $params -vf scale=640:360 -c:v libx264 -b:v 800k -c:a aac -b:a 96k -f hls -hls_time 4 -hls_playlist_type vod -hls_segment_filename 360p_%03d.ts 360p.m3u8 -vf scale=842:480 -c:v libx264 -b:v 1400k -c:a aac -b:a 128k -f hls -hls_time 4 -hls_playlist_type vod -hls_segment_filename 480p_%03d.ts 480p.m3u8 -vf scale=1280:720 -c:v libx264 -b:v 2800k -c:a aac -b:a 128k -f hls -hls_time 4 -hls_playlist_type vod -hls_segment_filename 720p_%03d.ts 720p.m3u8
