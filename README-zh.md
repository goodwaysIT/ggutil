# ggutil

**ggutil** 是一款面向企业级 Oracle GoldenGate (OGG) 多实例环境的命令行管理工具，支持多 OGG Home 并发批量操作，极大提升 OGG 日常运维、监控、配置和数据采集的自动化与效率。  
本工具完全开源，欢迎贡献与交流！

本工具的开发源于多年的实践经验，旨在解决单一服务器或集群上管理大量 OGG 软件实例的复杂性。它利用 GoldenGate 12c 及更高版本中提供的、支持大多数数据库的远程操作能力，提供了一个强大的集中式维护工具集。目前该工具已经过多年的生产环境检验，稳定支持 GoldenGate for Oracle, MySQL, DB2 LUW, DB2 z/OS, 和 Big Data。

- 开源地址: [https://github.com/goodwaysIT/ggutil](https://github.com/goodwaysIT/ggutil)

---

## 主要特性

- **多 OGG Home 并发管理**：支持配置多个 OGG Home 路径，所有命令自动并发处理，极大提升效率。
- **丰富的子命令体系**：
  - `tasks`：批量查询所有 OGG Home 下的 `SOURCEISTABLE` 任务，分组展示，支持无任务友好提示。
  - `mon`：批量获取所有 OGG Home 的版本、路径与 `info all` 运行态信息。
  - `info <process_name>`：查询指定进程（Extract/Replicat）在所有 OGG Home 下的详细信息。
  - `param <process_name>`：批量读取指定进程的参数文件内容。
  - `config`：批量展示所有 OGG Home 下主要进程的参数、配置表格。
  - `backup`：一键备份所有 OGG Home 的关键配置、日志、报告文件等，归档至时间戳目录并自动清理。
  - `stats <process_name>`：按总量、每日、每小时维度统计并展示指定进程的业务表操作数。
  - `collect <process_name>`：批量采集指定进程的 info/detail/showch/status 等所有相关文件，自动归档。
- **输出美观**：所有表格输出采用 gotabulate，结构清晰，适合直接用于运维报告。
- **并发与性能**：所有批量操作均为并发执行，充分利用多核资源。
- **强大的参数校验与错误提示**：所有参数、环境变量、路径均有详细校验与提示，支持 debug 模式输出详细日志。
- **高度可扩展与二次开发友好**：核心逻辑高度模块化，便于集成到更大运维平台或定制开发。
- **详细英文注释与专业代码风格**：便于团队协作、代码审查和国际化。

---

## 架构总览

ggutil CLI 并发工作流架构如下图所示：

![ggutil CLI 并发架构](./ggutil-cli-concurrent-workflow-architecture.svg)

---

## 安装与环境要求

- **操作系统**：Linux (建议 Oracle Linux/RedHat/CentOS)
- **依赖**：
  - Go 1.18 及以上
  - Oracle GoldenGate 已安装并配置（支持多 Home）
  - 依赖 Go 第三方库：`urfave/cli/v2`、`bndr/gotabulate`、`mholt/archiver/v3`
- **安装方式**：

  **方式一：直接下载安装（推荐）**
  
  无需编译和 Go 环境，直接前往 [Releases 页面](https://github.com/goodwaysIT/ggutil/releases) 下载对应平台的二进制包（如 `ggutil-x86_64`、`ggutil-arm64` 等），解压后赋予可执行权限即可：
  ```bash
  wget https://github.com/goodwaysIT/ggutil/releases/download/v1.0.0/ggutil-x86_64
  chmod +x ggutil-x86_64
  ./ggutil-x86_64 -h
  ```
  > ARM 架构请下载 `ggutil-arm64`，其它平台选择对应版本。

  **方式二：源码编译安装**
  
  需本地具备 Go 环境，适用于需要定制或二次开发场景。
  ```bash
  git clone https://github.com/goodwaysIT/ggutil.git
  cd ggutil
  go build -o ggutil main.go
  # 或直接 go run main.go <命令>
  ```

---

> 详细构建说明请见 [BUILD-zh.md](./BUILD-zh.md)

## 快速上手

### 1. 配置 OGG Home 路径

- 推荐通过环境变量 `GG_HOMES` 或 `-g/--gghomes` 参数指定多个 OGG Home（用英文逗号或分号分隔）：
  ```bash
  export GG_HOMES="/ogg1,/ogg2,/ogg3"
  ./ggutil tasks
  # 或
  ./ggutil -g "/ogg1,/ogg2" info extorcl
  ```

### 2. 查看所有命令与帮助

```bash
./ggutil -h
./ggutil <子命令> -h
```

### 3. 常用命令示例

- **显示帮助信息**

  ```bash
  $ ggutil -h
  NAME:
     ggutil - Oracle GoldenGate multi-instance management tool
              Open Source: https://github.com/goodwaysIT/ggutil

  USAGE:
     ggutil [global options] command [command options]

  COMMANDS:
     version  Show application version and open source repository
     tasks    List all OGG SOURCEISTABLE tasks under all homes.
     mon      Get version and path information for all OGG instances, print 'info all' results for each.
     info     Get information for OGG processes (iterates over all configured OGG Homes).
     param    Get parameter configuration for OGG processes (iterates over all configured OGG Homes).
     config   View process configuration details within OGG instances (iterates over all configured OGG Homes).
     backup   Backup configuration, log, report files, etc., for OGG instances (iterates over all configured OGG Homes).
     stats    View statistics for a specific OGG process (total, daily, hourly) (iterates over all configured OGG Homes).
     collect  Collect information for a specific OGG process (info, infodetail, showch, status) (iterates over all configured OGG Homes).
     help, h  Shows a list of commands or help for one command

  GLOBAL OPTIONS:
     --gghomes value, -g value  Specify one or more OGG Home paths, comma-separated. If not specified, attempts to read from GG_HOMES environment variable. [$GG_HOMES]
     --debug                    Enable debug output (show errors, warnings, exceptions) (default: false)
     --help, -h                 show help
  ```

- **监控所有 OGG 实例 (`mon`)**

  ```bash
  $ ggutil mon

  ==== Home: /acfsogg/oggb, OGG for Big Data, Version 19.1.0.0.200714 OGGCORE_19.1.0.0.0OGGBP_PLATFORMS_200628.2141

  Program     Status      Group       Lag at Chkpt  Time Since Chkpt

  MANAGER     RUNNING
  REPLICAT    RUNNING     RKAFKA      00:00:00      00:00:05


  --------------------------------------------------------------------------------


  ==== Home: /acfsogg/oggm, OGG for MySQL, Version 19.1.0.0.230418 OGGCORE_19.1.0.0.0OGGBP_PLATFORMS_230413.1325

  Program     Status      Group       Lag at Chkpt  Time Since Chkpt

  MANAGER     RUNNING
  REPLICAT    RUNNING     REPMYSQL    00:00:00      00:00:05


  --------------------------------------------------------------------------------


  ==== Home: /acfsogg/oggp, OGG for PostgreSQL, Version 21.14.0.0.0 OGGCORE_21.14.0.0.0OGGRU_PLATFORMS_240404.1108

  Program     Status      Group       Lag at Chkpt  Time Since Chkpt

  MANAGER     RUNNING
  EXTRACT     RUNNING     EXT_PG      00:00:00      00:00:07
  REPLICAT    RUNNING     REP_PG      00:00:00      00:00:01


  --------------------------------------------------------------------------------


  ==== Home: /acfsogg/oggo, OGG for Oracle, Version 19.1.0.0.4 OGGCORE_19.1.0.0.0_PLATFORMS_191017.1054_FBO

  Program     Status      Group       Lag at Chkpt  Time Since Chkpt

  MANAGER     RUNNING
  EXTRACT     RUNNING     DPORA       00:00:00      00:00:06
  EXTRACT     RUNNING     EXTORA      00:00:00      00:00:03


  --------------------------------------------------------------------------------
  ```

- **查看进程配置 (`config`)**

  ```bash
  $ ggutil config

  ==== Home: /acfsogg/oggm, OGG for MySQL, Version 19.1.0.0.230418 OGGCORE_19.1.0.0.0OGGBP_PLATFORMS_230413.1325

  Program    Status     Group      TabNo(prm) TabNo(rpt) Source                                       Target
  ---------- ---------- ---------- ---------- ---------- -------------------------------------------- --------------------------------------------
  REPLICAT   RUNNING    REPMYSQL   1          1          ./dirdat/my000000000(4539575)                ogg_target_db@mysqldb:3306


  ==== Home: /acfsogg/oggb, OGG for Big Data, Version 19.1.0.0.200714 OGGCORE_19.1.0.0.0OGGBP_PLATFORMS_200628.2141

  Program    Status     Group      TabNo(prm) TabNo(rpt) Source                                       Target
  ---------- ---------- ---------- ---------- ---------- -------------------------------------------- --------------------------------------------
  REPLICAT   RUNNING    RKAFKA     1          1          AdapterExamples/trail/tr000000000(5660)      Kafka


  ==== Home: /acfsogg/oggp, OGG for PostgreSQL, Version 21.14.0.0.0 OGGCORE_21.14.0.0.0OGGRU_PLATFORMS_240404.1108

  Program    Status     Group      TabNo(prm) TabNo(rpt) Source                                       Target
  ---------- ---------- ---------- ---------- ---------- -------------------------------------------- --------------------------------------------
  EXTRACT    RUNNING    EXT_PG     1          1          testpdb,                                     ./dirdat/pg000000001(1719914)
  REPLICAT   RUNNING    REP_PG     1          1          ./dirdat/pg000000001(1719914)                testpdb


  ==== Home: /acfsogg/oggo, OGG for Oracle, Version 19.1.0.0.4 OGGCORE_19.1.0.0.0_PLATFORMS_191017.1054_FBO

  Program    Status     Group      TabNo(prm) TabNo(rpt) Source                                       Target
  ---------- ---------- ---------- ---------- ---------- -------------------------------------------- --------------------------------------------
  EXTRACT*p  RUNNING    DPORA      1          1          /acfsogg/oggo/dirdat/or000000001(2700)       mysqldb,
  EXTRACT    RUNNING    EXTORA     1          1          oracledb/orcl,                               ./dirdat/or000000001(2700)
  ```

- **查看参数文件 (`param`)**

  ```bash
  $ ggutil param extora

  ==== OGG Process [ EXTORA ] Under Home: [ /acfsogg/oggo ] ====

  Param file [ /acfsogg/oggo/dirprm/extora.prm ] content for 'EXTORA':

  EXTRACT extora
  USERID c##ogguser@oracledb/orcl, PASSWORD ogguser2025
  FETCHOPTIONS FETCHPKUPDATECOLS
  discardfile ./dirrpt/extora.dsc, append, megabytes 1000
  exttrail ./dirdat/or
  sourcecatalog orclpdb
  DDL INCLUDE MAPPED
  TABLE TUSER.TTAB1;
  ```

- **查看进程统计信息 (`stats`)**

  ```bash
  $ ggutil stats rep_pg

  ==== OGG Process [ REP_PG ] Under Home: [ /acfsogg/oggp ] ====

  ========================================[total stats]========================================

  *** Total statistics since 2025-07-05 14:53:19 ***
  +-------------------------------+------------+------------+------------+------------+------------+-------------+---------------+
  | Table Name                    | Insert     | Updates    | Befores    | Deletes    | Upserts    | Discards    | Operations    |
  +===============================+============+============+============+============+============+=============+===============+
  | source_schema.source_table    | 5000.00    | 4000.00    |            | 3000.00    | 0.00       | 0.00        | 12000.00      |
  +-------------------------------+------------+------------+------------+------------+------------+-------------+---------------+

  ========================================[daily stats]========================================

  *** Daily statistics since 2025-07-05 14:53:19 ***
  +-------------------------------+------------+------------+------------+------------+------------+-------------+---------------+
  | Table Name                    | Insert     | Updates    | Befores    | Deletes    | Upserts    | Discards    | Operations    |
  +===============================+============+============+============+============+============+=============+===============+
  | source_schema.source_table    | 5000.00    | 4000.00    |            | 3000.00    | 0.00       | 0.00        | 12000.00      |
  +-------------------------------+------------+------------+------------+------------+------------+-------------+---------------+

  ========================================[hourly stats/sec]========================================

  *** Hourly statistics since 2025-07-05 14:53:19 ***
  +-------------------------------+-----------+------------+------------+------------+------------+-------------+---------------+
  | Table Name                    | Insert    | Updates    | Befores    | Deletes    | Upserts    | Discards    | Operations    |
  +===============================+===========+============+============+============+============+=============+===============+
  | source_schema.source_table    | 0.01      | 0.01       |            | 0.01       | 0.00       | 0.00        | 0.02          |
  +-------------------------------+-----------+------------+------------+------------+------------+-------------+---------------+
  ```

- **备份关键文件 (`backup`)**

  ```bash
  $ ggutil backup

  Please refer to gz file /tmp/oggbackup_xugu01_20250711_195536.tar.gz
  ```

---

## 设计理念与架构说明

- **并发优先**：所有 OGG Home 操作均采用 goroutine 并发，极大提升批量处理性能。
- **输出友好**：所有命令输出均为结构化表格，便于人工阅读和自动化采集。
- **代码可维护性**：所有核心模块、工具函数均有详细英文注释，便于团队协作和二次开发。
- **异常健壮**：参数、环境变量、路径等均有详细校验与错误提示，debug 模式可追踪详细日志。
- **安全性**：所有归档、删除、文件操作均有异常处理，避免误删、误覆盖。

---

## 适用场景

- 多 OGG Home 环境的日常批量运维与监控
- OGG 配置、日志、报告文件的自动归档与备份
- OGG 进程状态、参数、统计信息的批量采集与报表输出
- 适合 DBA、数据同步平台、自动化运维团队使用

---

## 贡献与支持

- 欢迎 issue、PR、建议与交流
- 企业级定制开发与技术支持请联系作者或通过 github issue 留言

---

## 开源协议

本项目采用 MIT 协议，详见 [LICENSE](LICENSE)。

---

如需功能扩展或企业定制，请联系维护者或提交 issue！
