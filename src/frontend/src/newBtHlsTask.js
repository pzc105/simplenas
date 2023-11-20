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

export default function BtHlsTaskPanel({ downloadReq, onCreate }) {
  const parentIdRef = useRef(-1)

  const newBtHlsTask = () => {
    if (!downloadReq) {
      return
    }
    console.log(downloadReq)
    var req = new User.NewBtHlsTaskReq()
    req.setReq(downloadReq)
    req.setCategoryParentId(parentIdRef.current)
    userService.newBtHlsTask(req, {}, (err, respone) => {
      if (onCreate) {
        onCreate()
      }
      if (err != null) {
        console.log(err)
        return
      }
    })
  }

  return (
    <>
      <Grid container >
        <Grid item xs={10} className='draggableWindow'>
          <Typography sx={{ userSelect: 'none', ml: "1em" }}>
            {"选择保存目录"}
          </Typography>
        </Grid>
      </Grid>
      <Grid container sx={{ display: 'flex' }} alignItems="center" justify="center">
        <FolderSelector select={(id) => parentIdRef.current = id} />
        <Button variant="outlined" color="primary" onClick={newBtHlsTask}>创建任务</Button>
      </Grid>
    </>
  )
}