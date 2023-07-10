import React, { useState, useEffect } from 'react';
import { Container, Grid, Link, TextField, Button, InputAdornment, CssBaseline } from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";
import Typography from '@mui/material/Typography';
import { useDispatch } from 'react-redux';

import SideUtils from './sideUtils.js';
import { btSlice, eventSlice } from './store.js'
import * as Bt from './prpc/bt_pb.js'
import userService from './rpcClient.js'
import FileUpload from './uploadTorrent.js'
import { ProgressLists } from './downloadlist.js'

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

  useEffect(() => {
    const statusRequest = new Bt.StatusRequest()
    var stream = userService.onStatus(statusRequest)
    stream.on('data', function (sResponse) {
      const trs = sResponse.getStatusArrayList()
      trs.map((t) => {
        dispatch(btSlice.actions.updateTorrent(t.toObject()))
        return null
      })
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
    dispatch(eventSlice.actions.setDownloadPageMouse(true))
  }
  const handleListItemMouseUp = () => {
    dispatch(eventSlice.actions.setDownloadPageMouse(false))
  }

  return (
    <DownloadContainer
      onMouseDown={handleListItemMouseDown}
      onMouseUp={handleListItemMouseUp}>
      <CssBaseline />
      <SideUtils name="下载" child={DownloadRequest()} />
      <LeftColumn>
        <TopLeftArea sx={{ backgroundColor: 'background.default' }}>
          <TorrentNavigation />
        </TopLeftArea>
      </LeftColumn>
      <RightColumn sx={{ backgroundColor: 'background.default' }}>
        <ProgressLists />
      </RightColumn>
    </DownloadContainer>
  )
}

const TorrentNavigation = () => {

  return (
    <Container >
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://yts.mx/" target="_blank" rel="noopener" >
          YTS
        </Link>
        可下载国外的电影与连续剧，需要科学上网
      </Typography>
      <Typography sx={{ marginTop: '1em' }}>
        <Link href="https://bt.acgzero.com/" target="_blank" rel="noopener" >
          零度动漫
        </Link>
        可以下载番剧
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

function DownloadRequest(props) {
  const [magnetUri, setMagnetUri] = useState('')
  function handleChange(e) {
    setMagnetUri(e.target.value)
  }

  function requestDownload(e) {
    e.preventDefault()
    var dr = new Bt.DownloadRequest()
    dr.setType(Bt.DownloadRequest.ReqType.MagnetUri)
    const encoder = new TextEncoder();
    dr.setContent(encoder.encode(magnetUri))
    userService.download(dr, {}, (err, dRespone) => {
      if (err != null) {
        console.log(err)
        return
      }
      console.log(dRespone)
    })
  }

  const onUploadFile = (fileData) => {
    var dr = new Bt.DownloadRequest()
    dr.setType(Bt.DownloadRequest.ReqType.TORRENT)
    dr.setContent(fileData)
    userService.download(dr, {}, (err, dRespone) => {
      if (err != null) {
        console.log(err)
        return
      }
      console.log(dRespone)
    })
  }

  return (
    <Container sx={{ width: "30vw" }}>
      <Grid >
        <Grid item>
          <TextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="uri"
            label="address"
            name="address"
            onChange={handleChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <CloudDownloadIcon />
                </InputAdornment>
              ),
            }}
            autoFocus />
          <Grid >
            <Button
              type="submit"
              fullWidth
              variant="contained"
              color="primary"
              disabled={magnetUri === ""}
              onClick={requestDownload}>
              通过磁力下载
            </Button>
          </Grid>
        </Grid>
        <Grid item sx={{ marginTop: "1em" }}>
          <FileUpload title={"选择种子文件"} onUpload={onUploadFile} />
        </Grid>
      </Grid>
    </Container >
  )
}