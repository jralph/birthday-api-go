package utils

import "birthdays-api/internal/birthdaysApi/userStore"

type MockUserStore struct {
	Users    map[string]*userStore.User
	PutError error
	GetError error
}

func (m *MockUserStore) Put(user *userStore.User) error {
	if m.PutError != nil {
		return m.PutError
	}

	m.Users[user.Username] = user

	return nil
}

func (m *MockUserStore) Get(username string) (*userStore.User, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}

	if user, exists := m.Users[username]; exists {
		return user, nil
	}

	return nil, nil
}
