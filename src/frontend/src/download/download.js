import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Container, Grid, Link, TextField, Button, InputAdornment,
  CssBaseline, Input, FormControl, FormLabel, FormHelperText, InputLabel, Box, Paper, List, ListItem
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";
import Typography from '@mui/material/Typography';
import { useDispatch, useSelector } from 'react-redux';

import SideUtils from '../sideManager.js';
import * as store from '../store.js'
import * as Bt from '../prpc/bt_pb.js'
import * as User from '../prpc/user_pb.js'
import * as utils from '../utils.js'
import userService from '../rpcClient.js'
import * as router from '../router.js'
import FileUpload from '../uploadTorrent.js'
import { ProgressLists } from './downloadlist.js'
import { FloatingChat, DraggableDialog } from '../chat/chat.js';
import BtHlsTaskPanel from '../newBtHlsTask.js'
import { DownloadRequest } from './downloadReq.js'

const DownloadContainer = styled('div')(({ theme }) => ({
  display: 'flex',
  height: '94vh',
}))

const LeftColumn = styled('div')(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
}))

const TopLeftArea = styled('div')(({ theme }) => ({
  flex: 1,
  backgroundColor: '#f2f2f2',
}))

const RightColumn = styled('div')(({ theme }) => ({
  width: '65%',
  backgroundColor: '#e5e5e5',
}))

export default function Download() {
  const dispatch = useDispatch()
  const showGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))
  const torrents = useSelector(state => store.selectTorrents(state))

  const queryTorrens = () => {
    var req = new User.GetTorrentsReq()
    userService.getTorrents(req, {}, (err, rsp) => {
      if (err != null) {
        console.log(err)
        return
      }
      let ts = {}
      rsp.getTorrentInfoList().map((t) => {
        let torrent = t.toObject()
        ts[torrent.infoHash.hash] = torrent
        return null
      })
      dispatch(store.btSlice.actions.updateTorrents(ts))
    })
  }

  useEffect(() => {
    queryTorrens()

    const statusRequest = new Bt.BtStatusRequest()
    var stream = userService.onBtStatus(statusRequest)
    stream.on('data', function (sResponse) {
      const trs = sResponse.getStatusArrayList()
      let needQueryTorrents = false
      trs.map((st) => {
        let status = st.toObject()
        let t = torrents[status.infoHash.hash]
        if (!t || status.state === Bt.BtStateEnum.DOWNLOADING && t.name.length === 0) {
          needQueryTorrents = true
        }
        dispatch(store.btSlice.actions.updateTorrentStatus(status))
        return null
      })
      if (needQueryTorrents) {
        queryTorrens()
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
  }, [dispatch])

  const handleListItemMouseDown = () => {
    dispatch(store.eventSlice.actions.setDownloadPageMouse(true))
  }
  const handleListItemMouseUp = () => {
    dispatch(store.eventSlice.actions.setDownloadPageMouse(false))
  }

  const closeGlobalChat = () => {
    dispatch(store.userSlice.actions.setOpenGlobalChat(false))
  }

  return (
    <DownloadContainer
      onMouseDown={handleListItemMouseDown}
      onMouseUp={handleListItemMouseUp}>
      <CssBaseline />
      <LeftColumn>
        <TopLeftArea sx={{ backgroundColor: 'background.default' }}>
          <DownloadRequest />
          <PeerInfoList />
        </TopLeftArea>
      </LeftColumn>
      <RightColumn sx={{ backgroundColor: 'background.default' }}>
        <ProgressLists />
      </RightColumn>
      {showGlobalChat ? <FloatingChat itemId={1} onClose={() => closeGlobalChat()} /> : null}
    </DownloadContainer>
  )
}

const TorrentNavigation = () => {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const navigate = useNavigate()
  const navigateToMagnetPage = () => {
    router.navigate2mgnetshares(navigate, userInfo.magnetRootId)
  }

  return (
    <Container >
      <Typography sx={{ marginTop: '1em' }}>
        <Link onClick={navigateToMagnetPage} target="_blank" rel="noopener" sx={{ cursor: 'pointer' }}>
          magnet uri
        </Link>
        分享中心
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://yts.mx/" target="_blank" rel="noopener" >
          YTS
        </Link>
        电影与连续剧
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://bt.acgzero.com/" target="_blank" rel="noopener" >
          零度动漫
        </Link>
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://bbs.xfsub.org/" target="_blank" rel="noopener" >
          动漫下载
        </Link>
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://bbs.acgrip.com/" target="_blank" rel="noopener" >
          字幕论坛
        </Link>
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://thepiratebay.org/index.html" target="_blank" rel="noopener" >
          https://thepiratebay.org/index.html
        </Link>
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://1337x.to/home/" target="_blank" rel="noopener" >
          https://1337x.to/home/
        </Link>
      </Typography>
    </Container >
  )
}

const PeerInfoList = () => {
  const dispatch = useDispatch()
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const [peerInfos, setPeerInfo] = useState([])
  const [sortedPeerInfos, setSortedPeerInfos] = useState([])

  useEffect(() => {
    if (!utils.isAdmin(userInfo)) {
      return
    }
    const request = () => {
      const req = new Bt.GetPeerInfoReq()
      userService.getPeerInfo(req, {}, (err, rsp) => {
        if (err != null) {
          return
        }
        let peerInfos = []
        rsp.getPeerInfosList().map((p) => {
          peerInfos.push(p.toObject())
          return null
        })
        dispatch(store.btSlice.actions.updateDisplayPeerInfos(peerInfos))
        setPeerInfo(peerInfos)
      })
    }

    let myInterval = setInterval(request, 3000)
    return () => {
      clearInterval(myInterval);
    };
  }, [userInfo])

  useEffect(() => {
    let tmp = [...peerInfos]
    tmp.sort((a, b) => {
      if (a.upSpeed != b.upSpeed) {
        return b.upSpeed - a.upSpeed
      }
      return b.peerAddr < a.peerAddr ? -1 : 1
    })
    setSortedPeerInfos(tmp)
  }, [peerInfos])

  return (
    <Container>
      {
        sortedPeerInfos.length > 0 ? <Paper style={{ maxHeight: "30vh", overflow: 'auto' }}>
          <List>
            {
              sortedPeerInfos.map((p, i) => {
                return (
                  <ListItem key={i}>
                    {"name:" + p.client + ", dr:" + `${(p.downSpeed / 1024 / 1024).toFixed(2)} MB` + ", ur:" + `${(p.upSpeed / 1024 / 1024).toFixed(2)} MB` + ", addr:" + p.peerAddr}
                  </ListItem>
                )
              })
            }
          </List>
        </Paper> : null
      }
    </Container>
  )
}