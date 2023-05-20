package utils

import (
	"path"
	"strings"
)

func FileNameFormat(fileName string) string {
	return strings.Replace(fileName, "\\", "/", -1)
}

func GetFileName(absFileName string) string {
	absFileName = strings.Replace(absFileName, "\\", "/", -1)
	filenameall := path.Base(absFileName)
	filesuffix := path.Ext(absFileName)
	ret := filenameall[0 : len(filenameall)-len(filesuffix)]
	return ret
}
