# FSAPC
FSAPC 是一个基于 Hyperledger Fabric 框架实现的异步分层联邦训练管理平台， 采用 ABAC 访问控制模型作为基础权限控制模型，用 golang 编写 chaincode 智能合约；提供了完整脚本，轻松启动 fabric 网络（如果你使用过 fabric 提供的 sample 网络，应该能体会到差别）。  

## 准备工作
### fabric 环境
[ github : hyperledger-fabric](https://github.com/hyperledger/fabric)  
[ doc 1.4 ](https://hyperledger-fabric.readthedocs.io/en/release-1.4/)
[ node-sdk 1.4 ](https://hyperledger.github.io/fabric-sdk-node/release-1.4/module-fabric-network.html)

### 操作系统 & 必要软件
OS|版本
-|-
ubuntu|16.04  

软件|版本
-|-
git|>=2.0.0
docker|>=19
docker-compose|>=1.24
node|>=12
golang|>=1.10
Hyperledger Fabric|=1.4.*

实验平台基于Fabric 1.4.3，默认下载的 Fabric 版本是 1.4.3，如果需要修改版本，请更改 `bootstrap.sh`
```
VERSION=1.4.3
CA_VERSION=1.4.3
```

## 部署和使用 FSAPC 平台

### 克隆仓库
使用 `git clone` 或下载后解压缩

### 下载并配置 Hyperledger Fabric
下载源码和 bin  
```bash
./bootstrap.sh
```
export PATH
```bash
export PATH=$PATH:$(pwd)/fabric-samples/bin
```

## 1.目录结构

|-chaincode....................................链码目录  
    |--go....................................................golang编写的chaincode代码

|-client............................................交互客户端  
    |--nodejs..............................................nodejs编写的客户端代码

|-network........................................fabric-iot网络配置和启动脚本  
    |--channel-artifacts..............................channel配置文件，创世区块目录  
    |--compose-files...................................docker-compose文件目录  
    |--conn-conf.........................................链接docker网络配置目录  
    |--crypto-config....................................存放证书、秘钥文件目录  
    |--scripts...............................................cli运行脚本目录  
    |--ccp-generate.sh...............................生成节点证书目录配置文件脚本  
    |--ccp-template.json.............................配置文件json模板  
    |--ccp-template.yaml.............................配置文件yaml模板  
    |--clear-docker.sh................................清理停用docker镜像脚本  
    |--configtx.yaml....................................configtxgen的配置文件  
    |--crypto-config.yaml...........................生成证书、秘钥的配置文件  
    |--down.sh............................................关闭网络shell脚本  
    |--env.sh...............................................其他shell脚本import的公共配置  
    |--generate.sh......................................初始化系统配置shell脚本  
    |--install-cc.sh......................................安装链码shell脚本  
    |--rm-cc.sh...........................................删除链码shell脚本  
    |--up.sh................................................启动网络shell脚本  

## 如何搭建网络

### 生成证书、配置文件等
进入目录  
```shell
cd network
```
运行脚本  
```shell
./generate.sh
```
证书保存在`crypto-config`中

### 启动网络
```shell
./up.sh
```

### 安装 chaincode
```shell
./cc-install.sh
```

### 更新 chaincode（本步骤在修改了chain code代码后再执行，[new version]要大于上一个版本）
```shell
./cc-upgrade.sh [new version]
```
*chaincode被保存在`/chaincode/go`目录中，目前只用`golang`实现

### 关闭网络
```shell
./down.sh
```
**每次关闭网络会删除所有docker容器和镜像，请谨慎操作**

## 与区块链网络交互

## 使用 shell 脚本调用链码 
进入shell脚本目录
```bash
cd network
```
调用链码(示例代码在`cc-test.sh`)
```bash
./cc.sh invoke pc 1.0 go/pc AddPolicy '"{\"AS\":{\"userId\":\"13800010001\",\"role\":\"u1\",\"group\":\"g1\"},\"AO\":{\"deviceId\":\"D100010001\",\"MAC\":\"00:11:22:33:44:55\"}}"'
./cc.sh invoke pc 1.0 go/pc QueryPolicy '"40db810e4ccb4cc1f3d5bc5803fb61e863cf05ea7fc2f63165599ef53adf5623"'
```

## 使用 Node JS 客户端调用链码 
**推荐用这种方式测试链码**

### 初始化代码
进入客户端代码目录  
*目前只实现了nodejs的客户端  

```shell
cd client/nodejs
```
安装依赖
```shell
npm install
```

### 创建用户
创建管理员账户
```shell
node ./enrollAdmin.js
```
用管理员账户创建子用户
```shell
node ./registerUser.js
```
### 调用 chaincode
**如果这一步出错，尝试执行一下初始化脚本！**
```shell
node ./invoke.js [chaincode_name] [function_name] [args]
```

