package main

import (
	"fmt"
	"net/http"

	// strconv パッケージは、基本データ型の文字列表現との変換を実装
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func main() {
	// データベースへの接続を開始
	db := sqlConnect()
	// Userモデルのマイグレーションを行う。これにより、データベースに User テーブルが存在しない場合、新たに作成されます。
	db.AutoMigrate(&User{})
	// 関数の終了時にデータベースの接続を閉じるよう指示
	defer db.Close()

	// ginフレームワークのデフォルトのルータを初期化
	router := gin.Default()
	// templates ディレクトリのHTMLテンプレートを全て読み込む
	router.LoadHTMLGlob("templates/*.html")

	// ユーザーの一覧を表示するルートです。データベースからユーザーの情報を取得して、index.html というテンプレートに渡しています。
	router.GET("/", func(ctx *gin.Context) {
		db := sqlConnect()
		var users []User
		db.Order("created_at asc").Find(&users)
		defer db.Close()

		ctx.HTML(200, "index.html", gin.H{
			"users": users,
		})
	})
	// 'ctx'はGinのオブジェクト。HTTPリクエストとレスポンスの橋渡し役
	router.POST("/new", func(ctx *gin.Context) {
		db := sqlConnect()
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")
		fmt.Println("create user " + name + " with email " + email)
		db.Create(&User{Name: name, Email: email})
		defer db.Close()

		// 保存後、ルートページにリダイレクト
		ctx.Redirect(302, "/")
	})

	// IDを取得して、そのIDのユーザーをデータベースから削除
	router.POST("/delete/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("id is not a number")
		}
		var user User
		db.First(&user, id)
		db.Delete(&user)
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	router.GET("/user/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		// URLの :id パートに該当する部分を取得
		n := ctx.Param("id")
		//  取得したID(変数n)を変換して'id'と'err'変数に代入しています。 ’Atoi’ は文字列を整数に変換する関数です。
		id, err := strconv.Atoi(n)
		if err != nil {
			ctx.String(http.StatusBadRequest, "IDが不正です")
			return
		}
		var user User
		db.First(&user, id)
		if user.ID == 0 {
			ctx.String(http.StatusNotFound, "ユーザが見つかりません")
			return
		}
		defer db.Close()
		ctx.HTML(http.StatusOK, "user_show.html", gin.H{
			"user": user,
		})
	})

	// Webサーバーを起動
	router.Run()
}

func sqlConnect() (database *gorm.DB) {
	DBMS := "mysql"
	USER := "go_test"
	PASS := "password"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "go_database"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

	count := 0
	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		for {
			if err == nil {
				fmt.Println("")
				break
			}
			fmt.Print(".")
			time.Sleep(time.Second)
			count++
			if count > 5 {
				fmt.Println("")
				fmt.Println("DB接続失敗")
				panic(err)
			}
			// DBに接続する
			db, err = gorm.Open(DBMS, CONNECT)
		}
	}
	// fmt.Println("DB接続成功")

	return db
}
