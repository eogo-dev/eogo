# Windows Development Guide

This guide helps you set up and run ZGO on Windows.

## 📋 Prerequisites

1. **Go 1.22+**: Download from [go.dev/dl](https://go.dev/dl/)
2. **Git**: Download from [git-scm.com](https://git-scm.com/download/win)
3. **PostgreSQL** (optional): For database features
4. **Make Alternative**: We provide `make.bat` and `make.ps1` scripts

## 🚀 Quick Start

### Option 1: Using PowerShell (Recommended)

```powershell
# Clone the repository
git clone https://github.com/zgiai/zgo.git
cd zgo

# Copy environment file
copy .env.example .env

# Setup development environment
.\make.ps1 setup

# Build and install globally
.\make.ps1 install

# Run the server
.\make.ps1 dev
```

### Option 2: Using Command Prompt

```cmd
# Clone the repository
git clone https://github.com/zgiai/zgo.git
cd zgo

# Copy environment file
copy .env.example .env

# Setup development environment
make.bat setup

# Build and install globally
make.bat install

# Run the server
make.bat dev
```

### Option 3: Direct Go Commands

```cmd
# Build CLI
go build -o zgo.exe cmd/zgo/main.go

# Run server
go run cmd/server/main.go
```

## 🛠️ Common Commands

### PowerShell

```powershell
.\make.ps1 build        # Build the CLI tool
.\make.ps1 build-server # Build the server
.\make.ps1 test         # Run tests
.\make.ps1 lint         # Code linting
.\make.ps1 wire         # Generate DI code
.\make.ps1 docs         # Generate API docs
.\make.ps1 air          # Run with hot-reload
.\make.ps1 help         # Show all commands
```

### Command Prompt

```cmd
make.bat build        # Build the CLI tool
make.bat build-server # Build the server
make.bat test         # Run tests
make.bat lint         # Code linting
make.bat wire         # Generate DI code
make.bat docs         # Generate API docs
make.bat help         # Show all commands
```

## ⚙️ Configuration

Edit `.env` file to configure your environment:

```env
# Server
APP_NAME=ZGO
APP_ENV=development
APP_PORT=8025

# Database
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=zgo
DB_USERNAME=postgres
DB_PASSWORD=your_password

# JWT
JWT_SECRET=your_secret_key_here
JWT_EXPIRATION=3600
```

## 🔧 IDE Setup

### Visual Studio Code

1. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
2. Open the project folder
3. Press `F5` to debug

Recommended extensions:
- Go (golang.go)
- Go Test Explorer
- GitLens

### GoLand

1. Open the project
2. Ensure Go SDK is configured (File → Settings → Go → GOROOT)
3. Right-click `cmd/server/main.go` → Run

## 🐛 Troubleshooting

### Issue: Command not found after install

**Solution**: Add `%GOPATH%\bin` to your PATH:

```powershell
# PowerShell (run as Administrator)
$env:Path += ";$env:GOPATH\bin"
[Environment]::SetEnvironmentVariable("Path", $env:Path, "User")
```

```cmd
# Command Prompt (run as Administrator)
setx PATH "%PATH%;%GOPATH%\bin"
```

### Issue: Wire not working

**Solution**: Install Wire manually:

```cmd
go install github.com/google/wire/cmd/wire@latest
```

### Issue: PowerShell execution policy

If you get "cannot be loaded because running scripts is disabled", run:

```powershell
# Run PowerShell as Administrator
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Issue: golangci-lint not working on Windows

**Solution**: Download pre-built binary from [golangci-lint releases](https://github.com/golangci/golangci-lint/releases) and add to PATH.

### Issue: Database connection failed

**Solution**: 
1. Ensure PostgreSQL is running
2. Check `.env` configuration
3. Create database: 
   ```sql
   CREATE DATABASE zgo;
   ```

## 🏃‍♂️ Development Workflow

### Standard Development

```powershell
# 1. Setup (first time only)
.\make.ps1 setup

# 2. Start development server with hot-reload
.\make.ps1 air

# 3. In another terminal, run tests
.\make.ps1 test

# 4. Before committing
.\make.ps1 lint
.\make.ps1 test
```

### Building for Production

```powershell
# Build both CLI and server
.\make.ps1 build-all

# This creates:
# - zgo.exe (CLI tool)
# - server.exe (HTTP server)
```

### Running Migrations

```cmd
# Build CLI first
.\make.ps1 build

# Run migrations
.\zgo.exe migrate

# Rollback
.\zgo.exe migrate:rollback
```

## 📝 Notes

- **Line Endings**: Windows uses CRLF, while Linux/Mac use LF. Configure Git:
  ```cmd
  git config --global core.autocrlf true
  ```

- **Path Separators**: Go handles both `/` and `\`, but prefer `/` in code for cross-platform compatibility.

- **Environment Variables**: Use `.env` file instead of system environment variables for easier configuration.

## 🔗 Additional Resources

- [Go on Windows](https://go.dev/doc/install/windows)
- [Git for Windows](https://gitforwindows.org/)
- [PostgreSQL on Windows](https://www.postgresql.org/download/windows/)
- [VS Code Go Development](https://code.visualstudio.com/docs/languages/go)

## 💡 Tips

1. **Use Windows Terminal** for a better command-line experience
2. **WSL2** is an alternative if you prefer Linux commands
3. **Docker Desktop** can be used for containerized development
4. Enable **Developer Mode** in Windows 11 for better dev experience

---

Need help? Check the main [README.md](README.md) or open an issue on GitHub.
