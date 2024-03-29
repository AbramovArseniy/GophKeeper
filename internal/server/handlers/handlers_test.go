package handlers

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/config"
	"github.com/AbramovArseniy/GophKeeper/internal/server/utils/storage/mockstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*func Example() {
	cfg := config.Config{
		Address: "locashost:8080",
	}
	s := NewServer(cfg)
	server := httptest.NewServer(s.Route())
	defer server.Close()
	requests := []struct {
		Name   string
		URL    string
		Method string
		Body   string
	}{
		{
			Name:   "Post ",
			URL:    server.URL + "/update/gauge/Alloc/200.10",
			Method: http.MethodPost,
		},
		{
			Name:   "Post counter",
			URL:    server.URL + "/update/counter/PollCount/5",
			Method: http.MethodPost,
		},
		{
			Name:   "Get gauge",
			URL:    server.URL + "/value/gauge/Alloc",
			Method: http.MethodGet,
		},
		{
			Name:   "Get counter",
			URL:    server.URL + "/value/counter/PollCount",
			Method: http.MethodGet,
		},
		{
			Name:   "Get all metrics",
			URL:    server.URL + "/",
			Method: http.MethodGet,
		},
		{
			Name:   "Post gauge",
			URL:    server.URL + "/update/",
			Method: http.MethodPost,
			Body: `{
				"id":"Alloc",
				"type":"gauge",
				"value":400
			}`,
		},
		{
			Name:   "Post counter",
			URL:    server.URL + "/update/counter/PollCount/5",
			Method: http.MethodPost,
			Body: `{
				"id":"PollCount",
				"type":"counter",
				"value":100
			}`,
		},
		{
			Name:   "Get counter",
			URL:    server.URL + "/value/",
			Method: http.MethodPost,
			Body: `{
				"id":"PollCount",
				"type":"counter"
			}`,
		},
		{
			Name:   "Get gauge",
			URL:    server.URL + "/value/",
			Method: http.MethodPost,
			Body: `{
				"id":"Alloc",
				"type":"gauge"
			}`,
		},
		{
			Name:   "Get all metrics",
			URL:    server.URL + "/",
			Method: http.MethodGet,
		},
	}
	for _, v := range requests {
		if v.Method == http.MethodPost {
			rdr := strings.NewReader(v.Body)
			resp, err := http.DefaultClient.Post(v.URL, "application/json", rdr)
			if err != nil {
				fmt.Println("error while getting response from server", err)
				return
			}
			body, _ := io.ReadAll(resp.Body)
			stringBody := string(body)
			fmt.Println(stringBody)
			resp.Body.Close()
		} else {
			resp, err := http.DefaultClient.Get(v.URL)
			if err != nil {
				fmt.Println("error while getting response from server", err)
				return
			}
			body, _ := io.ReadAll(resp.Body)
			stringBody := string(body)
			fmt.Println(stringBody)
			resp.Body.Close()
		}
	}
	// Output:
	// 	200.1
	// 5
	// PollCount: 5
	// Alloc: 200.100000
	//
	// {"id":"Alloc","type":"gauge","value":400}
	//
	// {"id":"PollCount","type":"counter","delta":10}
	// {"id":"Alloc","type":"gauge","value":400}
	// PollCount: 10
	// Alloc: 400.000000
}
*/

