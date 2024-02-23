package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUserTx(t *testing.T) UserTxResult {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserTxParams{
		CreateUserParams: CreateUserParams{
			Account:  util.RandomAccount(),
			Email:    util.RandomEmail(),
			FullName: util.RandomName(),
			GenderID: pgtype.Int4{
				Int32: util.RandomGender(),
				Valid: true,
			},
			Cellphone:       util.RandomCellPhone(),
			Address:         util.RandomAddress(),
			ShippingAddress: util.RandomAddress(),
			PostCode:        util.RandomPostCode(),
			HashedPassword:  hashedPassword,
			Status:          1,
			AvatarUrl:       util.RandomString(20),
		},
		AfterCreate: func(user User) error {
			return nil
		},
	}

	txResult, err := testStore.CreateUserTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, txResult)

	require.Equal(t, arg.CreateUserParams.Account, txResult.User.Account)
	require.Equal(t, arg.CreateUserParams.Email, txResult.User.Email)
	require.Equal(t, arg.CreateUserParams.FullName, txResult.User.FullName)
	require.Equal(t, arg.CreateUserParams.GenderID.Int32, txResult.User.GenderID.Int32)
	require.Equal(t, arg.CreateUserParams.Cellphone, txResult.User.Cellphone)
	require.Equal(t, arg.CreateUserParams.Address, txResult.User.Address)
	require.Equal(t, arg.CreateUserParams.ShippingAddress, txResult.User.ShippingAddress)
	require.Equal(t, arg.CreateUserParams.PostCode, txResult.User.PostCode)
	require.Equal(t, arg.CreateUserParams.HashedPassword, txResult.User.HashedPassword)
	require.Equal(t, arg.CreateUserParams.Status, txResult.User.Status)
	require.Equal(t, arg.CreateUserParams.AvatarUrl, txResult.User.AvatarUrl)

	return txResult
}

func TestCreateUserTx(t *testing.T) {
	createRandomUserTx(t)
}
