# Copilot Instructions

## リポジトリ概要

Go製のシンプルなBook CRUD APIを持つ実験用リポジトリ。
AI agent（GitHub Copilot / Claude Code）を使ったissue駆動開発の試し撃ち場として使用する。

## 技術スタック

- Language: Go 1.22+
- HTTP: `net/http`（標準ライブラリのみ）
- Storage: インメモリ（外部DBなし）

## ドメインモデル

`internal/book/` に Book CRUD を実装。

```go
type Book struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    ISBN   string `json:"isbn"`
}
```

## API エンドポイント

| Method | Path        | 説明         |
|--------|-------------|--------------|
| GET    | /books      | 一覧取得     |
| POST   | /books      | 新規作成     |
| GET    | /books/{id} | 1件取得      |
| PUT    | /books/{id} | 更新         |
| DELETE | /books/{id} | 削除         |

## Copilot agentへの期待役割

- シンプルな機能追加・修正issueを担当する
- 実装は `internal/book/` 配下に集約する
- 外部依存ライブラリは追加しない（標準ライブラリのみ）
- テストを追加する場合は `_test.go` ファイルに記述する

## 設計が必要な複雑なissue

設計の壁打ちが必要な場合は Copilot ではなく Claude Code（issue-design スキル）で対応する。

## 実装原則

詳細は `.github/skills/go-design-principles/SKILL.md` を参照。以下の原則を遵守する:

- **Always-Valid Domain**: ドメイン型はコンストラクタ経由でのみ生成し `Validate()` で検証。生 struct リテラル禁止。
- **panic と error の役割分離**: 契約違反は panic、業務エラー・外部入力エラーは error で返す。
- **Tell, Don't Ask**: 業務判定はドメインに持たせ、ハンドラは薄いオーケストレーターに留める。
- **Sentinel Error と詳細エラー型**: `Err*` 変数（`errors.Is` 用）と struct 型（`errors.As` 用）を対で定義し `Unwrap()` でチェーン維持。
- **層間エラー変換**: repository のエラーをそのまま handler に返さず、自層エラーに変換してから返す。

## PRレビュー観点

**PRレビューを行う際は、必ず `.github/skills/code-review-lenses/SKILL.md` を読み込み、そこに記載された指示にすべて従うこと。**

詳細は `.github/skills/code-review-lenses/SKILL.md` を参照。レビューは以下の4観点と重大度分類で行う。

### 4観点
1. **バグ**: 境界値・nil・競合状態・エラー握りつぶし
2. **セキュリティ**: 入力検証・機微情報のログ漏れ・インジェクション
3. **パフォーマンス**: 不要なアロケーション・N+1・mutex の粒度
4. **保守性**: 単一責任・命名・テスト容易性・過剰な抽象化

### 重大度
- **Critical**: 本番障害・データ損失・脆弱性 → マージブロック推奨
- **Important**: 品質劣化・将来のバグ温床 → マージブロック推奨
- **Suggestion**: 改善提案
- **Informational**: 知識共有

### Go 固有追加観点
- `context.Context` を下位層まで伝搬しているか
- `errors.Is` / `errors.As` が使えるようにエラーチェーンを維持しているか
- interface は呼び出し側に必要なメソッドだけを持つ小さな単位か
