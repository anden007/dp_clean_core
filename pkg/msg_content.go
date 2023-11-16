package pkg

type MsgContent struct {
	ChannelId string `json:"channelId"`
	MsgType   string `json:"msgType"`
	Data      string `json:"data"`
}
