import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper, List, ListItem, Link, Dialog
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import SearchIcon from '@mui/icons-material/Search';
import { styled } from "@mui/material/styles";
import Draggable from 'react-draggable';
import CloseIcon from '@mui/icons-material/Close';
import { CopyToClipboard } from 'react-copy-to-clipboard';

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'
import SideUtils from './sideManager.js';
import ChatPanel from './chat/chat.js';
import SubtitleUploader from './uploadSubtitle.js';
import * as category from './category.js'
import * as router from './router.js'
import { FloatingChat } from './chat/chat.js';
import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { serverAddress } from './rpcClient.js'
import UnifiedPage from './page.js'


export default function MagnetSharesPage() {
  const dispatch = useDispatch()
  const showGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))
  const location = useLocation()
  const searchParams = new URLSearchParams(location.search)
  const itemId = searchParams.get('itemid') ? Number(searchParams.get('itemid')) : -1
  const pageRows = 10
  const pageNum = useSelector((state) => store.selectMagnetSharesPageNum(state, itemId))
  const totalRows = useSelector((state) => store.selectMagnetSharesTotalRows(state, itemId))
  console.log("pageNum", totalRows, pageNum)

  useEffect(() => {
    if (!pageNum && pageNum !== 0) {
      return
    }
    queryMagnet(dispatch, itemId, "", pageNum, pageRows)
  }, [itemId, pageNum, pageRows])

  const closeGlobalChat = () => {
    dispatch(store.userSlice.actions.setOpenGlobalChat(false))
  }

  return (
    <MagnetContainer>
      <CssBaseline />
      <SideUtils
        name="管理"
        child={Manager({ parentId: itemId })}
      />
      <Grid container sx={{ display: 'flex' }} alignItems="center" justify="center">
        <Grid item xs={12}>
          <MagnetItems parentId={itemId} pageNum={pageNum - 1} rows={pageRows} />
        </Grid>
        <Grid item xs={12}>
          <Container>
            <UnifiedPage PageTotalCount={parseInt(totalRows / pageRows + 0.5)} PageNum={parseInt(pageNum + 1)} onPage={(n) => dispatch(store.categorySlice.actions.updateMagnetSharesPageNum({ num: n - 1, parentId: itemId }))} />
          </Container>
        </Grid>
      </Grid>
      {showGlobalChat ? <FloatingChat itemId={1} onClose={closeGlobalChat} /> : null}
    </MagnetContainer>
  )
}

const Manager = ({ parentId, pageNum, rows }) => {
  const dispatch = useDispatch()
  const [magnetUri, setMagnetUri] = useState('')
  const [magnetUriIntroduce, setMagnetUriIntroduce] = useState('')
  const [newCategory, setNewCategory] = useState('')
  const [newCategoryIntroduce, setNewCategoryIntroduce] = useState('')

  const handleMagnetUriChange = (e) => {
    setMagnetUri(e.target.value)
  }
  const handleMagnetUriIntroduceChangle = (e) => {
    setMagnetUriIntroduce(e.target.value)
  }

  const handleNewCategory = (e) => {
    setNewCategory(e.target.value)
  }
  const handleNewCategoryIntroduce = (e) => {
    setNewCategoryIntroduce(e.target.value)
  }

  const saveMagnetUri = (e) => {
    e.stopPropagation()
    if (magnetUri.length === 0) {
      return
    }
    var req = new User.AddMagnetUriReq()
    req.setCategoryId(parentId)
    req.setMagnetUri(magnetUri)
    req.setIntroduce(magnetUriIntroduce)
    userService.addMagnetUri(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
      queryMagnet(dispatch, parentId, "", pageNum, rows)
    })
  }

  const addMagnetCategory = (e) => {
    e.stopPropagation()
    if (newCategory.length === 0) {
      return
    }
    var req = new User.AddMagnetCategoryReq()
    req.setParentId(parentId)
    req.setCategoryName(newCategory)
    req.setIntroduce(newCategoryIntroduce)
    userService.addMagnetCategory(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
      queryMagnet(dispatch, parentId, "", pageNum, rows)
    })
  }

  return (
    <Container maxWidth="xs">
      <Container>
        <Grid container>
          <Grid item xs={12}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="magnet uri"
                margin="normal"
                onChange={handleMagnetUriChange}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <CloudDownloadIcon />
                    </InputAdornment>
                  ),
                }}
                autoFocus />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="introduce"
                margin="normal"
                onChange={handleMagnetUriIntroduceChangle}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <CloudDownloadIcon />
                    </InputAdornment>
                  ),
                }} />
            </Grid>
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
      <Container sx={{ mt: '1em' }}>
        <Grid item xs={12}>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="category name"
              margin="normal"
              onChange={handleNewCategory}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <CloudDownloadIcon />
                  </InputAdornment>
                ),
              }} />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="introduce"
              margin="normal"
              onChange={handleNewCategoryIntroduce}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <CloudDownloadIcon />
                  </InputAdornment>
                ),
              }} />
          </Grid>
          <Grid item xs={12}>
            <Button
              fullWidth
              color="primary"
              onClick={addMagnetCategory}
              variant="contained">
              新增分类
            </Button>
          </Grid>
        </Grid>
      </Container >
    </Container >
  )
}

