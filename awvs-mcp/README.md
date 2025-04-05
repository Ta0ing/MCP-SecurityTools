# AWVS MCP

基于[mcp-go](https://github.com/mark3labs/mcp-go)实现的Acunetix Web Vulnerability Scanner (AWVS)扫描器MCP。

## 功能特点

- 支持Stdio和SSE两种模式
- 批量添加URL到AWVS扫描器并进行扫描
- 支持多种扫描类型：完全扫描、高风险漏洞扫描、XSS漏洞扫描、SQL注入漏洞扫描等
- 支持清空扫描任务和目标
- 自定义扫描参数

## 安装

```bash
go get github.com/taoing/awvs-mcp
```

## 使用方法

### 配置

使用前需要在config.json中配置AWVS API的地址和API密钥：

```json
{
  "api_url": "https://localhost:3443/api/v1",
  "api_key": "your_api_key_here",
  "verify_ssl": false
}
```

### 启动服务

#### Stdio模式

```bash
awvs-mcp stdio
```

#### SSE模式

```bash
awvs-mcp sse --port 8080
```

## API工具

本MCP实现提供以下工具：

- `scan` - 添加URL并开始扫描
- `list_targets` - 列出所有扫描目标
- `list_scans` - 列出所有扫描任务
- `delete_all` - 删除所有目标和扫描任务
- `delete_scans` - 仅删除扫描任务
- `scan_existing` - 对已有目标开始新的扫描
