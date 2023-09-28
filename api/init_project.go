package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
	"github.com/juker1141/shopping-mall-go/util"
)

var permssions = []string{"後台帳號管理", "會員管理", "商品管理", "訂單管理", "優惠卷管理", "最新消息管理"}

var orderStatuses = []db.CreateOrderStatusParams{
	{
		Name:        "待付款",
		Description: "客戶已提交付款信息，但尚未確認款項已成功付款。",
	},
	{
		Name:        "待確認付款",
		Description: "訂單已確認付款，但尚未開始處理或出貨。",
	},
	{
		Name:        "處理中",
		Description: "訂單正在處理中，商品正被包裝、準備出貨，或者正在進行其他相關的處理工作。",
	},
	{
		Name:        "已出貨",
		Description: "訂單中的商品已經出貨，並且正在運送給客戶。",
	},
	{
		Name:        "已送達",
		Description: "商品已成功送達客戶指定的送貨地址。",
	},
	{
		Name:        "已取消",
		Description: "客戶或系統管理員已取消訂單，訂單不再有效。",
	},
	{
		Name:        "退貨處理中",
		Description: "客戶申請退貨，並且退貨處理程序正在進行中。",
	},
	{
		Name:        "已退貨",
		Description: "商品已經被客戶退回，退款程序可能已經完成。",
	},
	{
		Name:        "已完成",
		Description: "訂單已經完成，包括付款、處理、出貨、送達等所有步驟。",
	},
	{
		Name:        "問題訂單",
		Description: "訂單可能存在某種問題，需要進一步調查或處理，例如庫存不足、付款問題等。",
	},
	{
		Name:        "等待回覆",
		Description: "系統或客服正在等待客戶的回覆或進一步信息，以解決某些問題。",
	},
}

var payMethods = []string{"貨到付款", "信用卡", "銀行轉帳"}

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

	for _, orderStatus := range orderStatuses {
		arg := db.CreateOrderStatusParams{
			Name:        orderStatus.Name,
			Description: orderStatus.Description,
		}

		_, err := server.store.CreateOrderStatus(context.Background(), arg)
		if err != nil {
			return err
		}
	}

	for _, payMethod := range payMethods {
		_, err := server.store.CreatePayMethod(context.Background(), payMethod)
		if err != nil {
			return err
		}
	}

	return nil
}
