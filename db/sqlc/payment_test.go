package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatePayment(t *testing.T) {
ctx := context.Background()

	// First create a user
	createUserParams := CreateUserParams{
		Email:          "payment@example.com",
		PhoneNumber:    "+1234567898",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	// Test create payment
	params := CreatePaymentParams{
		FromAccount: createdUser.ID,
		ToAccount:   "company_account_123",
		Amount:      10000,
	}

	payment, err := testQueries.CreatePayment(ctx, params)
	require.NoError(t, err)
	assert.NotZero(t, payment.ID)
	assert.Equal(t, params.FromAccount, payment.FromAccount)
	assert.Equal(t, params.ToAccount, payment.ToAccount)
	assert.Equal(t, params.Amount, payment.Amount)
}

func TestUpdatePayment(t *testing.T) {
ctx := context.Background()

	// First create a user and payment
	createUserParams := CreateUserParams{
		Email:          "updatepayment@example.com",
		PhoneNumber:    "+1234567899",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	createPaymentParams := CreatePaymentParams{
		FromAccount: createdUser.ID,
		ToAccount:   "company_account_123",
		Amount:      10000,
	}

	payment, err := testQueries.CreatePayment(ctx, createPaymentParams)
	require.NoError(t, err)

	// Test update payment
	updateParams := UpdatePaymentParams{
		Type:   "TICKET",    // Assuming PaymentType enum
		Status: "COMPLETED", // Assuming PaymentStatus enum
		ID:     payment.ID,
	}

	updatedPayment, err := testQueries.UpdatePayment(ctx, updateParams)
	require.NoError(t, err)
	assert.Equal(t, updateParams.Type, updatedPayment.Type)
	assert.Equal(t, updateParams.Status, updatedPayment.Status)
}

func TestUpdatePaymentStatus(t *testing.T) {
ctx := context.Background()

	// First create a user and payment
	createUserParams := CreateUserParams{
		Email:          "paymentstatus@example.com",
		PhoneNumber:    "+1234567800",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	createPaymentParams := CreatePaymentParams{
		FromAccount: createdUser.ID,
		ToAccount:   "company_account_123",
		Amount:      10000,
	}

	payment, err := testQueries.CreatePayment(ctx, createPaymentParams)
	require.NoError(t, err)

	// Test update payment status
	updateParams := UpdatePaymentStatusParams{
		Status: "FAILED", // Assuming PaymentStatus enum
		ID:     payment.ID,
	}

	updatedPayment, err := testQueries.UpdatePaymentStatus(ctx, updateParams)
	require.NoError(t, err)
	assert.Equal(t, updateParams.Status, updatedPayment.Status)
}

func TestUpdatePaymentAmount(t *testing.T) {
ctx := context.Background()

	// First create a user and payment
	createUserParams := CreateUserParams{
		Email:          "paymentamount@example.com",
		PhoneNumber:    "+1234567801",
		HashedPassword: "hashedpassword123",
		EmailVerified:  false,
		PhoneVerified:  false,
	}

	createdUser, err := testQueries.CreateUser(ctx, createUserParams)
	require.NoError(t, err)

	createPaymentParams := CreatePaymentParams{
		FromAccount: createdUser.ID,
		ToAccount:   "company_account_123",
		Amount:      10000,
	}

	payment, err := testQueries.CreatePayment(ctx, createPaymentParams)
	require.NoError(t, err)

	// Test update payment amount
	updateParams := UpdatePaymentAmountParams{
		Amount: 1000,
		ID:     payment.ID,
	}

	err = testQueries.UpdatePaymentAmount(ctx, updateParams)
	require.NoError(t, err)
}
