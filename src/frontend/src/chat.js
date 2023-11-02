import React, { useState, useEffect, useRef } from 'react';
import ReactDOM from 'react-dom';
import {
  Container, CssBaseline, Grid, List, Paper, Button, Box, TextField, ListItem,
  Typography, Popper, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle
} from '@mui/material';
import { useSelector, useDispatch } from 'react-redux';
import Draggable from 'react-draggable';
import CloseIcon from '@mui/icons-material/Close';

import * as User from './prpc/user_pb.js'
import * as store from './store.js'
import userService from './rpcClient.js'
import './chat.css'

const ChatPanel = ({ itemId }) => {
  const maxMaxCount = 1000
  const [messages, setMessages] = useState([]);
  const [inputValue, setInputValue] = useState('');

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const handleSendMessage = () => {
    if (inputValue !== '') {
      const chatMsg = new User.ChatMessage()
      chatMsg.setMsg(inputValue)
      const req = new User.SendMsg2ChatRoomReq()
      req.setItemId(itemId)
      req.setChatMsg(chatMsg)
      userService.sendMsg2ChatRoom(req, {}, (err, res) => {
      })
      setInputValue('');
    }
  };

  const msgsRef = useRef([])

  useEffect(() => {
    const req = new User.JoinChatRoomReq()
    req.setItemId(itemId)
    var stream = userService.joinChatRoom(req)
    msgsRef.current = []
    setMessages([])
    stream.on('data', function (res) {
      const chatMsgs = res.getChatMsgsList()
      msgsRef.current.push(...chatMsgs)
      if (msgsRef.current.length > maxMaxCount) {
        msgsRef.current = msgsRef.current.slice(-maxMaxCount)
      }
      setMessages([...msgsRef.current])
    })
    stream.on('status', function (status) {
    });
    stream.on('end', function (end) {
      stream.cancel()
    });

    return () => {
      stream.cancel()
    };
  }, [itemId])

  const chatAreaRef = useRef(null);
  useEffect(() => {
    chatAreaRef.current.scrollTop = chatAreaRef.current.scrollHeight;
  }, [messages]);

  return (
    <Container>
      <Paper style={{ maxHeight: '50vh', overflow: 'auto' }} ref={chatAreaRef}>
        <List>
          {messages.map((message, i) => (
            <ListItem key={i}>
              <Box mb={1}>
                <Typography variant="subtitle2" color="textSecondary">
                  {message.getUserName()} - {new Date(message.getSentTime()).toLocaleString()}
                </Typography>
                <Typography
                  variant="body2"
                  style={{ wordWrap: 'break-word', wordBreak: 'break-all' }}
                >
                  {message.getMsg()}
                </Typography>
              </Box>
            </ListItem>
          ))}
        </List>
      </Paper>
      <div>
        <TextField
          label="输入"
          value={inputValue}
          onChange={handleInputChange}
        />
        <Button variant="contained" color="primary" onClick={handleSendMessage}>
          发送
        </Button>
      </div>
    </Container>
  )
}

export const FloatingChat = ({ itemId, onClose, defaultPosition }) => {
  console.log(defaultPosition)
  const handleClose = () => {
    if (onClose) {
      onClose()
    }
  }

  useEffect(() => {
    document.body.style.overflow = 'hidden';
    return () => {
      document.body.style.overflow = 'auto';
    }
  }, [])

  return (
    ReactDOM.createPortal(
      <Draggable handle='.draggableWindow' positionOffset={{ x: '-50%', y: '-50%' }}>
        <div className='myElement'>
          <Paper >
            <Grid container className='draggableWindow'>
              <Grid item xs={6}>
                <Typography>
                  聊天室
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Box sx={{ display: 'flex', justifyContent: 'flex-end', pr: 1 }}>
                  <Button size="small" color="secondary" onClick={handleClose}>
                    <CloseIcon />
                  </Button>
                </Box>
              </Grid>
            </Grid>
            <ChatPanel itemId={itemId} />
          </Paper>
        </div>
      </Draggable>,
      document.getElementById('portal-root') // 这是一个在public/index.html中定义的元素
    )
  )
}