// TestHandlers tests handlers
func TestHandlers(t *testing.T) {
	type want struct {
		code int
		body []string
	}
	tests := []struct {
		name   string
		URL    string
		method string
		body   string
		want   want
	}{
		{
			name:   "200 Success register",
			URL:    "/user/auth/register/",
			method: http.MethodPost,
			body:   `{"login":"test_login", "password":"test_password"}`,
			want:   want{code: 200},
		},
		{
			name:   "400 Bad request register empty login",
			URL:    "/user/auth/register/",
			method: http.MethodPost,
			body:   `{"login":"", "password":"test_password"}`,
			want:   want{code: 400},
		},
		{
			name:   "400 Bad request register empty password",
			URL:    "/user/auth/register/",
			method: http.MethodPost,
			body:   `{"login":"err_login", "password":""}`,
			want:   want{code: 400},
		},
		{
			name:   "400 Bad request login empty login",
			URL:    "/user/auth/login/",
			method: http.MethodPost,
			body:   `{"login":"", "password":"test_password"}`,
			want:   want{code: 400},
		},
		{
			name:   "400 Bad request login empty password",
			URL:    "/user/auth/login/",
			method: http.MethodPost,
			body:   `{"login":"test_login", "password":""}`,
			want:   want{code: 400},
		},
		{
			name:   "200 Success login",
			URL:    "/user/auth/login/",
			method: http.MethodPost,
			body:   `{"login":"test_login", "password":"test_password"}`,
			want:   want{code: 200},
		},
		{
			name:   "200 Success upload text data",
			URL:    "/user/add-data/",
			method: http.MethodPost,
			body:   `{"text":"some_text","type":"text","name":"text_data"}`,
			want:   want{code: 200},
		},
		{
			name:   "200 Success upload login_password",
			URL:    "/user/add-data/",
			method: http.MethodPost,
			body:   `{"login":"some_login","password":"some_password", "type":"login-password","name":"login-password_data"}`,
			want:   want{code: 200},
		},
		{
			name:   "200 Success upload card",
			URL:    "/user/add-data/",
			method: http.MethodPost,
			body:   `{"number":"123456","holder":"user","exp_date":"10/22","cvc":"123", "type":"card","name":"card_data"}`,
			want:   want{code: 200},
		},
		{
			name:   "400 Bad Request upload wrong type",
			URL:    "/user/add-data/",
			method: http.MethodPost,
			body:   `{"type":"wrong_type","name":"wrong-type_data"}`,
			want:   want{code: http.StatusBadRequest},
		},
		{
			name:   "200 Success get card data",
			URL:    "/user/get-data-by-name/",
			method: http.MethodPost,
			body:   `{"type":"card","name":"card_data"}`,
			want: want{code: http.StatusNotFound,
				body: []string{`{"number":"123456","holder":"user","exp_date":"10/22","cvc":"123"}`}},
		},
		{
			name:   "200 Success get text data",
			URL:    "/user/get-data-by-name/",
			method: http.MethodPost,
			body:   `{"type":"card","name":"text_data"}`,
			want: want{code: http.StatusNotFound,
				body: []string{`{"text":"some_text"}`}},
		},
		{
			name:   "200 Success get text data",
			URL:    "/user/get-data-by-name/",
			method: http.MethodPost,
			body:   `{"type":"card","name":"login-password_data"}`,
			want: want{code: http.StatusNotFound,
				body: []string{`{"login":"some_login","password":"some_password"}`}},
		},
		{
			name:   "404 Not Found get new name",
			URL:    "/user/get-data-by-name/",
			method: http.MethodPost,
			body:   `{"type":"card","name":"new-name_data"}`,
			want:   want{code: http.StatusNotFound},
		},
		{
			name:   "404 Not Found get wrong type",
			URL:    "/user/get-data-by-name/",
			method: http.MethodPost,
			body:   `{"type":"wrong-type","name":"new-name_data"}`,
			want:   want{code: http.StatusNotFound},
		},
	}
	cfg := config.Config{
		Address:   "locashost:8080",
		SecretKey: "secretKeyReallyy",
	}
	s := NewServer(cfg)
	ms := mockstorage.NewMockStorage()
	s.Storage = ms
	s.Auth = NewAuth(context.Background(), ms, "jwt_secret")
	server := httptest.NewServer(s.Route())
	var auth = ""
	defer server.Close()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, authHeader, body := RunRequest(t, server, tt.method, tt.URL, tt.body, "application/json", auth)
			auth = authHeader
			defer resp.Body.Close()
			assert.Equal(t, tt.want.code, resp.StatusCode)
			for _, s := range tt.want.body {
				assert.Equal(t, body, s)
			}
			assert.Equal(t, tt.want.code, resp.StatusCode)
		})
	}
}

// RunRequest does request to a server
func RunRequest(t *testing.T, ts *httptest.Server, method string, query string, body string, contentType string, authorization string) (*http.Response, string, string) {
	reader := strings.NewReader(body)
	req, err := http.NewRequest(method, ts.URL+query, reader)
	req.Header.Add("Content-Type", contentType)
	if query != "/user/auth/login/" {
		req.Header.Add("Authorization", "Bearer "+authorization)
	}
	log.Printf(`query: %s
	auth: %s
	body: %s`, query, req.Header.Get("Authorization"), body)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	auth := resp.Header.Get("Authorization")
	if query != "/user/auth/login/" || auth == "" || auth == "Bearer " {
		auth = authorization
	}
	RespBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, auth, string(RespBody)
}
