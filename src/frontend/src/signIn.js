import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import { Container, Typography, Avatar, Button, CssBaseline, TextField, FormControlLabel, Checkbox, Link, Grid, Box } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import MuiAlert from '@mui/lab/Alert';
import { styled } from "@mui/material/styles";
import { useNavigate } from 'react-router-dom';

import * as utils from './utils.js'
import userService from './rpcClient.js'
import * as User from './prpc/user_pb.js'
import { userSlice, btSlice } from './store.js'
import * as test from './test.js'

function Alert(props) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

function Copyright() {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {'Copyright © '}
      <Link color="inherit" href="">
        PNAS
      </Link>{' '}
      {new Date().getFullYear()}
      {'.'}
    </Typography>
  );
}

const PaperDiv = styled("div")(({ theme }) => ({
  marginTop: theme.spacing(8),
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
}))

const MyAvatar = styled(Avatar)(({ theme }) => ({
  margin: theme.spacing(1),
  backgroundColor: theme.palette.secondary.main,
}))

const MyForm = styled("form")(({ theme }) => ({
  width: '100%',
  marginTop: theme.spacing(3),
}))

const MySubmit = styled(Button)(({ theme }) => ({
  margin: theme.spacing(3, 0, 2),
}))

export function checkLogined() {
  if (utils.isLoginFailed()) {
    return
  }
  var loginInfo = new User.LoginInfo()
  userService.isLogined(loginInfo, {}, (err, loginRet) => {
    if (err != null) {
      utils.enterNeedLogin()
      return
    }
  })
}

export function FastLogin(navigate, dispatch, onFailed) {
  var loginInfo = new User.LoginInfo()
  let re = localStorage.getItem("rememberMe") === "true"
  loginInfo.setEmail("")
  loginInfo.setPasswd("")
  loginInfo.setRememberMe(re)
  utils.enterLogining()
  userService.fastLogin(loginInfo, {}, (err, loginRet) => {
    if (err != null) {
      if (err.name !== "RpcError") {
        utils.enterLoginFailed()
      } else {
        utils.enterNeedLogin()
      }
      if (onFailed) {
        onFailed(err)
      }
      return
    }
    let re = loginRet.getRememberMe()
    localStorage.setItem("rememberMe", re)
    const userInfo = loginRet.getUserInfo().toObject()
    dispatch(userSlice.actions.setUserInfo(userInfo))
    utils.enterLogined()
  })
}

export default function SignIn(props) {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [showEmailWarn, setShowEmailWarn] = useState(false)
  const [showPwdWarn, setShowPwdWarn] = useState(false)
  const [emailWarnMsg, setEmailWarnMsg] = useState('')
  const [pwdWarnMsg, setPwdWarnMsg] = useState('')
  const [rememberMe, setRememberMe] = useState(false)
  const [rawUserInfo, setRawUserInfo] = useState(new Map())

  function handleChange(e) {
    setRawUserInfo(rawUserInfo.set(e.target.name, e.target.value))
    if (e.target.name === 'email' && utils.isEmail(e.target.value)) {
      setShowEmailWarn(false)
    }
    else if (e.target.name === 'password' && e.target.value !== '') {
      setShowPwdWarn(false)
    }
  }

  function handleCheck(e) {
    setRememberMe(e.target.checked)
  }

  function handleSubmit(e) {
    e.preventDefault()

    var haveError = false
    if (rawUserInfo.get('email') === undefined || rawUserInfo.get('email') === '') {
      setShowEmailWarn(true)
      setEmailWarnMsg('邮箱不能为空')
      haveError = true
    }
    if (rawUserInfo.get('password') === undefined || rawUserInfo.get('password') === '') {
      setShowPwdWarn(true)
      setPwdWarnMsg('密码不能为空')
      haveError = true
    }
    if (haveError) {
      return
    }

    const CryptoJS = require("crypto-js")
    var email = rawUserInfo.get('email')
    var passwd = CryptoJS.MD5(rawUserInfo.get('password')).toString()

    var loginInfo = new User.LoginInfo()
    loginInfo.setEmail(email)
    loginInfo.setPasswd(passwd)
    loginInfo.setRememberMe(rememberMe)
    utils.enterLogining()
    userService.login(loginInfo, {}, (err, loginRet) => {
      if (err != null) {
        console.log(err)
        alert("登录失败")
        utils.enterLoginFailed()
        return
      }
      utils.enterLogined()
      let re = loginRet.getRememberMe()
      localStorage.setItem("rememberMe", re)
      const userInfo = loginRet.getUserInfo().toObject()
      dispatch(userSlice.actions.setUserInfo(userInfo))
      dispatch(btSlice.actions.removeAllTorrent())
      navigate("/download")
    })
  }

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <PaperDiv>
        <MyAvatar>
          <LockOutlinedIcon />
        </MyAvatar>
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <MyForm
          onSubmit={handleSubmit}
          method="post"
          noValidate>
          <TextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email Address"
            name="email"
            autoComplete="email"
            onChange={handleChange}
            autoFocus />
          {showEmailWarn ? <Alert severity="warning">{emailWarnMsg}</Alert> : null}
          <TextField
            variant="outlined"
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            onChange={handleChange} />
          {showPwdWarn ? <Alert severity="warning">{pwdWarnMsg}</Alert> : null}
          <FormControlLabel
            control={<Checkbox value="remember" color="primary" />}
            label="Remember me"
            onChange={handleCheck} />
          <MySubmit
            type="submit"
            fullWidth
            variant="contained"
            color="primary">
            登录
          </MySubmit>
          <Grid container>
            <Grid item xs>
              <Link href="#" variant="body2">
                Forgot password?
              </Link>
            </Grid>
            <Grid item>
              <Link href="/signup" variant="body2">
                {"没有账号? 注册吧"}
              </Link>
            </Grid>
          </Grid>
        </MyForm>
      </PaperDiv>
      <Box mt={8}>
        <Copyright />
      </Box>
    </Container>
  );
}