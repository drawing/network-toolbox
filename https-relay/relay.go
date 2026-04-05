package main

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

const (
	// 证书相关常量
	caCertFile     = "ca.crt"
	caKeyFile      = "ca.key"
	serverCertFile = "server.crt"
	serverKeyFile  = "server.key"
	dnsCacheFile   = "dns_cache.conf"
	listenAddress  = ":443"
)

var (
	dnsCache = make(map[string]string)
)

// init 初始化函数，加载配置和证书
func init() {
	// 加载 DNS 缓存配置
	loadDNSCache(dnsCacheFile)
	
	// 生成或加载证书
	if err := setupCertificates(); err != nil {
		log.Fatalf("[E] Failed to setup certificates: %v", err)
	}
}

// loadDNSCache 从配置文件加载 DNS 缓存
func loadDNSCache(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("[W] Failed to open DNS cache file: %v", err)
		log.Println("[W] Using empty DNS cache")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("[W] Invalid line format: %s", line)
			continue
		}

		domain := strings.TrimSpace(parts[0])
		address := strings.TrimSpace(parts[1])
		if domain != "" && address != "" {
			dnsCache[domain] = address
			log.Printf("[I] Loaded DNS cache: %s -> %s", domain, address)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[E] Error reading DNS cache file: %v", err)
	}

	if len(dnsCache) == 0 {
		log.Println("[W] DNS cache is empty")
	} else {
		log.Printf("[I] Loaded %d DNS cache entries", len(dnsCache))
	}
}

// setupCertificates 生成或加载证书
func setupCertificates() error {
	// 检查 CA 证书是否存在，不存在则生成
	if !fileExists(caCertFile) || !fileExists(caKeyFile) {
		log.Println("[I] CA certificate not found, generating...")
		if err := generateCACertificate(); err != nil {
			return fmt.Errorf("failed to generate CA certificate: %v", err)
		}
		log.Println("[I] CA certificate generated successfully")
	}

	// 每次都重新生成服务器证书，因为 DNS 缓存可能发生变化
	log.Println("[I] Generating server certificate...")
	if err := generateServerCertificate(); err != nil {
		return fmt.Errorf("failed to generate server certificate: %v", err)
	}
	log.Println("[I] Server certificate generated successfully")

	return nil
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// generateCACertificate 生成自签名 CA 证书
func generateCACertificate() error {
	// 生成 CA 私钥
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate CA private key: %v", err)
	}

	// 创建 CA 证书模板
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"HTTPS Relay CA"},
			Country:       []string{"CN"},
			Province:      []string{"Beijing"},
			Locality:      []string{"Beijing"},
			StreetAddress: []string{"Relay Street"},
			PostalCode:    []string{"100000"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	// 自签名 CA 证书
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("failed to create CA certificate: %v", err)
	}

	// 保存 CA 证书
	caCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})
	if err := os.WriteFile(caCertFile, caCertPEM, 0644); err != nil {
		return fmt.Errorf("failed to write CA certificate: %v", err)
	}

	// 保存 CA 私钥
	caKeyDER, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		return fmt.Errorf("failed to marshal CA private key: %v", err)
	}
	caKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: caKeyDER})
	if err := os.WriteFile(caKeyFile, caKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write CA private key: %v", err)
	}

	return nil
}

// generateServerCertificate 使用 CA 证书签发服务器证书
func generateServerCertificate() error {
	// 加载 CA 证书和私钥
	caCertPEM, err := os.ReadFile(caCertFile)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caKeyPEM, err := os.ReadFile(caKeyFile)
	if err != nil {
		return fmt.Errorf("failed to read CA private key: %v", err)
	}

	// 解析 CA 证书
	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		return fmt.Errorf("failed to decode CA certificate")
	}
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	// 解析 CA 私钥
	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		return fmt.Errorf("failed to decode CA private key")
	}
	caKey, err := x509.ParseECPrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA private key: %v", err)
	}

	// 生成服务器私钥
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate server private key: %v", err)
	}

	// 收集所有域名作为 SAN
	var dnsNames []string
	for domain := range dnsCache {
		dnsNames = append(dnsNames, domain)
	}
	// 如果没有配置域名，添加默认域名
	if len(dnsNames) == 0 {
		dnsNames = append(dnsNames, "localhost", "127.0.0.1")
	}

	// 创建服务器证书模板
	serverTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"HTTPS Relay Server"},
			Country:      []string{"CN"},
			Province:     []string{"Beijing"},
			Locality:     []string{"Beijing"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("0.0.0.0")},
	}

	// 使用 CA 签发服务器证书
	serverCertDER, err := x509.CreateCertificate(rand.Reader, &serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("failed to create server certificate: %v", err)
	}

	// 保存服务器证书
	serverCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertDER})
	if err := os.WriteFile(serverCertFile, serverCertPEM, 0644); err != nil {
		return fmt.Errorf("failed to write server certificate: %v", err)
	}

	// 保存服务器私钥
	serverKeyDER, err := x509.MarshalECPrivateKey(serverKey)
	if err != nil {
		return fmt.Errorf("failed to marshal server private key: %v", err)
	}
	serverKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: serverKeyDER})
	if err := os.WriteFile(serverKeyFile, serverKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write server private key: %v", err)
	}

	return nil
}

