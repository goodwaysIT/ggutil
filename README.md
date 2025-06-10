[简体中文](./README-zh.md) | [日本語](./README-ja.md)

# ggutil

**ggutil** is a command-line management tool designed for enterprise-level Oracle GoldenGate (OGG) multi-instance environments. It supports concurrent batch operations across multiple OGG Homes, significantly enhancing the automation and efficiency of daily OGG operations, monitoring, configuration, and data collection. This tool is fully open-source; contributions and discussions are welcome!

- Open Source Repository: [https://github.com/goodwaysIT/ggutil](https://github.com/goodwaysIT/ggutil)

---

## Key Features

- **Concurrent Management of Multiple OGG Homes**: Supports configuration of multiple OGG Home paths, with all commands processed concurrently to greatly improve efficiency.
- **Rich Subcommand System**:
  - `tasks`: Batch query `SOURCEISTABLE` tasks across all OGG Homes, grouped display, with friendly prompts for no tasks.
  - `mon`: Batch retrieve version, path, and `info all` runtime status for all OGG Homes.
  - `info <process_name>`: Query detailed information for a specified process (Extract/Replicat) across all OGG Homes.
  - `param <process_name>`: Batch read the parameter file content for a specified process.
  - `config`: Batch display parameters and configuration tables for major processes in all OGG Homes.
  - `backup`: One-click backup of key configurations, logs, report files, etc., for all OGG Homes, archived to a timestamped directory and automatically cleaned up.
  - `stats <process_name>`: Collect and display business table operation counts for a specified process by total, daily, and hourly dimensions.
  - `collect <process_name>`: Batch collect all relevant files (info/detail/showch/status, etc.) for a specified process, automatically archived.
- **Elegant Output**: All table outputs use `gotabulate` for clear structure, suitable for direct use in operational reports.
- **Concurrency and Performance**: All batch operations are executed concurrently, fully utilizing multi-core resources.
- **Robust Parameter Validation & Error Prompts**: Detailed validation and prompts for all parameters, environment variables, and paths, with debug mode for detailed log output.
- **Highly Extensible & Developer-Friendly**: Core logic is highly modular, facilitating integration into larger operations platforms or custom development.
- **Detailed English Comments & Professional Code Style**: Facilitates team collaboration, code review, and internationalization.

---

## Installation and Requirements

- **Operating System**: Linux (Oracle Linux/RedHat/CentOS recommended)
- **Dependencies**:
  - Go 1.18 or higher
  - Oracle GoldenGate installed and configured (multi-Home support)
  - Third-party Go libraries: `urfave/cli/v2`, `bndr/gotabulate`, `mholt/archiver/v3`
- **Installation Methods**:

  **Method 1: Direct Download (Recommended)**

  No compilation or Go environment needed. Directly go to the [Releases Page](https://github.com/goodwaysIT/ggutil/releases) to download the binary package for your platform (e.g., `ggutil-x86_64`, `ggutil-arm64`), then grant executable permissions:
  ```bash
  wget https://github.com/goodwaysIT/ggutil/releases/download/v1.0.0/ggutil-x86_64
  chmod +x ggutil-x86_64
  ./ggutil-x86_64 -h
  ```
  > For ARM architecture, download `ggutil-arm64`. Choose the corresponding version for other platforms.

  **Method 2: Compile from Source**

  Requires a local Go environment. Suitable for customization or secondary development scenarios.
  ```bash
  git clone https://github.com/goodwaysIT/ggutil.git
  cd ggutil
  go build -o ggutil main.go
  # Or run directly: go run main.go <command>
  ```

---

> For detailed build instructions, see [BUILD.md](./BUILD.md)

## Quick Start

### 1. Configure OGG Home Paths

- Recommended: Specify multiple OGG Homes (separated by commas or semicolons) via the `GG_HOMES` environment variable or the `-g/--gghomes` parameter:
  ```bash
  export GG_HOMES="/ogg1,/ogg2,/ogg3"
  ./ggutil tasks
  # Or
  ./ggutil -g "/ogg1,/ogg2" info extorcl
  ```

### 2. View All Commands and Help

```bash
./ggutil -h
./ggutil <subcommand> -h
```

### 3. Common Command Examples

- Query table-level tasks in all OGG Homes
  ```bash
  ./ggutil tasks
  ```
- View OGG version and `info all` for all Homes
  ```bash
  ./ggutil mon
  ```
- Query detailed information for a specific process
  ```bash
  ./ggutil info extorcl
  ```
- View process parameter file content
  ```bash
  ./ggutil param extorcl
  ```
- Backup all key configuration/log/report files
  ```bash
  ./ggutil backup
  ```
- Collect statistics on business table operations for a process
  ```bash
  ./ggutil stats extorcl
  ```
- Collect and archive all relevant files for a process
  ```bash
  ./ggutil collect extorcl
  ```

---

## Design Philosophy and Architecture

- **Concurrency First**: All OGG Home operations use goroutines for concurrency, greatly improving batch processing performance.
- **User-Friendly Output**: All command outputs are structured tables, easy for manual reading and automated collection.
- **Code Maintainability**: All core modules and utility functions have detailed English comments, facilitating team collaboration and secondary development.
- **Robustness**: Detailed validation and error prompts for parameters, environment variables, paths, etc. Debug mode allows tracking detailed logs.
- **Safety**: All archiving, deletion, and file operations include exception handling to prevent accidental deletion or overwriting.

---

## Use Cases

- Daily batch operations and monitoring for multi-OGG Home environments.
- Automated archiving and backup of OGG configuration, log, and report files.
- Batch collection and report output of OGG process status, parameters, and statistical information.
- Suitable for DBAs, data synchronization platforms, and automated operations teams.

---

## Contribution and Support

- Issues, PRs, suggestions, and discussions are welcome.
- For enterprise-level custom development and technical support, please contact the author or leave a message via GitHub issues.

---

## Open Source License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

For feature extensions or enterprise customization, please contact the maintainers or submit an issue!
