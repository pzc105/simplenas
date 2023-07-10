import { useState } from "react";
import { Box, Button, FormControl, FormLabel, Input, FormHelperText } from "@mui/material";

import * as User from './prpc/user_pb.js'
import userService from './rpcClient.js'

export default function SubtitleUploader({ itemId, onClose }) {
  const [fileData, setFileData] = useState(new Map());

  const handleFileSelect = (event) => {
    const files = event.target.files
    for (const file of files) {
      if (file) {
        const reader = new FileReader();
        reader.onload = () => {
          setFileData(fileData.set(file.name, new Uint8Array(reader.result)))
        }
        reader.readAsArrayBuffer(file)
      }
    }
  }

  const handleUpload = () => {
    if (fileData.size === 0) return
    var req = new User.UploadSubtitleReq()
    req.setItemId(itemId)
    var fileList = []
    console.log(fileData)
    for (const [key, value] of fileData.entries()) {
      var f = new User.SubtitleFile()
      f.setName(key)
      f.setContent(value)
      fileList.push(f)
    }
    console.log(fileList)
    req.setSubtitlesList(fileList)
    userService.uploadSubtitle(req, {}, (err, res) => {
      if (err != null) {
      }
      onClose()
    })
  }

  return (
    <Box>
      <FormControl>
        <FormLabel>{'字幕上传'}</FormLabel>
        <Input type="file" onChange={handleFileSelect} accept="*/*" inputProps={{
          multiple: true
        }} />
        <FormHelperText>{'字幕上传'}</FormHelperText>
      </FormControl>
      <Button fullWidth variant="contained" color="primary" onClick={handleUpload} disabled={!fileData}>
        字幕上传
      </Button>
    </Box>
  );
}