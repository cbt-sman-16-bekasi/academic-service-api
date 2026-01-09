# Dapodik Proxy Server

Proxy server untuk mengatasi masalah IP Address yang didaftarkan pada Web Service Dapodik berbeda dengan IP aplikasi utama.

## Problem

Dapodik Web Service memvalidasi IP Address yang melakukan request. Jika IP tidak sesuai dengan yang terdaftar, akan muncul error:

```json
{"success":false,"http_code":403,"status_code":"Forbidden","message":"IP Address yang didaftarkan pada Web Service Dapodik berbeda"}
```

## Solution

Install proxy ini di komputer/server yang IP-nya sudah terdaftar di Dapodik. Kemudian aplikasi utama akan request melalui proxy ini.

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│  Academic API   │ ───► │  Dapodik Proxy  │ ───► │  Dapodik Server │
│  (Any Server)   │      │  (Registered IP)│      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

## Installation (Windows)

1. Download file berikut ke komputer Windows yang IP-nya terdaftar di Dapodik:
   - `dapodik-proxy-windows-amd64.exe` (untuk Windows 64-bit)
   - `config.yaml`

2. Edit `config.yaml` sesuai kebutuhan:
   ```yaml
   server:
     port: 8888  # Port yang digunakan
   
   security:
     api_key: "your-secret-key"  # Optional: untuk keamanan
   ```

3. Jalankan proxy:
   ```cmd
   dapodik-proxy-windows-amd64.exe
   ```

   Atau dengan config file custom:
   ```cmd
   dapodik-proxy-windows-amd64.exe -config C:\path\to\config.yaml
   ```

   Atau override port:
   ```cmd
   dapodik-proxy-windows-amd64.exe -port 9999
   ```

## Usage

### Method 1: Query Parameter (Recommended)

```
GET http://proxy-server:8888/proxy?target=http://dapodik-server:5774/WebService/getSekolah?npsn=20275048
Authorization: Bearer <dapodik_token>
```

### Method 2: Header

```
GET http://proxy-server:8888/proxy
X-Proxy-Target: http://dapodik-server:5774/WebService/getSekolah?npsn=20275048
Authorization: Bearer <dapodik_token>
```

### Method 3: Path-based (dengan base_url di config)

Jika `dapodik.base_url` di config diset ke `http://dapodik-server:5774`:

```
GET http://proxy-server:8888/proxy/WebService/getSekolah?npsn=20275048
Authorization: Bearer <dapodik_token>
```

## Security

### API Key Protection

Set `security.api_key` di config untuk require authentication:

```yaml
security:
  api_key: "my-secret-proxy-key"
```

Client harus kirim header:
```
X-API-Key: my-secret-proxy-key
```

### IP Whitelist

Set `server.allowed_ips` untuk membatasi siapa yang bisa akses proxy:

```yaml
server:
  allowed_ips:
    - "192.168.1.100"
    - "10.0.0.1"
```

## Health Check

```
GET http://proxy-server:8888/health
```

Response:
```json
{"status":"ok","service":"dapodik-proxy"}
```

## Building from Source

```bash
# Install dependencies
go mod tidy

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o dapodik-proxy.exe main.go

# Or use build script
chmod +x build.sh
./build.sh
```

## Running as Windows Service

Untuk menjalankan sebagai Windows Service, gunakan tools seperti [NSSM](https://nssm.cc/):

```cmd
nssm install DapodikProxy C:\path\to\dapodik-proxy.exe
nssm set DapodikProxy AppDirectory C:\path\to
nssm set DapodikProxy AppParameters -config config.yaml
nssm start DapodikProxy
```

## Configuration Reference

| Setting | Description | Default |
|---------|-------------|---------|
| `server.port` | Port to listen on | 8888 |
| `server.read_timeout` | Request timeout (seconds) | 60 |
| `server.write_timeout` | Response timeout (seconds) | 120 |
| `server.allowed_ips` | List of allowed client IPs | [] (all) |
| `dapodik.base_url` | Default Dapodik server URL | "" |
| `logging.level` | Log level (debug/info/warn/error) | info |
| `logging.pretty` | Human readable logs | true |
| `security.api_key` | Required API key for clients | "" |

## Troubleshooting

### Port sudah digunakan
```
Error: listen tcp :8888: bind: address already in use
```
Solution: Ganti port di config atau gunakan flag `-port 9999`

### Timeout
```
Error: Proxy request failed: context deadline exceeded
```
Solution: Naikkan `write_timeout` di config

### Connection refused
Pastikan:
1. Dapodik server bisa diakses dari komputer proxy
2. Firewall tidak memblokir port proxy
3. URL target benar

