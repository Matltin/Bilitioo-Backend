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

func TestGetUserProfile(t *testing.T) {
	ctx := context.Background()

	// First create a user and initial profile
	createUserParams := CreateUserParams{
		Email:          "userprofile@example.com",
		PhoneNumber:    "+1234567802",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	err = testQueries.InitialProfile(ctx, createdUser.ID)
	require.NoError(t, err)

	// Test get user profile
	profile, err := testQueries.GetUserProfile(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, profile.UserID)
	assert.Equal(t, createdUser.Email, profile.Email)
	assert.Equal(t, createdUser.PhoneNumber, profile.PhoneNumber)
}

func TestUpdateProfile(t *testing.T) {
	ctx := context.Background()

	// First create a user and initial profile
	createUserParams := CreateUserParams{
		Email:          "updateprofile@example.com",
		PhoneNumber:    "+1234567803",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	err = testQueries.InitialProfile(ctx, createdUser.ID)
	require.NoError(t, err)

	// Test update profile
	updateParams := UpdateProfileParams{
		UserID:       createdUser.ID,
		FirstName:    sql.NullString{String: "John", Valid: true},
		LastName:     sql.NullString{String: "Doe", Valid: true},
		NationalCode: sql.NullString{String: "1234567890", Valid: true},
	}

	updatedProfile, err := testQueries.UpdateProfile(ctx, updateParams)
	require.NoError(t, err)
	assert.Equal(t, updateParams.FirstName.String, updatedProfile.FirstName)
	assert.Equal(t, updateParams.LastName.String, updatedProfile.LastName)
	assert.Equal(t, updateParams.NationalCode.String, updatedProfile.NationalCode)
}

func TestAddToUserWallet(t *testing.T) {
	ctx := context.Background()

	// First create a user and initial profile
	createUserParams := CreateUserParams{
		Email:          "wallet@example.com",
		PhoneNumber:    "+1234567804",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	err = testQueries.InitialProfile(ctx, createdUser.ID)
	require.NoError(t, err)

	// Test add to user wallet
	addParams := AddToUserWalletParams{
		Wallet: 5000,
		UserID: createdUser.ID,
	}

	err = testQueries.AddToUserWallet(ctx, addParams)
	require.NoError(t, err)
}
