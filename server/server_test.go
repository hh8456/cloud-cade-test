package main

import (
	"cloud-cade-test/server/gameApp"
	"testing"
)

// 模拟客户端在控制台输入用户名登录
func TestLogin(t *testing.T) {
	game := gameApp.NewGameApp()

	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	if false == game.ClientLogin("player2") {
		t.Fatal("fail")
	}

	if false == game.ClientLogin("player3") {
		t.Fatal("fail")
	}

	if game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	if game.ClientLogin("player2") {
		t.Fatal("fail")
	}

	if game.ClientLogin("player3") {
		t.Fatal("fail")
	}
}

// 工会创建, 加入, 解散
func TestAlliance(t *testing.T) {
	game := gameApp.NewGameApp()

	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player2 登录
	if false == game.ClientLogin("player2") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	t.Log("/createAlliance player1:")
	t.Logf("%s\n\n", game.CreateAlliance("player1", "alliance1"))

	// player2 输入 /joinAlliance 加入工会 alliance1
	game.JoinAlliance("player2", "alliance1")

	// player1 输入 /whichAlliance
	t.Log("/whichAlliance player1:")
	t.Logf("%s\n\n", game.WhichAlliance("player1"))

	// player2 输入 /whichAlliance
	t.Log("/whichAlliance player2:")
	t.Logf("%s\n\n", game.WhichAlliance("player2"))

	// 任何玩家输入 /allianceList 查看所有工会
	t.Logf("/allianceList: %s\n\n", game.AllianceList(""))

	// player3 登录
	if false == game.ClientLogin("player3") {
		t.Fatal("fail")
	}

	// player3 输入 /createAlliance 创建工会 alliance3
	game.CreateAlliance("player3", "alliance3")
	t.Logf("/allianceList: %s\n\n", game.AllianceList(""))

	// player3 输入 /dismissAlliance 解散工会
	game.DismissAlliance("player3", "alliance3")
	// 任何玩家输入 /allianceList 查看所有工会
	t.Logf("/allianceList: %s\n\n", game.AllianceList(""))
}

// 工会仓库扩容
func TestIncreaseCapacity(t *testing.T) {
	game := gameApp.NewGameApp()

	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	game.CreateAlliance("player1", "alliance1")

	t.Logf("alliance1 capacity: %d\n\n", game.AllianceCapacity("alliance1"))

	game.IncreaseCapacity("player1")
	t.Logf("alliance1 capacity( after increate ): %d\n\n", game.AllianceCapacity("alliance1"))
}

// 打印工会仓库初始化
func TestInitWarehouse(t *testing.T) {
	game := gameApp.NewGameApp()
	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	// 默认有 12 个道具1, 7个道具2, 3个道具3, 5个道具4, 1个道具5; 由于每个格子只能堆叠 5 个同类道具
	// 那么初始化仓库时:
	// 0, 1, 2 下标的格子保存道具1,数量分别是 5,5,2
	// 3, 4 下标的格子保存道具2,数量分别是 5,2
	// 5 下标的格子保存道具3,数量是 3
	// 6 下标的格子保存道具4,数量是 4
	// 7 下标的格子保存道具5,数量是 1
	game.CreateAlliance("player1", "alliance1")

	// 0, 1, 2 下标的格子保存道具1,数量分别是 5,5,2
	item0 := game.WarehouseItem("alliance1", 0)
	if item0.GetNumber() != 5 || item0.GetItemType() != 1 {
		t.Fatal("fail")
	}

	item1 := game.WarehouseItem("alliance1", 1)
	if item1.GetNumber() != 5 || item1.GetItemType() != 1 {
		t.Fatal("fail")
	}

	item2 := game.WarehouseItem("alliance1", 2)
	if item2.GetNumber() != 2 || item2.GetItemType() != 1 {
		t.Fatal("fail")
	}

	// 3, 4 下标的格子保存道具2,数量分别是 5,2
	item3 := game.WarehouseItem("alliance1", 3)
	if item3.GetNumber() != 5 || item3.GetItemType() != 2 {
		t.Fatal("fail")
	}

	item4 := game.WarehouseItem("alliance1", 4)
	if item4.GetNumber() != 2 || item4.GetItemType() != 2 {
		t.Fatal("fail")
	}

	// 5 下标的格子保存道具3,数量是 3
	item5 := game.WarehouseItem("alliance1", 5)
	if item5.GetNumber() != 3 || item5.GetItemType() != 3 {
		t.Fatal("fail")
	}

	// 6 下标的格子保存道具4,数量是 4
	item6 := game.WarehouseItem("alliance1", 6)
	if item6.GetNumber() != 5 || item6.GetItemType() != 4 {
		t.Fatal("fail")
	}

	// 7 下标的格子保存道具5,数量是 1
	item7 := game.WarehouseItem("alliance1", 7)
	if item7.GetNumber() != 1 || item7.GetItemType() != 5 {
		t.Fatal("fail")
	}
}

