import { useState } from "react";
import { Box, Button, FormControl, FormLabel, Input, FormHelperText } from "@mui/material";

export default function FileUpload({ title, onUpload }) {
  const [fileData, setFileData] = useState(null);

  const handleFileSelect = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        setFileData(new Uint8Array(reader.result));
      };
      reader.readAsArrayBuffer(file);
    }
  };

  const handleUpload = () => {
    if (!fileData) return
    onUpload(fileData)
  }

  return (
    <Box>
      <FormControl>
        <FormLabel>{title}</FormLabel>
        <Input type="file" onChange={handleFileSelect} accept="image/*" />
        <FormHelperText>{title}</FormHelperText>
      </FormControl>
      <Button fullWidth variant="contained" color="primary" onClick={handleUpload} disabled={!fileData}>
        上传种子下载
      </Button>
    </Box>
  );
}