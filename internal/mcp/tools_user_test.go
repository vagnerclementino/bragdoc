package mcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vagnerclementino/bragdoc/internal/domain"
	"pgregory.net/rapid"
)

// Feature: mcp-server-mode, Property 15: User Retrieve Round-Trip
//
// **Validates: Requirements 5.1, 5.3**

func TestHandleUserGet_Property_RoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		fix := newTestServer()
		ctx := context.Background()

		// Generate user data
		id := rapid.Int64Range(1, 100000).Draw(t, "userID")
		name := rapid.StringMatching(`[A-Za-z ]{3,30}`).Draw(t, "name")
		email := rapid.StringMatching(`[a-z]{3,10}@[a-z]{3,8}\.[a-z]{2,4}`).Draw(t, "email")
		jobTitle := rapid.StringMatching(`[A-Za-z ]{3,20}`).Draw(t, "jobTitle")
		company := rapid.StringMatching(`[A-Za-z ]{3,20}`).Draw(t, "company")

		user := &domain.User{
			ID:        id,
			Name:      name,
			Email:     email,
			JobTitle:  jobTitle,
			Company:   company,
			Locale:    domain.Locale("en-US"),
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}

		// Mock both retrieval paths
		fix.userRepo.On("Select", mock.Anything, id).Return(user, nil)
		fix.userRepo.On("SelectByEmail", mock.Anything, email).Return(user, nil)

		// Call handleUserGet by ID
		resultByID, _, err := fix.server.handleUserGet(ctx, nil, UserGetParams{ID: id})
		if err != nil {
			t.Fatalf("handleUserGet returned error: %v", err)
		}
		assert.False(t, resultByID.IsError)

		var respByID UserResponse
		if err := unmarshalResult(resultByID, &respByID); err != nil {
			t.Fatalf("unmarshal by ID: %v", err)
		}

		// Call handleUserGetByEmail
		resultByEmail, _, err := fix.server.handleUserGetByEmail(ctx, nil, UserGetByEmailParams{Email: email})
		if err != nil {
			t.Fatalf("handleUserGetByEmail returned error: %v", err)
		}
		assert.False(t, resultByEmail.IsError)

		var respByEmail UserResponse
		if err := unmarshalResult(resultByEmail, &respByEmail); err != nil {
			t.Fatalf("unmarshal by email: %v", err)
		}

		// Property: both retrieval paths return the same profile
		assert.Equal(t, respByID, respByEmail)
		assert.Equal(t, id, respByID.ID)
		assert.Equal(t, name, respByID.Name)
		assert.Equal(t, email, respByID.Email)
		assert.Equal(t, jobTitle, respByID.JobTitle)
		assert.Equal(t, company, respByID.Company)
	})
}

func TestHandleUserGet_NotFound(t *testing.T) {
	fix := newTestServer()
	ctx := context.Background()

	fix.userRepo.On("Select", mock.Anything, int64(999)).Return(nil, fmt.Errorf("user not found"))

	result, _, err := fix.server.handleUserGet(ctx, nil, UserGetParams{ID: 999})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, extractText(result), "not found")
}

func TestHandleUserGetByEmail_NotFound(t *testing.T) {
	fix := newTestServer()
	ctx := context.Background()

	fix.userRepo.On("SelectByEmail", mock.Anything, "nobody@example.com").Return(nil, fmt.Errorf("user not found"))

	result, _, err := fix.server.handleUserGetByEmail(ctx, nil, UserGetByEmailParams{Email: "nobody@example.com"})
	assert.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, extractText(result), "not found")
}