func TestStoreItem(t *testing.T) {
	game := gameApp.NewGameApp()
	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	// 初始化仓库时:
	// 0, 1, 2 下标的格子保存道具1,数量分别是 5,5,2
	// 3, 4 下标的格子保存道具2,数量分别是 5,2
	// 5 下标的格子保存道具3,数量是 3
	// 6 下标的格子保存道具4,数量是 4
	// 7 下标的格子保存道具5,数量是 1
	game.CreateAlliance("player1", "alliance1")
	t.Logf("创建仓库后,使用的格子数: %d\n\n", game.WarehouseUsed("alliance1"))

	// 从第10个格子(下标 9 ), 保存 3 个 道具1
	game.StoreItem("player1", "0", "3", "9")
	//t.Logf("已使用的格子数: %d\n\n", game.WarehouseUsed("alliance1"))
	item9 := game.WarehouseItem("alliance1", 9)
	// 9 下标的格子保存道具1,数量是 3
	if item9.GetNumber() != 3 || item9.GetItemType() != 1 {
		t.Fatal("fail")
	}

	// 默认 30 个格子, 还有 30 - 9 = 21 个空格, 下标 2 的格子还可以保存 3 个道具 1,
	// 下标 9 的格子还可以保存 2 个道具 1, 那么仓库还可以保存 21*5 + 3 + 2 = 110 个道具 1

	// 下标 9 的格子保存 110 个道具 1
	t.Logf("%s\n\n", game.StoreItem("player1", "0", "110", "20")) // 道具1 的 itemid = 0, ItemType = 1
	t.Logf("已使用的格子数: %d\n\n", game.WarehouseUsed("alliance1"))

	items := game.Warehouse("alliance1")
	// 因为初始化仓库占用了 8 个格子, 所以从第9个格子(下标 8 )开始检查
	for i := 8; i < len(items); i++ {
		item := items[i]
		if item.GetId() != 0 || item.GetNumber() != 5 {
			//t.Logf("index: %d, id: %d, number: %d\n\n", i, item.GetId(), item.GetNumber())
			t.Fatal("fail")
		}
	}
	//下标 2 的格子可以堆叠满 5 个
	item := items[2]
	if item.GetId() != 0 || item.GetNumber() != 5 {
		t.Fatal("fail")
	}
}

func TestDestroyItem(t *testing.T) {
	game := gameApp.NewGameApp()
	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	game.CreateAlliance("player1", "alliance1")

	// 初始化仓库时:
	// 0, 1, 2 下标的格子保存道具1,数量分别是 5,5,2
	// 3, 4 下标的格子保存道具2,数量分别是 5,2
	// 5 下标的格子保存道具3,数量是 3
	// 6 下标的格子保存道具4,数量是 4
	// 7 下标的格子保存道具5,数量是 1

	// player1 作为工会会长, 删除工会仓库第3个格子(下标 2)的道具
	game.DestroyItem("player1", "2")
	item := game.WarehouseItem("alliance1", 2)
	if item != nil {
		t.Fatal("fail")
	}
}

func TestClearup(t *testing.T) {
	game := gameApp.NewGameApp()
	// player1 登录
	if false == game.ClientLogin("player1") {
		t.Fatal("fail")
	}

	// player1 输入 /createAlliance 创建工会, 工会名 alliance1
	// 初始化仓库时:
	// 0, 1, 2 下标的格子保存道具1,数量分别是 5,5,2
	// 3, 4 下标的格子保存道具2,数量分别是 5,2
	// 5 下标的格子保存道具3,数量是 3
	// 6 下标的格子保存道具4,数量是 4
	// 7 下标的格子保存道具5,数量是 1
	game.CreateAlliance("player1", "alliance1")

	// 9, 19, 23, 25 格子添加道具 1, 总共占用 5 格
	// 9 下标格子, 保存 4 个 道具5( 第二个参数表示 itemid = 4)
	game.StoreItem("player1", "4", "4", "9")

	// 19 下标格子, 保存 3 个 道具1(第二个参数表示 itemid = 0)
	game.StoreItem("player1", "0", "3", "19")

	// 2 下标格子, 保存 6 个 道具1(第二个参数表示 itemid = 0)
	game.StoreItem("player1", "0", "6", "2")

	// 25 下标格子, 保存 4 个 道具1(第二个参数表示 itemid = 0)
	game.StoreItem("player1", "0", "4", "25")

	// 26, 27, 28 格子添加道具, 总共占用 3 格
	// 26 下标格子, 保存 3 个 道具2(第二个参数表示 itemid = 1)
	game.StoreItem("player1", "1", "3", "26")

	// 27 下标格子, 保存 2 个 道具3(第二个参数表示 itemid = 2)
	game.StoreItem("player1", "2", "2", "27")

	// 28 下标格子, 保存 1 个 道具4(第二个参数表示 itemid = 3)
	game.StoreItem("player1", "3", "1", "28")

	//t.Logf("已使用的格子数: %d\n\n", game.WarehouseUsed("alliance1"))

	t.Logf("整理仓库前的道具存放情况:\n")

	for i := 0; i < 30; i++ {
		item := game.WarehouseItem("alliance1", i)
		if item != nil {
			t.Logf("index: %d, itemType: %d, itemNumber: %d\n", i, item.GetItemType(), item.GetNumber())
		}
	}

	game.Clearup("player1")

	//t.Logf("已使用的格子数: %d\n\n", game.WarehouseUsed("alliance1"))

	t.Logf("\n\n整理仓库后的道具存放情况:\n")
	// 打印仓库中的道具 itemType 和 itemNumber
	for i := 0; i < 30; i++ {
		item := game.WarehouseItem("alliance1", i)
		if item != nil {
			t.Logf("index: %d, itemType: %d, itemNumber: %d\n", i, item.GetItemType(), item.GetNumber())
		}
	}
}
