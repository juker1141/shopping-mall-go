package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
	// 開始單元測試，通過 m.Run()
	// m.Run() 會回傳退出代碼，讓 os.Exit() 退出
}
