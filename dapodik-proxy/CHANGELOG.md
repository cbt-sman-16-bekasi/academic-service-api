# Changelog

All notable changes to Dapodik Proxy will be documented in this file.

## [Unreleased]

### Added
- **Windows Installer (installer.exe)**: Automated installer untuk Windows dengan fitur:
  - Install/Uninstall service otomatis menggunakan Windows Service Manager
  - Konfigurasi firewall otomatis (open port 8888)
  - Auto-start saat Windows boot
  - Interactive menu untuk manage service
- **GitHub Actions CI/CD**: Automated build workflow
  - Build untuk Windows (32/64-bit), Linux, macOS
  - Automatic release creation dengan tag versioning
  - Pre-packaged distribution files (.zip untuk Windows, .tar.gz untuk Linux)
- **Enhanced build script**: 
  - Build installer bersamaan dengan proxy
  - Create distribution packages otomatis
  - Support macOS ARM64 (Apple Silicon)

### Changed
- Updated README.md dengan installation instructions untuk automated installer
- Enhanced build.sh untuk compile dan package installer

## [1.0.0] - Initial Release

### Added
- Basic HTTP proxy server untuk Dapodik Web Service
- Support multiple target URL methods (query param, header, path-based)
- Configuration via YAML file
- Security features (API Key, IP Whitelist)
- CORS support
- Health check endpoint
- Logging dengan zerolog
- Cross-platform builds (Windows, Linux, macOS)
