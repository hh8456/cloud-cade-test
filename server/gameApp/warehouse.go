package gameApp

import (
	"cloud-cade-test/server/gameApp/pb"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"

	"github.com/golang/protobuf/proto"
)

func (g *GameApp) initWarehouseNonLock(allianceName string) {
	warehouse := make([]*pb.TestItem, 30)
	// 刻意每次从文件读取,便于支持仓库初始化的热更新
	f, err := ioutil.ReadFile("testItem.data")
	if err != nil {
		fmt.Printf("read testItem.data fail, %v\n", err)
		g.mapAllianceItem[allianceName] = warehouse
		return
	}

	itemInfoArr := &pb.TestItem_Array{}
	proto.Unmarshal(f, itemInfoArr)
	for _, v := range itemInfoArr.GetItems() {
		g.mapItemInfo[v.GetId()] = v
	}

	itemArr := &pb.TestItem_Array{}
	proto.Unmarshal(f, itemArr)

	items := itemArr.GetItems()
	i := 0
	for _, v := range items {
		if v.GetNumber() <= 5 {
			warehouse[i] = v
			i++
		} else {
			for {
				if v.GetNumber() > 0 {
					if v.GetNumber() <= 5 {
						p := proto.Clone(v).(*pb.TestItem)
						v.Number = 0
						warehouse[i] = p
					} else {
						p := proto.Clone(v).(*pb.TestItem)
						p.Number = 5
						v.Number -= 5
						warehouse[i] = p
					}
					i++
				} else {
					break
				}
			}
		}
	}

	g.mapAllianceItem[allianceName] = warehouse
}

func (g *GameApp) AllianceCapacity(allianceName string) int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return len(g.mapAllianceItem[allianceName])
}

// 工会仓库使用了多少个格子
func (g *GameApp) WarehouseUsed(allianceName string) int {
	g.lock.RLock()
	defer g.lock.RUnlock()
	items := g.mapAllianceItem[allianceName]
	i := 0
	for _, v := range items {
		if v != nil {
			i++
		}
	}

	return i
}

// 查询仓库格子上的道具信息
func (g *GameApp) WarehouseItem(allianceName string, index int) *pb.TestItem {
	g.lock.RLock()
	defer g.lock.RUnlock()
	items := g.mapAllianceItem[allianceName]
	if items != nil {
		return items[index]
	}

	return nil
}

func (g *GameApp) Warehouse(allianceName string) []*pb.TestItem {
	g.lock.RLock()
	defer g.lock.RUnlock()
	return g.mapAllianceItem[allianceName]
}

func (g *GameApp) IncreaseCapacity(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	// 如果玩家是工会创建者,才允许扩容
	if allianceName, ok := g.isAllianceCreatorNonLock(name); ok {
		newWarehouse := make([]*pb.TestItem, 40)
		items := g.mapAllianceItem[allianceName]
		for i, v := range items {
			newWarehouse[i] = v
		}

		g.mapAllianceItem[allianceName] = newWarehouse
		msg := []byte("increaseCapacity success")
		return msg
	} else {
		msg := []byte("increaseCapacity fail")
		return msg
	}
}

// 获得玩家所在行会的仓库
func (g *GameApp) getWarehoueNonLock(name string) []*pb.TestItem {
	if allianceName, ok := g.mapClientAlliance[name]; ok {
		return g.mapAllianceItem[allianceName]
	}

	return nil
}

func (g *GameApp) StoreItem(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	items := g.getWarehoueNonLock(name)
	itemid, err1 := strconv.ParseInt(param[0], 10, 32)
	itemNum, err2 := strconv.ParseInt(param[1], 10, 32)
	index, err3 := strconv.ParseInt(param[2], 10, 32)
	if err1 != nil || err2 != nil || err3 != nil {
		msg := []byte("storeItem fail")
		return msg

	}

	// index 错误
	if index < 0 || index >= int64(len(items)) {
		msg := []byte("storeItem fail")
		return msg
	}

	tip, ok := g.appendItemNonLock(int32(itemid), int32(itemNum), int32(index), items)
	if ok {
		msg := []byte("storeItem success")
		return msg
	} else {
		msg := []byte(tip)
		return msg
	}
}

