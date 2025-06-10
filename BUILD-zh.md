# 构建指南（Build Instructions, 中文）

本项目支持在主流 Linux x86_64 及 ARM (aarch64) 架构下编译和运行。以下为标准构建步骤。

## 1. 环境准备

- 操作系统：Linux（建议 Oracle Linux、CentOS、RedHat、Ubuntu 等主流发行版）
- Go 版本：建议 Go 1.18 及以上
- 已安装并配置好 Oracle GoldenGate（多 Home 支持）

## 2. 获取源码

```bash
git clone https://github.com/goodwaysIT/ggutil.git
cd ggutil
```

## 3. 依赖安装

Go 会自动拉取依赖，首次编译时无需手动安装。若需提前拉取：

```bash
go mod tidy
```

## 4. 本地编译（默认当前平台架构）

```bash
go build -o ggutil main.go
```

## 5. 跨平台编译

### 构建 Linux x86_64（64位 Intel/AMD）

```bash
GOOS=linux GOARCH=amd64 go build -o ggutil_x86_64 main.go
```

### 构建 Linux ARM64（64位 ARM，如鲲鹏、飞腾、树莓派等）

```bash
GOOS=linux GOARCH=arm64 go build -o ggutil_arm64 main.go
```

> 说明：如需其它架构（如 32 位 ARM），可将 `GOARCH` 设置为 `arm`，并根据实际需求调整 `GOARM` 环境变量。

## 6. 验证构建

```bash
file ggutil*
# 应显示对应的 ELF 64-bit LSB executable, x86-64 或 aarch64
./ggutil_x86_64 -h
./ggutil_arm64 -h
```

## 7. 常见问题

- **依赖拉取失败**：请检查网络或配置 Go 代理（如 `GOPROXY=https://goproxy.cn`）。
- **权限问题**：如遇执行权限问题，使用 `chmod +x ggutil*`。
- **OGG 相关命令不可用**：请确保目标机器已正确安装并授权 Oracle GoldenGate。

---

如有特殊平台、静态编译、容器化需求，欢迎通过 issue 或 PR 反馈！
