import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper, FormControlLabel, Switch, Dialog
} from '@mui/material';
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import { styled } from "@mui/material/styles";

import { useSelector, useDispatch } from 'react-redux';
import * as store from '../store.js'
import SideUtils from '../sideManager.js';
import { FloatingChat } from '../chat/chat.js';
import SubtitleUploader from '../uploadSubtitle.js';
import UnifiedPage from '../page.js'

import * as User from '../prpc/user_pb.js'
import * as Category from '../prpc/category_pb.js'
import userService from '../rpcClient.js'
import { serverAddress } from '../rpcClient.js'
import { CategoryItems } from './categoryitems.js'
import { CategoryCreatorPanel } from './creator.js';
import { queryItem, querySubItems, navigateToItem, navigateToVideo } from './utils.js'


const CategoryItemMgr = ({ parentId, onRefresh }) => {
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
      if (onRefresh) {
        onRefresh()
      }
    })
  }

  const [desc, setDesc] = useState(useSelector((state) => store.selectCategoryDesc(state)))

  return (
    <Container maxWidth="xs">
      <Container>
        <FormControlLabel
          control={
            <Switch
              checked={desc}
              onClick={
                (e) => {
                  let v = !desc
                  setDesc(v)
                  dispatch(store.categorySlice.actions.setDesc(v))
                }
              }
              color="primary"
              inputProps={{ 'aria-label': 'controlled' }}
            />
          }
          label={'逆序'}
        />
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
  const showGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))
  const thisItem = useSelector((state) => store.selectCategoryItem(state, itemId))
  const categoryDesc = useSelector((state) => store.selectCategoryDesc(state))

  const pageRows = 20
  const [totalRows, setTotalRows] = useState(0)
  const pageNum = useRef(0)
  const [pageNumState, setPageNumState] = useState(0)
  const refresh = () => {
    queryItem(itemId, shareid, dispatch, (item) => {
      setTotalRows(item.subItemIdsList.length)
    })
    querySubItems({
      itemId, shareid, dispatch, pageNum: pageNum.current, pageRows: pageRows, desc: categoryDesc, callback: (items) => {
        dispatch(store.categorySlice.actions.updateDisplayItems(items))
      }
    })
  }

  const closeChatPanel = () => {
    dispatch(store.userSlice.actions.setShowChatPanel(false))
  }
  const closeGlobalChat = () => {
    dispatch(store.userSlice.actions.setOpenGlobalChat(false))
  }

  useEffect(() => {
    pageNum.current = 0
    setPageNumState(0)
  }, [itemId])

  useEffect(() => {
    if (thisItem && thisItem.typeId === Category.CategoryItem.Type.VIDEO) {
      navigateToVideo(navigate, { replace: true }, thisItem.id, shareid)
    }
  }, [thisItem])

  useEffect(() => {
    if (showGlobalChat) {
      closeChatPanel()
    }
  }, [showGlobalChat])

  const dispatch = useDispatch()

  useEffect(() => {
    refresh()
  }, [itemId, dispatch, navigate, shareid, categoryDesc])

  const [MenuPos, setMenuPos] = useState({ left: 0, top: 0 });
  const [MenuOpen, setMenuOpen] = useState(false);
  const handleContextMenu = (event) => {
    event.preventDefault();
    setMenuPos({ left: event.clientX, top: event.clientY });
    setMenuOpen(true);
  };
  const handleMenuClose = () => {
    setMenuOpen(false);
  };

  const [CreatorOpen, setCreatorOpen] = useState(false)

  return (
    <CategoryContainer>
      <CssBaseline />
      <SideUtils
        name="管理"
        child={CategoryItemMgr({ parentId: itemId, onRefresh: refresh })}
      />
      <Container style={{ flex: 1, display: 'flex', flexDirection: 'column' }}>
        <Paper style={{ height: "85vh" }} onContextMenu={(e) => handleContextMenu(e)}>
          <CategoryItems shareid={shareid} onRefresh={refresh} />
          <Dialog open={CreatorOpen} onClose={() => setCreatorOpen(false)}>
            <CategoryCreatorPanel parentId={itemId} onRefresh={refresh} onClose={() => setCreatorOpen(false)} />
          </Dialog>
          <Menu
            anchorReference="anchorPosition"
            anchorPosition={MenuPos}
            open={MenuOpen}
            onClose={handleMenuClose} >
            <MenuItem onClick={() => { setCreatorOpen(true); handleMenuClose(); }}>创建分类</MenuItem>
          </Menu>
        </Paper>
        <UnifiedPage
          PageTotalCount={Math.ceil(totalRows / pageRows)}
          PageNum={parseInt(pageNumState + 1)}
          onPage={(n) => { pageNum.current = n - 1; setPageNumState(pageNum.current); refresh() }} />
      </Container>
      {shownChatPanel && !showGlobalChat ? <FloatingChat itemId={itemId} onClose={closeChatPanel} /> : null}
      {showGlobalChat ? <FloatingChat itemId={1} onClose={closeGlobalChat} /> : null}
    </CategoryContainer>
  );
}

