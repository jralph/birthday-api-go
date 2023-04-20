package handlers

import (
	"birthdays-api/internal/birthdaysapi/userstore"
	"birthdays-api/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// GetHelloUsernameResponse contains the response format from the GET /hello/:username request
type GetHelloUsernameResponse struct {
	Message string `json:"message"`
}

// ErrorResponse is for use with common errors to be returned from any request
type ErrorResponse struct {
	Error string `json:"error"`
}

// PutHelloUsernameRequest contains the request format for the PUT /hello/:username request
type PutHelloUsernameRequest struct {
	DateOfBirth userstore.DateOfBirth `json:"dateOfBirth"`
}

// HelloHandler is the primary handler for dealing with /hello/ based requests
type HelloHandler struct {
	UserStore userstore.UserStore
}

// PutHelloUsername handles any PUT /hello/:username request and deals with validating and then saving the user and its date of birth
func (h *HelloHandler) PutHelloUsername(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Read the request body to its byte representation
	body, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("error reading request body: %w", err))
		return
	}

	// Unmarshal the request body into the required struct
	requestData := PutHelloUsernameRequest{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("error decoding request: %w", err))

		response := ErrorResponse{Error: "unable to decode request body"}
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			log.Println(fmt.Errorf("error encoding response: %w", err))
			return
		}
		return
	}

	// If the time is in the future, or neither future nor past (the exact same time), throw an error
	if utils.IsTimeValid(time.Time(requestData.DateOfBirth)) {
		writer.WriteHeader(http.StatusBadRequest)
		response := ErrorResponse{Error: "unable to use date of birth in the future"}
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			log.Println(fmt.Errorf("error encoding response: %w", err))
			return
		}
		return
	}

	// Create a new user object
	user := &userstore.User{
		Username:    params.ByName("username"),
		DateOfBirth: requestData.DateOfBirth,
	}

	// Validate the new user object
	if err = validateUser(user); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("error validating user struct: %w", err))

		response := ErrorResponse{Error: "username must be alpha (letters only)"}
		err = json.NewEncoder(writer).Encode(response)
		if err != nil {
			log.Println(fmt.Errorf("error encoding response: %w", err))
			return
		}
		return
	}

	// Store the user in the provided store
	err = h.UserStore.Put(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("error putting user: %w", err))
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

// GetHelloUsername handles GET /hello/:username requests
func (h *HelloHandler) GetHelloUsername(writer http.ResponseWriter, _ *http.Request, params httprouter.Params) {
	// For simplicity, we just create a user object here with the username from the request in
	// We can then use this user object to quickly perform validation on the username format and return a 400
	user := &userstore.User{
		Username: params.ByName("username"),
	}
	if err := validateUser(user); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println(fmt.Errorf("error validating user struct: %w", err))

		response := ErrorResponse{Error: "username must be alpha (letters only)"}
		err := json.NewEncoder(writer).Encode(response)
		if err != nil {
			log.Println(fmt.Errorf("error encoding response: %w", err))
			return
		}
		return
	}

	// Fetch the user (and override the existing user object we already have)
	user, err := h.UserStore.Get(user.Username)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("error fetching user: %w", err))
		return
	}

	// If user and error are both nil, the user was simply not found
	if user == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Fetch the required message for the user based on its date of birth
	response := GetHelloUsernameResponse{
		Message: getUserMessage(user),
	}

	// Write a json response back to the user
	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		log.Println(fmt.Errorf("error encoding response: %w", err))
		return
	}
}

// getUserMessage handles converting the date of birth of a user into a message to display to them
func getUserMessage(user *userstore.User) string {
	dob := time.Time(user.DateOfBirth)
	now := time.Now()
	y, m1, d1 := now.Date()
	_, m2, d2 := dob.Date()

	if m1 == m2 && d1 == d2 {
		return fmt.Sprintf("Hello, %s! Happy birthday!", user.Username)
	}

	// Calculate the next birthday
	birthday, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", y, m2, d2))
	birthday = birthday.Add((time.Hour * 23) + (time.Minute + 59) + (time.Second + 59))
	until := time.Until(birthday)
	if until < 0 {
		birthday = birthday.AddDate(1, 0, 0)
		until = time.Until(birthday)
	}
	return fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", user.Username, int(math.Max(until.Hours()/24, 1)))
}

func validateUser(user *userstore.User) error {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return fmt.Errorf("error validating user struct: %w", err)
	}

	return nil
}
