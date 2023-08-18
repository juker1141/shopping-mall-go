package db

import (
	"context"
	"math"
	"testing"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAdminUserTx(t *testing.T) AdminUserTxResult {
	account := util.RandomAccount()
	fullName := util.RandomName()
	status := int32(1)
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	var rolesID []int64

	n := 5
	for i := 0; i < n; i++ {
		roleTxResult := createRandomRoleTx(t)
		rolesID = append(rolesID, roleTxResult.Role.ID)
	}

	adArg := CreateAdminUserTxParams{
		Account:        account,
		FullName:       fullName,
		HashedPassword: hashedPassword,
		Status:         status,
		RolesID:        rolesID,
	}

	result, err := testStore.CreateAdminUserTx(context.Background(), adArg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	createdAdminUser := result.AdminUser
	require.NotEmpty(t, createdAdminUser)
	require.Equal(t, account, createdAdminUser.Account)
	require.Equal(t, fullName, createdAdminUser.FullName)
	require.Equal(t, hashedPassword, createdAdminUser.HashedPassword)
	require.Equal(t, status, createdAdminUser.Status)
	require.NotZero(t, createdAdminUser.CreatedAt)

	permissions := result.PermissionList

	targetLen := int(math.Pow(float64(n), 2))
	require.Len(t, permissions, targetLen)

	return result
}

func TestCreateAdminUserTx(t *testing.T) {
	createRandomAdminUserTx(t)
}
