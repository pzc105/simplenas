import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Container, Grid, Link, TextField, Button, InputAdornment,
  CssBaseline, Input, FormControl, FormLabel, FormHelperText, InputLabel, Box, Paper, List, ListItem, FormControlLabel, RadioGroup, Radio
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


export function DownloadRequest(props) {
  const [magnetUri, setMagnetUri] = useState('')
  const [downloadReq, setDownloadReq] = useState(null)
  const magneturiTag = "1"
  const torrentfileTag = "2"
  const [selectedTagValue, setSelectedTagValue] = useState(magneturiTag);

  function handleChange(e) {
    let uri = e.target.value
    setMagnetUri(uri)
    if (uri.length === 0) {
      return
    }
    const fileInput = document.getElementById('fileInput')
    if (fileInput) {
      fileInput.value = ''
    }
    var dr = new Bt.DownloadRequest()
    dr.setType(Bt.DownloadRequest.ReqType.MagnetUri)
    const encoder = new TextEncoder();
    dr.setContent(encoder.encode(uri))
    setDownloadReq(dr)
  }

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        setMagnetUri('')
        var dr = new Bt.DownloadRequest()
        dr.setType(Bt.DownloadRequest.ReqType.TORRENT)
        dr.setContent(new Uint8Array(reader.result))
        setDownloadReq(dr)
      };
      reader.readAsArrayBuffer(file);
    }
  };

  return (
    <Container>
      <FormControl component="fieldset">
        <RadioGroup
          style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', }}
          value={selectedTagValue}
          onChange={(e) => setSelectedTagValue(e.target.value)}>
          <FormControlLabel value={magneturiTag} control={<Radio />} label="MagnetURI" />
          <FormControlLabel value={torrentfileTag} control={<Radio />} label="Torrent文件" />
        </RadioGroup>
      </FormControl>

      {
        selectedTagValue == magneturiTag ?
          <Grid >
            <Grid item>
              <TextField
                margin="normal"
                fullWidth
                id="uri"
                label="magnet uri"
                value={magnetUri}
                onChange={handleChange}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <CloudDownloadIcon />
                    </InputAdornment>
                  ),
                }}
                autoFocus />
            </Grid>
            <Grid item mt={'3em'}>
              <BtHlsTaskPanel downloadReq={downloadReq} />
            </Grid>
          </Grid> :
          <Grid >
            <Grid item sx={{ marginTop: "1em" }}>
              <FormControl fullWidth>
                <Input id="fileInput" type="file" onChange={handleFileSelect} accept="image/*" />
              </FormControl>
            </Grid>
            <Grid item mt={'3em'}>
              <BtHlsTaskPanel downloadReq={downloadReq} />
            </Grid>
          </Grid>
      }
    </Container >
  )
}