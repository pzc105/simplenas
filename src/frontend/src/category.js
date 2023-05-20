import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper,
  Typography, Tooltip, Card, CardContent, CardActions, CardMedia, InputAdornment
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'
import SideUtils from './sideUtils.js';

import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { serverAddress } from './rpcClient.js'
import * as utils from './utils.js'

const CategoryItems = ({ parentId }) => {
  const navigate = useNavigate()
  const items = useSelector((state) => store.selectCategorySubItems(state, parentId))
  const dispatch = useDispatch()

  const onClick = (item) => {
    if (item.typeId === Category.CategoryItem.Type.VIDEO) {
      navigate("/video/" + item.id)
      return
    }
    navigate("/citem/" + item.id)
  }

  const refreshSubtitle = (item) => {
    var req = new User.RefreshSubtitleReq()
    req.setItemId(item.id)
    userService.refreshSubtitle(req, {}, (err, res) => {
      if (err != null) {
        console.log(err)
        return
      }
    })
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


  return (
    <Paper style={{ width: "100em", maxHeight: '90vh', overflow: 'auto' }}>
      <Grid container spacing={2} sx={{ display: "flex" }}>
        <Grid item xs={12}>
          <Grid container spacing={2}>
            {items ?
              items.map((item) => (
                <Grid key={item.id} item xs={2}>
                  <Card onContextMenu={(e) => handleContextMenu(e, item.id)}
                    sx={{ width: "10em", marginBottom: "0", mt: "0" }}>
                    <CardMedia
                      component="img"
                      alt="Movie Poster"
                      width="100em"
                      image={serverAddress + "/poster/item/" + item.id}
                    />
                    <CardContent>
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
                    <CardActions>
                      {
                        item.typeId === Category.CategoryItem.Type.VIDEO ?
                          <div>
                            <Button
                              onClick={() => refreshSubtitle(item)}
                              size="small">
                              刷新字幕
                            </Button>
                            <Button>
                              举报
                            </Button>
                          </div>
                          : null
                      }
                    </CardActions>
                    <Menu
                      anchorReference="anchorPosition"
                      anchorPosition={anchorPosition}
                      open={open[item.id] ? open[item.id] : false}
                      onClose={() => handleClose(item.id)}
                    >
                      <MenuItem onClick={(e) => DelCategoryItem(item)}>删除</MenuItem>
                    </Menu>
                  </Card>
                </Grid>
              )) : null
            }
          </Grid>
        </Grid>
      </Grid>
    </Paper>
  )
}

const CategoryItemCreator = ({ parentId }) => {
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
  </Container>
  )
}

const CategoryContainer = styled('div')({
  display: 'flex', /* 将子元素布局为行 */
  height: '94vh', /* 页面铺满整个视窗 */
})

export default function CategoryItemPage() {
  const { itemId } = useParams()

  const dispatch = useDispatch()

  useEffect(() => {
    if (!utils.isNumber(itemId)) {
      return
    }
    var req = new User.QuerySubItemsReq()
    req.setParentId(itemId)
    userService.querySubItems(req, {}, (err, respone) => {
      if (err == null) {
        dispatch(store.categorySlice.actions.updateItem(respone.getParentItem().toObject()))
        respone.getItemsList().map((i) => {
          dispatch(store.categorySlice.actions.updateItem(i.toObject()))
          return null
        })
      } else {
        console.log(err)
      }
    })
  }, [itemId, dispatch])

  return (
    <CategoryContainer>
      <CssBaseline />
      <SideUtils name="管理" child={CategoryItemCreator(itemId)} />
      <CategoryItems parentId={itemId} />
    </CategoryContainer>
  );
}
