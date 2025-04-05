package awvs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Config AWVS配置
type Config struct {
	APIURL    string
	APIKey    string
	VerifySSL bool
}

// Client AWVS API客户端
type Client struct {
	config  *Config
	httpCli *http.Client
}

// NewClient 创建一个新的AWVS客户端
func NewClient(config *Config) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.VerifySSL},
	}

	httpCli := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &Client{
		config:  config,
		httpCli: httpCli,
	}
}

// request 执行HTTP请求
func (c *Client) request(method, path string, body interface{}) ([]byte, error) {
	// 添加API版本号路径
	apiPath := fmt.Sprintf("/api/v1%s", path)
	url := fmt.Sprintf("%s%s", c.config.APIURL, apiPath)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body failed: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth", c.config.APIKey)

	log.Printf("请求方法：%s，请求地址：%s，请求头：%v，请求体：%v\n", method, url, req.Header, body)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}

	log.Printf("响应状态码：%d，响应头：%v，响应体：%s\n", resp.StatusCode, resp.Header, string(respBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// get 执行GET请求
func (c *Client) get(path string) ([]byte, error) {
	return c.request(http.MethodGet, path, nil)
}

// post 执行POST请求
func (c *Client) post(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodPost, path, body)
}

// delete 执行DELETE请求
func (c *Client) delete(path string) ([]byte, error) {
	return c.request(http.MethodDelete, path, nil)
}
