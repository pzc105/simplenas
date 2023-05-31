import React, { useState, useEffect } from 'react';
import { useSelector } from 'react-redux';
import { Container, Typography, Paper, Button, Grid } from '@mui/material';

import userService from './rpcClient.js'
import * as User from './prpc/user_pb.js'
import * as store from './store.js'

export default function UserInfoPage() {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const shownUsrInfo = { 名称: userInfo["name"], Email: userInfo["email"] }
  return (
    <Container component="main" style={{ maxWidth: "50%" }}>
      <UserInfoList userInfo={shownUsrInfo} />
      <SharedItems />
    </Container>
  )
}

const UserInfoList = ({ userInfo }) => {
  return (
    <Container style={{ maxWidth: "100%" }}>
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
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }}>
      <Grid container sx={{ display: 'flex' }} >
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={2}>
              <Typography whiteSpace={'pre'}>
                {name + ":"}
              </Typography>
            </Grid>
            <Grid item xs={10}>
              <Typography>
                {value}
              </Typography>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Paper>
  )
}

const SharedItems = () => {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const [sharedItems, setShareItems] = useState([])

  const querySharedItems = () => {
    let req = new User.QuerySharedItemsReq()
    req.setUserId(userInfo.id)
    userService.querySharedItems(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      let sharedItemsTmp = []
      res.getSharedItemsList().map((si) => {
        sharedItemsTmp.push(si.toObject())
        return null
      })
      setShareItems(sharedItemsTmp)
    })
  }

  useEffect(() => {
    querySharedItems()
  }, [userInfo])

  const delShareItem = (shareid) => {
    let req = new User.DelSharedItemReq()
    req.setShareId(shareid)
    userService.delSharedItem(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      querySharedItems()
    })
  }

  return (
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }}>
      <Grid container sx={{ display: 'flex' }} >
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={2}>
              <Typography >
                分享链接:
              </Typography>
            </Grid>
            <Grid item xs={10}>
              {
                sharedItems.map((si, i) => {
                  const urlPrefix = "https://" + window.location.hostname + ":" + window.location.port + "/citem?"
                  return (
                    <Grid container spacing={2} key={i}>
                      <Grid item xs={8}>
                        <Typography>
                          {(i + 1) + "、" + urlPrefix + "itemid=" + si.itemId + "&shareid=" + si.shareId}
                        </Typography>
                      </Grid>
                      <Grid item xs={4}>
                        <Button onClick={() => delShareItem(si.shareId)} >
                          删除
                        </Button>
                      </Grid>
                    </Grid>
                  )
                })
              }
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </Paper>
  )
}