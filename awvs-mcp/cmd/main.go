package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/taoing/awvs-mcp/awvs"
	"github.com/taoing/awvs-mcp/models"

	"os/signal"
)

func main() {
	// 定义命令行参数
	var (
		mode       string
		port       int
		configPath string
	)

	flag.StringVar(&configPath, "config", "config.json", "配置文件路径")

	// 解析命令行参数
	flag.Parse()

	// 判断运行模式
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("请指定运行模式: stdio 或 http")
		os.Exit(1)
	}

	mode = args[0]

	// 处理http模式的端口参数
	if mode == "http" {
		portFlag := flag.NewFlagSet("http", flag.ExitOnError)
		portFlag.IntVar(&port, "port", 8080, "HTTP服务器端口")
		if err := portFlag.Parse(args[1:]); err != nil {
			fmt.Printf("解析HTTP参数出错: %v\n", err)
			os.Exit(1)
		}
	}

	// 处理配置文件路径
	if !filepath.IsAbs(configPath) {
		absPath, err := filepath.Abs(configPath)
		if err != nil {
			fmt.Printf("获取配置文件绝对路径失败: %v\n", err)
			os.Exit(1)
		}
		configPath = absPath
	}

	// 读取配置文件
	configData, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 解析配置
	var config models.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		fmt.Printf("解析配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 创建AWVS客户端
	awvsClient := awvs.NewClient(&awvs.Config{
		APIURL:    config.APIURL,
		APIKey:    config.APIKey,
		VerifySSL: config.VerifySSL,
	})

	// 创建MCP服务器
	mcpServer := server.NewMCPServer(
		"AWVS Scanner", // 服务器名称
		"1.0.0",       // 版本
		server.WithLogging(),
		server.WithToolCapabilities(true),
	)

	// 注册AWVS工具
	registerAWVSTool(mcpServer, awvsClient)

	// 初始化上下文
	ctx := context.Background()

	// 根据模式启动服务器
	switch mode {
	case "stdio":
		fmt.Println("启动AWVS扫描器服务器 (Stdio模式)...")
		log.Println("Starting AWVS Scanner in Stdio mode...")
		// 启动终端模式服务
		if err := server.ServeStdio(mcpServer); err != nil {
			fmt.Printf("服务器错误: %v\n", err)
			os.Exit(1)
		}
	case "http":
		fmt.Printf("启动AWVS扫描器服务器 (HTTP模式，端口: %d)...\n", port)
		log.Printf("Starting AWVS Scanner in HTTP mode on port %d...", port)
		// 创建SSE服务器并启动
		sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL(fmt.Sprintf("http://localhost:%d", port)))
		// 启动服务器
		if err := sseServer.Start(fmt.Sprintf(":%d", port)); err != nil {
			fmt.Printf("服务器错误: %v\n", err)
			os.Exit(1)
		}
		
		// 等待中断信号
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		
		// 优雅关闭服务器
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()
		sseServer.Shutdown(shutdownCtx)
	default:
		fmt.Printf("不支持的模式 '%s'，请使用 'stdio' 或 'http'\n", mode)
		os.Exit(1)
	}
}

// 注册AWVS扫描工具
func registerAWVSTool(mcpServer *server.MCPServer, awvsClient *awvs.Client) {
	// 创建扫描站点工具
	scanTool := mcp.NewTool("scan_website",
		mcp.WithDescription("扫描网站漏洞"),
		mcp.WithString("url",
			mcp.Description("要扫描的目标URL"),
			mcp.Required(),
		),
		mcp.WithString("scan_type",
			mcp.Description("要执行的扫描类型"),
			mcp.Enum("full", "high_risk", "xss", "sqli", "weak_password", "crawl_only", "malware", "log4j", "bug_bounty", "known_vuln", "spring4shell"),
			mcp.Required(),
		),
		mcp.WithString("cookies",
			mcp.Description("扫描时使用的Cookie")),
		mcp.WithObject("headers",
			mcp.Description("扫描时使用的HTTP头"),
			mcp.AdditionalProperties(true)),
	)

	// 创建列出目标工具
	listTargetsTool := mcp.NewTool("list_targets",
		mcp.WithDescription("列出所有扫描目标"),
	)

	// 创建列出扫描工具
	listScansTool := mcp.NewTool("list_scans",
		mcp.WithDescription("列出所有扫描任务"),
	)

	// 创建删除所有目标工具
	deleteAllTool := mcp.NewTool("delete_all",
		mcp.WithDescription("删除所有目标和扫描"),
	)

	// 添加扫描工具到服务器
	mcpServer.AddTool(scanTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 从请求中获取参数
		url, _ := request.Params.Arguments["url"].(string)
		scanType, _ := request.Params.Arguments["scan_type"].(string)
		cookies, _ := request.Params.Arguments["cookies"].(string)
		headersObj, _ := request.Params.Arguments["headers"].(map[string]interface{})

		// 转换headers为map[string]string格式
		headersMap := make(map[string]string)
		if headersObj != nil {
			for key, value := range headersObj {
				if strValue, ok := value.(string); ok {
					headersMap[key] = strValue
				}
			}
		}

		// 添加目标并开始扫描
		scan, target, err := awvsClient.AddAndScan(url, scanType, cookies, headersMap)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("扫描失败: %v", err),
					},
				},
			}, nil
		}

		// 构建响应
		responseData := map[string]interface{}{
			"target_id": target.TargetID,
			"scan_id":   scan.ScanID,
			"url":       url,
			"scan_type": scanType,
		}

		// 转换为JSON
		responseJSON, _ := json.Marshal(responseData)

		// 返回结果
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(responseJSON),
				},
			},
		}, nil
	})

	// 添加列出目标工具到服务器
	mcpServer.AddTool(listTargetsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取所有目标
		targets, err := awvsClient.ListTargets()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("获取目标失败: %v", err),
					},
				},
			}, nil
		}

		// 构建响应
		responseData := map[string]interface{}{
			"targets": targets,
			"count":   len(targets),
		}

		// 转换为JSON
		responseJSON, _ := json.Marshal(responseData)

		// 返回目标列表
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(responseJSON),
				},
			},
		}, nil
	})

	// 添加列出扫描工具到服务器
	mcpServer.AddTool(listScansTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 获取所有扫描
		scans, err := awvsClient.ListScans()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("获取扫描失败: %v", err),
					},
				},
			}, nil
		}

		// 构建响应
		responseData := map[string]interface{}{
			"scans": scans,
			"count": len(scans),
		}

		// 转换为JSON
		responseJSON, _ := json.Marshal(responseData)

		// 返回扫描列表
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(responseJSON),
				},
			},
		}, nil
	})

	// 注册删除所有目标工具
	mcpServer.AddTool(deleteAllTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// 删除所有目标
		err := awvsClient.DeleteAllTargets()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: fmt.Sprintf("删除所有目标失败: %v", err),
					},
				},
			}, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "所有目标和扫描已成功删除",
				},
			},
		}, nil
	})
}
