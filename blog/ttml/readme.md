# TTML (stpp)

<https://wikipedia.org/wiki/Timed_Text_Markup_Language>

go.sum 0 LOC:

https://github.com/zmalltalker/ttml2vtt

go.sum 2 LOC:

https://github.com/wargarblgarbl/libgosubs

go.sum 28 LOC:

https://github.com/asticode/go-astisub

## VLC

VLC support it, but the timing is off. cant even file a fucking issue:

> Your account is pending approval from your GitLab administrator and hence
> blocked. Please contact your GitLab administrator if you think this is an
> error.

https://code.videolan.org/videolan/vlc/-/issues

## WebVTT

yeah, seems like support for TTML is pretty crappy at this time, missing from FFmpeg and MPV. as a workaround, it seems WebVTT is pretty common with better support:

- https://ffmpeg.org/ffmpeg-all.html#Subtitle-Formats
- https://wikipedia.org/wiki/WebVTT

## based on FFmpeg or MPV

- https://github.com/BazzaCuda/MinimalistMediaPlayerX/issues/22
- https://github.com/media-kit/media-kit/issues/728
- https://github.com/mpc-qt/mpc-qt/issues/112
- https://github.com/mpvnet-player/mpv.net/issues/667
- https://github.com/pkoshevoy/aeyae/issues/6
- https://github.com/smplayer-dev/smplayer/issues/920
- https://github.com/tsl0922/ImPlay/issues/80
- https://github.com/zaps166/QMPlay2/issues/685

## FFmpeg

FFmpeg only supports encoding not decoding

- https://ffmpeg.org/ffmpeg-all.html#Subtitle-Formats
- https://trac.ffmpeg.org/ticket/10902
- https://trac.ffmpeg.org/ticket/4859
- https://github.com/GyanD/codexffmpeg/issues/115

## MPC-HC

refuse to fix

https://github.com/clsid2/mpc-hc/issues/2665

## MPV

MPV is based on FFmpeg

https://github.com/mpv-player/mpv/issues/7573
