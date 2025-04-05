package models

// Config 表示AWVS MCP配置
type Config struct {
	APIURL    string `json:"api_url"`    // AWVS API URL
	APIKey    string `json:"api_key"`    // AWVS API 密钥
	VerifySSL bool   `json:"verify_ssl"` // 是否验证SSL证书
}
