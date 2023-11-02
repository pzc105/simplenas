import React, { useState } from 'react';
import { Avatar, Button, CssBaseline, TextField, FormControlLabel, Checkbox, Link, Grid, Box, Typography, Container } from '@mui/material';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import MuiAlert from '@mui/lab/Alert';
import { styled } from "@mui/material/styles";
import { useNavigate } from 'react-router-dom';

import userService from './rpcClient.js'
import * as User from './prpc/user_pb.js'
import * as utils from './utils.js'

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
  width: '100%', // Fix IE 11 issue.
  marginTop: theme.spacing(3),
}))

const MySubmit = styled(Button)(({ theme }) => ({
  margin: theme.spacing(3, 0, 2),
}))

export default function SignUp(props) {
  let navigate = useNavigate()

  const [showEmailWarn, setShowEmailWarn] = useState(false)
  const [showPwdWarn, setShowPwdWarn] = useState(false)
  const [showNameWarn, setShowNameWarn] = useState(false)
  const [agreementChecked, setAgreementChecked] = useState(false)
  const [showAgreementWarn, setShowAgreementWarn] = useState(false)
  const [emailWarnMsg, setEmailWarnMsg] = useState('')
  const [pwdWarnMsg, setPwdWarnMsg] = useState('')
  const [nameWarnMsg, setNameWarnMsg] = useState('')
  const [agreementWarnMsg, setAgreementWarnMsg] = useState('')

  const [rawUserInfo, setRawUserInfo] = useState(new Map())

  function handleChange(e) {
    setRawUserInfo(rawUserInfo.set(e.target.name, e.target.value))
    if (e.target.name === 'email' && utils.isEmail(e.target.value)) {
      setShowEmailWarn(false)
    }
    else if (e.target.name === 'name' && e.target.value !== '') {
      setShowNameWarn(false)
    }
    else if (e.target.name === 'password' && e.target.value !== '') {
      setShowPwdWarn(false)
    }
  }

  function handleCheck(e) {
    setAgreementChecked(e.target.checked)
    if (e.target.checked) {
      setShowAgreementWarn(false)
    }
  }

  // focus has left
  function emailOnBlur(e) {

    const email = e.target.value

    if (!utils.isEmail(email)) {
      setShowEmailWarn(true)
      setEmailWarnMsg("请输入有效的邮箱!")
      return
    }

    var emailInfo = new User.EmailInfo()
    emailInfo.setEmail(email)

    userService.isUsedEmail(emailInfo, {}, (err, isUsedEmailRet) => {
      if (err != null) {
        console.log(err)
        setShowEmailWarn(true)
        setEmailWarnMsg(err.message)
        return
      }

      if (isUsedEmailRet == null) {
        console.log("respone is null")
        return
      }
    })
  }

  function handleSubmit(e) {
    e.preventDefault()

    var haveError = false
    if (rawUserInfo.get('email') === undefined || rawUserInfo.get('email') === '') {
      setShowEmailWarn(true)
      setEmailWarnMsg('邮箱不能为空')
      haveError = true
    }
    if (rawUserInfo.get('name') === undefined || rawUserInfo.get('name') === '') {
      setShowNameWarn(true)
      setNameWarnMsg('用户名不能为空')
      haveError = true
    }
    if (rawUserInfo.get('password') === undefined || rawUserInfo.get('password') === '') {
      setShowPwdWarn(true)
      setPwdWarnMsg('密码不能为空')
      haveError = true
    }
    if (!agreementChecked) {
      setShowAgreementWarn(true)
      setAgreementWarnMsg("同意并勾选才能注册")
      haveError = true
    }
    if (haveError) {
      return
    }

    const CryptoJS = require("crypto-js")
    const email = rawUserInfo.get('email')
    const name = rawUserInfo.get('name')
    const password = CryptoJS.MD5(rawUserInfo.get('password')).toString()

    var registerInfo = new User.RegisterInfo()
    var userInfo = new User.UserInfo()
    userInfo.setEmail(email)
    userInfo.setName(name)
    userInfo.setPasswd(password)
    registerInfo.setUserInfo(userInfo)


    userService.register(registerInfo, {}, (err, registerRet) => {
      if (err != null) {
        console.log(err)
        alert("注册失败")
        return
      }

      navigate('/signin')
    })
  }

  return (
    <Container component="main" maxWidth="xs" >
      <CssBaseline />
      <PaperDiv>
        <MyAvatar>
          <LockOutlinedIcon />
        </MyAvatar>
        <Typography component="h1" variant="h5">
          Sign up
        </Typography>
        <MyForm
          onSubmit={handleSubmit}
          method="post"
          noValidate>
          <Grid container spacing={2} direction='row'>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                onChange={handleChange}
                onBlur={emailOnBlur} />
              {showEmailWarn ? <Alert severity="warning">{emailWarnMsg}</Alert> : null}
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                id="name"
                label="name"
                name="name"
                autoComplete="name"
                onChange={handleChange} />
              {showNameWarn ? <Alert severity="warning">{nameWarnMsg}</Alert> : null}
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                name="password"
                label="password"
                type="password"
                id="password"
                autoComplete="current-password"
                onChange={handleChange} />
              {showPwdWarn ? <Alert severity="warning">{pwdWarnMsg}</Alert> : null}
            </Grid>
            <Grid item xs={12}>
              <FormControlLabel
                control={<Checkbox value="allowExtraEmails" color="primary" />}
                label="I want to receive inspiration, marketing promotions and updates via email."
                onChange={handleCheck} />
              {showAgreementWarn ? <Alert severity="warning">{agreementWarnMsg}</Alert> : null}
            </Grid>
          </Grid>
          <MySubmit
            type="submit"
            fullWidth
            variant="contained"
            color="primary">
            注册
          </MySubmit>
          <Grid container justifyContent="flex-end">
            <Grid item>
              <Link onClick={() => navigate("/signin")} variant="body2">
                已经有账号? 点我登录
              </Link>
            </Grid>
          </Grid>
        </MyForm>
      </PaperDiv>
      <Box mt={5}>
        <Copyright />
      </Box>
    </Container>
  );
}