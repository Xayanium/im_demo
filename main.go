package main

import "im_demo/server/router"

func main() {
	engine := router.Router()
	_ = engine.Run(":8080")
}
