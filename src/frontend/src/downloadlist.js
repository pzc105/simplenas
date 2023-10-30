import React, { useState } from 'react';
import { Menu, MenuItem, List, Paper, Button, Box, Typography } from '@mui/material';
import { styled } from "@mui/material/styles";
import LinearProgress from '@mui/material/LinearProgress';
import { useSelector, useDispatch } from 'react-redux';

import { btSlice, selectTorrent, selectInfoHashs } from './store.js'
import * as Bt from './prpc/bt_pb.js'
import userService from './rpcClient.js'
import BtVideosHandler from './btVideoMetadata.js'


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
  const torrent = useSelector(state => selectTorrent(state, infoHash))
  const [showVideos, setShowVideos] = useState(false)
  const onClick = () => {
    setShowVideos(!showVideos)
  }

  const RemoveTorrent = () => {
    var req = new Bt.RemoveTorrentReq()
    var i = new Bt.InfoHash()
    i.setVersion(infoHash.version)
    i.setHash(infoHash.hash)
    req.setInfoHash(i)
    userService.removeTorrent(req, {}, (err, dRespone) => {
      if (err != null) {
        console.log(err)
        return
      }
      dispatch(btSlice.actions.removeTorrent(infoHash))
    })
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
      <Box p={2} boxShadow={3} borderRadius={6} sx={{ display: 'flex', flexDirection: 'column' }}>
        <Name onClick={onClick}>{torrent.name}</Name>
        <Box sx={{ alignItems: 'center' }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', width: '100%' }}>
            <DownloadedSize>{`${(torrent.totalDone / 1024 / 1024).toFixed(2)} MB`}</DownloadedSize>
            <TotalSize>{`${(torrent.total / 1024 / 1024).toFixed(2)} MB`}</TotalSize>
          </Box>
          <Box>
            <LinearProgress variant="determinate" value={torrent.progress * 100} />
          </Box>
        </Box>
        <DownloadSpeed>{`${(torrent.downloadPayloadRate / 1000).toFixed(2)} KB/s`}</DownloadSpeed>
      </Box>
      {showVideos ? <BtVideosHandler infoHash={torrent.infoHash} /> : null}

      <Menu
        anchorReference="anchorPosition"
        anchorPosition={anchorPosition}
        open={open}
        onClose={handleClose}
      >
        <MenuItem onClick={RemoveTorrent}>删除</MenuItem>
      </Menu>
    </Box>
  );
}

export function ProgressLists() {
  const infoHashs = useSelector(state => selectInfoHashs(state))

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