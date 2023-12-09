import React, { useEffect, useRef, useState } from 'react';
import { Container, Grid, CssBaseline, List, ListItem, Button, Typography, Tooltip, Switch, FormControlLabel, Paper } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';

import Plyr from 'plyr_p';
import 'plyr_p/dist/plyr.css';
import Hls from 'hls.js'

import { queryItem, querySubItems, navigateToItem } from '../category.js'
import { serverAddress } from '../rpcClient.js'
import * as store from '../store.js'
import * as utils from '../utils.js';
import { FloatingChat } from '../chat/chat.js';
import PlayList from './videoList.js'


export default function PlyrWrap() {
  const dispatch = useDispatch()
  const location = useLocation()
  const navigate = useNavigate()
  const searchParams = new URLSearchParams(location.search)
  const shareid = searchParams.get('shareid')
  const itemId = Number(searchParams.get('itemid'))
  const videoInfo = useSelector((state) => store.selectItemVideoInfo(state, itemId))
  const item = useSelector((state) => store.selectCategoryItem(state, itemId))
  const parentItemId = item ? item.parentId : null
  const [items, setItems] = useState([])
  const [videoItemList, setVideoItemList] = useState([])
  const videoItemListRef = useRef([])

  const plyr = useRef(null)
  const hls = useRef(null)
  const [url, setUrl] = useState('')
  const [subtitles, setSubtitles] = useState([])
  const videoRef = useRef(null);
  const vidRef = useRef(-1);
  const selectedAudio = useSelector((state) => store.getSelectedAudio(state, vidRef.current))

  useEffect(() => {
    queryItem(itemId, shareid, dispatch)
  }, [itemId, shareid, dispatch])

  useEffect(() => {
    if (parentItemId) {
      querySubItems({
        itemId: parentItemId, shareid, dispatch, callback: (subItems) => {
          setItems(subItems)
        }
      })
    }
  }, [parentItemId, shareid, dispatch])

  useEffect(() => {
    if (!items) {
      return
    }
    let vl = []
    items.map((item) => {
      if (utils.isVideoItem(item)) {
        vl.push(item)
      }
      return null
    })
    vl.sort((a, b) => {
      if (a.name < b.name) {
        return -1;
      }
      if (a.name > b.name) {
        return 1;
      }
      return 0;
    })
    setVideoItemList(vl)
    videoItemListRef.current = vl;
  }, [items])

  useEffect(() => {
    if (!videoInfo) {
      return
    }
    const vid = videoInfo.id
    vidRef.current = vid
    let urlPath = serverAddress + "/video/" + vid
    if (shareid) {
      urlPath += "?shareid=" + shareid + "&itemid=" + itemId
    }
    setUrl(urlPath)
    let cs = []
    videoInfo.subtitlePathsList.map((c) => {
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
  }, [itemId, shareid, videoInfo]);

  let lastOffsetTime = useRef(0.0)
  useEffect(() => {
    if (shareid) {
      return
    }
    const saveStartOffset = (offset) => {
      if (!videoRef.current) {
        return
      }
      if (lastOffsetTime.current === videoRef.current.currentTime) {
        return
      }
      lastOffsetTime.current = offset ? offset : videoRef.current.currentTime
      fetch(serverAddress + "/video/" + vidRef.current + "/set_offsettime/" + lastOffsetTime.current, {
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
    let myInterval = setInterval(saveStartOffset, 3000)
    return () => {
      clearInterval(myInterval);
    };
  }, [videoRef, shareid])

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
    hls.current.audioTrack = newTrack
    dispatch(store.playerSlice.actions.updateSelectedAudio({ vid: vidRef.current, aid: newTrack }))
  }

  var availableQualities = useRef([])

  useEffect(() => {
    if (url.length === 0) {
      return
    }

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
          if (hls.current && utils.isNumber(data)) {
            videoRef.current.currentTime = Number(data)
          }
        }).catch(error => {
          console.log(error)
          if (hls.current) {
            videoRef.current.currentTime = 0
          }
        });
    }

    var hlsconfig = {
      debug: false,
      autoStartLoad: false,
      xhrSetup: function (xhr) {
        xhr.withCredentials = true; // do send cookies
      },
    }
    if (hls.current) {
      let t = new Hls(hlsconfig);
      hls.current.detachMedia()
      hls.current.attachMedia(videoRef.current);
      hls.current.loadSource(url)
      return
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
      let defaultAudioId = audioTracks[0]
      if (selectedAudio) {
        defaultAudioId = Number(selectedAudio)
      }
      defaultOptions.audio = {
        default: defaultAudioId,
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
      defaultOptions.keyboard = {
        focused: true, global: true
      }
      plyr.current = new Plyr(videoRef.current, defaultOptions);
      plyr.current.on('enterfullscreen', event => {
        if (window.screen.orientation.lock)
          window.screen.orientation.lock('landscape')
      });
      plyr.current.on('ready', event => {
        updateAudioTrack(defaultAudioId)
        if (autoContinuedPlay) {
          plyr.current.play()
        }
      });
      plyr.current.on('ended', event => {
        if (autoContinuedPlay) {
          for (let i = 0; i < videoItemListRef.current.length; i++) {
            if (videoItemListRef.current[i].id === itemId && i < videoItemListRef.current.length - 1) {
              navigateToItem(navigate, {}, videoItemListRef.current[i + 1].id, shareid)
            }
          }
        }
      });

      plyr.current.on('exitfullscreen', event => {
        if (window.screen.orientation.lock)
          window.screen.orientation.lock('portrait');
      });
    })
    hls.current.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
      availableQualities.current = hls.current.levels.map((l) => l.height)
      hls.current.startLevel = availableQualities.current.length - 1
      availableQualities.current.unshift(0) //prepend 0 to quality array
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
    }
  }, [url, shareid, videoRef])

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

  const autoContinuedPlay = useSelector((state) => store.selectAutoPlayVideo(state))
  const showGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))
  const closeGlobalChat = () => {
    dispatch(store.userSlice.actions.setOpenGlobalChat(false))
  }

  return (
    <Container onTouchStart={touchstart} onTouchMove={touchmove} onTouchEnd={touchend} sx={{ backgroundColor: 'background.default' }}>
      <CssBaseline />
      <Grid container alignItems="center" justify="center" spacing={2}>
        <Grid item xs={12} sx={{ display: "flex" }}>
          <Grid item xs={8} >
            <Typography variant="button" component="div" noWrap>
              {item ? item.name : ""}
            </Typography>
            <video ref={videoRef} crossOrigin="use-credentials">
              {
                subtitles.map((s, i) => (
                  <track key={i} kind={s.find} label={s.srcLang} src={s.src} srcLang={s.srcLang} />
                ))
              }
            </video>
          </Grid>
          <Grid item xs={4} >
            <PlayList videoItemList={videoItemList} shareid={shareid} />
          </Grid>
        </Grid>
      </Grid>
      {showGlobalChat ? <FloatingChat itemId={1} onClose={closeGlobalChat} /> : null}
    </Container>
  );
}

