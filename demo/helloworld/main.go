package main

import (
	"fmt"
	"github.com/PPMac/gocn"
)

func main() {
	engine := gocn.New()
	g := engine.Group("user")
	g.Get("/", func(ctx *gocn.Context) {
		fmt.Fprintf(ctx.W, "%s hello world", "gocn")
	})
	g.Get("/list", func(ctx *gocn.Context) {
		fmt.Fprintf(ctx.W, "users：%s, %s", "LiLei", "HanMeimei")
	})
	g.Post("/add", func(ctx *gocn.Context) {
		fmt.Fprintf(ctx.W, "%s: is a %s", "LiLei", "boy")
	})
	engine.Run("8080")
}
