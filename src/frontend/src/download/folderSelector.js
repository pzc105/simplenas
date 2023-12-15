import React, { useState, useEffect, useRef } from 'react';
import { FormControl, InputLabel, MenuItem, Select, Checkbox }
  from "@mui/material";

import { useDispatch, useSelector } from 'react-redux';
import * as store from '../store.js'

import * as category from '../prpc/category_pb.js'
import * as User from '../prpc/user_pb.js'
import userService from '../rpcClient.js'

export default function FolderSelector({ style, select }) {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const lastUsedParentDirId = useSelector((state) => store.selectlastUsedParentDirId(state))
  const lastUsedDirId = useSelector((state) => store.selectlastUsedDirId(state))
  const [nowPathItemId, setNowPathItemId] = useState(lastUsedParentDirId > 0 ? lastUsedParentDirId : userInfo.homeDirectoryId)
  const [subDirectories, setSubDirectories] = useState([])
  const [pathItem, setPathItem] = useState(null)
  const [selectedValue, setSelectedValue] = useState(lastUsedDirId);
  const [isSelectOpen, setIsSelectOpen] = useState(false);
  const selectRef = useRef(null);
  const dispatch = useDispatch()

  useEffect(() => {
    if (nowPathItemId <= 0) {
      return
    }
    let req = new User.QuerySubItemsReq()
    req.setParentId(nowPathItemId)
    userService.querySubItems(req, {}, (err, respone) => {
      if (err == null) {
        select(selectedValue)
        setPathItem(respone.getParentItem().toObject())
        let ds = []
        respone.getItemsList().map((i) => {
          let item = i.toObject()
          if (item.typeId === category.CategoryItem.Type.DIRECTORY) {
            ds.push(item)
          }
          return null
        })
        setSubDirectories(ds)
      } else {
        if (nowPathItemId != userInfo.homeDirectoryId) {
          setNowPathItemId(userInfo.homeDirectoryId)
        }
      }
    })
  }, [])

  useEffect(() => {
    let req = new User.QuerySubItemsReq()
    req.setParentId(nowPathItemId)
    userService.querySubItems(req, {}, (err, respone) => {
      if (err == null) {
        setPathItem(respone.getParentItem().toObject())
        let ds = []
        respone.getItemsList().map((i) => {
          let item = i.toObject()
          if (item.typeId === category.CategoryItem.Type.DIRECTORY) {
            ds.push(item)
          }
          return null
        })
        setSubDirectories(ds)
      }
    })
  }, [nowPathItemId])

  const boxOnClick = (id) => {
    setSelectedValue(id)
    dispatch(store.userSlice.actions.setlastUsedParentDirId(nowPathItemId))
    dispatch(store.userSlice.actions.setlastUsedDirId(id))
    select(id)
  }

  const handleSelectOpen = () => {
    setIsSelectOpen(!isSelectOpen);
  };

  return (
    <div style={style}>
      {
        pathItem ? <FormControl fullWidth>
          <InputLabel id="select-directory">选择保存目录</InputLabel>
          <Select
            labelId="select-directory"
            value={selectedValue}
            open={isSelectOpen}
            onClick={handleSelectOpen}
            inputProps={{ "aria-label": "Without label" }}
            ref={selectRef}
          >
            {nowPathItemId !== userInfo.homeDirectoryId && (
              <MenuItem key="back"
                value={0}
                onClick={
                  (e) => {
                    setNowPathItemId(pathItem.parentId)
                    e.stopPropagation()
                  }
                }
              >
                返回上一级
              </MenuItem>
            )}
            {subDirectories.map((dir) => (
              <MenuItem
                key={dir.id}
                value={dir.id}
                onClick={
                  (e) => {
                    setNowPathItemId(dir.id)
                    e.stopPropagation()
                  }
                }>
                <Checkbox
                  checked={selectedValue == dir.id}
                  onClick={
                    (e) => {
                      boxOnClick(dir.id)
                      e.stopPropagation()
                    }
                  } />
                <label >
                  {dir.name}
                </label>
              </MenuItem>
            ))}
          </Select>
        </FormControl> : null
      }
    </div>
  )
}
