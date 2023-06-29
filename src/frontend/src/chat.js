import React, { useState, useEffect, useRef } from 'react';
import { Menu, MenuItem, List, Paper, Button, Box, TextField, ListItem, ListItemText, Typography } from '@mui/material';
import { styled } from "@mui/material/styles";
import LinearProgress from '@mui/material/LinearProgress';
import { useSelector, useDispatch } from 'react-redux';

import * as store from './store.js'
import * as User from './prpc/user_pb.js'
import userService from './rpcClient.js'
import { Container } from '@mui/system';

export default function ChatPanel({ itemId }) {

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
      <Paper style={{ maxHeight: '50vh', width: '30vw', overflow: 'auto' }} ref={chatAreaRef}>
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