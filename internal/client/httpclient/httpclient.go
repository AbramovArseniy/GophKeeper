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
	auth    string
}

func NewHTTPClient(address string) *HTTPClient {
	return &HTTPClient{
		client:  *http.DefaultClient,
		address: "http://" + address,
	}
}

func (c *HTTPClient) SaveData(ctx context.Context, req storage.Info, infoType storage.InfoType, infoName string) error {
	byteBody, err := json.Marshal(req)
	if err != nil {
		log.Println("error, while marshalling json body:", err)
		return err
	}
	byteBody = []byte(string(byteBody[:len(byteBody)-1]) + fmt.Sprintf(", \"type\":\"%s\", \"name\":\"%s\"}", infoType, infoName))
	reqBody := bytes.NewBuffer(byteBody)
	httpReq, err := http.NewRequest("Post", c.address+"/user/add-data/", reqBody)
	if err != nil {
		log.Println("error, while creating http request:", err)
		return err
	}
	httpReq.Header.Set("Authorization", c.auth)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error while doing request: %w", err)
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
	httpReq, err := http.NewRequest("Post", c.address+"/user/get-data-by-name/", reqBody)
	if err != nil {
		log.Println("error, while creating http request:", err)
		return nil, err
	}
	httpReq.Header.Set("Authorization", c.auth)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error while doing request: %w", err)
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

func (c *HTTPClient) Register(ctx context.Context, req types.AuthRequest) error {
	byteBody, err := json.Marshal(req)
	if err != nil {
		log.Println("error, while marshalling json body:", err)
		return err
	}
	reqBody := bytes.NewBuffer(byteBody)
	httpReq, err := http.NewRequest("Post", c.address+"/user/auth/register/", reqBody)
	if err != nil {
		log.Println("error, while creating http request:", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error while doing request: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned status %d, error", resp.StatusCode)
	}
	c.auth = "Bearer " + resp.Header.Get("Authorization")
	return nil
}

func (c *HTTPClient) Login(ctx context.Context, req types.AuthRequest) error {
	byteBody, err := json.Marshal(req)
	if err != nil {
		log.Println("error, while marshalling json body:", err)
		return err
	}
	reqBody := bytes.NewBuffer(byteBody)
	httpReq, err := http.NewRequest("Post", c.address+"/user/auth/login/", reqBody)
	if err != nil {
		log.Println("error, while creating http request:", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error while doing request: %w", err)
	}
	resp.Body.Close()
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned status %d, error", resp.StatusCode)
	}
	c.auth = "Bearer " + resp.Header.Get("Authorization")
	return nil
}
