package main

import (
	"birthdays-api/internal/birthdaysApi/userStore"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPServerHelloEndpoints(t *testing.T) {
	mockStore := &userStore.MockUserStore{
		Users: map[string]*userStore.User{},
	}

	s := httptest.NewServer(CreateHandler(mockStore))
	defer s.Close()

	var err error
	resp := &http.Response{}

	// Try and fetch user that does not (yet) exist
	resp, err = getHelloRequest(s.URL, "joe")
	if err != nil {
		t.Errorf("error making get user request when user does not exist: %s", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf("expected status code 404, got %d", resp.StatusCode)
	}

	// Put the user into the system with a date of birth
	resp, err = putHelloRequest(s.URL, "joe", "2000-05-05")
	if err != nil {
		t.Errorf("error putting new user: %s", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status code 204, got %d", resp.StatusCode)
	}

	// Fetch the newly created user
	resp, err = getHelloRequest(s.URL, "joe")
	if err != nil {
		t.Errorf("error making get user request when user exists: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	// Try and create an invalid user
	resp, err = putHelloRequest(s.URL, "joe123", "2000-05-05")
	if err != nil {
		t.Errorf("error making put user request when username is invalid: %s", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status code 400, got %d", resp.StatusCode)
	}

	// Update the existing users data of birth
	resp, err = putHelloRequest(s.URL, "joe", "2000-05-06")
	if err != nil {
		t.Errorf("error making put user request when user exists: %s", err)
	}
	if resp.StatusCode != 204 {
		t.Errorf("expected status code 204, got %d", resp.StatusCode)
	}
}

func urlForUser(host string, user string) string {
	return fmt.Sprintf("%s/hello/%s", host, user)
}

func getHelloRequest(host string, user string) (*http.Response, error) {
	return http.Get(urlForUser(host, user))
}

func putHelloRequest(host string, user string, dob string) (*http.Response, error) {
	jsonBody := []byte(fmt.Sprintf(`{"dateOfBirth": "%s"}`, dob))
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPut, urlForUser(host, user), bodyReader)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
