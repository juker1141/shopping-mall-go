package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(8))
	require.NoError(t, err)

	arg := CreateUserParams{
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
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Account, user.Account)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.GenderID, user.GenderID)
	require.Equal(t, arg.Cellphone, user.Cellphone)
	require.Equal(t, arg.Address, user.Address)
	require.Equal(t, arg.ShippingAddress, user.ShippingAddress)
	require.Equal(t, arg.PostCode, user.PostCode)
	require.Equal(t, arg.Status, user.Status)
	require.Equal(t, arg.AvatarUrl, user.AvatarUrl)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Account, user2.Account)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.GenderID, user2.GenderID)
	require.Equal(t, user1.Cellphone, user2.Cellphone)
	require.Equal(t, user1.Address, user2.Address)
	require.Equal(t, user1.ShippingAddress, user2.ShippingAddress)
	require.Equal(t, user1.PostCode, user2.PostCode)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.AvatarUrl, user2.AvatarUrl)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestGetUserByAccount(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testStore.GetUserByAccount(context.Background(), user1.Account)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Account, user2.Account)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.GenderID, user2.GenderID)
	require.Equal(t, user1.Cellphone, user2.Cellphone)
	require.Equal(t, user1.Address, user2.Address)
	require.Equal(t, user1.ShippingAddress, user2.ShippingAddress)
	require.Equal(t, user1.PostCode, user2.PostCode)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.AvatarUrl, user2.AvatarUrl)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserAllField(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		ID: oldUser.ID,
		FullName: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
		Cellphone: pgtype.Text{
			String: util.RandomCellPhone(),
			Valid:  true,
		},
		Address: pgtype.Text{
			String: util.RandomAddress(),
			Valid:  true,
		},
		ShippingAddress: pgtype.Text{
			String: util.RandomAddress(),
			Valid:  true,
		},
		PostCode: pgtype.Text{
			String: util.RandomPostCode(),
			Valid:  true,
		},
		AvatarUrl: pgtype.Text{
			String: util.RandomString(20),
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		Status: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
	}

	newUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.ID, newUser.ID)
	require.Equal(t, oldUser.Account, newUser.Account)
	require.Equal(t, oldUser.GenderID, newUser.GenderID)

	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.NotEqual(t, oldUser.Cellphone, newUser.Cellphone)
	require.NotEqual(t, oldUser.Address, newUser.Address)
	require.NotEqual(t, oldUser.ShippingAddress, newUser.ShippingAddress)
	require.NotEqual(t, oldUser.PostCode, newUser.PostCode)
	require.NotEqual(t, oldUser.Status, newUser.Status)
	require.NotEqual(t, oldUser.AvatarUrl, newUser.AvatarUrl)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)

	require.False(t, newUser.PasswordChangedAt.IsZero())
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	arg := UpdateUserParams{
		ID: oldUser.ID,
		FullName: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
	}

	newUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.ID, newUser.ID)
	require.Equal(t, oldUser.Account, newUser.Account)
	require.Equal(t, oldUser.GenderID, newUser.GenderID)
	require.Equal(t, oldUser.Cellphone, newUser.Cellphone)
	require.Equal(t, oldUser.Address, newUser.Address)
	require.Equal(t, oldUser.ShippingAddress, newUser.ShippingAddress)
	require.Equal(t, oldUser.PostCode, newUser.PostCode)
	require.Equal(t, oldUser.Status, newUser.Status)
	require.Equal(t, oldUser.AvatarUrl, newUser.AvatarUrl)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)

	require.NotEqual(t, oldUser.FullName, newUser.FullName)

	require.True(t, newUser.PasswordChangedAt.IsZero())
}

func TestUpdateUserOnlyStatus(t *testing.T) {
	oldUser := createRandomUser(t)

	arg := UpdateUserParams{
		ID: oldUser.ID,
		Status: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
	}

	newUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.ID, newUser.ID)
	require.Equal(t, oldUser.Account, newUser.Account)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.GenderID, newUser.GenderID)
	require.Equal(t, oldUser.Cellphone, newUser.Cellphone)
	require.Equal(t, oldUser.Address, newUser.Address)
	require.Equal(t, oldUser.ShippingAddress, newUser.ShippingAddress)
	require.Equal(t, oldUser.PostCode, newUser.PostCode)
	require.Equal(t, oldUser.AvatarUrl, newUser.AvatarUrl)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)

	require.NotEqual(t, oldUser.Status, newUser.Status)

	require.True(t, newUser.PasswordChangedAt.IsZero())
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		ID: oldUser.ID,
		HashedPassword: pgtype.Text{
			String: newHashedPassword,
			Valid:  true,
		},
		PasswordChangedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	newUser, err := testStore.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, oldUser.ID, newUser.ID)
	require.Equal(t, oldUser.Account, newUser.Account)
	require.Equal(t, oldUser.GenderID, newUser.GenderID)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Cellphone, newUser.Cellphone)
	require.Equal(t, oldUser.Address, newUser.Address)
	require.Equal(t, oldUser.ShippingAddress, newUser.ShippingAddress)
	require.Equal(t, oldUser.PostCode, newUser.PostCode)
	require.Equal(t, oldUser.Status, newUser.Status)
	require.Equal(t, oldUser.AvatarUrl, newUser.AvatarUrl)

	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)

	require.False(t, newUser.PasswordChangedAt.IsZero())
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)

	err := testStore.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testStore.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testStore.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}

func TestGetUsersCount(t *testing.T) {
	createRandomUser(t)

	count, err := testStore.GetUsersCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)
}
