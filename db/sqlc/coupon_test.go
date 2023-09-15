package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomCoupon(t *testing.T, name string, code string, startAt time.Time) Coupon {
	arg := CreateCouponParams{
		Title:     name,
		Code:      code,
		Percent:   int32(util.RandomInt(1, 100)),
		CreatedBy: util.RandomName(),
		StartAt:   startAt,
		ExpiresAt: startAt.Add(1 * time.Minute),
	}

	coupon, err := testStore.CreateCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, coupon)

	require.NotZero(t, coupon.ID)
	require.Equal(t, arg.Title, coupon.Title)
	require.Equal(t, arg.Code, coupon.Code)
	require.Equal(t, arg.Percent, coupon.Percent)
	require.Equal(t, arg.CreatedBy, coupon.CreatedBy)

	require.WithinDuration(t, arg.StartAt, coupon.StartAt, time.Second)
	require.WithinDuration(t, arg.ExpiresAt, coupon.ExpiresAt, time.Second)

	require.NotZero(t, coupon.CreatedAt)
	return coupon
}

func TestCreateCoupon(t *testing.T) {
	createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
}

func TestGetCoupon(t *testing.T) {
	coupon1 := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	coupon2, err := testStore.GetCoupon(context.Background(), coupon1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, coupon2)

	require.Equal(t, coupon1.ID, coupon2.ID)
	require.Equal(t, coupon1.Title, coupon2.Title)
	require.Equal(t, coupon1.Code, coupon2.Code)
	require.Equal(t, coupon1.Percent, coupon2.Percent)
	require.Equal(t, coupon1.CreatedBy, coupon2.CreatedBy)
	require.WithinDuration(t, coupon1.StartAt, coupon2.StartAt, time.Second)
	require.WithinDuration(t, coupon1.ExpiresAt, coupon2.ExpiresAt, time.Second)

	require.WithinDuration(t, coupon1.CreatedAt, coupon2.CreatedAt, time.Second)
}

func TestUpdateCouponAllField(t *testing.T) {
	oldCoupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	arg := UpdateCouponParams{
		ID: oldCoupon.ID,
		Title: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
		Code: pgtype.Text{
			String: util.RandomString(10),
			Valid:  true,
		},
		Percent: pgtype.Int4{
			Int32: int32(util.RandomInt(1, 100)),
			Valid: true,
		},
		StartAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	newCoupon, err := testStore.UpdateCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCoupon)

	require.Equal(t, oldCoupon.ID, newCoupon.ID)
	require.Equal(t, oldCoupon.CreatedBy, newCoupon.CreatedBy)
	require.WithinDuration(t, oldCoupon.CreatedAt, newCoupon.CreatedAt, time.Second)

	require.NotEqual(t, oldCoupon.Title, newCoupon.Title)
	require.NotEqual(t, oldCoupon.Code, newCoupon.Code)
	require.NotEqual(t, oldCoupon.Percent, newCoupon.Percent)

	require.False(t, oldCoupon.StartAt.Equal(newCoupon.StartAt))
	require.False(t, oldCoupon.ExpiresAt.Equal(newCoupon.ExpiresAt))
}

func TestUpdateCouponOnlyTitle(t *testing.T) {
	oldCoupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	arg := UpdateCouponParams{
		ID: oldCoupon.ID,
		Title: pgtype.Text{
			String: util.RandomName(),
			Valid:  true,
		},
	}

	newCoupon, err := testStore.UpdateCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCoupon)

	require.Equal(t, oldCoupon.ID, newCoupon.ID)
	require.Equal(t, oldCoupon.CreatedBy, newCoupon.CreatedBy)
	require.WithinDuration(t, oldCoupon.CreatedAt, newCoupon.CreatedAt, time.Second)
	require.Equal(t, oldCoupon.Code, newCoupon.Code)
	require.Equal(t, oldCoupon.Percent, newCoupon.Percent)
	require.True(t, oldCoupon.StartAt.Equal(newCoupon.StartAt))
	require.True(t, oldCoupon.ExpiresAt.Equal(newCoupon.ExpiresAt))

	require.NotEqual(t, oldCoupon.Title, newCoupon.Title)
}

