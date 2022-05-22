package gameApp

// 查询 playerName 是不是工会创建者, 返回值包含了所创建的工会名
func (g *GameApp) isAllianceCreatorNonLock(playerName string) (string, bool) {
	// 查找玩家所属工会
	if allianceName, ok := g.mapClientAlliance[playerName]; ok {
		// 查找工会创建者名字
		if creator, find := g.mapAlliance[allianceName]; find {
			if creator == playerName {
				return allianceName, true
			}
		}
	}

	return "", false
}

func (g *GameApp) WhichAlliance(name string, param ...string) []byte {
	g.lock.RLock()
	defer g.lock.RUnlock()
	msg := []byte("alliance name: ")
	allianceName := g.mapClientAlliance[name]
	msg = append(msg, []byte(allianceName)...)
	msg = append(msg, []byte("\nalliance members: ")...)
	if members, ok := g.mapAllianceMember[allianceName]; ok {
		for _, v := range members {
			v += " "
			msg = append(msg, []byte(v)...)
		}
	}
	return msg
}

func (g *GameApp) CreateAlliance(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	allianceName := param[0]
	if _, ok := g.mapAlliance[allianceName]; ok {
		msg := []byte("alliance name is already use, please input a new alliance name ")
		return msg
	}

	if g.mapClientAlliance[name] == "" {
		g.mapAlliance[allianceName] = name
		g.mapClientAlliance[name] = allianceName
		g.mapAllianceMember[allianceName] = append(g.mapAllianceMember[allianceName], name)
		g.initWarehouseNonLock(allianceName)
		msg := []byte("create alliance sucess, alliance name is: ")
		msg = append(msg, []byte(allianceName)...)
		return msg
	} else {
		msg := []byte("you already have alliance, can not create alliance again ")
		return msg
	}
}

func (g *GameApp) AllianceList(name string, param ...string) []byte {
	g.lock.RLock()
	defer g.lock.RUnlock()
	var alliancsList string
	for k, _ := range g.mapAlliance {
		alliancsList += k
		alliancsList += " "
	}

	return []byte(alliancsList)
}

func (g *GameApp) JoinAlliance(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	allianceName := param[0]
	if g.mapClientAlliance[name] == "" {
		if _, ok := g.mapAlliance[allianceName]; ok {
			g.mapClientAlliance[name] = allianceName
			g.mapAllianceMember[allianceName] = append(g.mapAllianceMember[allianceName], name)
			msg := []byte("join the alliance success, alliance name: ")
			msg = append(msg, []byte(allianceName)...)
			return msg
		} else {
			// 工会不存在
			msg := []byte("error, the alliance is not exist")
			return msg
		}
	} else {
		// 如果已经加入工会了
		msg := []byte("error, you already have alliance, can not join another alliance ")
		return msg
	}
}

func (g *GameApp) DismissAlliance(name string, param ...string) []byte {
	g.lock.Lock()
	defer g.lock.Unlock()
	if allianceName, find := g.mapClientAlliance[name]; find {
		if creator, ok := g.mapAlliance[allianceName]; ok {
			if creator == name {
				delete(g.mapAlliance, allianceName)

				members := g.mapAllianceMember[allianceName]
				for _, memberName := range members {
					delete(g.mapClientAlliance, memberName)
				}

				delete(g.mapAllianceMember, allianceName)
				delete(g.mapClientAlliance, allianceName)
				delete(g.mapAllianceItem, allianceName)

				// 解散行会成功
				msg := []byte("dismiss alliance success")
				return msg
			}

			// 需要会长才能解散工会
			msg := []byte("only the creator can dismiss alliance ")
			return msg

		}
	}
	// 行会不存在
	msg := []byte("the alliance is not exist")
	return msg
}
