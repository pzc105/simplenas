import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";
import Draggable from 'react-draggable';
import CloseIcon from '@mui/icons-material/Close';

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'
import SideUtils from './sideUtils.js';
import ChatPanel from './chat.js';
import SubtitleUploader from './uploadSubtitle.js';

import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { serverAddress } from './rpcClient.js'
import * as utils from './utils.js'

const CategoryItems = ({ parentId, shareid }) => {
  const navigate = useNavigate()
  const items = useSelector((state) => store.selectCategorySubItems(state, parentId))
  const dispatch = useDispatch()

  const onClick = (item) => {
    if (item.typeId === Category.CategoryItem.Type.VIDEO) {
      navigateToVideo(navigate, {}, item.id, shareid)
    } else {
      navigateToItem(navigate, {}, item.id, shareid)
    }
  }

  const DelCategoryItem = (item) => {
    var req = new User.DelCategoryItemReq()
    req.setItemId(item.id)
    userService.delCategoryItem(req, {}, (err, res) => {
      if (err != null) {
        console.log(err)
        return
      }
      dispatch(store.categorySlice.actions.deleteItem(item.id))
    })
    handleClose(item.id)
  }

  const ShareCategoryItem = (item) => {
    let req = new User.ShareItemReq()
    req.setItemId(item.id)
    userService.shareItem(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      const shareid = res.getShareId()
      alert("复制此共享URL: https://" + window.location.hostname + ":" + window.location.port + "/citem?itemid=" + item.id + "&shareid=" + shareid)
    })
    handleClose(item.id)
  }

  const [anchorPosition, setAnchorPosition] = useState({ left: 0, top: 0 });
  const [open, setOpen] = useState({});

  const handleContextMenu = (event, itemId) => {
    event.preventDefault();
    setAnchorPosition({ left: event.clientX, top: event.clientY });
    setOpen({ ...open, [itemId]: true });
  };

  const handleClose = (itemId) => {
    setOpen({ ...open, [itemId]: false });
  };

  const uploadSubtitleAnchorElRef = useRef(null)
  const [popoverOpen, setPopoverOpen] = useState(false)
  const handlePopoverClose = () => {
    setPopoverOpen(false)
  }
  const [subtitleUploadItemId, setSubtitleUploadItemId] = useState(-1)
  const UploadSubtitle = (item) => {
    setSubtitleUploadItemId(item.id)
    setPopoverOpen(true)
    handleClose(item.id)
  }

  return (
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }} ref={uploadSubtitleAnchorElRef}>
      <Grid container spacing={2} sx={{ display: "flex" }}>
        <Grid item xs={12}>
          <Grid container spacing={2}>
            {items ?
              items.map((item) => (
                <Grid key={item.id} item xs={10} sm={5} lg={2} sx={{ ml: "0.5em", mt: "0.5em" }}>
                  <Card onContextMenu={(e) => handleContextMenu(e, item.id)}>
                    <Box sx={{ display: "flex", justifyContent: "center", height: "4.3em" }}>
                      <img style={{ maxHeight: "5em" }} alt="Movie Poster"
                        src={serverAddress + "/poster/item?itemid=" + item.id + (shareid ? "&shareid=" + shareid : "")} />
                    </Box>
                    <CardContent sx={{ display: "flex", justifyContent: "center" }}>
                      <Tooltip title={item.name}>
                        <Typography variant="button" component="div" noWrap>
                          <Button onClick={() => onClick(item)}>
                            {item.name}
                          </Button>
                        </Typography>
                      </Tooltip>
                      <Typography variant="body2" color="text.secondary">
                        {item.introduce}
                      </Typography>
                    </CardContent>
                    <CardActions sx={{ display: "flex", justifyContent: "center" }}>
                    </CardActions>
                    <Menu
                      anchorReference="anchorPosition"
                      anchorPosition={anchorPosition}
                      open={open[item.id] ? open[item.id] : false}
                      onClose={() => handleClose(item.id)}
                    >
                      <MenuItem onClick={(e) => DelCategoryItem(item)}>删除</MenuItem>
                      <MenuItem onClick={(e) => ShareCategoryItem(item)}>分享</MenuItem>
                      <MenuItem onClick={(e) => UploadSubtitle(item)}>上传字幕</MenuItem>
                    </Menu>
                  </Card>
                </Grid>
              )) : null
            }
          </Grid>
        </Grid>
      </Grid>
      <Popover
        id={"id"}
        open={popoverOpen}
        anchorEl={uploadSubtitleAnchorElRef.current}
        onClose={handlePopoverClose}
        anchorOrigin={{
          vertical: 'center',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'center',
          horizontal: 'center',
        }}
      >
        <SubtitleUploader itemId={subtitleUploadItemId} onClose={() => setPopoverOpen(false)} />
      </Popover>
    </Paper>
  )
}

