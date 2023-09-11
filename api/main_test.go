package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

var (
	emptyPermission    = []int64{}
	accountPermissions = []int64{accountPermissionCode}
	// productPermissions = []int64{productPermissionCode}
	// orderPermissions   = []int64{orderPermissionCode}
	// couponPermissions  = []int64{couponPermissionCode}
	// newsPermissions    = []int64{newsPermissionCode}
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
	// 開始單元測試，通過 m.Run()
	// m.Run() 會回傳退出代碼，讓 os.Exit() 退出
}
