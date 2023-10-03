package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomVerifyEmail(t *testing.T) VerifyEmail {
	user := createRandomUser(t)

	arg := CreateVerifyEmailParams{
		UserID: pgtype.Int4{
			Int32: int32(user.ID),
			Valid: true,
		},
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}

	verifyEmail, err := testStore.CreateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, verifyEmail)

	require.Equal(t, int32(user.ID), verifyEmail.UserID.Int32)
	require.Equal(t, user.Email, verifyEmail.Email)
	require.Equal(t, arg.SecretCode, verifyEmail.SecretCode)
	require.False(t, verifyEmail.IsUsed)

	require.NotZero(t, verifyEmail.CreatedAt)
	require.True(t, verifyEmail.ExpiresAt.After(verifyEmail.CreatedAt))

	return verifyEmail
}

func TestCreateVerifyEmail(t *testing.T) {
	createRandomVerifyEmail(t)
}

func TestUpdateVerifyEmail(t *testing.T) {
	oldVerifyEmail := createRandomVerifyEmail(t)

	arg := UpdateVerifyEmailParams{
		ID:         oldVerifyEmail.ID,
		SecretCode: oldVerifyEmail.SecretCode,
	}

	newVerifyEmail, err := testStore.UpdateVerifyEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newVerifyEmail)

	require.Equal(t, oldVerifyEmail.ID, newVerifyEmail.ID)
	require.Equal(t, oldVerifyEmail.Email, newVerifyEmail.Email)
	require.Equal(t, oldVerifyEmail.SecretCode, newVerifyEmail.SecretCode)
	require.Equal(t, oldVerifyEmail.UserID.Int32, newVerifyEmail.UserID.Int32)

	require.NotEqual(t, oldVerifyEmail.IsUsed, newVerifyEmail.IsUsed)

	require.WithinDuration(t, oldVerifyEmail.CreatedAt, newVerifyEmail.CreatedAt, time.Second)
	require.WithinDuration(t, oldVerifyEmail.ExpiresAt, newVerifyEmail.ExpiresAt, time.Second)
}
