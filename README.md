# MCP-SecurityTools 简介

**MCP-SecurityTools 是一个专注于收录和更新网络安全领域 MCP 的开源项目，旨在汇总、整理和优化各类与 MCP 相关的安全工具、技术及实战经验。**

## 一：uncover-MCP(使AI具有调用FOFA/SHODAN能力)

**项目地址：** `https://github.com/Co5mos/uncover-mcp`

### 1.1 编译说明

```
# 克隆仓库
git clone https://github.com/Co5mos/uncover-mcp
cd uncover-mcp

# 构建项目(需在 Go 1.21 及更高版本中)
go build -o uncover-mcp ./cmd/uncover-mcp
```

![image-20250331145154372](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331145154372.png)

### 1.2 使用方法

**作为MCP服务运行实例：**

```json
{
    "mcpServers": {
        "uncover-mcp": {
            "command": "./uncover-mcp",
            "args": [],
            "env": {
                "SHODAN_API_KEY": "",
                "FOFA_EMAIL": "",
                "FOFA_KEY": ""
            }
        }
    }
}
```

**Cherry Studio中使用**

> Cherry StudioMCP 使用教程详见：`https://docs.cherry-ai.com/advanced-basic/mcp`

```json
{
  "mcpServers": {
    "uncover-mcp": {
      "name": "uncover-mcp",
      "isActive": true,
      "command": "You compile the uncover-mcp binary file",
      "args": [
        "uncover-mcp"
      ],
      "env": {
        "SHODAN_API_KEY": "xxxxxxxxxx",
        "FOFA_EMAIL": "xxxxxxxxxx",
        "FOFA_KEY": "xxxxxxxxxx"
      }
    }
  }
}
```

![image-20250331145745743](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331145745743.png)

![image-20250331150333597](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331150333597.png)

**cline中使用**

```json
{
  "mcpServers": {
    "uncover-mcp": {
      "command": "You compile the uncover-mcp binary file",
      "args": [],
      "env": {
        "SHODAN_API_KEY": "xxxxxxxxxx",
        "FOFA_EMAIL": "xxxxxxxxxx",
        "FOFA_KEY": "xxxxxxxxxx"
      },
      "autoApprove": [
        "fofa"
      ]
    }
  }
}
```

![image-20250331151122685](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331151122685.png)

## 二：ENScan_GO-MCP(使AI具有一键信息收集能力)

**项目地址：**  `https://github.com/wgpsec/ENScan_GO`

### 2.1 使用MCP

开启MCP服务器，将监听本地的 `http://localhost:8080`

```
./enscan --mcp
```

** 以 Cherry Studio 配置为例 **

![图像-20250329160425571](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250329160425571.png)

![图像-20250329160556011](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250329160556011.png)

## 三：VirusTotal-MCP(使AI具有VirusTotal的安全分析能力)

**项目地址：** ` https://github.com/BurtTheCoder/mcp-virustotal·`

### 3.1 编译说明

```txt
# 要有node环境
git clone https://github.com/BurtTheCoder/mcp-virustotal.git
cd mcp-virustotal
npm install
npm run build
```

## 3.2 使用方法

**Cherry Studio中使用**

> Cherry StudioMCP 使用教程详见：`https://docs.cherry-ai.com/advanced-basic/mcp`

```json
{
  "mcpServers": {
    "s4Q9KPP86Ec_MWVfGURLI": {
      "isActive": true,
      "name": "virustotal-mcp",
      "description": "virustotal-mcp",
      "command": "node",
      "args": [
        "--experimental-modules",
        "You compile the uncover-mcp binary file"
      ],
      "env": {
        "VIRUSTOTAL_API_KEY": "xxxxxxxxxx"
      }
    }
  }
}
```

![image-20250331214038166](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331214038166.png)

**cline中使用**

```json
  {
    "mcpServers": {
      "virustotal": {
        "command": "node",
        "args": [
          "--experimental-modules",
          "You compile the uncover-mcp binary file"
        ],
        "env": {
          "VIRUSTOTAL_API_KEY": "xxxxxxxxxx"
        },
        "autoApprove": [
          "get_url_report",
          "get_file_report",
          "get_url_relationship",
          "get_ip_report",
          "get_domain_report",
          "get_url_relationship",
          "get_file_relationship"
        ]
      }
    }
  }
```

![image-20250331214857434](https://imges-1255470970.cos.ap-nanjing.myqcloud.com/img/image-20250331214857434.png)
