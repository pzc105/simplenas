import * as Category from './prpc/category_pb.js'

export function isEmail(s) {
  return (/^[a-zA-Z0-9.!#$%&'*+\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/.test(s))
}

export function getTorrentStatus(t) {
  const info_hash = t.getInfoHash()
  return {
    infoHash: {
      version: info_hash.getVersion(),
      hash: Array.from(info_hash.getHash())
        .map(byte => ('0' + byte.toString(16)).slice(-2))
        .join(''),
    },
    name: t.getName(),
    downloadPayloadRate: t.getDownloadPayloadRate(),
    total: t.getTotal(),
    totalDone: t.getTotalDone(),
    progress: t.getProgress(),
  }
}

export function enterNeedLogin() {
  sessionStorage.setItem("login_state", "need_fast")
}

export function enterLogining() {
  sessionStorage.setItem("login_state", "logining")
}

export function enterLoginFailed() {

  sessionStorage.setItem("login_state", "failed")
}

export function enterLogined() {
  sessionStorage.setItem("login_state", "logined")
}

export function enterMnaullyLogin() {
  sessionStorage.setItem("login_state", "manully_login")
}

export function needFastLogin() {
  const isLogined = sessionStorage.getItem("login_state")
  if (isLogined === undefined || isLogined === null || isLogined === "need_fast") {
    return true
  }
  return false
}

export function isLogining() {
  const isLogined = sessionStorage.getItem("login_state")
  if (isLogined === "logining") {
    return true
  }
  return false
}

export function isLoginFailed() {
  const isLogined = sessionStorage.getItem("login_state")
  if (isLogined === "failed") {
    return true
  }
  return false
}

export function isLogined() {
  const isLogined = sessionStorage.getItem("login_state")
  if (isLogined === "logined") {
    return true
  }
  return false
}

export function isManuallyLogin() {
  const isLogined = sessionStorage.getItem("login_state")
  if (isLogined === "manully_login") {
    return true
  }
  return false
}

export function isNumber(str) {
  // 使用正则表达式匹配数字格式
  // const pattern = /^[0-9]+$/;
  // return pattern.test(str);
  return !isNaN(str)
}

export function isVideoItem(item) {
  return item.typeId === Category.CategoryItem.Type.VIDEO
}