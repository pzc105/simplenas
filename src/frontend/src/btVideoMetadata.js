import React, { useState, useEffect, useRef } from 'react';
import { Container, Button, FormControl, InputLabel, MenuItem, Select, FormControlLabel, FormGroup, Checkbox }
  from "@mui/material";

import { useSelector, useDispatch } from 'react-redux';
import * as store from './store.js'
import * as utils from './utils.js'

import * as category from './prpc/category_pb'
import * as User from './prpc/user_pb.js'
import * as Bt from './prpc/bt_pb.js'
import userService from './rpcClient.js'


const FolderSelector = ({ select }) => {
  const userInfo = useSelector((state) => store.selectUserInfo(state))
  const [nowPathItemId, setNowPathItemId] = useState(userInfo.homeDirectoryId)
  const [subDirectories, setSubDirectories] = useState([])
  const [pathItem, setPathItem] = useState(null)
  const [selectedValue, setSelectedValue] = useState('');
  const [isSelectOpen, setIsSelectOpen] = useState(false);
  const selectRef = useRef(null);

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
      } else {
        console.log(err)
      }
    })
  }, [nowPathItemId])

  const boxOnClick = (id) => {
    setSelectedValue(id)
    select(id)
  }

  const handleSelectOpen = () => {
    console.log("open")
    setIsSelectOpen(!isSelectOpen);
  };

  return (
    <FormControl sx={{ m: 1, minWidth: 120 }}>
      <InputLabel id="select-directory">Select</InputLabel>
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
              checked={selectedValue === dir.id}
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
    </FormControl>
  )
}

const FileListHandler = ({ infoHash }) => {
  const btViodeFiles = useSelector((state) => store.selectBtVideoFiles(state, infoHash))
  const [selectedVideoFiles, selectVideoFiles] = useState({})
  const [selectedDirId, selectDir] = useState(-1)
  const isMouseDown = useSelector((state) => store.isDownloadPageMouseDown(state))

  const handleChange = (e, index) => {
    selectVideoFiles({ ...selectedVideoFiles, [index]: e.target.checked })
  }

  const saveVideos = () => {
    var req = new User.AddBtVideosReq()
    var i = new Bt.InfoHash()
    i.setVersion(infoHash.version)
    i.setHash(infoHash.hash)
    req.setInfoHash(i)
    req.setCategoryItemId(selectedDirId)
    for (let fileIndex in selectedVideoFiles) {
      if (selectedVideoFiles[fileIndex]) {
        req.addFileIndexes(fileIndex)
      }
    }
    userService.addBtVideos(req, {}, (err, res) => {

    })
  }

  const [selectedIndexes, selectIndexes] = useState({})

  const onMouseEnter = (event, fileIndex) => {
    if (isMouseDown) {
      const f = selectedIndexes[fileIndex] ? false : true
      selectIndexes({ ...selectedIndexes, [fileIndex]: f })
      selectVideoFiles({ ...selectedVideoFiles, [fileIndex]: f })
    }
  }
  const onCheckBoxClick = (event, fileIndex) => {
    const f = selectedIndexes[fileIndex] ? false : true
    selectIndexes({ ...selectedIndexes, [fileIndex]: f })
  }

  return (
    <Container>
      <FolderSelector
        select={(id) => selectDir(id)} />
      <Button
        variant="contained"
        onClick={saveVideos}
        color="primary" >
        保存
      </Button>
      <FormGroup>{
        btViodeFiles ?
          btViodeFiles.map((f) => {
            return (
              <FormControlLabel
                onMouseEnter={(e) => onMouseEnter(e, f.fileIndex)}
                key={f.fileIndex}
                control={<Checkbox
                  checked={selectedIndexes[f.fileIndex] ? true : false}
                  onClick={(e) => onCheckBoxClick(e, f.fileIndex)}
                  onChange={(e) => handleChange(e, f.fileIndex)} name="gilad" />}
                label={"[" + utils.secondsToHHMMSS(f.meta.format.duration) + "] " + f.meta.format.filename} />
            )
          }) : null
      }
      </FormGroup>
    </Container>
  )
}

export default function BtVideosHandler({ infoHash }) {
  const userInfo = useSelector((state) => store.selectUserInfo(state))

  const dispatch = useDispatch()

  useEffect(() => {
    var req = new User.QueryBtVideosReq()
    var i = new Bt.InfoHash()
    i.setVersion(infoHash.version)
    i.setHash(infoHash.hash)
    req.setInfoHash(i)
    userService.queryBtVideos(req, {}, (err, respone) => {
      if (respone != null) {
        var localBtVideoMetadata = []
        var data = respone.getDataList()
        data.sort((a, b) => a.getFileIndex() - b.getFileIndex())
        data.map((d) => {
          localBtVideoMetadata.push(d.toObject())
          return null
        })
        localBtVideoMetadata.sort((a, b) => a.index - b.index)
        var payload = {
          infoHash: infoHash,
          btVideoMetadat: localBtVideoMetadata
        }
        dispatch(store.btSlice.actions.updateVideoFiles(payload))
      }
    })

    req = new User.QuerySubItemsReq()
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
  }, [userInfo, dispatch, infoHash])

  return (
    <Container>
      <FileListHandler infoHash={infoHash} />
    </Container >
  )
}