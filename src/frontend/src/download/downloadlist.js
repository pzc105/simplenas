import React, { useState, useEffect } from 'react';
import {
  Menu, MenuItem, List, Paper, Button, Box, Typography, Dialog, Grid, TextField, InputAdornment,
  FormControl, FormLabel, RadioGroup, FormControlLabel, Radio
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { styled } from "@mui/material/styles";
import LinearProgress from '@mui/material/LinearProgress';
import { useSelector, useDispatch } from 'react-redux';
import { CopyToClipboard } from 'react-copy-to-clipboard';

import * as store from '../store.js'
import * as Bt from '../prpc/bt_pb.js'
import userService from '../rpcClient.js'
import BtVideosHandler from './btFileList.js'


const Name = styled(Typography)({
  marginRight: (props) => props.theme.spacing(1),
  textOverflow: 'ellipsis',
  cursor: 'pointer',
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
        <CopyToClipboard text={magnetUri}>
          <Button>复制</Button>
        </CopyToClipboard>
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
  const torrents = useSelector(state => store.selectTorrents(state))
  const btstatus = useSelector(state => store.selectAllBtStatus(state))
  const [sortedTorrents, setSortedTorrents] = useState([])
  const [searchText, setSearchText] = useState('')
  const downloadingTag = "1"
  const downloadedTag = "2"
  const [selectedTagValue, setSelectedTagValue] = useState(downloadingTag);

  const onSearchText = (e) => {
    setSearchText(e.target.value)
  }

  const isDownloading = (st) => {
    if (!st) {
      return false
    }
    if (st.state === Bt.BtStateEnum.DOWNLOADING
      || st.state === Bt.BtStateEnum.DOWNLOADING_METADATA
      || st.state === Bt.BtStateEnum.CHECKING_FILES) {
      return true
    }
    return false
  }

  useEffect(() => {
    if (!torrents) {
      return
    }
    let tmp = []
    let emptyNameTs = []
    for (let t of Object.values(torrents)) {
      if (searchText.length > 0) {
        let existedWords = searchText.split(" ")
        let notfound = false
        for (const text of existedWords) {
          if (t.name.indexOf(text) == -1) {
            notfound = true
          }
        }
        if (notfound) {
          continue
        }
      }
      if (t.name.length === 0) {
        emptyNameTs.push(t)
      } else {
        let status = btstatus[t.infoHash.hash]
        if (selectedTagValue == downloadingTag && isDownloading(status)) {
          tmp.push(t)
        } else if (selectedTagValue == downloadedTag && !isDownloading(status)) {
          tmp.push(t)
        }
      }
    }
    tmp.sort((a, b) => {
      if (a.name < b.name) {
        return -1;
      }
      if (a.name > b.name) {
        return 1;
      }
      return 0;
    })
    emptyNameTs.sort((a, b) => {
      if (a.infoHash.hash < b.infoHash.hash) {
        return -1;
      }
      if (a.infoHash.hash > b.infoHash.hash) {
        return 1;
      }
      return 0;
    })
    if (selectedTagValue == downloadingTag) {
      tmp.push(...emptyNameTs)
    }
    setSortedTorrents(tmp)
  }, [torrents, btstatus, searchText, selectedTagValue])

  return (
    <Paper style={{ maxHeight: '90vh', overflow: 'auto', marginLeft: "1em" }}>
      <TextField
        onChange={onSearchText}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon />
            </InputAdornment>
          ),
        }} />
      <FormControl component="fieldset">
        <RadioGroup
          style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', }}
          value={selectedTagValue}
          onChange={(e) => setSelectedTagValue(e.target.value)}>
          <FormControlLabel value={downloadingTag} control={<Radio />} label="下载中" />
          <FormControlLabel value={downloadedTag} control={<Radio />} label="已下载" />
        </RadioGroup>
      </FormControl>
      {
        sortedTorrents.map((t) =>
          <Box key={t.infoHash.hash}>
            <ProgressBar infoHash={t.infoHash} key={t.infoHash.hash} />
          </Box>
        )
      }
    </Paper >
  )
}
