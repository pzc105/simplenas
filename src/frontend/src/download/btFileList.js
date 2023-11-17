import React, { useState, useEffect, useRef } from 'react';
import { Container, Button, FormControl, InputLabel, MenuItem, Select, FormControlLabel, FormGroup, Checkbox }
  from "@mui/material";

import { useSelector, useDispatch } from 'react-redux';
import * as store from '../store.js'
import * as utils from '../utils.js'
import FolderSelector from './folderSelector.js'
import * as category from '../prpc/category_pb.js'
import * as User from '../prpc/user_pb.js'
import * as Bt from '../prpc/bt_pb.js'
import userService from '../rpcClient.js'


const FileListHandler = ({ infoHash }) => {
  let btViodeFiles = useSelector((state) => store.selectBtVideoFiles(state, infoHash))
  const [selectedVideoFiles, selectVideoFiles] = useState({})
  const selectedDirId = useRef(-1)
  const isMouseDown = useSelector((state) => store.isDownloadPageMouseDown(state))
  let sortedVideos = []

  const sortVideos = () => {
    if (!btViodeFiles) {
      return
    }
    let tmp = Object.values(btViodeFiles)
    tmp.sort((a, b) => {
      if (a.meta.format.filename < b.meta.format.filename) {
        return -1;
      }
      if (a.meta.format.filename > b.meta.format.filename) {
        return 1;
      }
      return 0;
    })
    sortedVideos = tmp
  }
  sortVideos()

  const handleChange = (e, index) => {
    selectVideoFiles({ ...selectedVideoFiles, [index]: e.target.checked })
  }

  const saveVideos = () => {
    var req = new User.AddBtVideosReq()
    var i = new Bt.InfoHash()
    i.setVersion(infoHash.version)
    i.setHash(infoHash.hash)
    req.setInfoHash(i)
    req.setCategoryItemId(selectedDirId.current)
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
      <FolderSelector select={(id) => selectedDirId.current = id} />
      <Button
        variant="contained"
        onClick={saveVideos}
        color="primary" >
        保存
      </Button>
      <FormGroup>{
        sortedVideos.map((f) => {
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
        })
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

