package misc

import (
	"time"

	"github.com/btnguyen2k/consu/olaf"
)

var olafInstance *olaf.Olaf

func NewId() (result uint64) {
	result = olafInstance.Id64()
	return
}

func NewHexId() (result string) {
	result = olafInstance.Id64Hex()
	return
}

func init() {
	olafInstance = olaf.NewOlaf(time.Now().Unix() / 1000)
}
