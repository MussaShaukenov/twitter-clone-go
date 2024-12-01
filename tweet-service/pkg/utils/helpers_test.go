package utils

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestValidateId(t *testing.T) {
	tc := []struct {
		name string
		id   int
		err  error
	}{
		{
			name: "valid id",
			id:   1,
			err:  nil,
		},
		{
			name: "invalid id",
			id:   0,
			err:  ErrInvalidId,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			err := validateID(tt.id)
			if err != tt.err {
				t.Errorf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestGetIdFromQueryParam(t *testing.T) {
	// Define test cases
	tc := []struct {
		name    string
		id      string
		wantID  int
		wantErr error
	}{
		{
			name:    "valid id",
			id:      "1",
			wantID:  1,
			wantErr: nil,
		},
		{
			name:    "invalid id",
			id:      "0",
			wantID:  0,
			wantErr: ErrInvalidId,
		},
		{
			name:    "non-numeric id",
			id:      "abc",
			wantID:  0,
			wantErr: strconv.ErrSyntax,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock router to inject the path parameter
			r := chi.NewRouter()
			r.Get("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
				id, err := GetIdFromQueryParam(w, r)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				if id != tt.wantID {
					t.Errorf("expected id %d, got %d", tt.wantID, id)
				}
			})

			// Create a test request with the id as a path parameter
			req := httptest.NewRequest(http.MethodGet, "/test/"+tt.id, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
		})
	}
}

func TestGenerateRandomStringOfLength(t *testing.T) {
	// Generate a random string of length 10
	randomString := GenerateRandomStringOfLength(10)

	// Check if the length of the generated string is 10
	if len(randomString) != 10 {
		t.Errorf("expected string length 10, got %d", len(randomString))
	}
}

func TestHashPassword(t *testing.T) {
	// Define a plaintext password
	plaintextPassword := "password123"

	// Hash the plaintext password
	hashedPassword, err := HashPassword(plaintextPassword)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	// Check if the hashed password is not equal to the plaintext password
	if hashedPassword == plaintextPassword {
		t.Errorf("expected hashed password to be different from plaintext password")
	}
}

func TestMatchesPassword(t *testing.T) {
	// Define a plaintext password
	plaintextPassword := "password123"

	// Hash the plaintext password
	hashedPassword, err := HashPassword(plaintextPassword)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	// Check if the plaintext password matches the hashed password
	match, err := MatchesPassword(plaintextPassword, []byte(hashedPassword))
	if err != nil {
		t.Errorf("error matching password: %v", err)
	}

	// Check if the password match is true
	if !match {
		t.Errorf("expected password match to be true")
	}
}
