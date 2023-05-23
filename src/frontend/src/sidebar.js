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



const HomeItems = ({ }) => {
  const navigate = useNavigate()
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const items = useSelector((state) => store.selectCategorySubItems(state, userInfo.homeDirectoryId))

  const onClick = (itemId) => {
    navigate("/citem/" + itemId)
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
        {utils.isLogined() && item.subComponent != null ? item.subComponent(item.subComponentParams) : null}
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

  const handleNavigation = (path) => {
    navigate(path)
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

  const menuItems = [
    {
      icon: <PersonIcon />,
      title: "登录",
      onClick: () => handleNavigation("/signin"),
    },
    {
      icon: <FileDownloadIcon />,
      title: "下载",
      onClick: () => handleNavigation("/download"),
    },
    {
      icon: <PersonalVideoIcon />,
      title: "Home",
      subComponent: HomeItems,
      subComponentParams: {},
      onClick: () => handleNavigation("/citem/" + userInfo.homeDirectoryId),
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
