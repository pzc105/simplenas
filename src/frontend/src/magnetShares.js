import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  CssBaseline, Button, TextField, Menu, MenuItem, Container, Grid, Paper, Box,
  Typography, Tooltip, Card, CardContent, CardActions, InputAdornment, Popover, Popper, List, ListItem
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
import * as category from './category.js'

import * as User from './prpc/user_pb.js'
import * as Category from './prpc/category_pb.js'
import userService from './rpcClient.js'
import { serverAddress } from './rpcClient.js'

const MagnetItems = ({ parentId }) => {
  const navigate = useNavigate()
  const items = useSelector((state) => store.selectCategorySubItems(state, parentId))
  console.log(parentId, items)

  return (
    <Container>
      <List>
        {items ?
          items.map((item) => (
            <ListItem key={item.id}>
              {item.name}
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
  const thisItem = useSelector((state) => store.selectCategoryItem(state, itemId))

  useEffect(() => {
    var req = new User.QueryMagnetReq()
    req.setParentId(itemId)
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
  }, [itemId])

  return (
    <Container>
      <CssBaseline />
      <MagnetItems parentId={itemId} />
    </Container>
  )
}
