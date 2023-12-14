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

import * as User from '../prpc/user_pb.js'
import * as Category from '../prpc/category_pb.js'
import userService from '../rpcClient.js'
import { serverAddress } from '../rpcClient.js'
import { CategoryItems } from './categoryitems.js'
import { queryItem, querySubItems, navigateToItem, navigateToVideo } from './utils.js'


export function CategoryCreatorPanel({ parentId, onRefresh, onClose }) {

  const [ItemName, setItemName] = useState('')
  const [Introduce, setIntroduce] = useState('')

  function NewCategoryItem(e) {
    e.stopPropagation()
    if (ItemName.length === 0) {
      return
    }
    var req = new User.NewCategoryItemReq()
    req.setName(ItemName)
    req.setTypeId(Category.CategoryItem.Type.DIRECTORY)
    req.setParentId(parentId)
    req.setIntroduce(Introduce)
    userService.newCategoryItem(req, {}, (err, res) => {
      if (err != null) {
        console.log(err)
        if (onClose) {
          onClose()
        }
        return
      }
      if (onRefresh) {
        onRefresh()
      }
      if (onClose) {
        onClose()
      }
    })
  }

  return (
    <Container sx={{ mt: "0.5em", mb: "0.5em" }}>
      <Grid container spacing={2} direction='row'>
        <Grid item xs={12}>
          <TextField
            variant="outlined"
            required
            fullWidth
            id="categoryName"
            label="名称"
            name="categoryName"
            autoComplete="categoryName"
            onChange={(e) => setItemName(e.target.value)} />
        </Grid>
        <Grid item xs={12}>
          <TextField
            multiline
            variant="outlined"
            fullWidth
            label="介绍"
            onChange={(e) => setIntroduce(e.target.value)} />
        </Grid>
        <Grid item xs={12}>
          <Button
            fullWidth
            variant="contained"
            color="primary"
            onClick={NewCategoryItem}
          >
            创建
          </Button>
        </Grid>
      </Grid>
    </Container>
  )
}