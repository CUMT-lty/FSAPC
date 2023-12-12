package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/newham/fabric-iot/chaincode/go/m"
	"time"
)

// Model Management Contract
type ModelManagementContract interface {
	Init(shim.ChaincodeStubInterface) sc.Response
	Invoke(shim.ChaincodeStubInterface) sc.Response
	AddModel(shim.ChaincodeStubInterface, []string) sc.Response
	RequestModel(shim.ChaincodeStubInterface, []string) sc.Response
	Synchro() sc.Response
}

// Define the Contract structure
type ChainCode struct {
	ModelManagementContract
}

func NewModelManagementContract() ModelManagementContract {
	return new(ChainCode)
}

// 在这里进行模型初始化
func (cc *ChainCode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	//initialModel := m.Model{ // 初始化模型
	//	ModelType:      1,  // 0: 层聚合模型，1: 全局模型
	//	ModelTier:      -1, // -1: 表示全局模型，[0,4]表示层聚合模型
	//	TierRound:      -1, // -1: 表示是全局模型，0 之后表示是层模型的层更新轮次
	//	GlobalRound:    0,  // 全局通信轮次
	//	ModelEncryAddr: "",
	//}
	//modelID := initialModel.GetID()
	//err := APIstub.PutState(modelID, initialModel.ToBytes()) // 将初始化模型保存至状态数据库，模型 ID : 模型信息
	//if err != nil {                                          // 如果出错
	//	shim.Error("Error while initializing model!")
	//}
	return shim.Success(m.OK)
}

func (cc *ChainCode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "AddModel" {
		return cc.AddModel(APIstub, args)
	} else if function == "RequestModel" {
		return cc.RequestModel(APIstub, args)
	} else if function == "Synchro" {
		return cc.Synchro()
	}
	return shim.Error("Invalid Smart Contract function name.")
}

// 向链上添加模型，需要保存完整的模型信息，参数是 Model 类型的 json 字符串
func (cc *ChainCode) AddModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 检查参数数量
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 将 json 字符串参数解析为 Model 类型
	model, err := m.ParseModel(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// TODO: 加个时间
	ts := time.Date(2023, time.October, 15, 9, 47, 30, 22, time.UTC).Unix()
	model.Timestamp = ts
	// 将模型保存到状态数据库
	// 注意，这里以 modelID 为 key，也就等价于保存了模型的其他信息，因为 modelID 是根据其他必要信息哈希生成的
	err = APIstub.PutState(model.GetID(), model.ToBytes()) // ModelID : 模型的信息（里面有 IPFS 的 CID）
	if err != nil {                                        // 如果出错
		shim.Error("Error while storing model!")
	}
	// 如果保存成功
	return shim.Success([]byte("The model information has been on-chain successfully!"))
}

// 请求模型
func (cc *ChainCode) RequestModel(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 检查参数个数（参数应该是一个 json 字符串）
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 将 json 字符串参数解析为 Model 类型
	model, err := m.ParseModel(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// 在状态数据库中查询模型
	modelBytes, err := APIstub.GetState(model.GetID())
	if err != nil || modelBytes == nil { // 如果出错或查询不到
		return shim.Error("Error in model information!")
	}
	modelOnChain, err := m.ParseModel(string(modelBytes))
	return shim.Success([]byte(modelOnChain.ModelEncryAddr))
}

func (cc *ChainCode) Synchro() sc.Response {
	return shim.Success(m.OK)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(NewModelManagementContract())
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
