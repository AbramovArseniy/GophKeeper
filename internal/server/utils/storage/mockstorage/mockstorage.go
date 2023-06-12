package mockstorage

import (
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/types"
)

type MockData struct {
	Data  []byte
	Type  storage.InfoType
	Name  string
	Login string
}

type MockStorage struct {
	Storage []MockData
	Users   []types.User
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Storage: make([]MockData, 100),
	}
}

func (ms *MockStorage) SaveData(encryptedData []byte, metadata storage.InfoMeta) error {
	ms.Storage = append(ms.Storage, MockData{
		Data: encryptedData,
		Type: metadata.Type,
		Name: metadata.Name,
	})
	return nil
}

func (ms *MockStorage) GetData(metadata storage.InfoMeta) ([]byte, error) {
	var data []byte
	exists := false
	for _, md := range ms.Storage {
		if md.Type == metadata.Type && md.Name == metadata.Name {
			data = md.Data
			exists = true
		}
	}
	if !exists {
		return nil, storage.ErrDataNotFound
	}
	return data, nil
}

func (ms *MockStorage) Close() {}

func (ms *MockStorage) FindUser(login string) (*types.User, error) {
	for _, user := range ms.Users {
		if user.Login == login {
			return &user, nil
		}
	}
	return nil, storage.ErrDataNotFound
}

func (ms *MockStorage) RegisterNewUser(login string, password string) (types.User, error) {
	user := types.User{
		ID:           len(ms.Users) + 1,
		Login:        login,
		HashPassword: password,
	}
	for _, user := range ms.Users {
		if user.Login == login {
			return types.User{}, storage.ErrUserExists
		}
	}
	ms.Users = append(ms.Users, user)
	return user, nil
}

func (ms *MockStorage) GetUserData(login string) (types.User, error) {
	for _, user := range ms.Users {
		if user.Login == login {
			return user, nil
		}
	}
	return types.User{}, storage.ErrDataNotFound
}
