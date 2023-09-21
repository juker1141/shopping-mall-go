package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
)

var permssions = []string{"後台帳號管理", "會員管理", "商品管理", "訂單管理", "優惠卷管理", "最新消息管理"}

func (server *Server) InitProject() error {
	adminUser, err := server.store.GetAdminUserByAccount(context.Background(), server.config.TestAccount)
	if err == nil && adminUser.Account == server.config.TestAccount {
		return nil
	}

	var permissions_id []int64
	for _, permission := range permssions {
		permission, err := server.store.CreatePermission(context.Background(), permission)
		if err != nil {
			return err
		}
		permissions_id = append(permissions_id, permission.ID)
	}

	roleArg := db.CreateRoleTxParams{
		Name:          "最高管理者",
		PermissionsID: permissions_id,
	}

	result, err := server.store.CreateRoleTx(context.Background(), roleArg)
	if err != nil {
		return err
	}

	hashedPassword, err := util.HashPassword(server.config.TestPassword)
	if err != nil {
		return err
	}

	userArg := db.CreateAdminUserParams{
		Account:        server.config.TestAccount,
		FullName:       "測試管理者",
		Status:         1,
		HashedPassword: hashedPassword,
		RoleID: pgtype.Int4{
			Int32: int32(result.Role.ID),
			Valid: true,
		},
	}

	_, err = server.store.CreateAdminUser(context.Background(), userArg)
	if err != nil {
		return err
	}

	return nil
}
