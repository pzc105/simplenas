import React, { useEffect, useRef, useState } from 'react';
import { Container, Grid, CssBaseline, List, ListItem, Button, Typography, Tooltip, Switch, FormControlLabel, Paper } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';

import Plyr from 'plyr';
import 'plyr/dist/plyr.css';
import Hls from 'hls.js'

import { queryItem, querySubItems, navigateToItem } from './category.js'
import * as Category from './prpc/category_pb.js'
import { serverAddress } from './rpcClient.js'
import * as store from './store.js'
import { isNumber } from './utils';


export default function PlyrWrap() {
  const dispatch = useDispatch()
  const location = useLocation()
  const navigate = useNavigate()
  const searchParams = new URLSearchParams(location.search)
  const shareid = searchParams.get('shareid')
  const itemId = Number(searchParams.get('itemid'))
  const videoInfo = useSelector((state) => store.selectItemVideoInfo(state, itemId))
  const item = useSelector((state) => store.selectCategoryItem(state, itemId))
  const items = useSelector((state) => store.selectCategorySubItems(state, item ? item.parentId : -1))
  const [videoItemList, setVideoItemList] = useState([])

  const player = useRef(null)
  const hls = useRef(null)
  const [url, setUrl] = useState('')
  const [subtitles, setSubtitles] = useState([])
  const videoRef = useRef(null);
  const vidRef = useRef(-1);
  const selectedAudio = useSelector((state) => store.getSelectedAudio(state, vidRef.current))

  useEffect(() => {
    queryItem(itemId, shareid, dispatch)
    querySubItems(item.parentId, shareid, dispatch)
  }, [itemId, shareid, dispatch])

  useEffect(() => {
    if (!items) {
      return
    }
    let vl = []
    items.map((item) => {
      if (item.typeId === Category.CategoryItem.Type.VIDEO) {
        vl.push(item)
      }
      return null
    })
    setVideoItemList(vl)
  }, [])

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

      player.current = new Plyr(videoRef.current, defaultOptions);
      player.current.on('enterfullscreen', event => {
        if (window.screen.orientation.lock)
          window.screen.orientation.lock('landscape')
      });
      player.current.on('ready', event => {
        updateAudioTrack(defaultAudioId)
      });
      player.current.on('ended', event => {
        if (autoContinuedPlay) {
          for (let i = 0; i < videoItemList.length; i++) {
            if (videoItemList[i].id === itemId && i < videoItemList.length - 1) {
              navigateToItem(navigate, {}, videoItemList[i + 1].id, shareid)
            }
          }
        }
      });

      player.current.on('exitfullscreen', event => {
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
      hls.current.destroy()
      if (player.current) {
        player.current.destroy()
      }
    }
  }, [url, shareid])

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

  const [autoContinuedPlay, setAutoContinuedPlay] = useState(useSelector((state) => store.selectAutoContinuedPlayVideo(state)));

  return (
    <Container onTouchStart={touchstart} onTouchMove={touchmove} onTouchEnd={touchend} sx={{ backgroundColor: 'background.default' }}>
      <CssBaseline />
      <Grid container spacing={2}>
        <Grid item xs={12} sx={{ display: "flex" }}>
          <Grid item xs={8} >
            <video style={{ height: '50vh' }} ref={videoRef} crossOrigin="use-credentials">
              {
                subtitles.map((s, i) => (
                  <track key={i} kind={s.find} label={s.srcLang} src={s.src} srcLang={s.srcLang} />
                ))
              }
            </video>
          </Grid>
          <Grid item xs={4} >
            <Container>
              <Typography variant="button" component="div" noWrap>
                播放列表
              </Typography>
              <FormControlLabel
                control={
                  <Switch
                    checked={autoContinuedPlay}
                    onClick={
                      (e) => {
                        let v = !autoContinuedPlay
                        setAutoContinuedPlay(v)
                        dispatch(store.playerSlice.actions.setAutoContinuedPlayVideo(v))
                      }
                    }
                    color="primary"
                    inputProps={{ 'aria-label': 'controlled' }}
                  />
                }
                label={'自动连播'}
              />
              <Paper style={{ maxHeight: '50vh', overflow: 'auto' }}>
                <List>
                  {
                    videoItemList.map((item) => {
                      return (
                        <ListItem
                          key={item.id} >
                          <Tooltip title={item.name}>
                            <Typography variant="button" component="div" noWrap>
                              <Button onClick={() => navigateToItem(navigate, {}, item.id, shareid)}>
                                {item.name}
                              </Button>
                            </Typography>
                          </Tooltip>
                        </ListItem>
                      )
                    })
                  }
                </List>
              </Paper>
            </Container>
          </Grid>
        </Grid>
      </Grid>

    </Container>
  );
}
