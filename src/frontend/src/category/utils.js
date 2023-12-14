import * as store from '../store.js'

import * as User from '../prpc/user_pb.js'
import * as Category from '../prpc/category_pb.js'
import userService from '../rpcClient.js'

export function navigateToItem(navigate, navigateParams, itemId, shareid) {
  let path = "/citem?itemid=" + itemId
  if (shareid) {
    path += "&shareid=" + shareid
  }
  navigate(path, navigateParams)
}

export function navigateToVideo(navigate, navigateParams, itemId, shareid) {
  let path = "/video?itemid=" + itemId
  if (shareid) {
    path += "&shareid=" + shareid
  }
  navigate(path, navigateParams)
}

export const querySubItems = ({ itemId, shareid, dispatch, callback, pageNum, pageRows, desc }) => {
  var req = new User.QuerySubItemsReq()
  req.setParentId(itemId)
  req.setPageNum(pageNum)
  req.setRows(pageRows)
  req.setDesc(desc)
  if (shareid) {
    req.setShareId(shareid)
  }
  userService.querySubItems(req, {}, (err, respone) => {
    if (err == null) {
      let subItems = []
      respone.getItemsList().map((i) => {
        let obj = i.toObject()
        subItems.push(obj)
        dispatch(store.categorySlice.actions.updateItem(obj))
        return null
      })
      if (callback) {
        callback(subItems)
      }
    } else {
      console.log(err)
    }
  })
}

export const queryItem = (itemId, shareId, dispatch, callback) => {
  var req = new User.QueryItemInfoReq()
  req.setItemId(itemId)
  if (shareId) {
    req.setShareId(shareId)
  }
  userService.queryItemInfo(req, {}, (err, res) => {
    if (err != null || !res) {
      return
    }
    const itemInfo = res.getItemInfo()
    if (itemInfo.getTypeId() === Category.CategoryItem.Type.VIDEO) {
      dispatch(store.categorySlice.actions.updateVideoInfo({ itemId: itemInfo.getId(), videoInfo: res.getVideoInfo().toObject() }))
    }
    dispatch(store.categorySlice.actions.updateItem(itemInfo.toObject()))
    if (callback) {
      callback(itemInfo.toObject())
    }
  })
}

export const isDirectory = (item) => {
  if (item == null) {
    return false
  }
  return item.typeId === Category.CategoryItem.Type.DIRECTORY || item.typeId === Category.CategoryItem.Type.HOME
}