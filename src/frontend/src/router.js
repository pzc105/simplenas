import { React, useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { purple, green } from '@mui/material/colors';
import { CssBaseline, Container, Grid, Typography, Box, Button, Paper } from '@mui/material';

import SignIn from './user/signIn';
import SignUp from './user/signUp';
import Download from './download/download';
import GlobalSidebar from './globalSidebar.js'
import CategoryItemPage from './category/category.js'
import MagnetSharesPage from './magnetShares';
import PlyrWrap from './player/player.js';
import Player from './player/player2.js';
import CheckLoginHandler from './user/checklogin.js'
import UserInfoPage from './user/user.js'
import * as test from './test.js'
import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'

export default function App() {
  const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
  const defaultTheme = createTheme({
    palette: {
      mode: darkThemeMq && darkThemeMq.matches ? 'dark' : 'light',
      primary: {
        main: purple[500],
      },
      secondary: {
        main: green[500],
      },
    },
  });

  const [myTheme, setMyTheme] = useState(defaultTheme);

  const openGlobalChat = useSelector((state) => store.selectOpenGlobalChat(state))

  useEffect(() => {
    darkThemeMq.onchange = e => {
      if (e.matches) {
        setMyTheme(createTheme({
          palette: {
            mode: 'dark',
            primary: {
              main: purple[500],
            },
            secondary: {
              main: green[500],
            },
          },
        }))
      } else {
        setMyTheme(createTheme({
          palette: {
            mode: 'light',
            primary: {
              main: purple[500],
            },
            secondary: {
              main: green[500],
            },
          },
        }))
      }
    };
  })

  return (
    <ThemeProvider theme={myTheme}>
      <Router>
        <GlobalSidebar />
        <CheckLoginHandler />
        <Routes>
          <Route path="/" element={<Download />} />
          <Route path="/test" element={<test.Test />} />
          <Route path="/user" element={<UserInfoPage />} />
          <Route path="/signin" element={<SignIn />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/download" element={<Download />} />
          <Route path="/citem" element={<CategoryItemPage />} />
          <Route path="/video" element={<Player />} />
          <Route path="/mgnetshares" element={<MagnetSharesPage />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export function navigate2mgnetshares(navigate, itemId) {
  let path = "/mgnetshares"
  if (itemId != null) {
    path += "?itemid=" + itemId
  }
  navigate(path)
}
