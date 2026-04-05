# curl

一个支持 HTTP/1.1、HTTP/2 和 HTTP/3 协议的简单 HTTP 客户端工具，类似于 curl 命令行工具。

## 功能特点

- 📡 **多协议支持**：支持 HTTP/1.1、HTTP/2 和 HTTP/3 协议
- 🔒 **证书验证**：可选择忽略 SSL 证书验证
- 📤 **POST 请求**：支持发送 POST 数据
- ⏱️ **延迟发送**：支持设置 POST 数据发送前的延迟时间
- 📋 **头部请求**：支持仅获取 HTTP 头部信息
- 🖨️ **详细输出**：显示响应状态码和响应体

## 技术实现

- **语言**：Go 语言
- **HTTP 客户端**：标准库 `net/http`
- **HTTP/2 支持**：`golang.org/x/net/http2`
- **HTTP/3 支持**：`github.com/quic-go/quic-go`

## 快速开始

### 前置要求

- Go 1.16 或更高版本
- Git

### 安装和运行

1. **克隆项目**

```bash
git clone https://github.com/yourusername/network-toolbox.git
cd network-toolbox/curl
```

2. **安装依赖**

```bash
go get github.com/quic-go/quic-go
go get golang.org/x/net/http2
```

3. **构建和运行**

```bash
# 构建
go build -o curl .

# 运行
./curl https://example.com
```

## 使用方法

### 基本用法

```bash
# 使用默认 HTTP/1.1 协议获取网页
./curl https://example.com

# 使用 -url 参数指定 URL
./curl -url https://example.com
```

### 协议选择

```bash
# 使用 HTTP/2 协议
./curl -http2 https://example.com

# 使用 HTTP/3 协议
./curl -http3 https://example.com
```

### POST 请求

```bash
# 发送 POST 数据
./curl -d "key=value&another=value" https://example.com/api

# 发送 POST 数据并设置延迟（1000ms）
./curl -d "key=value" -delay 1000 https://example.com/api
```

### 其他选项

```bash
# 忽略 SSL 证书验证
./curl -k https://self-signed.example.com

# 仅获取头部信息
./curl -I https://example.com
```

## 命令行参数

| 参数 | 描述 | 默认值 |
|------|------|--------|
| `-k` | 忽略 SSL 证书验证 | `false` |
| `-I` | 仅获取头部信息 | `false` |
| `-http3` | 使用 HTTP/3 协议 | `false` |
| `-http2` | 使用 HTTP/2 协议 | `false` |
| `-url` | 指定要获取的 URL | `""` |
| `-d` | POST 数据 | `""` |
| `-delay` | POST 数据发送前的延迟时间（毫秒） | `0` |

## 示例

### 1. 基本 GET 请求

```bash
./curl https://example.com
```

输出：
```
REQUEST: https://example.com, Data=
Response status: 200 OK
Response body:
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
    ...
</html>
```

### 2. POST 请求

```bash
./curl -d "name=test&value=123" https://httpbin.org/post
```

输出：
```
REQUEST: https://httpbin.org/post, Data=name=test&value=123
Response status: 200 OK
Response body:
{
  "args": {},
  "data": "",
  "files": {},
  "form": {
    "name": "test",
    "value": "123"
  },
  ...
}
```

### 3. 使用 HTTP/3 协议

```bash
./curl -http3 https://www.google.com
```

输出：
```
REQUEST: https://www.google.com, Data=
Response status: 200 OK
Response body:
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Google</title>
    ...
</html>
```

## 项目结构

```
curl/
├── curl.go          # 主程序文件
└── README.md        # 项目说明文档
```

## 核心功能实现

### 协议选择

- 根据命令行参数选择使用 HTTP/1.1、HTTP/2 或 HTTP/3 协议
- 对于 HTTP/3，使用 quic-go 库实现 QUIC 协议

### 延迟发送

- 实现了 `delayedBodyReader` 结构体，支持在发送 POST 数据前添加延迟
- 延迟仅在第一次读取数据时应用

### 错误处理

- 提供详细的错误信息，包括 URL 验证、传输创建和请求发送等环节
- 添加了 HTTP 客户端超时设置，防止请求无限期挂起

## 注意事项

- HTTP/3 功能依赖于 quic-go 库，可能需要额外的系统依赖
- 目前 HTTP/3 的 qlog 功能仅创建日志文件路径，尚未实现完整的日志记录
- 对于 HTTPS 请求，默认会验证 SSL 证书，可以使用 `-k` 参数忽略验证

## 扩展建议

- 添加更多 HTTP 方法支持（PUT、DELETE 等）
- 实现文件上传功能
- 添加 HTTP 头部自定义功能
- 实现代理支持
- 增加更多输出格式选项（如 JSON、安静模式等）

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！