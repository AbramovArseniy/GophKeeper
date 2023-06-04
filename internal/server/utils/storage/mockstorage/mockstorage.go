package mockstorage

import (
	"fmt"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
)

type MockData struct {
	Data  []byte
	Type  storage.InfoType
	Name  string
	Login string
}

type MockUser struct {
	Login    string
	Password string
}

type MockStorage struct {
	Storage []MockData
	Users   []MockUser
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

func (ms *MockStorage) GetData(metadata storage.InfoMeta) (storage.Info, error) {
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
	info := storage.NewInfo(metadata.Type)
	err := info.DecodeBinary(data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding binary: %w", err)
	}
	return info, nil
}

func (ms *MockStorage) Close() {
	return
}
