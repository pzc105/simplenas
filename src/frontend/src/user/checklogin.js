
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import * as utils from '../utils.js'
import { FastLogin, checkLogined } from './signIn.js';

export default function CheckLoginHandler() {
  const navigate = useNavigate()
  const dispatch = useDispatch()

  useEffect(() => {
    const check = () => {
      if (utils.isLogined()) {
        checkLogined()
      }
      if (utils.needFastLogin()) {
        FastLogin(navigate, dispatch)
      }
      if (utils.isLoginFailed()) {
        alert("需要重新登录")
        utils.enterMnaullyLogin()
        navigate("/signin")
      }
    }

    let myInterval = setInterval(check, 30000)
    return () => {
      clearInterval(myInterval);
    };
  }, [dispatch, navigate])

  return null
}