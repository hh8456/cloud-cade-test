## 服务端设计思路

服务端对于每个客户端 tcp 连接,都有一个接收协程和一个发送协程  
所有客户端对象,共同操作几个 map 对象, 这几个 map 对象保存了下面几种关系  
map1: 公会名 - 公会会长  
map2: 公会名 - 公会成员列表  
map3: 玩家   - 所属公会名字  
map4: 公会名 - 公会仓库列表  

玩家的所有操作,都是对上面几个 map 进行数据的读取和写入  

### TCP 私有协议和粘包处理

server - client 使用私有 TCP 协议通信, 协议格式: 前面4个字节 | 字符串  
前面 4 个字节表示后面字符串的长度, 比如玩家发送指令 /createAlliance 长度总共是 15 个字节  
那么发给服务端的逻辑包,总长度是 4 + 15 = 19 个字节  
前面 4  个字节格式化成 int 类型后, 值是 15, 即指令 /createAlliance 的长度  
后面 15 个字节就是字符串 "/createAlliance"  
程序中逻辑包最大长度是 4096 字节  

### 代码地图

server/gameApp/gameApp.go  
// 接受客户端的 tcp 连接,并生成一个 Client 对象  
func (g *GameApp) Listen(addr string) 

server/gameApp/client/client.go  
// 收到客户端发来一个逻辑包后, 网络模块会回调这个函数处理玩家登录和下面几个指令  
// /createAlliance /whichAlliance /allianceList /joinAlliance /dismissAlliance 
// /increaseCapacity /storeItem /destroyItem /clearup  
func (c *Client) recvData(binData []byte)  

// 客户端断线时,网络模块会回调这个函数  
func (c *Client) exit()  

//公会逻辑 /createAlliance /whichAlliance /allianceList /joinAlliance /dismissAlliance  
server/gameApp/alliance.go 

//公会仓库逻辑 /increaseCapacity /storeItem /destroyItem /clearup  
server/gameApp/warehouse.go 

## 公会命令测试
启动 3 个 client, 第一次输入的字符串是名字, 随后就可以输入命令进行测试  
  
client1 启动后依次输入  
player1  
/createAlliance alliance1  

client2 启动后依次输入  
player2  
/joinAlliance alliance1  

client3 启动后依次输入   
player3  
/createAlliance alliance3  

然后, 在任何一个 client 输入 /allianceList, 都会打印出 alliance1 alliance3  
client1 和 client2 输入 /whichAlliance, 会打印出 alliance1 和 player1 player2  
client3 输入 /whichAlliance, 会打印出 alliance3 和 player3  

其他命令测试省略  

## 单元测试
server_test.go 包含了服务端所有逻辑的单元测试, 使用 go test -run TestXXX -v 可以进行查看打印的逻辑信息  


## 部署
客户端和服务器最好在同一台机器上运行, 个人运行环境是 ubuntu 18  
执行 dockerfile/run.sh 可以一键部署服务端 docker  
服务端默认监听 0.0.0.0:4567  

编译好的客户端是 bin/client, 默认连接 127.0.0.1:4567   
在终端输入 bin/client -h 可以查看客户端的帮助  

如果服务端docker 部署在另外一台机器上, 假设ip 地址是 192.168.1.100:5678  
那么客户端需要带参数执行, 终端中输入 bin/client -a 192.168.1.100:5678  



## 完毕; 身体健康,谢谢阅读






