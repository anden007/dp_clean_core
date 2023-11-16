package misc

import (
	"fmt"
	"os"
)

func WriteWithFileWrite(filePath, fileName, content string) {
	targetFile := filePath + fileName
	if exist, _ := PathExists(filePath); !exist {
		if err := os.Mkdir(filePath, os.ModePerm); err != nil {
			fmt.Println("× 生成目录失败，请检查文件路径", err.Error())
			return
		}
	}
	if existOldFile, _ := PathExists(targetFile); existOldFile {
		targetFile += ".new"
	}

	if fileObj, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err == nil {
		defer fileObj.Close()
		if fileObj != nil {
			if _, err := fileObj.WriteString(content); err == nil {
				fmt.Printf("√ %s  --生成成功\n", targetFile)
				return
			}
		}
	} else {
		fmt.Println("× 生成文件失败，请检查文件路径", err.Error())
	}
	fmt.Printf("× %s  --生成失败\n", targetFile)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
