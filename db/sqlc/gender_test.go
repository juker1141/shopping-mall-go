package db

import (
	"context"
	"testing"

	"github.com/juker1141/shopping-mall-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomGender(t *testing.T) Gender {
	name := util.RandomName()
	gender, err := testStore.CreateGender(context.Background(), name)

	require.NoError(t, err)
	require.NotZero(t, gender.ID)
	require.Equal(t, name, gender.Name)
	return gender
}

func TestCreateGender(t *testing.T) {
	createRandomGender(t)
}

func TestGetGender(t *testing.T) {
	gender1 := createRandomGender(t)

	gender2, err := testStore.GetGender(context.Background(), gender1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gender2)
	require.Equal(t, gender1.ID, gender2.ID)
	require.Equal(t, gender1.Name, gender2.Name)
}

func TestListGenders(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomGender(t)
	}

	genders, err := testStore.ListGenders(context.Background())
	require.NoError(t, err)
	require.NotZero(t, genders)

	for _, gender := range genders {
		require.NotEmpty(t, gender)
	}
}
