import React, { useEffect, useRef, useState } from 'react';
import { Container, Grid, CssBaseline, List, ListItem, Button, Typography, Tooltip, Switch, FormControlLabel, Paper } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';

import Plyr from 'plyr_p';
import 'plyr_p/dist/plyr.css';
import Hls from 'hls.js'

import { queryItem, querySubItems, navigateToItem } from '../category/utils.js'
import { serverAddress } from '../rpcClient.js'
import * as store from '../store.js'
import * as utils from '../utils.js';
import { FloatingChat } from '../chat/chat.js';


export default function PlayList({ videoItemList, shareid }) {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [autoContinuedPlay, setAutoContinuedPlay] = useState(useSelector((state) => store.selectAutoPlayVideo(state)));

  return (
    <Container>
      <Grid container sx={{ display: 'flex' }} alignItems="center" justify="center">
        <Grid item xs={5}>
          <Typography variant="button" component="div" noWrap>
            播放列表
          </Typography>
        </Grid>
        <Grid item xs={7} sx={{ display: 'flex', justifyContent: 'flex-end', pr: 1 }}>
          <FormControlLabel
            control={
              <Switch
                checked={autoContinuedPlay}
                onClick={
                  (e) => {
                    let v = !autoContinuedPlay
                    setAutoContinuedPlay(v)
                    dispatch(store.playerSlice.actions.setAutoContinuedPlayVideo(v))
                  }
                }
                color="primary"
                inputProps={{ 'aria-label': 'controlled' }}
              />
            }
            label={'自动连播'}
          />
        </Grid>
      </Grid>
      <Paper style={{ maxHeight: '50vh', overflow: 'auto' }}>
        <List>
          {
            videoItemList.map((item) => {
              if (!utils.isVideoItem(item)) {
                return null
              }
              return (
                <ListItem
                  key={item.id} >
                  <Tooltip title={item.name}>
                    <Typography variant="button" component="div" noWrap>
                      <Button onClick={() => navigateToItem(navigate, {}, item.id, shareid)}>
                        {item.name}
                      </Button>
                    </Typography>
                  </Tooltip>
                </ListItem>
              )
            })
          }
        </List>
      </Paper>
    </Container>
  )
}