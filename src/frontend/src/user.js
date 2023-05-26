import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Container, Typography, Avatar, Button, CssBaseline, TextField, FormControlLabel, Checkbox, Link, Grid, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import MuiAlert from '@mui/lab/Alert';
import { styled } from "@mui/material/styles";
import { useNavigate } from 'react-router-dom';

import * as utils from './utils.js'
import userService from './rpcClient.js'
import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import * as store from './store.js'

export default function UserInfoPage({ }) {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const shownUsrInfo = { 名称: userInfo["name"], Email: userInfo["email"] }
  return (
    <Container component="main" maxWidth="xs">
      <UserInfoList userInfo={shownUsrInfo} />
      <SharedItems />
    </Container>
  )
}

const UserInfoList = ({ userInfo }) => {
  return (
    <Container>
      {
        Object.keys(userInfo).map((k) => {
          return <UserListItem key={k} name={k} value={userInfo[k]} />
        })
      }
    </Container>
  )
}

const UserListItem = ({ name, value }) => {
  return (
    <Container sx={{ display: 'flex' }}>
      <Typography whiteSpace={'pre'}>
        {name + ":  "}
      </Typography>
      <Typography>
        {value}
      </Typography>
    </Container>
  )
}

const SharedItems = () => {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const [sharedUrls, setShareUrls] = useState([])

  useEffect(() => {
    let req = new User.QuerySharedItemsReq()
    req.setUserId(userInfo.id)
    userService.querySharedItems(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      const urlPrefix = "https://" + window.location.hostname + ":" + window.location.port + "/citem?"
      let sharedUrlsTmp = []
      res.getSharedItemsList().map((si) => {
        sharedUrlsTmp.push(urlPrefix + "itemid=" + si.getItemId() + "&shareid=" + si.getShareId())
      })

      setShareUrls(sharedUrlsTmp)
    })
  }, [userInfo])

  return (
    <Grid Container sx={{ display: 'flex' }} >
      <Grid item xs={12}>
        <Grid item xs={3}>
          <Typography >
            分享链接:
          </Typography>
        </Grid>
        <Grid item xs={9}>
          {
            sharedUrls.map((url, i) => {
              return (
                <Box key={i}>
                  <Typography>
                    {(i + 1) + ". " + url}
                  </Typography>
                </Box>
              )
            })
          }
        </Grid>
      </Grid>
    </Grid>
  )
}