# ggutil

**ggutil** は、エンタープライズレベルの Oracle GoldenGate (OGG) マルチインスタンス環境向けのコマンドライン管理ツールです。複数のOGGホームにまたがる同時バッチ操作をサポートし、日常のOGG運用、監視、設定、データ収集の自動化と効率を大幅に向上させます。このツールは完全にオープンソースであり、貢献と議論を歓迎します！

長年の実務経験に基づき開発された `ggutil` は、単一のサーバーまたはクラスター上で多数のOGGソフトウェアインスタンスを管理する際の複雑さに対応します。GoldenGate 12c以降のバージョンで利用可能になったリモート操作機能（ほとんどのデータベースをサポート）を活用し、強力で一元化されたメンテナンスツールキットを提供します。このツールは長年の本番環境での使用実績があり、GoldenGate for Oracle, MySQL, DB2 LUW, DB2 z/OS, Big Dataを安定してサポートしています。

- オープンソースリポジトリ: [https://github.com/goodwaysIT/ggutil](https://github.com/goodwaysIT/ggutil)

---

## 主な特徴

- **複数OGGホームの同時管理**: 複数のOGGホームパスの設定をサポートし、すべてのコマンドが自動的に同時処理されるため、効率が大幅に向上します。
- **豊富なサブコマンド体系**:
  - `tasks`: すべてのOGGホーム下の `SOURCEISTABLE` タスクを一括クエリし、グループ化して表示します。タスクがない場合はフレンドリーなプロンプトを表示します。
  - `mon`: すべてのOGGホームのバージョン、パス、および `info all` 実行時情報を一括取得します。
  - `info <プロセス名>`: 指定されたプロセス（Extract/Replicat）のすべてのOGGホームでの詳細情報をクエリします。
  - `param <プロセス名>`: 指定されたプロセスのパラメータファイルの内容を一括読み取りします。
  - `config`: すべてのOGGホームの主要プロセスのパラメータと設定テーブルを一括表示します。
  - `backup`: すべてのOGGホームの主要な設定、ログ、レポートファイルなどをワンクリックでバックアップし、タイムスタンプ付きディレクトリにアーカイブして自動的にクリーンアップします。
  - `stats <プロセス名>`: 指定されたプロセスの業務テーブル操作数を合計、日別、時間別で統計し表示します。
  - `collect <プロセス名>`: 指定されたプロセスのすべての関連ファイル（info/detail/showch/statusなど）を一括収集し、自動的にアーカイブします。
- **美しい出力**: すべてのテーブル出力は `gotabulate` を使用し、明確な構造で、運用レポートに直接使用するのに適しています。
- **並行性とパフォーマンス**: すべてのバッチ操作は並行して実行され、マルチコアリソースを最大限に活用します。
- **強力なパラメータ検証とエラープロンプト**: すべてのパラメータ、環境変数、パスには詳細な検証とプロンプトがあり、デバッグモードでは詳細なログ出力がサポートされます。
- **高度な拡張性と二次開発フレンドリー**: コアロジックは高度にモジュール化されており、より大きな運用プラットフォームへの統合やカスタム開発が容易です。
- **詳細な英語コメントとプロフェッショナルなコードスタイル**: チームコラボレーション、コードレビュー、国際化を促進します。

---

## インストールと環境要件

- **オペレーティングシステム**: Linux (Oracle Linux/RedHat/CentOS推奨)
- **依存関係**:
  - Go 1.18 以上
  - Oracle GoldenGate がインストールされ設定済みであること（マルチホームサポート）
  - サードパーティGoライブラリ: `urfave/cli/v2`、`bndr/gotabulate`、`mholt/archiver/v3`
- **インストール方法**:

  **方法1：直接ダウンロード（推奨）**

  コンパイルやGo環境は不要です。[リリースベージ](https://github.com/goodwaysIT/ggutil/releases) からプラットフォームに対応するバイナリパッケージ（例：`ggutil-x86_64`、`ggutil-arm64`など）を直接ダウンロードし、解凍後に実行権限を付与します。
  ```bash
  wget https://github.com/goodwaysIT/ggutil/releases/download/v1.0.0/ggutil-x86_64
  chmod +x ggutil-x86_64
  ./ggutil-x86_64 -h
  ```
  > ARMアーキテクチャの場合は `ggutil-arm64` をダウンロードしてください。他のプラットフォームの場合は対応するバージョンを選択してください。

  **方法2：ソースからのコンパイルインストール**

  ローカルにGo環境が必要です。カスタマイズや二次開発のシナリオに適しています。
  ```bash
  git clone https://github.com/goodwaysIT/ggutil.git
  cd ggutil
  go build -o ggutil main.go
  # または直接実行: go run main.go <コマンド>
  ```

---

> 詳細なビルド手順については、[BUILD-ja.md](./BUILD-ja.md) を参照してください。

## クイックスタート

### 1. OGGホームパスの設定

- 環境変数 `GG_HOMES` または `-g/--gghomes` パラメータを使用して、複数のOGGホーム（英語のカンマまたはセミコロンで区切る）を指定することを推奨します。
  ```bash
  export GG_HOMES="/ogg1,/ogg2,/ogg3"
  ./ggutil tasks
  # または
  ./ggutil -g "/ogg1,/ogg2" info extorcl
  ```

### 2. すべてのコマンドとヘルプの表示

```bash
./ggutil -h
./ggutil <サブコマンド> -h
```

### 3. 一般的なコマンド例

- **ヘルプ情報の表示**

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

- **すべてのOGGインスタンスの監視 (`mon`)**

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

- **プロセス設定の表示 (`config`)**

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

- **パラメータファイルの表示 (`param`)**

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

- **プロセス統計の表示 (`stats`)**

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

- **キーファイルのバックアップ (`backup`)**

  ```bash
  $ ggutil backup

  Please refer to gz file /tmp/oggbackup_xugu01_20250711_195536.tar.gz
  ```

---

## 設計思想とアーキテクチャ

- **並行性優先**: すべてのOGGホーム操作はgoroutineを使用して並行処理され、バッチ処理性能を大幅に向上させます。
- **ユーザーフレンドリーな出力**: すべてのコマンド出力は構造化されたテーブルであり、手動での読み取りや自動収集が容易です。
- **コードの保守性**: すべてのコアモジュールとユーティリティ関数には詳細な英語のコメントがあり、チームコラボレーションと二次開発を促進します。
- **堅牢性**: パラメータ、環境変数、パスなどには詳細な検証とエラープロンプトがあります。デバッグモードでは詳細なログを追跡できます。
- **安全性**: すべてのアーカイブ、削除、ファイル操作には例外処理が含まれており、誤った削除や上書きを防ぎます。

---

## ユースケース

- 複数OGGホーム環境の日常的なバッチ運用と監視。
- OGG設定、ログ、レポートファイルの自動アーカイブとバックアップ。
- OGGプロセスステータス、パラメータ、統計情報のバッチ収集とレポート出力。
- DBA、データ同期プラットフォーム、自動運用チームに適しています。

---

## 貢献とサポート

- issue、PR、提案、議論を歓迎します。
- エンタープライズレベルのカスタム開発や技術サポートについては、作者に連絡するか、GitHubのissue経由でメッセージを残してください。

---

## オープンソースライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細は [LICENSE](LICENSE) を参照してください。

---

機能拡張やエンタープライズカスタマイズについては、メンテナーに連絡するか、issueを提出してください！