// handleConnection 处理客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("[E] Connection is not a TLS connection")
		return
	}

	// 执行 TLS 握手
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("[E] TLS handshake failed: %v", err)
		return
	}

	// 获取连接状态
	state := tlsConn.ConnectionState()
	if state.ServerName == "" {
		log.Println("[E] SSL SNI is empty")
		return
	}

	// 查询目标地址
	targetAddr, exists := dnsCache[state.ServerName]
	if !exists {
		log.Printf("[E] Domain not found in DNS cache: %s", state.ServerName)
		return
	}

	log.Printf("[I] Connected with: %s -> %s", state.ServerName, targetAddr)

	// 连接到目标服务器
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 跳过目标服务器证书验证
	}
	targetConn, err := tls.Dial("tcp", targetAddr, tlsConfig)
	if err != nil {
		log.Printf("[E] Failed to connect to target: %s[%s]: %v", state.ServerName, targetAddr, err)
		return
	}
	defer targetConn.Close()

	// 双向数据转发
	go io.Copy(conn, targetConn)
	io.Copy(targetConn, conn)
}

// updateHostsFile 更新 /etc/hosts 文件
func updateHostsFile() {
	const hostsFile = "/etc/hosts"
	const localhostIP = "127.0.0.1"
	const startMarker = "# === HTTPS Relay START - DO NOT EDIT MANUALLY ==="
	const endMarker = "# === HTTPS Relay END ==="

	// 检查是否有 root 权限
	file, err := os.OpenFile(hostsFile, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("[W] Permission denied: Cannot modify %s. Please manually add the following entries:", hostsFile)
		for domain := range dnsCache {
			log.Printf("[W]   %s %s", localhostIP, domain)
		}
		return
	}
	defer file.Close()

	// 读取现有 hosts 文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[E] Failed to read %s: %v", hostsFile, err)
		return
	}

	// 处理文件内容，移除旧的自动添加区域
	lines := strings.Split(string(content), "\n")
	var newContent []string
	inAutoSection := false

	for _, line := range lines {
		if line == startMarker {
			inAutoSection = true
			continue
		}
		if line == endMarker {
			inAutoSection = false
			continue
		}
		if !inAutoSection {
			newContent = append(newContent, line)
		}
	}

	// 移除末尾空行
	for len(newContent) > 0 && newContent[len(newContent)-1] == "" {
		newContent = newContent[:len(newContent)-1]
	}

	// 如果没有域名需要添加，只清理旧内容
	if len(dnsCache) == 0 {
		if err := os.Truncate(hostsFile, 0); err != nil {
			log.Printf("[E] Failed to truncate %s: %v", hostsFile, err)
			return
		}
		if _, err := file.WriteString(strings.Join(newContent, "\n")); err != nil {
			log.Printf("[E] Failed to write cleaned content to %s: %v", hostsFile, err)
			return
		}
		log.Println("[I] Cleaned up old HTTPS Relay entries from hosts file")
		return
	}

	// 添加新的自动添加区域
	newContent = append(newContent, "")
	newContent = append(newContent, startMarker)
	
	for domain := range dnsCache {
		newContent = append(newContent, fmt.Sprintf("%s %s", localhostIP, domain))
		log.Printf("[I] Added to hosts file: %s %s", localhostIP, domain)
	}
	
	newContent = append(newContent, endMarker)

	// 写入更新后的内容
	if err := os.Truncate(hostsFile, 0); err != nil {
		log.Printf("[E] Failed to truncate %s: %v", hostsFile, err)
		return
	}

	if _, err := file.WriteString(strings.Join(newContent, "\n")); err != nil {
		log.Printf("[E] Failed to write updated content to %s: %v", hostsFile, err)
		return
	}

	log.Printf("[I] Successfully updated %s with %d domains", hostsFile, len(dnsCache))
}

// main 主函数
func main() {
	// 更新 hosts 文件
	updateHostsFile()

	// 加载证书
	cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("[E] Failed to load server certificate: %v", err)
	}

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// 启动 TLS 监听器
	listener, err := tls.Listen("tcp", listenAddress, tlsConfig)
	if err != nil {
		log.Fatalf("[E] Failed to start TLS listener: %v", err)
	}
	defer listener.Close()

	log.Printf("[I] HTTPS Relay server started, listening on %s", listenAddress)
	log.Printf("[I] CA certificate: %s (add this to trusted certificates)", caCertFile)
	log.Printf("[I] Server certificate: %s", serverCertFile)

	// 接受客户端连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[E] Accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
