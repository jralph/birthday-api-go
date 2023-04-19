package handlers

import (
	"birthdays-api/internal/birthdaysApi/userStore"
	"birthdays-api/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGETHelloUsername(t *testing.T) {
	testGetHello(t, "Joe", false, 0, http.StatusNotFound, false)
	testGetHello(t, "Joe", true, 5, http.StatusOK, false)
	testGetHello(t, "Joe", true, 2, http.StatusOK, false)
	testGetHello(t, "Joe", true, 100, http.StatusOK, false)
	testGetHello(t, "Joe", true, 0, http.StatusOK, false)
	testGetHello(t, "Joe123", false, 0, http.StatusBadRequest, false)
	testGetHello(t, "Joe", true, 0, http.StatusInternalServerError, true)
}

func TestPUTHelloUsername(t *testing.T) {
	testPutHello(t, "Joe", "hi", false, http.StatusBadRequest, false)
	testPutHello(t, "Joe", `{"dateOfBirth": "2020-05-08"}`, true, http.StatusInternalServerError, true)
	testPutHello(t, "Joe", `{"dateOfBirth": "2020-05-08"}`, true, http.StatusNoContent, false)
	testPutHello(t, "Joe", `{"dateOfBirth": "2020-05-08"}`, false, http.StatusNoContent, false)
	testPutHello(t, "Joe", `{"dateOfBirth": "2020-15-08"}`, false, http.StatusBadRequest, false)
	testPutHello(t, "Joe", `{"dateOfBirth": "2100-05-08"}`, false, http.StatusBadRequest, false)
	testPutHello(t, "Joe123", `{"dateOfBirth": "2020-05-08"}`, false, http.StatusBadRequest, false)
	testPutHello(t, "Joe", fmt.Sprintf(`{"dateOfBirth": "%s"}`, time.Now().Format("2006-01-02")), false, http.StatusBadRequest, false)
}

func testGetHello(t *testing.T, user string, create bool, days int, statusCode int, forceErr bool) {
	t.Run(fmt.Sprintf("get /hello/%s status code %d with days until birthday %d", user, statusCode, days), func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/hello/%s", user), nil)
		response := httptest.NewRecorder()
		params := httprouter.Params{
			{Key: "username", Value: user},
		}

		mockUserStore := utils.MockUserStore{
			Users: map[string]*userStore.User{},
		}
		if forceErr {
			mockUserStore.GetError = fmt.Errorf("mock store get error")
		}

		if create {
			storeErr := mockUserStore.Put(&userStore.User{
				Username:    user,
				DateOfBirth: userStore.DateOfBirth(time.Now().AddDate(0, 0, days)),
			})

			if storeErr != nil && !forceErr {
				t.Errorf("error with user store mock: %s", storeErr)
			}
		}

		handler := HelloHandler{UserStore: &mockUserStore}
		handler.GetHelloUsername(response, request, params)

		body := response.Body.String()
		got := GetHelloUsernameResponse{}
		err := json.Unmarshal([]byte(body), &got)
		if err != nil && response.Code == 200 {
			t.Error(err)
		}

		if response.Code != statusCode {
			t.Errorf("unexpected error code: got %d want %d", response.Code, statusCode)
		}

		if statusCode == 200 {
			want := GetHelloUsernameResponse{Message: fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", user, days)}
			if days == 0 {
				want = GetHelloUsernameResponse{Message: fmt.Sprintf("Hello, %s! Happy birthday!", user)}
			}

			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("GetHelloUsernameResponse{} mismatch (-want +got):\n%s", diff)
			}
		}
	})
}

func testPutHello(t *testing.T, user string, body string, create bool, statusCode int, forceErr bool) {
	t.Run(fmt.Sprintf("put /hello/%s", user), func(t *testing.T) {
		body := bytes.NewReader([]byte(body))
		request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/hello/%s", user), body)
		response := httptest.NewRecorder()
		params := httprouter.Params{
			{Key: "username", Value: user},
		}

		mockUserStore := utils.MockUserStore{
			Users: map[string]*userStore.User{},
		}
		if forceErr {
			mockUserStore.PutError = fmt.Errorf("mock store put error")
		}

		if create {
			storeErr := mockUserStore.Put(&userStore.User{
				Username:    user,
				DateOfBirth: userStore.DateOfBirth(time.Now().AddDate(0, 0, 10)),
			})

			if storeErr != nil && !forceErr {
				t.Errorf("error with user store mock: %s", storeErr)
			}
		}

		handler := HelloHandler{UserStore: &mockUserStore}
		handler.PutHelloUsername(response, request, params)

		want := statusCode
		if response.Code != want {
			t.Errorf("unexpected error code: got %d want %d", response.Code, want)
		}
	})
}
