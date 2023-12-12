package m

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
)

type Client struct {
	ClientID    int `json:"ClientID"`    // 只是客户端的编号，取值 [0, 39]
	ClientTier  int `json:"ClientTier"`  // 客户端所在通信层
	TierRound   int `json:"TierRound"`   // 客户端所在层的通信轮次
	GlobalRound int `json:"GlobalRound"` // 当前全局轮次
	//Timestamp   int64 `json:"Timestamp"`   // 发起请求的时间戳
}

type AC struct { // 表示客户端 CLient 在访问控制中的必要信息
	ClientTier  int `json:"ClientTier"`
	GlobalRound int `json:"GlobalRound"`
}

func (client *Client) GenAC() AC {
	return AC{
		ClientTier:  client.ClientTier,
		GlobalRound: client.GlobalRound,
	}
}

func (client *Client) GetID() string {
	return fmt.Sprintf("%x", sha256.Sum256(
		[]byte(strconv.Itoa(client.ClientID)+
			strconv.Itoa(client.ClientTier))))
}

func (client *Client) ToBytes() []byte {
	bs, err := json.Marshal(*client)
	if err != nil {
		return nil
	}
	return bs
}
