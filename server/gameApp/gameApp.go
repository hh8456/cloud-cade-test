package gameApp

import (
	"cloud-cade-test/base-library/base_net"
	"fmt"
	"net"
	"sync"
	"time"

	"cloud-cade-test/server/gameApp/client"
	"cloud-cade-test/server/gameApp/pb"
)

type GameApp struct {
	*base_net.Connection
	lock              sync.RWMutex
	mapName           map[string]struct{}       //注册的玩家名字集合
	mapAlliance       map[string]string         // 工会 - 工会会长
	mapAllianceMember map[string][]string       // 工会 - 工会成员
	mapClientAlliance map[string]string         // 玩家 - 工会
	mapAllianceItem   map[string][]*pb.TestItem // 工会名 - 工会仓库
	//mapLogicFunc      map[string]func(*client.Client, ...string)
	mapLogicFunc map[string]func(string, ...string) []byte
	mapItemInfo  map[int32]*pb.TestItem // 道具 id - 道具信息
}

func NewGameApp() *GameApp {
	g := &GameApp{
		mapLogicFunc:      map[string]func(string, ...string) []byte{},
		mapClientAlliance: map[string]string{},
		mapName:           map[string]struct{}{},
		mapAlliance:       map[string]string{},
		mapAllianceMember: map[string][]string{},
		mapAllianceItem:   map[string][]*pb.TestItem{},
		mapItemInfo:       map[int32]*pb.TestItem{}, // 道具 id - 道具信息
	}

	g.mapLogicFunc["/whichAlliance"] = g.WhichAlliance
	g.mapLogicFunc["/createAlliance"] = g.CreateAlliance
	g.mapLogicFunc["/allianceList"] = g.AllianceList
	g.mapLogicFunc["/joinAlliance"] = g.JoinAlliance
	g.mapLogicFunc["/dismissAlliance"] = g.DismissAlliance

	g.mapLogicFunc["/increaseCapacity"] = g.IncreaseCapacity
	g.mapLogicFunc["/storeItem"] = g.StoreItem
	g.mapLogicFunc["/destroyItem"] = g.DestroyItem
	g.mapLogicFunc["/clearup"] = g.Clearup

	return g
}

func (g *GameApp) ClientLogin(name string) bool {
	g.lock.Lock()
	defer g.lock.Unlock()
	_, ok := g.mapName[name]
	if ok {
		return false
	}
	g.mapName[name] = struct{}{}
	return true
}

func (g *GameApp) Listen(addr string) {
	onEstablish := func(conn net.Conn) {
		c := base_net.CreateSocket(conn, 4096) // 设定只能处理 4k 以内的逻辑包

		g.lock.RLock() // 加锁保护 g.mapLogicFunc
		cli := client.NewClient(c, g.ClientLogin, g.mapLogicFunc)
		g.lock.RUnlock()

		cli.Run()
		fmt.Printf("server: create a new client tcp connection, serialNumber: %d\n", cli.ConnId())
	}

	l := base_net.CreateListener()

	for {
		err := l.Start(addr, onEstablish)
		if err != nil {
			panic(fmt.Sprintf("listen tcp addr: %s, error:%v", addr, err))
		}

		time.Sleep(time.Second)
	}
}