const MagnetContainer = styled('div')({
  display: 'flex', /* 将子元素布局为行 */
  height: '94vh', /* 页面铺满整个视窗 */
})

const queryMagnet = (dispatch, id, searchWords, pageNum, rows) => {
  console.log("q", pageNum, rows)
  var req = new User.QueryMagnetReq()
  req.setParentId(id)
  req.setSearchCond(searchWords)
  req.setPageNum(pageNum)
  req.setRows(rows)
  userService.queryMagnet(req, {}, (err, respone) => {
    if (err != null) {
      console.log(err)
      return
    }
    let objs = []
    respone.getItemsList().map((i) => {
      let obj = i.toObject()
      if (obj.id != id) {
        objs.push(obj)
      }
      dispatch(store.categorySlice.actions.updateItem(obj))
      return null
    })
    dispatch(store.categorySlice.actions.updateMagnetSharesTotalRows({ totalRows: respone.getTotalRowCount(), parentId: id }))
    dispatch(store.categorySlice.actions.updateMagnetSharesItems({ items: objs, parentId: id }))
  })
}

const MagnetItems = ({ parentId, pageNum, rows }) => {
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const items = useSelector((state) => store.selectMagnetSharesItems(state))
  const [searchWords, setSearchWords] = useState('')

  const [copyDialogOpen, setCopyDialogOpen] = useState(false)

  const delItem = (id) => {
    var req = new User.DelMagnetCategoryReq()
    req.setId(id)
    userService.delMagnetCategory(req, {}, (err, respone) => {
      if (err != null) {
        console.log(err)
        return
      }
      queryMagnet(dispatch, parentId, searchWords, pageNum, rows)
    })
  }

  const onSearchText = (e) => {
    setSearchWords(e.target.value)
  }
  const search = (e) => {
    queryMagnet(dispatch, parentId, searchWords, pageNum, rows)
  }

  return (
    <Container>
      <Grid container sx={{ display: 'flex' }} alignItems="center" justify="center">
        <Grid item>
          <TextField
            label="搜索关键字"
            variant="outlined"
            margin="normal"
            onChange={onSearchText}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }} />
        </Grid>
        <Grid item>
          <Button
            color="primary"
            onClick={search}
            variant="contained">
            搜索
          </Button>
        </Grid>
      </Grid>
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
                          <Grid item xs={4}>
                            <Typography whiteSpace={'pre'}>
                              {"分类: "}
                            </Typography>
                          </Grid>
                          <Grid item xs={6}>
                            <Tooltip title={item.introduce}>
                              <Link onClick={() => router.navigate2mgnetshares(navigate, item.id)} sx={{ cursor: 'pointer' }}>
                                {item.name}
                              </Link>
                            </Tooltip>
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
                          <Grid item xs={4}>
                            <Typography whiteSpace={'pre'}>
                              {"magnet uri: "}
                            </Typography>
                          </Grid>
                          <Grid item xs={6}>
                            <CopyToClipboard text={item.other}>
                              <Tooltip title={item.introduce}>
                                <Typography variant="button" component="div" noWrap>
                                  <Button onClick={() => setCopyDialogOpen(true)}>{item.other}</Button>
                                </Typography>
                              </Tooltip>
                            </CopyToClipboard>
                            <Dialog open={copyDialogOpen} onClose={() => setCopyDialogOpen(false)}>
                              <div style={{ padding: '16px' }}>
                                已复制到剪贴板
                              </div>
                            </Dialog>
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
    </Container >
  )
}