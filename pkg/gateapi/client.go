package gateapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	apiKey    string
	apiSecret string
	baseURL   string
	client    *http.Client
}

func NewClient(apiKey, apiSecret, baseURL string) *Client {
	return &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// 生成签名
func (c *Client) generateSignature(method, path, body string, timestamp int64) string {
	message := fmt.Sprintf("%s\n%s\n%s\n%d\n", method, path, body, timestamp)
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// 添加认证头
func (c *Client) setAuthHeaders(req *http.Request, method, path, body string) {
	timestamp := time.Now().Unix()
	signature := c.generateSignature(method, path, body, timestamp)

	req.Header.Set("KEY", c.apiKey)
	req.Header.Set("SIGN", signature)
	req.Header.Set("Timestamp", strconv.FormatInt(timestamp, 10))
	req.Header.Set("Content-Type", "application/json")
}

// GET 请求
func (c *Client) Get(path string, result interface{}) error {
	fullURL := c.baseURL + path
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return err
	}

	c.setAuthHeaders(req, "GET", path, "")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp, result)
}

// POST 请求
func (c *Client) Post(path string, body, result interface{}) error {
	fullURL := c.baseURL + path

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return err
	}

	c.setAuthHeaders(req, "POST", path, string(jsonBody))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp, result)
}

// DELETE 请求
func (c *Client) Delete(path string, result interface{}) error {
	fullURL := c.baseURL + path
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return err
	}

	c.setAuthHeaders(req, "DELETE", path, "")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp, result)
}

// 解析响应
func (c *Client) parseResponse(resp *http.Response, result interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return json.Unmarshal(body, result)
}

// API 错误
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Gate API Error [%d]: %s", e.StatusCode, e.Message)
}
