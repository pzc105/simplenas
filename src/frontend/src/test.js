import { useState } from "react";
import { Box, Button, FormControl, FormLabel, Input, FormHelperText } from "@mui/material";

export function Test({ itemId }) {
  const [fileData, setFileData] = useState(new Map());

  const handleFileSelect = (event) => {
    const files = event.target.files
    for(const file of files) {
      console.log(file)
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
    if (!fileData) return
    
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