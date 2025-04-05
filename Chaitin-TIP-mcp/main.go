package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ChaitinIPResponse represents the response structure from Chaitin API
type ChaitinIPResponse struct {
	Code int             `json:"code"`
	Msg  string         `json:"msg"`
	Data map[string]any `json:"data"`
}

// ChaitinIPLookup implements the IP lookup functionality
func ChaitinIPLookup(ip string) (*ChaitinIPResponse, error) {
	sk := os.Getenv("CHAITIN_SK")
	if sk == "" {
		return nil, fmt.Errorf("CHAITIN_SK environment variable not set")
	}

	url := fmt.Sprintf("https://ip-0.rivers.chaitin.cn/api/share/s?sk=%s&ip=%s", sk, ip)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query Chaitin API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var result ChaitinIPResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &result, nil
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Chaitin IP Lookup",
		"1.0.0",
		server.WithLogging(),
	)

	// Add IP lookup tool
	ipLookupTool := mcp.NewTool("ip_lookup",
		mcp.WithDescription("Look up IP information using Chaitin Threat Intelligence"),
		mcp.WithString("ip",
			mcp.Required(),
			mcp.Description("The IP address to look up"),
		),
	)

	// Add the IP lookup handler
	s.AddTool(ipLookupTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ip := request.Params.Arguments["ip"].(string)
		
		result, err := ChaitinIPLookup(ip)
		if err != nil {
			return nil, err
		}

		// Convert result to JSON string for output
		jsonResult, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to format result: %v", err)
		}

		return mcp.NewToolResultText(string(jsonResult)), nil
	})

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}
