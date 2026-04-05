# https-relay

An HTTPS relay server that supports mapping domain names to target addresses via configuration file and automatically generates certificates.

## Features

- 🚀 **Automatic certificate generation**: Automatically generates self-signed CA certificate and server certificate
- 🔒 **Automatic domain signing**: Automatically signs certificates for domains in the configuration file
- 📝 **Configuration file support**: Manage domain-to-target mappings through configuration file
- 🔄 **Real-time updates**: Server certificate automatically regenerates based on configuration changes
- 🔐 **Secure communication**: Uses TLS encryption for secure data transmission

## Technical Implementation

- **Language**: Go
- **Encryption**: ECDSA P256 algorithm
- **TLS version**: TLS 1.2+
- **Certificate validity**: 1 year

## Quick Start

### Prerequisites

- Go 1.16 or higher
- Git

### Installation and Running

1. **Clone the project**

```bash
git clone https://github.com/yourusername/network-toolbox.git
cd network-toolbox/https-relay
```

2. **Configure domain mappings**

Edit the `dns_cache.conf` file to add domain and target address mappings:

```
# DNS Cache Configuration
# Format: domain=ip:port
# Example: example.com=192.168.1.1:443

example.com=192.168.1.1:443
```

3. **Run the service**

```bash
go run relay.go
```

The service will automatically generate:
- `ca.crt` - CA certificate (needs to be added to system trusted list)
- `ca.key` - CA private key
- `server.crt` - Server certificate
- `server.key` - Server private key

4. **Trust the CA certificate**

Add the generated `ca.crt` file to your system's trusted certificate list:

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
1. Double-click the `ca.crt` file
2. Click "Install Certificate"
3. Select "Local Machine"
4. Select "Place all certificates in the following store"
5. Browse and select "Trusted Root Certification Authorities"
6. Complete the installation

## Configuration File

`dns_cache.conf` file format:
```
# Comment lines start with #
domain1=ip1:port1
domain2=ip2:port2
```

## Usage

1. **Configure domain mappings**: Add domains and target addresses to `dns_cache.conf`
2. **Start the service**: Run `go run relay.go`
3. **Trust the CA certificate**: Add the generated `ca.crt` to your system's trusted list
4. **Access services**: Use the configured domains to access services, traffic will be relayed to target addresses

## Examples

### Example: Relay example.com

Configuration:
```
example.com=192.168.1.1:443
```

After starting the service, accessing `https://example.com` will relay traffic to `192.168.1.1:443`.

## Project Structure

```
https-relay/
├── relay.go          # Main program file
├── dns_cache.conf    # DNS cache configuration file (create this)
├── ca.crt            # CA certificate (auto-generated)
├── ca.key            # CA private key (auto-generated)
├── server.crt        # Server certificate (auto-generated)
├── server.key        # Server private key (auto-generated)
└── README.md         # Project documentation
```

## Notes

- The configuration file `dns_cache.conf` is not committed to Git
- Certificate files (`ca.crt`, `ca.key`, `server.crt`, `server.key`) are not committed to Git
- CA certificate only needs to be generated once, server certificate automatically regenerates based on configuration changes
- The server listens on port 443, requiring administrator privileges

## License

MIT License

## Contributing

Welcome to submit Issues and Pull Requests to improve this project!