// 提交物品到仓库; 要假设传递的参数都是正确的
func (g *GameApp) appendItemNonLock(itemId, itemNum, index int32, items []*pb.TestItem) (string, bool) {
	itemInfo := g.mapItemInfo[itemId]

	// 仓库还剩多少可用空间来保存 itemNum 个 itemId
	leftNumber := 0

	// 每个格子能堆叠 5 个同类道具
	for _, item := range items {
		if item == nil {
			leftNumber += 5
		} else {
			if item.GetItemType() == itemInfo.GetItemType() {
				n := 5 - int(item.GetNumber())
				leftNumber += n
			}
		}
	}

	// 超过了仓库容量
	if leftNumber < int(itemNum) {
		str := fmt.Sprintf("The warehouse capacity is exceeded, leftNumber :%d", leftNumber)
		return str, false
	}

	left := g.appendItemAtRangeNonLock(itemId, itemNum, index, int32(len(items))-1, items)

	// 如果还剩得有, 就放到前面的格子中
	if left > 0 {
		g.appendItemAtRangeNonLock(itemId, left, 0, index-1, items)
	}

	return "", true
}

// 在仓库区间 [startIndex, endIndex] 保存提交的物品; 返回值是还剩余多少个 itemId 没有被保存
func (g *GameApp) appendItemAtRangeNonLock(itemId, itemNum, startIndex, endIndex int32, items []*pb.TestItem) int32 {
	itemInfo := g.mapItemInfo[itemId]
	for i := startIndex; i <= endIndex; i++ {
		if items[i] == nil {
			item := &pb.TestItem{}
			item.Id = itemInfo.GetId()
			item.Name = itemInfo.GetName()
			item.ItemType = itemInfo.GetItemType()
			if itemNum <= 5 {
				item.Number = itemNum
				itemNum = 0
				items[i] = item
				return 0
			} else {
				item.Number = 5
				itemNum -= 5
				items[i] = item
			}
		} else {
			if items[i].GetItemType() == itemInfo.GetItemType() && items[i].GetNumber() < 5 {
				left := 5 - items[i].GetNumber()
				items[i].Number = 5
				itemNum -= left
				if itemNum == 0 {
					return 0
				}
			}
		}
	}

	return itemNum
}

func (g *GameApp) DestroyItem(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	if allianceName, find := g.mapClientAlliance[name]; find {
		if creator, ok := g.mapAlliance[allianceName]; ok {
			// 会长才能删除道具
			if creator == name {
				items := g.getWarehoueNonLock(name)
				index, err := strconv.ParseInt(param[0], 10, 32)
				// index 错误
				if err != nil || index < 0 || index >= int64(len(items)) {
					msg := []byte("storeItem fail")
					return msg
				}
				if items[index] != nil {
					items[index] = nil
					msg := []byte("storeItem success")
					return msg
				}

				msg := []byte("storeItem fail, there is no item")
				return msg
			} else {
				msg := []byte("storeItem fail, creator can DestroyItem")
				return msg
			}
		}
	}

	msg := []byte("storeItem fail")
	return msg
}

// 每种类型的道具,有多少个; 按 type 升序排列
func (g *GameApp) getItemTypeNum(name string) []*pb.TestItem {
	// 每种类型的道具,有多少个
	mapItemNumByType := map[int32]*pb.TestItem{}
	items := g.getWarehoueNonLock(name)
	for _, item := range items {
		if item != nil {
			if pItem, ok := mapItemNumByType[item.GetItemType()]; ok {
				pItem.Number += item.GetNumber()
			} else {
				p := proto.Clone(item).(*pb.TestItem)
				mapItemNumByType[item.GetItemType()] = p
			}
		}
	}

	var ss []*pb.TestItem
	for _, v := range mapItemNumByType {
		ss = append(ss, v)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].ItemType < ss[j].ItemType // 升序排列
	})

	return ss
}

func (g *GameApp) Clearup(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	if allianceName, ok := g.isAllianceCreatorNonLock(name); ok {
		// 得到按 type 升序排列的仓库道具集合
		storedItems := g.getItemTypeNum(name)

		newItems := make([]*pb.TestItem, len(g.mapAllianceItem[allianceName]))
		i := 0
		for _, item := range storedItems {
			for {
				if item.GetNumber() > 0 {
					if item.GetNumber() <= 5 {
						p := proto.Clone(item).(*pb.TestItem)
						newItems[i] = p
						i++
						item.Number = 0
					} else {
						p := proto.Clone(item).(*pb.TestItem)
						p.Number = 5
						item.Number -= 5
						newItems[i] = p
						i++
					}
				} else {
					break
				}
			}
		}

		g.mapAllianceItem[allianceName] = newItems
	}

	msg := []byte("clearup over")
	return msg
}
