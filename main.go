package main

import (
	"fmt"
	"net/http"

	// strconv パッケージは、基本データ型の文字列表現との変換を実装
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// 構造体のフィールドに対するバリデーションを行うためのパッケージ
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string `validate:"required"`       // Name は必須項目
	Email string `validate:"required,email"` // Email は必須項目で、メールアドレス形式
	Age   int    `validate:"gte=0,lte=130"`  // Age は0以上130以下
}

var validate *validator.Validate

func main() {
	// データベースへの接続を開始
	db := sqlConnect()
	// Userモデルのマイグレーションを行う。これにより、データベースに User テーブルが存在しない場合、新たに作成されます。
	db.AutoMigrate(&User{})
	// 関数の終了時にデータベースの接続を閉じるよう指示
	defer db.Close()

	// バリデーションの初期化
	validate = validator.New()
	// ginフレームワークのデフォルトのルータを初期化
	router := gin.Default()
	// templates ディレクトリのHTMLテンプレートを全て読み込む
	router.LoadHTMLGlob("templates/*.html")

	// ユーザーの一覧を表示するルートです。データベースからユーザーの情報を取得して、index.html というテンプレートに渡しています。
	router.GET("/", func(ctx *gin.Context) {

		// Custom-Headerを追加
		ctx.Header("Custom-Header", "some-value")

		db := sqlConnect()
		var users []User
		db.Order("created_at asc").Find(&users)
		defer db.Close()

		ctx.HTML(200, "index.html", gin.H{
			"header": "HelloWorld", // ヘッダーとして表示するテキスト
			"users":  users,
		})
	})

	// 'ctx'はGinのオブジェクト。HTTPリクエストとレスポンスの橋渡し役
	router.POST("/new", func(ctx *gin.Context) {
		db := sqlConnect()
		name := ctx.PostForm("name")
		email := ctx.PostForm("email")
		age, _ := strconv.Atoi(ctx.PostForm("age"))
		user := &User{Name: name, Email: email, Age: age}
		err := validate.Struct(user)
		if err != nil {
			// バリデーションエラーの処理
			ctx.String(http.StatusBadRequest, "入力が不正です")
			return
		}
		db.Create(user)
		defer db.Close()
		ctx.Redirect(302, "/")
	})

	// IDを取得して、そのIDのユーザーをデータベースから削除
	router.POST("/delete/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("IDがありません")
		}
		var user User
		db.Find(&user, id)
		db.Delete(&user)
		defer db.Close()

		ctx.Redirect(302, "/")
	})

	// IDを取得して、そのIDのユーザー詳細を表示
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
		db.Find(&user, id)
		if user.ID == 0 {
			ctx.String(http.StatusNotFound, "ユーザが見つかりません")
			return
		}
		defer db.Close()
		ctx.HTML(http.StatusOK, "user_show.html", gin.H{
			"user": user,
		})
	})

	// IDを取得して、そのIDのユーザー編集画面の表示を行う
	router.GET("/user/edit/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
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
		ctx.HTML(http.StatusOK, "user_edit.html", gin.H{
			"user": user,
		})
	})

	// IDを取得して、そのIDのユーザーを更新
	router.POST("/user/update/:id", func(ctx *gin.Context) {
		db := sqlConnect()
		n := ctx.Param("id")
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
		user.Name = ctx.PostForm("name")
		user.Email = ctx.PostForm("email")
		db.Save(&user)
		defer db.Close()
		ctx.Redirect(http.StatusSeeOther, "/user/"+n)
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