func TestUpdateCouponOnlyStartTime(t *testing.T) {
	oldCoupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	arg := UpdateCouponParams{
		ID: oldCoupon.ID,
		StartAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	newCoupon, err := testStore.UpdateCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCoupon)

	require.Equal(t, oldCoupon.ID, newCoupon.ID)
	require.Equal(t, oldCoupon.CreatedBy, newCoupon.CreatedBy)
	require.Equal(t, oldCoupon.Title, newCoupon.Title)
	require.Equal(t, oldCoupon.Code, newCoupon.Code)
	require.Equal(t, oldCoupon.Percent, newCoupon.Percent)
	require.WithinDuration(t, oldCoupon.CreatedAt, newCoupon.CreatedAt, time.Second)
	require.True(t, oldCoupon.ExpiresAt.Equal(newCoupon.ExpiresAt))

	require.False(t, oldCoupon.StartAt.Equal(newCoupon.StartAt))
}

func TestUpdateCouponOnlyExpiresTime(t *testing.T) {
	oldCoupon := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	arg := UpdateCouponParams{
		ID: oldCoupon.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	newCoupon, err := testStore.UpdateCoupon(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newCoupon)

	require.Equal(t, oldCoupon.ID, newCoupon.ID)
	require.Equal(t, oldCoupon.CreatedBy, newCoupon.CreatedBy)
	require.Equal(t, oldCoupon.Title, newCoupon.Title)
	require.Equal(t, oldCoupon.Code, newCoupon.Code)
	require.Equal(t, oldCoupon.Percent, newCoupon.Percent)
	require.WithinDuration(t, oldCoupon.CreatedAt, newCoupon.CreatedAt, time.Second)
	require.True(t, oldCoupon.StartAt.Equal(newCoupon.StartAt))

	require.False(t, oldCoupon.ExpiresAt.Equal(newCoupon.ExpiresAt))
}

func TestDeleteCoupon(t *testing.T) {
	coupon1 := createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	err := testStore.DeleteCoupon(context.Background(), coupon1.ID)
	require.NoError(t, err)

	coupon2, err := testStore.GetCoupon(context.Background(), coupon1.ID)
	require.Error(t, err)
	require.Empty(t, coupon2)
}

func TestListCoupons(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	}

	arg := ListCouponsParams{
		Limit:  5,
		Offset: 5,
	}
	coupons, err := testStore.ListCoupons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, coupons, 5)

	for _, coupon := range coupons {
		require.NotEmpty(t, coupon)
	}
}

func TestListCouponsSearchTitle(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomCoupon(t, name, util.RandomString(10), time.Now())
	}
	for i := 0; i < 10-n; i++ {
		createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	}

	arg := ListCouponsParams{
		Key:      KeyTitle,
		KeyValue: name,
		Limit:    10,
		Offset:   0,
	}
	coupons, err := testStore.ListCoupons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, coupons, n)

	for _, coupon := range coupons {
		require.NotEmpty(t, coupon)
	}
}

func TestListCouponsSearchTitleButNoKeyValue(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomCoupon(t, name, util.RandomString(10), time.Now())
	}
	for i := 0; i < 10-n; i++ {
		createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	}

	arg := ListCouponsParams{
		Key:    KeyTitle,
		Limit:  5,
		Offset: 5,
	}
	coupons, err := testStore.ListCoupons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, coupons, 5)

	for _, coupon := range coupons {
		require.NotEmpty(t, coupon)
	}
}

func TestListCouponsSearchTitleButNoKey(t *testing.T) {
	n := 3
	name := util.RandomName()
	for i := 0; i < n; i++ {
		createRandomCoupon(t, name, util.RandomString(10), time.Now())
	}
	for i := 0; i < 10-n; i++ {
		createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	}

	arg := ListCouponsParams{
		KeyValue: name,
		Limit:    5,
		Offset:   5,
	}
	coupons, err := testStore.ListCoupons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, coupons, 5)

	for _, coupon := range coupons {
		require.NotEmpty(t, coupon)
	}
}

func TestListCouponsSearchCode(t *testing.T) {
	n := 3
	code := util.RandomString(10)
	for i := 0; i < n; i++ {
		createRandomCoupon(t, util.RandomName(), code, time.Now())
	}
	for i := 0; i < 10-n; i++ {
		createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())
	}

	arg := ListCouponsParams{
		Key:      KeyCode,
		KeyValue: code,
		Limit:    10,
		Offset:   0,
	}
	coupons, err := testStore.ListCoupons(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, coupons, n)

	for _, coupon := range coupons {
		require.NotEmpty(t, coupon)
	}
}

func TestGetCouponsCount(t *testing.T) {
	createRandomCoupon(t, util.RandomName(), util.RandomString(10), time.Now())

	count, err := testStore.GetCouponsCount(context.Background())
	require.NoError(t, err)
	require.NotZero(t, count)
}
