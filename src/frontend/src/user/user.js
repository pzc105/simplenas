import React, { useState, useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Container, Typography, Paper, Button, Grid, CssBaseline } from '@mui/material';

import userService from '../rpcClient.js'
import * as User from '../prpc/user_pb.js'
import * as store from '../store.js'
import * as cateutils from '../category/utils.js'

export default function UserInfoPage() {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const shownUsrInfo = { 名称: userInfo["name"], Email: userInfo["email"] }
  return (
    <Container component="main" style={{ maxWidth: "90%" }} sx={{ backgroundColor: 'background.default' }}>
      <CssBaseline />
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
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }} sx={{ mt: "1em" }}>
      <Grid container sx={{ display: 'flex' }} >
        <Grid item xs={4}>
          <Typography variant="h6" whiteSpace={'pre'}>
            {name + ":"}
          </Typography>
        </Grid>
        <Grid item xs={8}>
          <Typography noWrap>
            {value}
          </Typography>
        </Grid>
      </Grid>
    </Paper>
  )
}

const SharedItems = () => {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const [sharedItems, setShareItems] = useState([])
  const [sharedIds, setSharedIds] = useState([])
  const items = useSelector((state) => store.selectCategoryItems(state, ...sharedIds))
  const dispatch = useDispatch()

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
        cateutils.queryItem(si.getItemId(), "", dispatch)
        let tmp = sharedIds
        tmp.push(si.getItemId())
        setSharedIds(tmp)
        return null
      })
      setShareItems(sharedItemsTmp)
    })
  }

  useEffect(() => {
    querySharedItems()
  }, [])

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
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }} sx={{ mt: "1em" }}>
      <Grid container sx={{ display: 'flex' }} >
        <Grid item xs={12}>
          <Grid container spacing={2}>
            <Grid item xs={4}>
              <Typography variant="h6" whiteSpace={'pre'}>
                分享链接:
              </Typography>
            </Grid>
            <Grid item xs={8}>
              {
                sharedItems.map((si, i) => {
                  const urlPrefix = "https://" + window.location.hostname + ":" + window.location.port + "/citem?"
                  const item = items[si.itemId]
                  return (
                    <Grid container spacing={2} key={i}>
                      <Grid item xs={8}>
                        <Typography>
                          {(i + 1) + "、" + (item ? item.name + ": " : "") + urlPrefix + "itemid=" + si.itemId + "&shareid=" + si.shareId}
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