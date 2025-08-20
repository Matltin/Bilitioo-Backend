package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	params := CreateUserParams{
		Email:          "test@example.com",
		PhoneNumber:    "+1234567890",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	user, err := testQueries.CreateUser(ctx, params)
	require.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, params.Email, user.Email)
	assert.Equal(t, params.PhoneNumber, user.PhoneNumber)
	assert.Equal(t, params.HashedPassword, user.HashedPassword)
	assert.Equal(t, params.EmailVerified, user.EmailVerified)
	assert.Equal(t, params.PhoneVerified, user.PhoneVerified)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()

	// First create a user
	createParams := CreateUserParams{
		Email:          "getuser@example.com",
		PhoneNumber:    "+1234567891",
		HashedPassword: "hashedpassword123",
		EmailVerified:  true,
		PhoneVerified:  true,
	}

	createdUser, err := testQueries.CreateUser(ctx, createParams)
	require.NoError(t, err)

	// Test get user by email
	getParams := GetUserParams{
		Email:       createParams.Email,
		PhoneNumber: "",
	}

	user, err := testQueries.GetUser(ctx, getParams)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Email, user.Email)
	assert.Equal(t, createdUser.PhoneNumber, user.PhoneNumber)
}

func TestGetUserByID(t *testing.T) {
	ctx := context.Background()
	
	// First create a user
	createParams := CreateUserParams{
		Email:          "getuserbyid@example.com",
		PhoneNumber:    "+1234567892",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}
	
	createdUser, err := testQueries.CreateUser(ctx, createParams)
	require.NoError(t, err)

	// Test get user by ID
	user, err := testQueries.GetUserByID(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Email, user.Email)
	assert.Equal(t, createdUser.PhoneNumber, user.PhoneNumber)
	assert.Equal(t, createdUser.HashedPassword, user.HashedPassword)
}