[构建指南](./BUILD-zh.md) | [ビルドガイド](./BUILD-ja.md)

# Build Guide

This project supports compilation and execution on mainstream Linux x86_64 and ARM (aarch64) architectures. The following are standard build steps.

## 1. Environment Preparation

- Operating System: Linux (Oracle Linux, CentOS, RedHat, Ubuntu, etc., mainstream distributions recommended)
- Go Version: Go 1.18 or higher recommended
- Oracle GoldenGate installed and configured (multi-Home support)

## 2. Get Source Code

```bash
git clone https://github.com/goodwaysIT/ggutil.git
cd ggutil
```

## 3. Install Dependencies

Go will automatically fetch dependencies during the first compilation; manual installation is not required. To fetch them beforehand:

```bash
go mod tidy
```

## 4. Local Compilation (Default to Current Platform Architecture)

```bash
go build -o ggutil main.go
```

## 5. Cross-Platform Compilation

### Build for Linux x86_64 (64-bit Intel/AMD)

```bash
GOOS=linux GOARCH=amd64 go build -o ggutil_x86_64 main.go
```

### Build for Linux ARM64 (64-bit ARM, e.g., Kunpeng, Phytium, Raspberry Pi)

```bash
GOOS=linux GOARCH=arm64 go build -o ggutil_arm64 main.go
```

> Note: For other architectures (e.g., 32-bit ARM), set `GOARCH` to `arm` and adjust the `GOARM` environment variable as needed.

## 6. Verify Build

```bash
file ggutil*
# Should display corresponding ELF 64-bit LSB executable, x86-64 or aarch64
./ggutil_x86_64 -h
./ggutil_arm64 -h
```

## 7. Common Issues

- **Dependency Fetch Failure**: Check your network or configure a Go proxy (e.g., `GOPROXY=https://goproxy.cn`).
- **Permission Issues**: If you encounter execution permission problems, use `chmod +x ggutil*`.
- **OGG Related Commands Unavailable**: Ensure Oracle GoldenGate is correctly installed and authorized on the target machine.

---

For special platforms, static compilation, or containerization needs, feel free to provide feedback via issues or PRs!
