package video

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

func GetMetadata(fileName string) (*Metadata, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-show_format", "-show_streams", "-print_format", "json", fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	var ret Metadata
	json.Unmarshal(output, &ret)
	return &ret, nil
}

func GetMetadataByStdin(content []byte) (*Metadata, error) {
	cmd := exec.Command("ffprobe", "-i", "-", "-v", "quiet", "-show_format", "-show_streams", "-print_format", "json")
	var b bytes.Buffer
	cmd.Stdout = &b
	var wb bytes.Buffer
	wb.Write(content)
	cmd.Stdin = &wb
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	err := cmd.Wait()
	if err != nil {
		return nil, err
	}
	var ret Metadata
	json.Unmarshal(b.Bytes(), &ret)
	return &ret, nil
}
