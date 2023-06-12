package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AbramovArseniy/GophKeeper/internal/client/utils/types"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage"
)

type HTTPClient struct {
	client  http.Client
	address string
}

func (c *HTTPClient) Connect(address string) error {
	c.client = http.Client{}
	c.address = address
	return nil
}

func (c *HTTPClient) SavePassword(ctx context.Context, req storage.InfoLoginPass, infoType storage.InfoType) error {
	byteBody, err := json.Marshal(req)
	if err != nil {
		log.Println("error, while marshalling json body:", err)
		return err
	}
	byteBody = []byte(string(byteBody[:len(byteBody)-1]) + fmt.Sprintf(", \"type\":\"%s\"}", infoType))
	reqBody := bytes.NewBuffer(byteBody)
	resp, err := c.client.Post(c.address, "application/json", reqBody)
	if err != nil {
		return fmt.Errorf("error while making request: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned status %d, error", resp.StatusCode)
	}
	return nil
}

func (c *HTTPClient) GetData(ctx context.Context, req types.GetRequest) (storage.Info, error) {
	byteBody, err := json.Marshal(req)
	if err != nil {
		log.Println("error, while marshalling json body:", err)
		return nil, err
	}
	reqBody := bytes.NewBuffer(byteBody)
	resp, err := c.client.Post(c.address, "application/json", reqBody)
	if err != nil {
		return nil, fmt.Errorf("error while making request: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}
	respInfo := storage.NewInfo(req.Type)
	err = json.Unmarshal(body, respInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal response body: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("server returned status %d, error", resp.StatusCode)
	}
	return respInfo, nil
}
