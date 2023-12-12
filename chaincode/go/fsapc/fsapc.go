package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/newham/fabric-iot/chaincode/go/m"
	"time"
)

// Federal Secure Access Control Smart Contracts
type FSAPContract interface {
	Init(shim.ChaincodeStubInterface) sc.Response
	Invoke(shim.ChaincodeStubInterface) sc.Response
	AddValidClient(shim.ChaincodeStubInterface, []string) sc.Response
	CheckClient(shim.ChaincodeStubInterface, m.Client) bool // TODO: 不会被客户端直接调用
	AddPolicy(shim.ChaincodeStubInterface, []string) sc.Response
	QueryPolicy(shim.ChaincodeStubInterface, []string) sc.Response
	CheckValid(shim.ChaincodeStubInterface, []string) sc.Response
	Synchro() sc.Response
}

// Define the Contract structure
type ChainCode struct {
	FSAPContract
}

func NewFSAPContract() FSAPContract {
	return new(ChainCode)
}

func (cc *ChainCode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(m.OK)
}

func (cc *ChainCode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "AddValidClient" {
		return cc.AddValidClient(APIstub, args)
	} else if function == "AddPolicy" {
		return cc.AddPolicy(APIstub, args)
	} else if function == "QueryPolicy" {
		return cc.QueryPolicy(APIstub, args)
	} else if function == "CheckValid" {
		return cc.CheckValid(APIstub, args)
	} else if function == "Synchro" {
		return cc.Synchro()
	}
	return shim.Error("Invalid Smart Contract function name.")
}

// 添加合法的客户端
func (cc *ChainCode) AddValidClient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	client, err := m.ParseClient(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(client.GetID(), client.ToBytes())
	if err != nil {
		shim.Error(err.Error())
	}
	return shim.Success([]byte("Client has been registered successfully!"))
}

// 检查是否是注册过的合法客户端
func (cc *ChainCode) CheckClient(APIstub shim.ChaincodeStubInterface, client m.Client) bool {
	clientBytes, err := APIstub.GetState(client.GetID())
	if err != nil {
		return false
	}
	if clientBytes == nil {
		return false
	}
	return true
}

// 添加访问控制策略，参数是 Policy 类型字符串
func (cc *ChainCode) AddPolicy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 检查参数个数，参数应当是一个 json 类型的字符串
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 将 json 字符串参数解析为 Policy 类型
	policy, err := m.ParsePolicy(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	// 检查策略是否重复
	policyBytes, err := APIstub.GetState(policy.GetID())
	if policyBytes != nil {
		return shim.Error("Policy already exists!")
	}
	// 如果不是重复策略，将策略添加到状态数据库
	policy.AE = time.Now().Unix()                            // 添加系统时间
	err = APIstub.PutState(policy.GetID(), policy.ToBytes()) // 策略ID : 策略
	if err != nil {                                          // 在添加时出错
		return shim.Error(err.Error())
	}
	// 添加成功
	return shim.Success([]byte("The policy has been on-chain successfully!"))
}

// 查询访问控制策略是否存在，参数是 Policy 类型字符串
func (cc *ChainCode) QueryPolicy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 { // 参数不合法返回 -2
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	policy, err := m.ParsePolicy(args[0])
	if err != nil {
		return shim.Error("Error in policy parameter format!")
	}
	policyID := policy.GetID()
	policyBytes, err := APIstub.GetState(policyID) // 在状态数据库中查询
	if err != nil {                                // 查询出错返回 -1
		return shim.Error("Error while the querying!")
	}
	if policyBytes == nil { // 策略不存在
		return shim.Error("Policy does not exist!")
	}
	// 查询成功
	return shim.Success([]byte(policyID))
}

// 检查是否是一次合法的访问，args 是一个 FSAPRequest 的 json 字符串
func (cc *ChainCode) CheckValid(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	// 检查参数个数，参数形式是一个嵌套的 json 字符串
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 解析请求参数
	fsap_req, err := m.ParseFSAPRequest(args[0])
	if err != nil {
		return shim.Error("Error in request parameter format！")
	}
	// 检查客户端是否合法
	clientBytes, err := APIstub.GetState(fsap_req.Client.GetID())
	if err != nil || clientBytes == nil {
		return shim.Error("Invalid Client!")
	}
	// 判断操作类型：上传本地模型 / 请求全局模型
	// 0: 层聚合模型，1: 全局模型，-1: 客户端本地模型
	op := fsap_req.OP.OPType
	if op == 0 { // OPType == 0 代表客户端向区块链请求模型（也可能是不合法的请求）
		// 验证访问控制策略
		attrs := fsap_req.GetAttrs()
		reqID := attrs.GetID()
		policyBytes, err := APIstub.GetState(reqID)
		if err != nil || policyBytes == nil { // 没有查询到访问控制请求
			return shim.Error("Policy does not exist, access is denied!")
		}
		return shim.Success(m.OK)
	} else { // op == 1，OPType == 1 代表客户端提交本地模型
		if fsap_req.Model.ModelType != -1 {
			return shim.Error("Wrong model type, access is denied!")
		}
		return shim.Success(m.OK)
	}
}

func (cc *ChainCode) Synchro() sc.Response {
	return shim.Success(m.OK)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(NewFSAPContract())
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
