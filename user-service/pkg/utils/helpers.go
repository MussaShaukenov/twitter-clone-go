package utils

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidId = errors.New("invalid ID value")

func ReadJson(w http.ResponseWriter, r *http.Request, in interface{}) error {
	body := r.Body
	defer body.Close()

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, status int, in interface{}, headers http.Header) error {
	response, err := json.Marshal(&in)
	if err != nil {
		http.Error(w, "error during marshaling", http.StatusInternalServerError)
		return err
	}
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)

	return nil
}

func GetIdFromQueryParam(w http.ResponseWriter, r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return 0, err
	}
	if err = validateID(id); err != nil {
		return 0, err
	}
	return id, nil
}

func validateID(id int) error {
	if id < 1 {
		return ErrInvalidId
	}
	return nil
}

func GenerateRandomStringOfLength(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a slice to store the random characters
	result := make([]byte, n)

	// Populate the slice with random characters
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(result)
}

func HashPassword(plaintextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func MatchesPassword(plaintextPassword string, hashedPassword []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
