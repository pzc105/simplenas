import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Container, Grid, Link, TextField, Button, Paper, CssBaseline } from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";
import Typography from '@mui/material/Typography';
import { useDispatch, useSelector } from 'react-redux';


import * as Bt from './prpc/bt_pb.js'
import * as User from './prpc/user_pb.js'
import userService from './rpcClient.js'
import FolderSelector from './download/folderSelector.js'

export default function BtHlsTaskPanel({ downloadReq }) {

  const [btInfo, setBtInfo] = useState(null)
  const [parentId, setParentId] = useState(-1)

  const newBtHlsTask = () => {
    var req = new User.NewBtHlsTaskReq()
    req.setReq(downloadReq)
    req.setCategoryParentId(parentId)
    userService.newBtHlsTask(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
    })
  }

  return (
    <>
      <CssBaseline />
      <Paper style={{}}>
        <FolderSelector select={(id) => setParentId(id)} />
        <Button onClick={newBtHlsTask}>创建</Button>
      </Paper>
    </>
  )
}