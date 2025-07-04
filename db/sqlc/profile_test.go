package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialProfile(t *testing.T) {
	ctx := context.Background()

	// First create a user
	createParams := CreateUserParams{
		Email:          "profile@example.com",
		PhoneNumber:    "+1234567893",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createParams)
	require.NoError(t, err)

	// Test initial profile creation
	err = testQueries.InitialProfile(ctx, createdUser.ID)
	require.NoError(t, err)
}

func TestUpdateUserContact(t *testing.T) {
	ctx := context.Background()

	// First create a user
	createParams := CreateUserParams{
		Email:          "update@example.com",
		PhoneNumber:    "+1234567894",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createParams)
	require.NoError(t, err)

	// Test update user contact
	updateParams := UpdateUserContactParams{
		ID:             createdUser.ID,
		Email:          sql.NullString{String: "updated@example.com", Valid: true},
		PhoneNumber:    sql.NullString{String: "+9876543210", Valid: true},
		HashedPassword: sql.NullString{String: "newhashedpassword", Valid: true},
	}

	updatedUser, err := testQueries.UpdateUserContact(ctx, updateParams)
	require.NoError(t, err)
	assert.Equal(t, updateParams.Email.String, updatedUser.Email)
	assert.Equal(t, updateParams.PhoneNumber.String, updatedUser.PhoneNumber)
	assert.Equal(t, updateParams.HashedPassword.String, updatedUser.HashedPassword)
}
