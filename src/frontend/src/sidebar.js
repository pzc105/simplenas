import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import MenuIcon from '@mui/icons-material/Menu';
import FileDownloadIcon from '@mui/icons-material/FileDownload';
import PersonalVideoIcon from '@mui/icons-material/PersonalVideo';
import PersonIcon from "@mui/icons-material/Person"
import { Container } from '@mui/system';
import { Drawer, Accordion, AccordionSummary, AccordionDetails, List, ListItemButton, Typography, IconButton } from "@mui/material";

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'

import * as utils from './utils.js'
import * as User from './prpc/user_pb.js'
import userService from './rpcClient.js'
import { navigateToItem } from './category.js'


const ListItemWithChildren = ({ item }) => {
  const [expanded, setExpanded] = useState(false);
  const handleExpand = () => {
    setExpanded(!expanded);
  };
  return (
    <Accordion
      key={item.title}
      expanded={expanded}
      onMouseEnter={item.subComponent != null ? handleExpand : null}
      onMouseLeave={item.subComponent != null ? handleExpand : null}>
      <AccordionSummary expandIcon={item.icon}>
        <Typography
          variant="subtitle1"
          onClick={item.onClick}>
          {item.title}
        </Typography>
      </AccordionSummary>
      <AccordionDetails>
        {item.subComponent != null ? item.subComponent(item.subComponentParams) : null}
      </AccordionDetails>
    </Accordion>
  );
};

export default function Sidebar() {
  const navigate = useNavigate()
  const dispatch = useDispatch()
  const userInfo = useSelector((state) => store.selectUserInfo(state))

  const [isOpen, setIsOpen] = useState(false)

  const handleToggleSidebar = () => {
    setIsOpen(!isOpen);
  }

  useEffect(() => {
    if (!utils.isLogined() || userInfo == null) {
      return
    }
    var req = new User.QuerySubItemsReq()
    req.setParentId(userInfo.homeDirectoryId)
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
  }, [userInfo, dispatch])

  const firstItem = userInfo ? {
    icon: <PersonIcon />,
    title: "个人信息",
    onClick: () => navigate("/user"),
    subComponent: Relogin,
    subComponentParams: {},
  } : {
    icon: <PersonIcon />,
    title: "登录",
    onClick: () => navigate("/signin"),
  }
  const menuItems = [
    firstItem,
    {
      icon: <FileDownloadIcon />,
      title: "下载",
      onClick: () => navigate("/download"),
    },
    {
      icon: <PersonalVideoIcon />,
      title: "Home",
      subComponent: HomeItems,
      subComponentParams: {},
      onClick: userInfo ? () => navigateToItem(navigate, {}, userInfo.homeDirectoryId, null) : null,
    },

  ];

  return (
    <Container sx={{ backgroundColor: 'background.default', ml: "0", width: "100%" }}>
      <IconButton onClick={handleToggleSidebar} edge="start" color="inherit" aria-label="menu">
        <MenuIcon />
      </IconButton>
      <Drawer anchor="left" open={isOpen} onClose={handleToggleSidebar}>
        <List>
          {menuItems.map((item, index) => {
            return (
              <Container key={index}>
                <ListItemWithChildren
                  item={item} />
              </Container>
            )
          })}
        </List>
      </Drawer>
    </Container>
  );
}

const HomeItems = ({ }) => {
  const navigate = useNavigate()
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const items = useSelector((state) => {
    if (userInfo == null) {
      return []
    }
    return store.selectCategorySubItems(state, userInfo.homeDirectoryId)
  })

  const onClick = (itemId) => {
    navigateToItem(navigate, {}, itemId, null)
  }

  return (
    <List>
      {
        items ?
          items.map((item) => {
            return (
              <ListItemButton
                key={item.id}
                onClick={() => onClick(item.id)}
              >
                {item.name}
              </ListItemButton>
            )
          }) : null
      }
    </List>
  );
};


const Relogin = ({ }) => {
  const navigate = useNavigate()
  const onClick = () => {
    navigate("/signin")
  }
  return (
    <List>
      <ListItemButton
        key={1}
        onClick={() => onClick()} >
        重新登录
      </ListItemButton>
    </List>
  );
}