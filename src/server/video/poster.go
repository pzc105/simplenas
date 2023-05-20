package video

import (
	"os"
	"os/exec"
	"path"
	"pnas/log"
)

type GenPosterParams struct {
	InputFileName  string
	OutputFileName string
}

func GenPoster(params *GenPosterParams) error {
	err := os.MkdirAll(path.Dir(params.InputFileName), 0755)
	if err != nil {
		return err
	}
	var cmdParams []string
	cmdParams = append(cmdParams, "-i")
	cmdParams = append(cmdParams, params.InputFileName)
	cmdParams = append(cmdParams, "-y")
	cmdParams = append(cmdParams, "-vf")
	cmdParams = append(cmdParams, "select=gt(scene\\,0.5)")
	cmdParams = append(cmdParams, "-frames:v")
	cmdParams = append(cmdParams, "1")
	cmdParams = append(cmdParams, "-vsync")
	cmdParams = append(cmdParams, "vfr")
	cmdParams = append(cmdParams, params.OutputFileName)
	cmd := exec.Command("ffmpeg", cmdParams...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Warn(cmd.String())
		log.Warnf("%+v, %v", string(out), err)
		return err
	}
	return nil
}
