package video

import (
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
