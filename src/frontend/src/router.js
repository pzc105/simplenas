import { React } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import { createTheme, ThemeProvider } from '@mui/material/styles';
import { purple, green } from '@mui/material/colors';
import SignIn from './signIn';
import SignUp from './signUp';
import Download from './download';
import Sidebar from './sidebar.js'
import CategoryItemPage from './category.js'
import PlyrWrap from './plyrwrap.js';
import CheckLoginHandler from './checklogin.js'
import * as test from './test.js'

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: purple[500],
    },
    secondary: {
      main: green[500],
    },
  },
});

export default function App() {
  return (
    <ThemeProvider theme={darkTheme}>
      <Router>
        <Sidebar />
        <CheckLoginHandler />
        <Routes>
          <Route path="/" element={<Download />} />
          <Route path="/test" element={<test.MultiSelectList />} />
          <Route path="/signin" element={<SignIn />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/download" element={<Download />} />
          <Route path="/citem/:itemId" element={<CategoryItemPage />} />
          <Route path="/video/:itemId" element={<PlyrWrap />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}