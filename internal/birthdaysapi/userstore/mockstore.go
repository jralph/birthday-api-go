package userstore

import "sync/atomic"

type MockUserStore struct {
	Users    map[string]*User
	PutError error
	GetError error
	GetCalls int32
	PutCalls int32
}

func (m *MockUserStore) Put(user *User) error {
	atomic.AddInt32(&m.PutCalls, 1)

	if m.PutError != nil {
		return m.PutError
	}

	m.Users[user.Username] = user

	return nil
}

func (m *MockUserStore) Get(username string) (*User, error) {
	atomic.AddInt32(&m.GetCalls, 1)

	if m.GetError != nil {
		return nil, m.GetError
	}

	if user, exists := m.Users[username]; exists {
		return user, nil
	}

	return nil, nil
}
