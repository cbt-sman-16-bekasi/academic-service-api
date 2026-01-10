# Quick Start Guide - Dapodik Proxy

## 🚀 Instalasi Windows (Recommended)

### Prasyarat
- Windows 7/8/10/11 (32-bit atau 64-bit)
- Komputer harus IP-nya sudah terdaftar di Dapodik Web Service
- Hak akses Administrator

### Langkah-langkah

1. **Download Package**
   - Kunjungi [Releases Page](https://github.com/YOUR_ORG/YOUR_REPO/releases)
   - Download `dapodik-proxy-windows-amd64.zip` (untuk 64-bit)
   - Extract file ZIP ke folder yang diinginkan (misalnya `C:\DapodikProxy`)

2. **Konfigurasi (Opsional)**
   - Edit file `config.yaml` jika perlu mengubah port atau security settings
   - Default port: **8888**

3. **Install Service**
   - Right-click `installer.exe`
   - Pilih **"Run as Administrator"**
   - Pilih menu **[1] Install Service**
   - Tunggu sampai proses selesai

4. **Verifikasi**
   - Buka browser
   - Akses: `http://localhost:8888/health`
   - Jika muncul `{"status":"ok","service":"dapodik-proxy"}` → **Berhasil!** ✅

## 🔧 Penggunaan

### Dari Aplikasi Academic Service

Edit konfigurasi Academic Service untuk menggunakan proxy:

```yaml
dapodik:
  use_proxy: true
  proxy_url: "http://IP_KOMPUTER_PROXY:8888/proxy"
  base_url: "http://IP_DAPODIK_SERVER:5774"
```

### Contoh Request API

**Method 1: Query Parameter (Recommended)**
```bash
curl -H "Authorization: Bearer TOKEN_DAPODIK" \
  "http://IP_PROXY:8888/proxy?target=http://IP_DAPODIK:5774/WebService/getSekolah?npsn=20275048"
```

**Method 2: Header**
```bash
curl -H "Authorization: Bearer TOKEN_DAPODIK" \
     -H "X-Proxy-Target: http://IP_DAPODIK:5774/WebService/getSekolah?npsn=20275048" \
  "http://IP_PROXY:8888/proxy"
```

## 🛠️ Manajemen Service

Jalankan `installer.exe` sebagai Administrator untuk:

| Menu | Fungsi |
|------|--------|
| [1]  | Install Service (termasuk firewall & autostart) |
| [2]  | Uninstall Service (hapus service & firewall rule) |
| [3]  | Start Service |
| [4]  | Stop Service |
| [5]  | Cek Status Service |

## ❓ Troubleshooting

### Service gagal start
```cmd
# Cek status dengan perintah:
sc query DapodikProxy

# Atau gunakan installer.exe → pilih [5]
```

### Port 8888 sudah digunakan
- Edit `config.yaml`, ubah port ke yang lain (misal 9999)
- Uninstall service lama → Install ulang

### Firewall memblokir koneksi
- Pastikan Windows Firewall rule sudah dibuat (otomatis saat install)
- Atau buat manual:
  ```cmd
  netsh advfirewall firewall add rule name="DapodikProxy-8888" dir=in action=allow protocol=TCP localport=8888
  ```

### Tidak bisa akses dari komputer lain
- Pastikan firewall di komputer proxy tidak memblokir koneksi dari jaringan
- Test koneksi: `telnet IP_PROXY 8888` dari komputer lain

## 📋 Uninstall

1. Jalankan `installer.exe` sebagai Administrator
2. Pilih menu **[2] Uninstall Service**
3. Hapus folder instalasi jika sudah tidak diperlukan

## 🆘 Butuh Bantuan?

- Baca dokumentasi lengkap: [README.md](README.md)
- Report issue: [GitHub Issues](https://github.com/YOUR_ORG/YOUR_REPO/issues)
