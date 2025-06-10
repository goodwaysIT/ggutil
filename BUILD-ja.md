# ビルドガイド

このプロジェクトは、主流のLinux x86_64およびARM (aarch64) アーキテクチャでのコンパイルと実行をサポートしています。以下は標準的なビルド手順です。

## 1. 環境準備

- オペレーティングシステム：Linux（Oracle Linux、CentOS、RedHat、Ubuntuなどの主流ディストリビューションを推奨）
- Goバージョン：Go 1.18以上を推奨
- Oracle GoldenGateがインストールされ、設定済みであること（マルチホームサポート）

## 2. ソースコードの取得

```bash
git clone https://github.com/goodwaysIT/ggutil.git
cd ggutil
```

## 3. 依存関係のインストール

Goは最初のコンパイル時に自動的に依存関係を取得するため、手動でのインストールは不要です。事前に取得する場合：

```bash
go mod tidy
```

## 4. ローカルコンパイル（現在のプラットフォームアーキテクチャをデフォルトとする）

```bash
go build -o ggutil main.go
```

## 5. クロスプラットフォームコンパイル

### Linux x86_64向けビルド（64ビット Intel/AMD）

```bash
GOOS=linux GOARCH=amd64 go build -o ggutil_x86_64 main.go
```

### Linux ARM64向けビルド（64ビット ARM、例：Kunpeng、Phytium、Raspberry Pi）

```bash
GOOS=linux GOARCH=arm64 go build -o ggutil_arm64 main.go
```

> 注意：他のアーキテクチャ（例：32ビットARM）の場合、`GOARCH`を`arm`に設定し、必要に応じて`GOARM`環境変数を調整してください。

## 6. ビルドの検証

```bash
file ggutil*
# 対応するELF 64ビットLSB実行可能ファイル、x86-64またはaarch64が表示されるはずです
./ggutil_x86_64 -h
./ggutil_arm64 -h
```

## 7. よくある問題

- **依存関係の取得失敗**：ネットワークを確認するか、Goプロキシを設定してください（例：`GOPROXY=https://goproxy.cn`）。
- **権限の問題**：実行権限の問題が発生した場合は、`chmod +x ggutil*`を使用してください。
- **OGG関連コマンドが利用不可**：ターゲットマシンにOracle GoldenGateが正しくインストールされ、承認されていることを確認してください。

---

特別なプラットフォーム、静的コンパイル、またはコンテナ化のニーズについては、issueやPRを通じてフィードバックをお寄せください！
