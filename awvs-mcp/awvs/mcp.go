package awvs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Server 表示AWVS控制服务器
type Server struct {
	client  *Client
	version string
	name    string
}

// NewServer 创建一个新的AWVS控制服务器
func NewServer(configPath string) (*Server, error) {
	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	// 解析配置
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("parse config file failed: %w", err)
	}

	// 创建AWVS客户端
	client := NewClient(&config)

	// 返回服务器实例
	return &Server{
		client:  client,
		version: "1.0.0",
		name:    "AWVS Scanner",
	}, nil
}

// 请求和响应结构
type Request struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Response struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ScanRequest struct {
	URL      string            `json:"url"`
	ScanType string            `json:"scan_type"`
	Cookies  string            `json:"cookies,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

type ListTargetsRequest struct{}

type ListScansRequest struct{}

type DeleteAllRequest struct{}

type DeleteScansRequest struct{}

type ScanExistingRequest struct {
	TargetID string `json:"target_id"`
	ScanType string `json:"scan_type"`
}

// ServeStdio 启动Stdio模式的服务器
func (s *Server) ServeStdio() error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	log.Println("Starting AWVS Scanner in Stdio mode...")

	for {
		// 读取一行输入
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read stdin failed: %w", err)
		}

		// 解析请求
		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			fmt.Fprintf(writer, "{\"error\":\"parse request failed: %s\"}\n", err)
			continue
		}

		// 处理请求
		resp, err := s.handleRequest(req)
		if err != nil {
			resp = Response{
				ID:   req.ID,
				Type: "error",
				Payload: ErrorResponse{
					Error: err.Error(),
				},
			}
		}

		// 发送响应
		respData, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(writer, "{\"error\":\"marshal response failed: %s\"}\n", err)
			continue
		}

		fmt.Fprintf(writer, "%s\n", respData)
	}

	return nil
}

// ServeSSE 启动SSE模式的服务器
func (s *Server) ServeSSE(port int) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "AWVS Scanner Server is running. Use /api endpoint for communication.")
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
			return
		}

		// 设置SSE响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 解析请求
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			json.NewEncoder(w).Encode(Response{
				ID:   "",
				Type: "error",
				Payload: ErrorResponse{
					Error: fmt.Sprintf("parse request failed: %s", err),
				},
			})
			return
		}

		// 处理请求
		resp, err := s.handleRequest(req)
		if err != nil {
			resp = Response{
				ID:   req.ID,
				Type: "error",
				Payload: ErrorResponse{
					Error: err.Error(),
				},
			}
		}

		// 发送响应
		respData, err := json.Marshal(resp)
		if err != nil {
			json.NewEncoder(w).Encode(Response{
				ID:   req.ID,
				Type: "error",
				Payload: ErrorResponse{
					Error: fmt.Sprintf("marshal response failed: %s", err),
				},
			})
			return
		}

		fmt.Fprintf(w, "data: %s\n\n", respData)
	})

	log.Printf("Starting AWVS Scanner in SSE mode on port %d...", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// 处理请求
func (s *Server) handleRequest(req Request) (Response, error) {
	switch req.Type {
	case "scan":
		return s.handleScan(req)
	case "list_targets":
		return s.handleListTargets(req)
	case "list_scans":
		return s.handleListScans(req)
	case "delete_all":
		return s.handleDeleteAll(req)
	case "delete_scans":
		return s.handleDeleteScans(req)
	case "scan_existing":
		return s.handleScanExisting(req)
	default:
		return Response{}, fmt.Errorf("unknown request type: %s", req.Type)
	}
}

// 处理scan请求
func (s *Server) handleScan(req Request) (Response, error) {
	// 解析请求载荷
	var scanReq ScanRequest
	if err := json.Unmarshal(req.Payload, &scanReq); err != nil {
		return Response{}, fmt.Errorf("parse scan request failed: %w", err)
	}

	// 验证URL
	if scanReq.URL == "" {
		return Response{}, fmt.Errorf("url must be a non-empty string")
	}

	// 验证扫描类型
	if scanReq.ScanType == "" {
		return Response{}, fmt.Errorf("scan_type must be a non-empty string")
	}

	// 添加目标
	target, err := s.client.AddTarget(scanReq.URL, scanReq.Cookies, scanReq.Headers)
	if err != nil {
		return Response{}, fmt.Errorf("add target failed: %w", err)
	}

	// 开始扫描
	scan, err := s.client.StartScan(target.TargetID, scanReq.ScanType)
	if err != nil {
		return Response{}, fmt.Errorf("start scan failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"target":  target,
		"scan":    scan,
		"message": fmt.Sprintf("成功添加目标 %s 并开始 %s 扫描", scanReq.URL, scanReq.ScanType),
	}

	return Response{
		ID:      req.ID,
		Type:    "scan_result",
		Payload: response,
	}, nil
}

// 处理list_targets请求
func (s *Server) handleListTargets(req Request) (Response, error) {
	// 获取所有目标
	targets, err := s.client.ListTargets()
	if err != nil {
		return Response{}, fmt.Errorf("list targets failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"targets": targets,
		"count":   len(targets),
	}

	return Response{
		ID:      req.ID,
		Type:    "targets_result",
		Payload: response,
	}, nil
}

// 处理list_scans请求
func (s *Server) handleListScans(req Request) (Response, error) {
	// 获取所有扫描任务
	scans, err := s.client.ListScans()
	if err != nil {
		return Response{}, fmt.Errorf("list scans failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"scans": scans,
		"count": len(scans),
	}

	return Response{
		ID:      req.ID,
		Type:    "scans_result",
		Payload: response,
	}, nil
}

// 处理delete_all请求
func (s *Server) handleDeleteAll(req Request) (Response, error) {
	// 删除所有目标（会级联删除所有扫描任务）
	if err := s.client.DeleteAllTargets(); err != nil {
		return Response{}, fmt.Errorf("delete all targets failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"message": "成功删除所有目标和扫描任务",
	}

	return Response{
		ID:      req.ID,
		Type:    "delete_all_result",
		Payload: response,
	}, nil
}

// 处理delete_scans请求
func (s *Server) handleDeleteScans(req Request) (Response, error) {
	// 删除所有扫描任务
	if err := s.client.DeleteAllScans(); err != nil {
		return Response{}, fmt.Errorf("delete all scans failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"message": "成功删除所有扫描任务",
	}

	return Response{
		ID:      req.ID,
		Type:    "delete_scans_result",
		Payload: response,
	}, nil
}

// 处理scan_existing请求
func (s *Server) handleScanExisting(req Request) (Response, error) {
	// 解析请求载荷
	var scanReq ScanExistingRequest
	if err := json.Unmarshal(req.Payload, &scanReq); err != nil {
		return Response{}, fmt.Errorf("parse scan_existing request failed: %w", err)
	}

	// 验证目标ID
	if scanReq.TargetID == "" {
		return Response{}, fmt.Errorf("target_id must be a non-empty string")
	}

	// 验证扫描类型
	if scanReq.ScanType == "" {
		return Response{}, fmt.Errorf("scan_type must be a non-empty string")
	}

	// 开始扫描
	scan, err := s.client.StartScan(scanReq.TargetID, scanReq.ScanType)
	if err != nil {
		return Response{}, fmt.Errorf("start scan failed: %w", err)
	}

	// 构建响应
	response := map[string]interface{}{
		"scan":    scan,
		"message": fmt.Sprintf("成功对目标 %s 开始 %s 扫描", scanReq.TargetID, scanReq.ScanType),
	}

	return Response{
		ID:      req.ID,
		Type:    "scan_existing_result",
		Payload: response,
	}, nil
}
