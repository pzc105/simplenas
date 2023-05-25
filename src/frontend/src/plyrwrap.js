import React, { useEffect, useRef, useState } from 'react';
import { Container, CssBaseline } from '@mui/material';
import { useParams, useLocation } from 'react-router-dom';

import Plyr from './lib/plyr/dist/plyr';
import './lib/plyr/dist/plyr.css';
import Hls from 'hls.js'

import { serverAddress } from './rpcClient.js'
import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { isNumber } from './utils';


export default function PlyrWrap() {
  const location = useLocation()
  const searchParams = new URLSearchParams(location.search)
  const shareid = searchParams.get('shareid')
  const itemId = searchParams.get('itemid')

  const player = useRef(null)
  const hls = useRef(null)
  const [url, setUrl] = useState('')
  const [subtitles, setSubtitles] = useState([])
  const videoRef = useRef(null);
  const vidRef = useRef(-1);

  const requestStartOffset = () => {
    if (shareid) {
      return
    }
    fetch(serverAddress + "/video/" + vidRef.current + "/get_offsettime", {
      method: 'GET',
      mode: 'cors',
      credentials: "include",
      headers: {
      },
    }).then(response => response.text())
      .then(data => {
        if (hls.current && isNumber(data)) {
          videoRef.current.currentTime = Number(data)
        }
      }).catch(error => {
        console.log(error)
        if (hls.current) {
          videoRef.current.currentTime = 0
        }
      });
  }

  let lastOffsetTime = useRef(0.0)
  useEffect(() => {
    if (shareid) {
      return
    }
    const saveStartOffset = () => {
      if (!videoRef.current) {
        return
      }
      if (lastOffsetTime.current === videoRef.current.currentTime) {
        return
      }
      lastOffsetTime.current = videoRef.current.currentTime
      fetch(serverAddress + "/video/" + vidRef.current + "/set_offsettime/" + videoRef.current.currentTime, {
        method: 'POST',
        mode: 'cors',
        credentials: "include",
        headers: {
        },
      })
    }

    videoRef.current.onseeking = () => {
      saveStartOffset()
    }

    setInterval(saveStartOffset, 3000)
    return () => {
      clearInterval(saveStartOffset);
    };
  }, [videoRef])

  function updateQuality(newQuality) {
    if (newQuality === 0) {
      hls.current.currentLevel = -1; //Enable AUTO quality if option.value = 0
    } else {
      hls.current.levels.forEach((level, levelIndex) => {
        if (level.height === newQuality) {
          hls.current.currentLevel = levelIndex;
        }
      });
    }
  }

  function updateAudioTrack(newTrack) {
    hls.current.audioTrack = newTrack;
  }

  var availableQualities = useRef([])

  useEffect(() => {
    if (url.length === 0) {
      return
    }
    var hlsconfig = {
      debug: false,
      autoStartLoad: false,
      xhrSetup: function (xhr) {
        xhr.withCredentials = true; // do send cookies
      },
    }
    hls.current = new Hls(hlsconfig);
    hls.current.attachMedia(videoRef.current);
    hls.current.on(Hls.Events.AUDIO_TRACKS_UPDATED, function (event, data) {
      let audioItemLabels = {}
      const audioTracks = hls.current.audioTracks.map((a, i) => {
        audioItemLabels[i] = a.lang
        return i
      })

      const defaultOptions = {
        debug: false,
        controls: ['play-large', 'play', 'progress', 'current-time', 'restart',
          'mute', 'volume', 'captions', 'settings', 'pip', 'airplay', 'fullscreen'],
        customSettings: ['audio'],
      }

      defaultOptions.quality = {
        default: availableQualities.current.length - 1,
        options: availableQualities.current,
        forced: true,
        onChange: (nv) => updateQuality(nv),
      }
      defaultOptions.audio = {
        default: audioTracks[0],
        options: audioTracks,
        onChange: (nv) => updateAudioTrack(nv),
      }
      defaultOptions.i18n = {
        qualityLabel: {
          0: 'Auto',
        },
        audio: "Audio",
        customItemLabel: {
          audio: audioItemLabels
        }
      }
      defaultOptions.seekTime = 3
      defaultOptions.clickToPlay = false
      defaultOptions.fullscreen = {
        enabled: true, fallback: true, iosNative: true, container: null
      }

      if (player.current) {
        player.current.destroy()
      }
      player.current = new Plyr(videoRef.current, defaultOptions);
      player.current.on('enterfullscreen', event => {
        if (window.screen.orientation.lock)
          window.screen.orientation.lock('landscape')
      });

      player.current.on('exitfullscreen', event => {
        if (window.screen.orientation.lock)
          window.screen.orientation.lock('portrait');
      });
    })
    hls.current.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
      availableQualities.current = hls.current.levels.map((l) => l.height)
      availableQualities.current.unshift(0) //prepend 0 to quality array

      hls.current.startLevel = availableQualities.current.length - 1
      hls.current.startLoad(0)
      requestStartOffset()
    })

    hls.current.on(Hls.Events.ERROR, function (event, data) {
      if (data.type === Hls.ErrorTypes.NETWORK_ERROR) {
        hls.current.startLoad()
      } else {
        console.log(event, data)
      }
    })

    hls.current.loadSource(url)

    return () => {
      hls.current.destroy()
      if (player.current) {
        player.current.destroy()
      }
    }
  }, [url])

  useEffect(() => {
    var req = new User.QueryItemInfoReq()
    req.setItemId(itemId)
    if (shareid) {
      req.setShareId(shareid)
    }
    userService.queryItemInfo(req, {}, (err, res) => {
      if (err != null || !res) {
        console.log(err)
        return
      }
      const itemInfo = res.getItemInfo()
      if (itemInfo.getTypeId() !== Category.CategoryItem.Type.VIDEO) {
        return
      }

      const videoInfo = res.getVideoInfo()
      const vid = videoInfo.getId()
      vidRef.current = vid
      let urlPath = serverAddress + "/video/" + vid
      if (shareid) {
        urlPath += "?shareid=" + shareid + "&itemid=" + itemId
      }
      setUrl(urlPath)
      let cs = []
      videoInfo.getSubtitlePathsList().map((c) => {
        let suffixes = c.split(".")
        let lang = "unknown"
        if (suffixes.length > 2) {
          suffixes.pop()
          lang = suffixes.pop()
        }
        let urlPath = serverAddress + "/video/" + vid + "/subtitle/" + c
        if (shareid) {
          urlPath += "?shareid=" + shareid + "&itemid=" + itemId
        }
        cs.push({
          kind: "subtitles",
          src: urlPath,
          srcLang: lang,
        })
        setSubtitles(cs)
        return null
      })
    })
  }, [itemId, shareid]);

  var touchStartX = useRef(0);
  var touchEndX = useRef(0);

  const touchstart = (e) => {
    touchStartX.current = e.touches[0].clientX;
    touchEndX.current = touchStartX.current
  }

  const touchmove = (e) => {
    touchEndX.current = e.touches[0].clientX;
  }

  const touchend = (e) => {
    var diffX = touchEndX.current - touchStartX.current;
    if (Math.abs(diffX) > 30) {
      videoRef.current.currentTime += diffX / 10
    }
  }

  return (
    <Container onTouchStart={touchstart} onTouchMove={touchmove} onTouchEnd={touchend} sx={{ backgroundColor: 'background.default' }}>
      <CssBaseline />
      <video ref={videoRef} crossOrigin="use-credentials">
        {
          subtitles.map((s, i) => (
            <track key={i} kind={s.find} label={s.srcLang} src={s.src} srcLang={s.srcLang} />
          ))
        }
      </video>
    </Container>
  );
}
