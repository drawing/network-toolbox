# https-relay

一个 HTTPS 中继服务器，支持通过配置文件指定域名和目标地址的映射关系，并自动生成证书。

## 功能特点

- 🚀 **自动证书生成**：自动生成自签名 CA 证书和服务器证书
- 🔒 **域名自动签发**：根据配置文件中的域名自动签发证书
- 📝 **配置文件支持**：通过配置文件管理域名和目标地址的映射
- 🔄 **实时更新**：服务器证书会根据配置文件的变化自动重新生成
- 🔐 **安全通信**：使用 TLS 加密通信，确保数据传输安全

## 技术实现

- **语言**：Go 语言
- **加密**：ECDSA P256 算法
- **TLS 版本**：TLS 1.2+
- **证书有效期**：1 年

## 快速开始

### 前置要求

- Go 1.16 或更高版本
- Git

### 安装和运行

1. **克隆项目**

```bash
git clone https://github.com/yourusername/network-toolbox.git
cd network-toolbox/https-relay
```

2. **配置域名映射**

编辑 `dns_cache.conf` 文件，添加域名和目标地址的映射关系：

```
# DNS Cache Configuration
# Format: domain=ip:port
# Example: example.com=192.168.1.1:443

example.com=192.168.1.1:443
```

3. **运行服务**

```bash
go run relay.go
```

服务会自动生成：
- `ca.crt` - CA 证书（需要添加到系统信任列表）
- `ca.key` - CA 私钥
- `server.crt` - 服务器证书
- `server.key` - 服务器私钥

4. **信任 CA 证书**

将生成的 `ca.crt` 文件添加到系统的可信证书列表中：

**macOS:**
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ca.crt
```

**Linux:**
```bash
sudo cp ca.crt /usr/local/share/ca-certificates/
sudo update-ca-certificates
```

**Windows:**
1. 双击 `ca.crt` 文件
2. 点击 "安装证书"
3. 选择 "本地计算机"
4. 选择 "将所有证书放入下列存储"
5. 浏览并选择 "受信任的根证书颁发机构"
6. 完成安装

## 配置文件说明

`dns_cache.conf` 文件格式：
```
# 注释行，以 # 开头
domain1=ip1:port1
domain2=ip2:port2
```

## 使用方法

1. **配置域名映射**：在 `dns_cache.conf` 中添加需要中继的域名和目标地址
2. **启动服务**：运行 `go run relay.go`
3. **信任 CA 证书**：将生成的 `ca.crt` 添加到系统信任列表
4. **访问服务**：使用配置的域名访问服务，流量会自动中继到目标地址

## 示例

### 示例：中继 example.com

配置文件：
```
example.com=192.168.1.1:443
```

启动服务后，访问 `https://example.com`，流量会被中继到 `192.168.1.1:443`。


## 项目结构

```
https-relay/
├── relay.go          # 主程序文件
├── dns_cache.conf    # DNS 缓存配置文件（需创建）
├── ca.crt            # CA 证书（自动生成）
├── ca.key            # CA 私钥（自动生成）
├── server.crt        # 服务器证书（自动生成）
├── server.key        # 服务器私钥（自动生成）
└── README.md         # 项目说明文档
```

## 注意事项

- 配置文件 `dns_cache.conf` 不会被提交到 Git
- 证书文件（`ca.crt`、`ca.key`、`server.crt`、`server.key`）不会被提交到 Git
- CA 证书只需要生成一次，服务器证书会根据配置文件的变化自动重新生成
- 服务器监听在 443 端口，需要管理员权限运行

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！