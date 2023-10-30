import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper, List, ListItem, Link
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";
import Draggable from 'react-draggable';
import CloseIcon from '@mui/icons-material/Close';
import { CopyToClipboard } from 'react-copy-to-clipboard';

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'
import SideUtils from './sideUtils.js';
import ChatPanel from './chat.js';
import SubtitleUploader from './uploadSubtitle.js';
import * as category from './category.js'
import * as router from './router.js'

import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { serverAddress } from './rpcClient.js'

const MagnetItems = ({ parentId }) => {
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const items = useSelector((state) => store.selectCategorySubItems(state, parentId))
  console.log(parentId, items)

  const delItem = (id) => {
    var req = new User.DelMagnetCategoryReq()
    req.setId(id)
    userService.delMagnetCategory(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
      queryMagnet(parentId, dispatch)
    })
  }

  return (
    <Container>
      <List>
        {items ?
          items.map((item) => (
            <ListItem key={item.id}>
              {
                category.isDirectory(item) ?
                  <Container>
                    <Grid container sx={{ display: 'flex' }} >
                      <Grid item xs={12}>
                        <Grid container spacing={2}>
                          <Grid item xs={2}>
                            <Typography whiteSpace={'pre'}>
                              {"分类: "}
                            </Typography>
                          </Grid>
                          <Grid item xs={8}>
                            <Link onClick={() => router.navigate2mgnetshares(navigate, item.id)} sx={{ cursor: 'pointer' }}>
                              {item.name}
                            </Link>
                          </Grid>
                          <Grid item xs={2}>
                            <Button onClick={() => delItem(item.id)}>删除</Button>
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  </Container> :
                  <Container>
                    <Grid container sx={{ display: 'flex' }} >
                      <Grid item xs={12}>
                        <Grid container spacing={2}>
                          <Grid item xs={2}>
                            <Typography whiteSpace={'pre'}>
                              {"magnet uri: "}
                            </Typography>
                          </Grid>
                          <Grid item xs={8}>
                            <CopyToClipboard text={item.resourcePath}>
                              <Button>{item.resourcePath}</Button>
                            </CopyToClipboard>
                          </Grid>
                          <Grid item xs={2}>
                            <Button onClick={() => delItem(item.id)}>删除</Button>
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  </Container>
              }
            </ListItem>
          )) : null
        }
      </List>
    </Container>
  )
}

export default function MagnetSharesPage() {
  const dispatch = useDispatch()

  const location = useLocation()
  const searchParams = new URLSearchParams(location.search)
  const itemId = searchParams.get('itemid') ? Number(searchParams.get('itemid')) : -1

  useEffect(() => {
    queryMagnet(itemId, dispatch)
  }, [itemId])

  return (
    <MagnetContainer>
      <CssBaseline />
      {
        <SideUtils
          name="管理"
          child={Manager({ parentId: itemId })}
        />
      }
      <MagnetItems parentId={itemId} />
    </MagnetContainer>
  )
}

const Manager = ({ parentId }) => {
  const dispatch = useDispatch()
  const [magnetUri, setMagnetUri] = useState('')
  function handleChange(e) {
    setMagnetUri(e.target.value)
  }
  function saveMagnetUri(e) {
    e.stopPropagation()
    if (magnetUri.length === 0) {
      return
    }
    var req = new User.AddMagnetUriReq()
    req.setCategoryId(parentId)
    req.setMagnetUri(magnetUri)
    userService.addMagnetUri(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
      queryMagnet(parentId, dispatch)
    })
  }

  return (<Container maxWidth="xs">
    <Container>
      <Grid container>
        <Grid item xs={12}>
          <TextField
            fullWidth
            onChange={handleChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <CloudDownloadIcon />
                </InputAdornment>
              ),
            }}
            autoFocus />
          <Grid item xs={12}>
            <Button
              fullWidth
              color="primary"
              onClick={saveMagnetUri}
              variant="contained">
              保存Magnet Uri
            </Button>
          </Grid>
        </Grid>
      </Grid>
    </Container>
  </Container>
  )
}

const MagnetContainer = styled('div')({
  display: 'flex', /* 将子元素布局为行 */
  height: '94vh', /* 页面铺满整个视窗 */
})

const queryMagnet = (id, dispatch) => {
  var req = new User.QueryMagnetReq()
  req.setParentId(id)
  userService.queryMagnet(req, {}, (err, respone) => {
    if (err != null) {
      console.log(err)
      return
    }
    respone.getItemsList().map((i) => {
      let obj = i.toObject()
      dispatch(store.categorySlice.actions.updateItem(obj))
      return null
    })
  })
}