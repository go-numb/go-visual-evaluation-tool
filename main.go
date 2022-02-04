package main

import (
	"go-visual-evaluation-tool/modules"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	fileroot = "./data/images"
)

func main() {
	// ファイルディレクトリからファイル群を読み込み必要な情報をリストに記載して保存する
	// リストパスは ./data/list.csv
	modules.CreateCSV(fileroot)

	// リストを読み込み表示
	// 評価後リスト更新
	// 次のファイルを表示する

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/data", "data")

	e.GET("/", func(c echo.Context) error {
		b, err := ioutil.ReadFile("./statics/sample.html")
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, string(b))
	})
	e.POST("/receive", modules.Receive)
	e.Logger.Fatal(e.Start(":1323"))
}
