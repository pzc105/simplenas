import React, { useState } from 'react';
import { Menu, MenuItem, List, Paper, Button, Box, Typography, Dialog } from '@mui/material';
import { styled } from "@mui/material/styles";
import LinearProgress from '@mui/material/LinearProgress';
import { useSelector, useDispatch } from 'react-redux';

import * as store from '../store.js'
import * as Bt from '../prpc/bt_pb.js'
import userService from '../rpcClient.js'
import BtVideosHandler from './btFileList.js'


const Name = styled(Button)({
  marginRight: (props) => props.theme.spacing(1),
  textOverflow: 'ellipsis',
});

const DownloadedSize = styled(Typography)({
  fontSize: '0.8rem',
  marginRight: (props) => props.theme.spacing(1),
});

const TotalSize = styled(Typography)({
  fontSize: '0.8rem',
});

const DownloadSpeed = styled(Typography)({
  fontSize: '0.8rem',
  marginLeft: 'auto',
});

function ProgressBar(props) {
  const infoHash = props.infoHash
  const dispatch = useDispatch()
  const torrent = useSelector(state => store.selectTorrent(state, infoHash))
  const torrentStatus = useSelector(state => store.selectTorrentStatus(state, infoHash))
  const [showVideos, setShowVideos] = useState(false)
  const [magnetUri, setMagnetUri] = useState("")
  const [sMagnetUri, setSMagnetUri] = useState(false)
  const onClick = () => {
    setShowVideos(!showVideos)
  }

  const RemoveTorrent = () => {
    var req = new Bt.RemoveTorrentReq()
    var ih = new Bt.InfoHash()
    ih.setVersion(infoHash.version)
    ih.setHash(infoHash.hash)
    req.setInfoHash(ih)
    userService.removeTorrent(req, {}, (err, dRespone) => {
      if (err != null) {
        console.log(err)
        return
      }
      dispatch(store.btSlice.actions.removeTorrent(infoHash))
    })
    setOpen(false)
  }

  const GetMagnetUri = () => {
    var req = new Bt.GetMagnetUriReq()
    var ih = new Bt.InfoHash()
    ih.setVersion(infoHash.version)
    ih.setHash(infoHash.hash)
    req.setType(Bt.GetMagnetUriReq.ReqType.INFOHASH)
    req.setInfoHash(ih)
    userService.getMagnetUri(req, {}, (err, rsp) => {
      if (err != null) {
        console.log(err)
        return
      }
      setSMagnetUri(true)
      setMagnetUri(rsp.getMagnetUri())
    })
    setOpen(false)
  }

  const setShowMagnetUri = (f) => {
    setSMagnetUri(false)
  }

  const [anchorPosition, setAnchorPosition] = useState({ left: 0, top: 0 });
  const [open, setOpen] = useState(false);

  const handleContextMenu = (event) => {
    event.preventDefault();
    setAnchorPosition({ left: event.clientX, top: event.clientY });
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  return (
    <Box onContextMenu={handleContextMenu}>
      {torrent ? <Box p={2} boxShadow={3} borderRadius={6} sx={{ display: 'flex', flexDirection: 'column' }}>
        <Name onClick={onClick}>{torrent.name}</Name>
        <Box sx={{ alignItems: 'center' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
            <DownloadedSize>{`${(torrentStatus.totalDone / 1024 / 1024).toFixed(2)} MB`}</DownloadedSize>
            <TotalSize>{`${(torrent.totalSize / 1024 / 1024).toFixed(2)} MB`}</TotalSize>
          </Box>
          <Box>
            <LinearProgress variant="determinate" value={torrentStatus.progress * 100} />
          </Box>
        </Box>
        <DownloadSpeed>{`${(torrentStatus.downloadPayloadRate / 1000).toFixed(2)} KB/s`}</DownloadSpeed>
      </Box> : null}
      {showVideos ? <BtVideosHandler infoHash={torrent.infoHash} /> : null}

      <Dialog open={sMagnetUri} onClose={() => setShowMagnetUri(false)}>
        <div style={{ padding: '16px' }}>
          {magnetUri}
        </div>
      </Dialog>

      <Menu
        anchorReference="anchorPosition"
        anchorPosition={anchorPosition}
        open={open}
        onClose={handleClose}
      >
        <MenuItem onClick={RemoveTorrent}>删除</MenuItem>
        <MenuItem onClick={GetMagnetUri}>获取magnet</MenuItem>
      </Menu>
    </Box>
  );
}

export function ProgressLists() {
  const infoHashs = useSelector(state => store.selectInfoHashs(state))

  return (
    <Paper style={{ maxHeight: '90vh', overflow: 'auto' }}>
      {
        infoHashs.map((infoHash) =>
          <List key={infoHash.hash}>
            <ProgressBar infoHash={infoHash} key={infoHash.hash} />
          </List>
        )
      }
    </Paper>
  )
}