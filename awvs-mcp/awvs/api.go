package awvs

import (
	"encoding/json"
	"fmt"
	"log"
)

// 扫描类型常量
const (
	ScanTypeFull          = "full"           // 完全扫描
	ScanTypeHighRisk      = "high_risk"      // 高风险漏洞扫描
	ScanTypeXSS           = "xss"            // XSS漏洞扫描
	ScanTypeSQLi          = "sqli"           // SQL注入漏洞扫描
	ScanTypeWeakPass      = "weak_password"  // 弱口令检测
	ScanTypeCrawlOnly     = "crawl_only"     // 仅爬行
	ScanTypeMalware       = "malware"        // 恶意软件扫描
	ScanTypeLog4j         = "log4j"          // Log4j漏洞扫描
	ScanTypeBugBounty     = "bug_bounty"     // Bug Bounty高频漏洞
	ScanTypeKnownVuln     = "known_vuln"     // 已知漏洞扫描
	ScanTypeSpring4Shell  = "spring4shell"   // Spring4Shell漏洞扫描
)

// 扫描配置文件映射
var scanProfileMap = map[string]string{
	ScanTypeFull:         "11111111-1111-1111-1111-111111111111", // Full Scan
	ScanTypeHighRisk:     "11111111-1111-1111-1111-111111111112", // High Risk Vulnerabilities
	ScanTypeXSS:          "11111111-1111-1111-1111-111111111116", // XSS
	ScanTypeSQLi:         "11111111-1111-1111-1111-111111111113", // SQL Injection
	ScanTypeWeakPass:     "11111111-1111-1111-1111-111111111115", // Weak Passwords
	ScanTypeCrawlOnly:    "11111111-1111-1111-1111-111111111117", // Crawl Only
	ScanTypeMalware:      "11111111-1111-1111-1111-111111111120", // Malware Scan
	// 下面这些是高级的扫描模式，可能需要根据具体AWVS版本来调整profile ID
	ScanTypeLog4j:        "log4j_scan_profile",                  // Log4j
	ScanTypeBugBounty:    "bug_bounty_profile",                 // Bug Bounty
	ScanTypeKnownVuln:    "known_vuln_profile",                 // Known Vulnerabilities
	ScanTypeSpring4Shell: "spring4shell_profile",               // Spring4Shell
}

// Target 表示AWVS扫描目标
type Target struct {
	TargetID  string `json:"target_id"`
	Address   string `json:"address"` 
	Criticity int    `json:"criticity"`
	Status    string `json:"status"` 
}

