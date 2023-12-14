import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper, FormControlLabel, Switch
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";

import { useSelector, useDispatch } from 'react-redux';
import * as store from '../store.js'
import SideUtils from '../sideManager.js';
import { FloatingChat } from '../chat/chat.js';
import SubtitleUploader from '../uploadSubtitle.js';
import UnifiedPage from '../page.js'
import { queryItem, querySubItems, navigateToItem, navigateToVideo } from './utils.js'

import * as User from '../prpc/user_pb.js'
import * as Category from '../prpc/category_pb.js'
import userService from '../rpcClient.js'
import { serverAddress } from '../rpcClient.js'

export function CategoryItems({ shareid, onRefresh }) {
  const navigate = useNavigate()
  const items = useSelector((state) => store.selectDisplayItems(state))
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
      if (onRefresh) {
        onRefresh()
      }
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
    event.stopPropagation();
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

  const RenameBtVideoName = (item) => {
    let req = new User.RenameBtVideoNameReq()
    req.setItemId(item.id)
    userService.renameBtVideoName(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      if (onRefresh) {
        onRefresh()
      }
    })
    handleClose(item.id)
  }

  const [shownRenameInput, setShownRenameInput] = useState({})
  const [renameInputFocus, setRenameInputFocus] = useState({})
  const showRenameImput = (item) => {
    setShownRenameInput({ ...shownRenameInput, [item.id]: true })
    handleClose(item.id)
  }
  const onFocusRenameInput = (e, item) => {
    e.target.select()
    setRenameInputFocus({ ...renameInputFocus, [item.id]: true })
  }
  const onBlurRenameInput = (e, item) => {
    if (!renameInputFocus[item.id]) {
      return
    }
    renameItem(item, e.target.value);
    setShownRenameInput({ ...shownRenameInput, [item.id]: false })
    setRenameInputFocus({ ...renameInputFocus, [item.id]: false })
  }
  const renameItem = (item, newName) => {
    if (item.name === newName) {
      return
    }
    let req = new User.RenameItemReq()
    req.setItemId(item.id)
    req.setNewName(newName)
    userService.renameItem(req, {}, (err, res) => {
      if (err != null) {
        return
      }
      if (onRefresh) {
        onRefresh()
      }
    })
  }

  return (
    <Paper style={{ width: "100%", maxHeight: '90vh', overflow: 'auto' }} ref={uploadSubtitleAnchorElRef}>
      <Grid container spacing={2} sx={{ display: "flex" }}>
        <Grid item xs={12}>
          <Grid container spacing={2}>
            {
              items.map((item) => {
                return (
                  <Grid key={item.id} item xs={10} sm={5} lg={2} sx={{ ml: "0.5em", mt: "0.5em" }}>
                    <Tooltip title={<div>{"Name:" + item.name}<br />{"介绍:" + item.introduce}</div>} >
                      <Card onContextMenu={(e) => handleContextMenu(e, item.id)}>
                        <Box sx={{ display: "flex", justifyContent: "center", height: "4.3em" }}>
                          <img style={{ maxHeight: "5em" }} alt="Movie Poster"
                            src={serverAddress + "/poster/item?itemid=" + item.id + (shareid ? "&shareid=" + shareid : "")} />
                        </Box>
                        <CardContent sx={{ display: "flex", justifyContent: "center" }}>
                          {
                            !shownRenameInput[item.id] ?
                              <Typography
                                onClick={() => onClick(item)}
                                noWrap
                                style={{ cursor: "pointer" }}>
                                {item.name}
                              </Typography> :
                              <TextField
                                autoFocus
                                defaultValue={item.name}
                                onFocus={(e) => onFocusRenameInput(e, item)}
                                onBlur={(e) => onBlurRenameInput(e, item)}
                              />
                          }
                          <Typography variant="body2" color="text.secondary">
                          </Typography>
                        </CardContent>
                        <CardActions sx={{ display: "flex", justifyContent: "center" }}>
                        </CardActions>
                        <Menu
                          anchorReference="anchorPosition"
                          anchorPosition={anchorPosition}
                          open={open[item.id] ? open[item.id] : false}
                          onClose={() => handleClose(item.id)} >
                          <MenuItem onClick={(e) => showRenameImput(item)}>重命名</MenuItem>
                          <MenuItem onClick={(e) => DelCategoryItem(item)}>删除</MenuItem>
                          <MenuItem onClick={(e) => ShareCategoryItem(item)}>分享</MenuItem>
                          <MenuItem onClick={(e) => UploadSubtitle(item)}>上传字幕</MenuItem>
                          <MenuItem onClick={(e) => RenameBtVideoName(item)}>智能整理BT视频名字</MenuItem>
                        </Menu>
                      </Card>
                    </Tooltip>
                  </Grid>
                )
              })
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
    </Paper >
  )
}