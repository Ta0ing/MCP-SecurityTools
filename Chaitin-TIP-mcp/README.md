# Chaitin IP Lookup MCP

这是一个用于查询长亭威胁情报 IP 数据库的 MCP（Machine Capability Provider）实现。

## 支持的平台

提供以下平台的预编译二进制文件：

- Linux
  - AMD64 (x86_64): `chaitin-mcp-linux-amd64`
  - ARM64: `chaitin-mcp-linux-arm64`

- macOS
  - Intel (AMD64): `chaitin-mcp-darwin-amd64`
  - Apple Silicon (ARM64): `chaitin-mcp-darwin-arm64`

- Windows
  - AMD64 (x86_64): `chaitin-mcp-windows-amd64.exe`
  - ARM64: `chaitin-mcp-windows-arm64.exe`

## 编译和运行

1. 设置长亭的密钥：

   ```bash
   # Linux/macOS
   export CHAITIN_SK="your_secret_key_here"
   
   # Windows (CMD)
   set CHAITIN_SK=your_secret_key_here
   
   # Windows (PowerShell)
   $env:CHAITIN_SK="your_secret_key_here"
   ```

2. 运行 MCP 服务器：

   ```bash
   # Linux/macOS
   ./chaitin-mcp-[os]-[arch]
   
   # Windows
   chaitin-mcp-windows-[arch].exe
   ```

   根据你的操作系统和架构选择对应的二进制文件。

### 从源码编译

如果需要自行编译，可以使用提供的编译脚本：

```bash
# 编译所有平台的二进制文件
./build.sh

# 编译结果将保存在 build/ 目录下
```

## 在 AI 中使用

在支持 MCP 的 AI 系统中，你可以这样使用这个功能：

### 连接到 MCP

连接到本地运行的 Chaitin IP Lookup MCP

### 使用 ip_lookup 工具

使用 ip_lookup 工具查询 IP，参数如下：
```json
{
  "ip": "8.8.8.8"  // 要查询的 IP 地址
}
```

### 返回结果示例

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    // IP 信息将在这里显示
  }
}
```

## 可用功能

### ip_lookup

- 功能：查询 IP 的威胁情报信息
- 参数：
  - `ip`：要查询的 IP 地址（必填）
- 返回：JSON 格式的查询结果

## 注意事项

1. 必须设置环境变量 `CHAITIN_SK` 为有效的长亭 API 密钥
2. MCP 服务器使用标准输入输出进行通信，不要直接在终端中运行
3. 选择运行与你的系统架构匹配的二进制文件
