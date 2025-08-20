package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserEmailVerified(t *testing.T) {
	ctx := context.Background()

	// First create a user
	createParams := CreateUserParams{
		Email:          "emailverify@example.com",
		PhoneNumber:    "+1234567895",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createParams)
	require.NoError(t, err)

	// Test update email verified
	updateParams := UpdateUserEmailVerifiedParams{
		ID:            createdUser.ID,
		EmailVerified: true,
	}

	err = testQueries.UpdateUserEmailVerified(ctx, updateParams)
	require.NoError(t, err)
}

func TestCreateVerifyEmail(t *testing.T) {
	ctx := context.Background()

	// First create a user
	createUserParams := CreateUserParams{
		Email:          "verifyemail@example.com",
		PhoneNumber:    "+1234567896",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	// Test create verify email
	params := CreateVerifyEmailParams{
		UserID:     createdUser.ID,
		Email:      createdUser.Email,
		SecretCode: "123456",
	}

	verifyEmail, err := testQueries.CreateVerifyEmail(ctx, params)
	require.NoError(t, err)
	assert.NotZero(t, verifyEmail.ID)
	assert.Equal(t, params.UserID, verifyEmail.UserID)
	assert.Equal(t, params.Email, verifyEmail.Email)
	assert.Equal(t, params.SecretCode, verifyEmail.SecretCode)
	assert.False(t, verifyEmail.IsUsed)
}

func TestUpdateVerifyEmail(t *testing.T) {
	ctx := context.Background()

	// First create a user and verify email
	createUserParams := CreateUserParams{
		Email:          "updateverify@example.com",
		PhoneNumber:    "+1234567897",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	createVerifyParams := CreateVerifyEmailParams{
		UserID:     createdUser.ID,
		Email:      createdUser.Email,
		SecretCode: "123456",
	}

	verifyEmail, err := testQueries.CreateVerifyEmail(ctx, createVerifyParams)
	require.NoError(t, err)

	// Test update verify email
	updateParams := UpdateVerifyEmailParams{
		ID:         verifyEmail.ID,
		SecretCode: "123456",
	}

	updatedVerifyEmail, err := testQueries.UpdateVerifyEmail(ctx, updateParams)
	require.NoError(t, err)
	assert.True(t, updatedVerifyEmail.IsUsed)
}
