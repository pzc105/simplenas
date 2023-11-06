import React, { useState, useEffect, useRef } from 'react';
import ReactDOM from 'react-dom';
import {
  Container, CssBaseline, Grid, List, Paper, Button, Box, TextField, ListItem,
  Typography, Popper, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, Popover
} from '@mui/material';
import { useSelector, useDispatch } from 'react-redux';
import Draggable from 'react-draggable';
import CloseIcon from '@mui/icons-material/Close';
import MoodIcon from '@mui/icons-material/Mood';

import * as User from '../prpc/user_pb.js'
import * as store from '../store.js'
import { emojiList } from './emojiList.js';
import userService from '../rpcClient.js'
import './chat.css'
import { navigateToItem } from '../category.js';


export const FloatingChat = ({ name, itemId, onClose }) => {
  const storedPosition = JSON.parse(localStorage.getItem('dragPosition'))
  const positionOffset = storedPosition ? undefined : { x: "-50%", y: "-50%" }
  const [position, setPosition] = useState(storedPosition ? storedPosition : { x: 0, y: 0 });

  useEffect(() => {
    document.body.style.overflow = 'hidden';

    return () => {
      document.body.style.overflow = 'auto';
    }
  }, [])

  const handleDrag = (e, data) => {
    const newPosition = { x: data.x, y: data.y };
    setPosition(newPosition);
  };

  const handleStop = (e, data) => {
    localStorage.setItem('dragPosition', JSON.stringify(position));
  }

  const handleClose = () => {
    if (onClose) {
      onClose()
    }
  }

  return (
    ReactDOM.createPortal(
      <Draggable
        onDrag={handleDrag}
        onStop={handleStop}
        handle='.draggableWindow'
        position={position}
        positionOffset={positionOffset}
      >
        <div className='floatingchat' >
          <Paper style={{ width: "20em", borderRadius: "20px", border: '2px solid #e178ce' }}>
            <Grid container>
              <Grid item xs={10} className='draggableWindow'>
                <Typography sx={{ userSelect: 'none', ml: "1em" }}>
                  {name ? name : "聊天室-" + itemId}
                </Typography>
              </Grid>
              <Grid item xs={2} sx={{ display: 'flex', justifyContent: 'flex-end', pr: 1 }}>
                <Box >
                  <Button size="small" color="secondary" onClick={handleClose}>
                    <CloseIcon />
                  </Button>
                </Box>
              </Grid>
            </Grid>
            <ChatPanel itemId={itemId} />
            <div className='draggableWindow' style={{ height: "1em" }}></div>
          </Paper>
        </div>
      </Draggable>,
      document.getElementById('portal-root') // 这是一个在public/index.html中定义的元素
    )
  )
}

const ChatPanel = ({ itemId }) => {
  const maxMaxCount = 1000
  const [messages, setMessages] = useState([]);
  const [inputValue, setInputValue] = useState('');

  var EmojiConvertor = require('emoji-js');
  var emoji = new EmojiConvertor();

  const handleInputChange = (e) => {
    setInputValue(e.target.value);
  };

  const onEmoji = (em) => {
    setInputValue(inputValue + em)
  }
  const translate2Emoji = (msg) => {
    emoji.replace_mode = 'unified';
    emoji.allow_native = true;
    return emoji.replace_colons(msg);
  }

  const handleSendMessage = () => {
    if (inputValue !== '') {
      const chatMsg = new User.ChatMessage()
      const room = new User.Room()
      room.setType(User.Room.Type.CATEGORY)
      room.setId(itemId)
      emoji.colons_mode = true
      let msg = emoji.replace_unified(inputValue)
      chatMsg.setMsg(msg)
      const req = new User.SendMsg2ChatRoomReq()
      req.setRoom(room)
      req.setChatMsg(chatMsg)
      userService.sendMsg2ChatRoom(req, {}, (err, res) => {
      })
      setInputValue('');
    }
  };

  const msgsRef = useRef([])

  useEffect(() => {
    const req = new User.JoinChatRoomReq()
    const room = new User.Room()
    room.setType(User.Room.Type.CATEGORY)
    room.setId(itemId)
    req.setRoom(room)
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
    <Container className='chatContainer'>
      <Paper style={{ maxHeight: '50vh', overflow: 'auto' }} ref={chatAreaRef}>
        <List >
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
                  {translate2Emoji(message.getMsg())}
                </Typography>
              </Box>
            </ListItem>
          ))}
        </List>
      </Paper>
      <div style={{ width: "100%" }}>
        <EmojiPicker emojiList={emojiList} onEmoji={onEmoji} />
        <Grid container alignItems="center" justify="center">
          <Grid item xs={10}>
            <TextField
              label="输入"
              value={inputValue}
              onChange={handleInputChange}
              autoFocus
            />
          </Grid>
          <Grid item xs={2}>
            <Button variant="contained" color="primary" onClick={handleSendMessage}>
              发送
            </Button>
          </Grid>
        </Grid>
      </div>
    </Container >
  )
}

const EmojiPicker = ({ emojiList, onEmoji }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open.current ? 'emoji-popover' : undefined;

  const onEmojiClick = (ce) => {
    if (onEmoji) {
      onEmoji(ce)
    }
  }

  return (
    <div>
      <IconButton
        color="primary"
        aria-describedby={id}
        onClick={handleClick}
      >
        <MoodIcon />
      </IconButton>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <div style={{ maxWidth: "6em" }}>
          {emojiList.map((e, i) => {
            return (
              <span key={i} className='pointer' onClick={() => { setAnchorEl(null); onEmojiClick(e); }}>
                {e}
              </span>
            )
          })}
        </div>
      </Popover>
    </div>
  );
};