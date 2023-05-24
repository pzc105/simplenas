package utils

import (
	"os"
	"path"
	"path/filepath"
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

func GetFilesByFileExtension(searchPath string, suffixes []string) []string {
	var ret []string
	filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != searchPath {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() && info.Size() > 0 {
			for _, s := range suffixes {
				if strings.HasSuffix(info.Name(), s) {
					ret = append(ret, info.Name())
					return nil
				}
			}
		}
		return nil
	})
	return ret
}
