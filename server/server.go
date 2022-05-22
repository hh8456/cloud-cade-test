package main

import "cloud-cade-test/server/gameApp"

func main() {
	game := gameApp.NewGameApp()
	game.Listen("0.0.0.0:4567")
}
