// package main

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"strings"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/assert/v2"
// 	_ "github.com/go-sql-driver/mysql"
// )

// // 仮のハンドラ関数
// func yourPostHandlerFunction(ctx *gin.Context) {
// 	name := ctx.PostForm("name")
// 	email := ctx.PostForm("email")
// 	if name == "" || email == "" {
// 		ctx.String(http.StatusBadRequest, "name or email is missing")
// 		return
// 	}
// 	// 以下、実際の処理（データベースへの保存など）
// 	ctx.String(http.StatusOK, "success")
// }

// func TestNameAndEmailPost(t *testing.T) {
// 	// GinとHTTPテストレコーダーの設定
// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	router.POST("/new", yourPostHandlerFunction)

// 	// 作成テスト用のHTTPリクエストの作成
// 	form := url.Values{}
// 	form.Add("name", "John")
// 	form.Add("email", "john@example.com")
// 	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
// 	if err != nil {
// 		t.Fatalf("Failed to make POST request: %v", err)
// 	}
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	// HTTPレスポンスの記録
// 	resp := httptest.NewRecorder()
// 	router.ServeHTTP(resp, req)

// 	// HTTPステータスコードの確認
// 	assert.Equal(t, http.StatusOK, resp.Code)

// 	// レスポンスボディの確認（オプション）
// 	assert.Equal(t, "success", resp.Body.String())
// }

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Database インターフェース
type Database interface {
	CreateUser(name string, email string) error
}

// MockDatabase の定義
type MockDatabase struct{}

func (db *MockDatabase) CreateUser(name string, email string) error {
	return nil
}

// ハンドラ関数
func yourPostHandlerFunction(ctx *gin.Context, db Database) {
	name := ctx.PostForm("name")
	email := ctx.PostForm("email")
	if err := db.CreateUser(name, email); err != nil {
		ctx.String(http.StatusInternalServerError, "Could not create user")
		return
	}
	ctx.String(http.StatusOK, "success")
}

// テストコード
func TestNameAndEmailPostWithMock(t *testing.T) {
	// GinとHTTPテストレコーダーの設定
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// MockDatabase のインスタンスを作成
	mockDB := new(MockDatabase)

	// ルーティングとハンドラ関数の設定
	router.POST("/new", func(ctx *gin.Context) {
		yourPostHandlerFunction(ctx, mockDB)
	})

	// HTTPリクエストの作成
	form := url.Values{}
	form.Add("name", "John")
	form.Add("email", "john@example.com")
	req, err := http.NewRequest("POST", "/new", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// HTTPレスポンスの記録
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// HTTPステータスコードとレスポンスボディの確認
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "success", resp.Body.String())
}
