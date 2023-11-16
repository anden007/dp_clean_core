package misc

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// PortCanUse 检测TCP端口是否可用
func PortCanUse(portNumber int) bool {
	findCount := 0
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr :%d", portNumber)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	resArray := strings.Split(resStr, "\r\n")
	for _, item := range resArray {
		ipWithPort := regexp.MustCompile(`\d+.\d+.\d+.\d+:\d+`).FindString(item)
		Port := regexp.MustCompile(`:\d+`).FindString(ipWithPort)
		if Port == fmt.Sprintf(":%d", portNumber) {
			findCount += 1
		}
	}

	return findCount == 0
}

// GetRndUnUsePortNumber 获取未使用随机端口
func GetRndUnUsePortNumber() (result int, err error) {
	tryCount := 0
	minP := 49152
	maxP := 65535
	rand.New(rand.NewSource(time.Now().UnixMicro()))
	port := rand.Intn(maxP-minP) + minP
	for {
		if PortCanUse(port) {
			result = port
			return
		} else {
			tryCount++
			if tryCount > 100 {
				result = -1
				err = errors.New("func GetRndUnUsePortNumber have exceeded the number attempts allowed(100)!")
				return
			}
		}
	}
}
