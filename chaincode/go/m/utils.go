package m

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// ok 信息
var OK = []byte("OK")

// 将 json 字符串解析为 Policy 类型
func ParsePolicy(arg string) (Policy, error) {
	policy := Policy{}
	err := json.Unmarshal([]byte(arg), &policy)
	return policy, err
}

// 将 json 字符串解析为 Client 类型
func ParseClient(arg string) (Client, error) {
	client := Client{}
	err := json.Unmarshal([]byte(arg), &client)
	return client, err
}

// 将 json 字符串解析为 Model 类型
func ParseModel(arg string) (Model, error) {
	model := Model{}
	err := json.Unmarshal([]byte(arg), &model)
	return model, err
}

// 将 json 字符串解析为 Model 类型
func ParseOP(arg string) (OP, error) {
	op := OP{}
	err := json.Unmarshal([]byte(arg), &op)
	return op, err
}

// 将 json 字符串解析为 FSAPRequest 类型
func ParseFSAPRequest(arg string) (FSAPRequest, error) {
	fsap_req := FSAPRequest{}
	err := json.Unmarshal([]byte(arg), &fsap_req)
	return fsap_req, err
}

func Sha256Addr(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}