// Scan 表示AWVS扫描任务
type Scan struct {
	ScanID    string `json:"scan_id"`
	TargetID  string `json:"target_id"` 
	ScanType  string `json:"scan_type"`
	ProfileID string `json:"profile_id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Severity  Severity `json:"severity"`
}

// Severity 表示漏洞严重性
type Severity struct {
	High    int `json:"high"`
	Medium  int `json:"medium"`
	Low     int `json:"low"`
	Info    int `json:"info"`
}

// 请求和响应的结构体
type addTargetRequest struct {
	Address   string   `json:"address"`
	Criticity int      `json:"criticity"`
	Description string  `json:"description"`
	Headers   []header `json:"custom_headers,omitempty"`
	Cookies   string   `json:"custom_cookies,omitempty"`
}

type header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type addTargetResponse struct {
	Target Target `json:"target"`
}

type startScanRequest struct {
	TargetID  string `json:"target_id"`
	ProfileID string `json:"profile_id"` 
	Schedule   struct {
		Disable    bool   `json:"disable"`
		StartDate  *string `json:"start_date,omitempty"`
		TimeZone   *string `json:"time_zone,omitempty"`
	} `json:"schedule"`
}

type startScanResponse struct {
	Scan Scan `json:"scan"`
}

type targetsResponse struct {
	Targets []Target `json:"targets"`
}

type scansResponse struct {
	Scans []Scan `json:"scans"`
}

// AddTarget 添加目标到AWVS
func (c *Client) AddTarget(url string, cookies string, headers map[string]string) (*Target, error) {
	// 构建请求体
	req := addTargetRequest{
		Address:     url,
		Criticity:   10, // Default criticity
		Description: "Added by AWVS MCP",
	}
	
	// 添加Cookie
	if cookies != "" {
		req.Cookies = cookies
	}
	
	// 添加自定义Header
	if len(headers) > 0 {
		for k, v := range headers {
			req.Headers = append(req.Headers, header{Key: k, Value: v})
		}
	}
	
	// 发送请求
	respBytes, err := c.post("/targets", req)
	if err != nil {
		return nil, fmt.Errorf("add target failed: %w", err)
	}
	
	// 解析响应
	var resp addTargetResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal add target response failed: %w", err)
	}
	
	return &resp.Target, nil
}

// StartScan 开始扫描目标
func (c *Client) StartScan(targetID, scanType string) (*Scan, error) {
	// 记录传入的targetID
	log.Printf("StartScan接收到的targetID: %s", targetID)
	
	// 获取扫描配置ID
	profileID, ok := scanProfileMap[scanType]
	if !ok {
		return nil, fmt.Errorf("invalid scan type: %s", scanType)
	}
	
	// 构建请求体
	req := startScanRequest{
		TargetID:  targetID,
		ProfileID: profileID,
		Schedule: struct {
			Disable    bool   `json:"disable"`
			StartDate  *string `json:"start_date,omitempty"`
			TimeZone   *string `json:"time_zone,omitempty"`
		}{
			Disable: true, // 禁用调度，立即开始扫描
		},
	}
	
	// 记录请求体
	jsonBytes, _ := json.Marshal(req)
	log.Printf("StartScan请求体: %s", string(jsonBytes))
	
	// 发送请求
	respBytes, err := c.post("/scans", req)
	if err != nil {
		return nil, fmt.Errorf("start scan failed: %w", err)
	}
	
	// 解析响应
	var resp startScanResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal start scan response failed: %w", err)
	}
	
	return &resp.Scan, nil
}

// AddAndScan 添加目标并开始扫描
func (c *Client) AddAndScan(url string, scanType string, cookies string, headers map[string]string) (*Scan, *Target, error) {
	// 首先尝试查找是否已经存在该URL的目标
	targets, err := c.ListTargets()
	if err == nil && len(targets) > 0 {
		// 查找匹配的目标
		var existingTarget *Target
		for _, t := range targets {
			if t.Address == url {
				existingTarget = &t
				log.Printf("找到已存在的目标: ID=%s, URL=%s", t.TargetID, t.Address)
				break
			}
		}
		
		// 如果找到匹配的目标，直接使用它
		if existingTarget != nil {
			log.Printf("使用已存在的目标: %s 开始扫描", existingTarget.TargetID)
			scan, err := c.StartScan(existingTarget.TargetID, scanType)
			return scan, existingTarget, err
		}
	}
	
	// 如果没有找到匹配的目标，添加新目标
	target, err := c.AddTarget(url, cookies, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("add target failed: %w", err)
	}
	
	log.Printf("成功添加新目标: ID=%s, URL=%s", target.TargetID, target.Address)
	
	// 开始扫描
	scan, err := c.StartScan(target.TargetID, scanType)
	if err != nil {
		return nil, target, fmt.Errorf("start scan failed: %w", err)
	}
	
	return scan, target, nil
}

// ListTargets 获取所有目标
func (c *Client) ListTargets() ([]Target, error) {
	respBytes, err := c.get("/targets")
	if err != nil {
		return nil, fmt.Errorf("list targets failed: %w", err)
	}
	
	var resp targetsResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal targets response failed: %w", err)
	}
	
	return resp.Targets, nil
}

// ListScans 获取所有扫描任务
func (c *Client) ListScans() ([]Scan, error) {
	respBytes, err := c.get("/scans")
	if err != nil {
		return nil, fmt.Errorf("list scans failed: %w", err)
	}
	
	var resp scansResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal scans response failed: %w", err)
	}
	
	return resp.Scans, nil
}

// DeleteTarget 删除指定目标
func (c *Client) DeleteTarget(targetID string) error {
	_, err := c.delete(fmt.Sprintf("/targets/%s", targetID))
	if err != nil {
		return fmt.Errorf("delete target failed: %w", err)
	}
	
	return nil
}

// DeleteScan 删除指定扫描任务
func (c *Client) DeleteScan(scanID string) error {
	_, err := c.delete(fmt.Sprintf("/scans/%s", scanID))
	if err != nil {
		return fmt.Errorf("delete scan failed: %w", err)
	}
	
	return nil
}

// DeleteAllTargets 删除所有目标
func (c *Client) DeleteAllTargets() error {
	targets, err := c.ListTargets()
	if err != nil {
		return fmt.Errorf("list targets failed: %w", err)
	}
	
	for _, target := range targets {
		if err := c.DeleteTarget(target.TargetID); err != nil {
			return fmt.Errorf("delete target %s failed: %w", target.TargetID, err)
		}
	}
	
	return nil
}

// DeleteAllScans 删除所有扫描任务
func (c *Client) DeleteAllScans() error {
	scans, err := c.ListScans()
	if err != nil {
		return fmt.Errorf("list scans failed: %w", err)
	}
	
	for _, scan := range scans {
		if err := c.DeleteScan(scan.ScanID); err != nil {
			return fmt.Errorf("delete scan %s failed: %w", scan.ScanID, err)
		}
	}
	
	return nil
}
