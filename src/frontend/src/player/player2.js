import React, { useEffect, useRef, useState } from 'react';
import { Container, Grid, CssBaseline, List, ListItem, Button, Typography, Tooltip, Switch, FormControlLabel, Paper, Hidden } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';

import Hls from 'hls.js'
import DPlayer from 'dplayer';

import { queryItem, querySubItems, navigateToItem } from '../category/utils.js'
import userService, { serverAddress } from '../rpcClient.js'
import * as User from '../prpc/user_pb.js'
import * as store from '../store.js'
import * as utils from '../utils.js';
import { FloatingChat } from '../chat/chat.js';
import PlayList from './videoList.js'

export default function Player() {
  const dispatch = useDispatch()
  const location = useLocation()
  const navigate = useNavigate()
  const searchParams = new URLSearchParams(location.search)
  const shareid = searchParams.get('shareid')
  const itemId = Number(searchParams.get('itemid'))
  const videoInfo = useSelector((state) => store.selectItemVideoInfo(state, itemId))
  const item = useSelector((state) => store.selectCategoryItem(state, itemId))
  const parentItemId = item ? item.parentId : null
  const categoryDesc = useSelector((state) => store.selectCategoryDesc(state))
  const [items, setItems] = useState([])
  const [videoItemList, setVideoItemList] = useState([])
  const videoItemListRef = useRef([])
  const autoContinuedPlay = useSelector((state) => store.selectAutoPlayVideo(state));
  const userInfo = useSelector((state) => store.selectUserInfo(state))

  const [urlRef, setUrl] = useState('')
  const subtitlesRef = useRef([])
  const dplayerRef = useRef(null);
  const vidRef = useRef(-1);
  const videoRef = useRef(null)

  const serverOffsetTime = useRef(undefined)
  const lastOffsetTime = useRef(0.0)
  const lastSaveTime = useRef(0)
  const afterRequestOffsetTime = useRef(false)

  const requestVideoTimeOffset = () => {
    if (shareid || serverOffsetTime.current != undefined) {
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
        if (dplayerRef.current && utils.isNumber(data)) {
          serverOffsetTime.current = Number(data)
          dplayerRef.current.seek(serverOffsetTime.current)
        }
        afterRequestOffsetTime.current = true;
      }).catch(error => {
        if (dplayerRef.current) {
          dplayerRef.current.seek(0)
        }
        afterRequestOffsetTime.current = true;
      });
  }

  const saveVideoTimeOffset = (offset, force) => {
    if (Math.abs(lastOffsetTime.current - offset) < 1.0 || !afterRequestOffsetTime.current) {
      return
    }
    lastOffsetTime.current = offset
    let now = Date.now()
    if (now - lastSaveTime.current < 2000 && !force) {
      return
    }
    lastSaveTime.current = now
    fetch(serverAddress + "/video/" + vidRef.current + "/set_offsettime/" + lastOffsetTime.current, {
      method: 'POST',
      mode: 'cors',
      credentials: "include",
      headers: {
      },
    })
  }

  const playNext = () => {
    for (let i = 0; i < videoItemListRef.current.length; i++) {
      if (videoItemListRef.current[i].id === itemId && i < videoItemListRef.current.length - 1) {
        navigateToItem(navigate, {}, videoItemListRef.current[i + 1].id, shareid)
      }
    }
  }

  const getDanmakuApi = () => {
    return serverAddress + "/video/" + vidRef.current + "/danmaku/"
  }

  useEffect(() => {
    queryItem(itemId, shareid, dispatch)
  }, [itemId, shareid, dispatch])

  useEffect(() => {
    if (parentItemId) {
      querySubItems({
        itemId: parentItemId, shareid, desc: categoryDesc, dispatch, callback: (subItems) => {
          setItems(subItems)
        }
      })
    }
  }, [parentItemId, shareid, dispatch, categoryDesc])

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
        url: urlPath,
        lang: lang,
        name: "zh-cn"
      })
      return null
    })
    subtitlesRef.current = cs
    setUrl(urlPath)
  }, [itemId, shareid, videoInfo]);

  useEffect(() => {
    if (urlRef.length === 0) {
      return
    }

    const req = new User.JoinChatRoomReq()
    const room = new User.Room()
    room.setType(User.Room.Type.DANMAKU)
    room.setId(vidRef.current)
    req.setRoom(room)
    var stream = userService.joinChatRoom(req)
    stream.on('data', function (res) {
      const chatMsgs = res.getChatMsgsList()
      let dans = chatMsgs.map((msg) => {
        try {
          let item = JSON.parse(msg.getMsg())
          return item
        } catch {
          return {}
        }
      })
      if (dplayerRef.current) {
        dplayerRef.current.onPushDanmaku(dans)
      }
    })
    stream.on('status', function (status) {
    });
    stream.on('end', function (end) {
      stream.cancel()
    });

    return () => {
      stream.cancel()
    };
  }, [urlRef])

  useEffect(() => {
    if (urlRef.length === 0) {
      return
    }

    window.Hls = Hls

    let hlsconfig = {
      debug: false,
      autoStartLoad: false,
      xhrSetup: function (xhr) {
        xhr.withCredentials = true; // do send cookies
      },
    }
    let danmakuUrl = getDanmakuApi()
    let options = {
      airplay: false,
      container: document.getElementById('dplayer'),
      video: {
        url: urlRef,
        type: 'hls',
      },
      pluginOptions: {
        hls: hlsconfig,
      },
      danmaku: {
        api: danmakuUrl,
        userId: userInfo ? userInfo.id : -1,
        userName: userInfo ? userInfo.userName : "",
        maximum: 10000,
        withCredentials: true,
      },
      nextButton: {
        show: true,
        callback: playNext,
      },
    }

    if (subtitlesRef.current.length > 0) {
      options.subtitle = {
        url: subtitlesRef.current,
        defaultSubtitle: 0,
        color: "orange",
        fontSize: "27px",
      }
    }

    const dp = new DPlayer({
      ...Object.assign({}, {
        lang: 'zh-cn',
        contextmenu: [
          {
            text: 'Author',
            link: 'https://github.com/pzc105/DPlayer'
          },
        ],
      }, options)
    });

    dp.on('canplay', (event) => {
      videoRef.current = event.target
      requestVideoTimeOffset()
      if (autoContinuedPlay) {
        dplayerRef.current.play()
      }
    })

    dp.on('timeupdate', (event) => {
      saveVideoTimeOffset(event.target.currentTime)
    })

    dp.on('fullscreen', () => {
      localStorage.setItem("fullscreen", "true")
      if (window.screen.orientation.lock) {
        window.screen.orientation.lock('landscape')
      }
    })

    dp.on('fullscreen_cancel', () => {
      localStorage.setItem("fullscreen", "false")
    })

    dp.on("ended", () => {
      saveVideoTimeOffset(0, true)
      if (autoContinuedPlay) {
        playNext()
      }
    })

    if (localStorage.getItem("fullscreen") === "true") {
      dp.fullScreen.request("browser")
    }

    dplayerRef.current = dp

    return (() => {
      dp.destroy()
    })
  }, [urlRef])

  const showGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))
  const closeGlobalChat = () => {
    dispatch(store.userSlice.actions.setOpenGlobalChat(false))
  }

  return (
    <Container >
      <CssBaseline />
      <Grid container spacing={2} sx={{ display: "flex" }}>
        <Grid item sm={12} md={8} >
          <Typography variant="button" component="div" noWrap>
            {item ? item.name : ""}
          </Typography>
          <div id="dplayer" />
        </Grid>
        <Hidden smDown>
          <Grid item sm={0} md={4} >
            <PlayList videoItemList={videoItemList} shareid={shareid} />
          </Grid>
        </Hidden>
      </Grid>
      {showGlobalChat ? <FloatingChat itemId={1} onClose={closeGlobalChat} /> : null}
    </Container>
  )
}
