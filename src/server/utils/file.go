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

func GetNotZeroFilesByFileExtension(searchPath string, suffixes []string) []string {
	var ret []string
	filepath.Walk(searchPath, func(pathstr string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if pathstr != searchPath {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() && info.Size() > 0 {
			for _, s := range suffixes {
				if strings.HasSuffix(info.Name(), s) {
					stat, err := os.Stat(pathstr)
					if err != nil || stat.Size() == 0 {
						return nil
					}
					ret = append(ret, info.Name())
					return nil
				}
			}
		}
		return nil
	})
	return ret
}