const CategoryItemCreator = ({ parentId, refreshParent }) => {
  const dispatch = useDispatch()
  const [itemName, setItemName] = useState('')
  function handleChange(e) {
    setItemName(e.target.value)
  }
  function NewCategoryItem(e) {
    e.stopPropagation()
    if (itemName.length === 0) {
      return
    }
    var req = new User.NewCategoryItemReq()
    req.setName(itemName)
    req.setTypeId(Category.CategoryItem.Type.DIRECTORY)
    req.setParentId(parentId)
    userService.newCategoryItem(req, {}, (err, res) => {
      if (err != null) {
        console.log(err)
        return
      }
      refreshParent()
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
              onClick={NewCategoryItem}
              variant="contained">
              新建分类
            </Button>
          </Grid>
        </Grid>
      </Grid>
    </Container>
    <Container sx={{ mt: '1em' }} >
      <Button
        fullWidth
        color="primary"
        variant="contained"
        onClick={() => { dispatch(store.userSlice.actions.setShowChatPanel(true)) }}>
        聊天室
      </Button>
    </Container>
  </Container>
  )
}

const CategoryContainer = styled('div')({
  display: 'flex', /* 将子元素布局为行 */
  height: '94vh', /* 页面铺满整个视窗 */
})

export default function CategoryItemPage() {
  const location = useLocation()
  const navigate = useNavigate()
  const searchParams = new URLSearchParams(location.search)
  const shareid = searchParams.get('shareid')
  const itemId = searchParams.get('itemid') ? Number(searchParams.get('itemid')) : -1
  const shownChatPanel = useSelector((state) => store.selectShownChatPanel(state))
  const thisItem = useSelector((state) => store.selectCategoryItem(state, itemId))

  useEffect(() => {
    if (thisItem && thisItem.typeId === Category.CategoryItem.Type.VIDEO) {
      navigateToVideo(navigate, { replace: true }, thisItem.id, shareid)
    }
  }, [thisItem])

  const dispatch = useDispatch()
  const refreshSubItems = () => {
    querySubItems()
  }

  useEffect(() => {
    queryItem(itemId, shareid, dispatch)
    querySubItems(itemId, shareid, dispatch)
  }, [itemId, dispatch, navigate, shareid])

  const anchorElRef = useRef(null)

  const handleClose = () => {
    dispatch(store.userSlice.actions.setShowChatPanel(false))
  };

  const [openPopper, setOpenPopper] = useState(false)
  useEffect(() => {
    setOpenPopper(Boolean(anchorElRef.current !== null && shownChatPanel))
  }, [anchorElRef, shownChatPanel])

  return (
    <CategoryContainer>
      <CssBaseline />
      {
        <SideUtils
          name="管理"
          child={CategoryItemCreator({ parentId: itemId, refreshParent: () => { refreshSubItems() } })}
        />
      }
      <CategoryItems parentId={itemId} shareid={shareid} />
      <Button
        ref={anchorElRef}
        style={{
          position: 'fixed',
          right: '20px',
          bottom: '20px',
          zIndex: 9999,
        }}>
      </Button>
      <Draggable >
        <Popper
          id={shownChatPanel ? 'floating-window' : undefined}
          open={openPopper}
          anchorEl={anchorElRef.current}
          placement='bottom'
          sx={{ backgroundColor: 'background.default' }}
        >
          <Box sx={{ display: 'flex', justifyContent: 'flex-end', pr: 1, backgroundColor: 'background.default' }}>
            <Button size="small" color="secondary" onClick={handleClose}>
              <CloseIcon />
            </Button>
          </Box>
          <ChatPanel itemId={itemId} />
        </Popper>
      </Draggable>
    </CategoryContainer>
  );
}

export function navigateToItem(navigate, navigateParams, itemId, shareid) {
  let path = "/citem?itemid=" + itemId
  if (shareid) {
    path += "&shareid=" + shareid
  }
  navigate(path, navigateParams)
}

export function navigateToVideo(navigate, navigateParams, itemId, shareid) {
  let path = "/video?itemid=" + itemId
  if (shareid) {
    path += "&shareid=" + shareid
  }
  navigate(path, navigateParams)
}

export const querySubItems = (itemId, shareid, dispatch, callback) => {
  var req = new User.QuerySubItemsReq()
  req.setParentId(itemId)
  if (shareid) {
    req.setShareId(shareid)
  }
  userService.querySubItems(req, {}, (err, respone) => {
    if (err == null) {
      let subItems = []
      respone.getItemsList().map((i) => {
        let obj = i.toObject()
        subItems.push(obj)
        dispatch(store.categorySlice.actions.updateItem(obj))
        return null
      })
      if (callback) {
        callback(subItems)
      }
    } else {
      console.log(err)
    }
  })
}

export const queryItem = (itemId, shareId, dispatch) => {
  var req = new User.QueryItemInfoReq()
  req.setItemId(itemId)
  if (shareId) {
    req.setShareId(shareId)
  }
  userService.queryItemInfo(req, {}, (err, res) => {
    if (err != null || !res) {
      return
    }
    const itemInfo = res.getItemInfo()
    if (itemInfo.getTypeId() === Category.CategoryItem.Type.VIDEO) {
      dispatch(store.categorySlice.actions.updateVideoInfo({ itemId: itemInfo.getId(), videoInfo: res.getVideoInfo().toObject() }))
    }
    dispatch(store.categorySlice.actions.updateItem(itemInfo.toObject()))
  })
}