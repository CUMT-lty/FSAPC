package m

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
)

// Define the model resource structure
type Model struct { // 模型
	// ModelID 根据其他信息生成，不需要保存
	ModelType      int    `json:"ModelType"`      // 0: 层聚合模型，1: 全局模型，-1: 客户端本地模型
	ModelTier      int    `json:"ModelTier"`      // -1: 表示时全局模型，[0,4]表示是层聚合模型
	TierRound      int    `json:"TierRound"`      // -1: 表示是全局模型，0 之后表示是层模型的层更新轮次
	GlobalRound    int    `json:"GlobalRound"`    // 全局通信轮次
	ModelEncryAddr string `json:"ModelEncryAddr"` // 模型的 IPFS 保存地址
	Timestamp      int64  `json:"Timestamp"`      // 模型保存到链上的时间
}

type AM struct { // 表示模型 Model 在访问控制中的必要信息
	ModelID string `json:"ModelID"`
}

func (model *Model) GenAM() AM {
	return AM{
		ModelID: model.GetID(),
	}
}

func (model *Model) GetID() string {
	return fmt.Sprintf("%x", sha256.Sum256(
		[]byte(strconv.Itoa(model.ModelType)+
			strconv.Itoa(model.ModelTier)+
			strconv.Itoa(model.TierRound)+
			strconv.Itoa(model.GlobalRound))))
}

func (model *Model) ToBytes() []byte {
	bs, err := json.Marshal(*model)
	if err != nil {
		return nil
	}
	return bs
}
