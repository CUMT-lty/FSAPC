package m

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type ABACRequest struct {
	AS AS
	AO AO
}

func (r *ABACRequest) ToBytes() []byte {
	b, err := json.Marshal(*r)
	if err != nil {
		return nil
	}
	return b
}

type Attrs struct {
	DeviceId  string
	UserId    string
	Timestamp int64
}

func (a Attrs) GetId() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(a.UserId+a.DeviceId)))
}

func (r ABACRequest) GetAttrs() Attrs {
	return Attrs{DeviceId: r.AO.DeviceId, UserId: r.AS.UserId, Timestamp: time.Now().Unix()}
}

// 用来映射请求 json 字符串：某客户端对某模型请求某个操作，应该是一个嵌套的 json
type FSAPRequest struct {
	Client Client
	Model  Model
	OP     OP // 0 代表客户端请求全局模型，1 代表客户端上传本地模型
}

//func (r *FSAPRequest) ToBytes() []byte {
//	bs, err := json.Marshal(*r)
//	if err != nil {
//		return nil
//	}
//	return bs
//}

// 用来存储生成访问控制请求所需的信息
type Attrs struct {
	ClientTier  int
	GlobalRound int
	ModelID     string
}

func (r *FSAPRequest) GetAttrs() Attrs { // 由 ABACRequest 来返回 Attrs
	return Attrs{
		ClientTier:  r.Client.ClientTier,
		GlobalRound: r.Client.GlobalRound,
		ModelID:     r.Model.GetID(),
	}
}

func (r *FSAPRequest) ToBytes() []byte { // 由 ABACRequest 来返回 Attrs
	bs, err := json.Marshal(*r)
	if err != nil {
		return nil
	}
	return bs
}

// 用生成的这个 ID 和 policyID 进行比较
func (a *Attrs) GetID() string {
	return fmt.Sprintf("%x", sha256.Sum256(
		[]byte(strconv.Itoa(a.ClientTier)+a.ModelID)))
}
