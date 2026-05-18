---
name: tell-dont-ask
description: Tell, Don't Ask 原則に基づくコードレビューと設計支援。getter後のif分岐、連鎖呼び出し（Train Wreck）、外部での集計、型/ステータス分岐などの違反を検知し、ドメインに振る舞いを集約する改善案を提示する。Go コードのPRレビュー時、設計レビュー時に参照する。
---

# Tell, Don't Ask

オブジェクトに問い合わせるな、命じよ。

## 核心原則

**オブジェクトの内部状態に基づく意思決定をし、その結果で該当オブジェクトを更新してはならない。**

| アプローチ | 特徴 | 問題 |
|-----------|------|------|
| Ask | 状態を取得→外部で判断→操作 | ロジックが散在、カプセル化破壊 |
| Tell | オブジェクトに直接命じる | 責任集約、変更に強い |

## 判断フロー

```
getter でフィールドを取得している
    ↓
その後 if / switch で判定している？
    ├─ YES → Askパターン（問題あり）
    └─ NO  → 表示・出力目的なら許容
```

## アンチパターン（Goでの例）

```go
// ❌ getter + if: 外部で状態を判断して更新
if order.Status == StatusPending && order.Total > 0 {
    order.Status = StatusProcessing
}

// ❌ 連鎖呼び出し (Train Wreck)
city := order.Customer.Address.City

// ❌ 外部での集計: コレクションの外でループ集計
total := 0.0
for _, item := range order.Items {
    total += item.Price
}

// ❌ 型/ステータス分岐: switch で型を見て処理を分ける
switch v := item.(type) {
case *Book:
    return v.Weight * 0.5
case *Electronics:
    return v.Weight*1.0 + 5.0
}
```

## 変換パターン（Go）

### 1. 状態判定の内部化

```go
// ✅ Tell: 判定と操作をメソッドに集約
func (o *Order) ProcessIfReady() error {
    if !o.canProcess() {
        return nil
    }
    o.status = StatusProcessing
    return nil
}

func (o *Order) canProcess() bool {
    return o.status == StatusPending && o.total > 0
}
```

### 2. インターフェースによる型分岐の解消

```go
// ✅ Tell: 各型に振る舞いを持たせる
type Shippable interface {
    ShippingCost() float64
}

func (b *Book) ShippingCost() float64        { return b.weight * 0.5 }
func (e *Electronics) ShippingCost() float64 { return e.weight*1.0 + 5.0 }
```

### 3. コレクション操作の委譲

```go
// ✅ Tell: 集計はコレクション所有者のメソッドに
func (o *Order) TotalPrice() float64 {
    total := 0.0
    for _, item := range o.items {
        total += item.price
    }
    return total
}
```

## 警戒シグナル早見表

| シグナル | 対応 |
|----------|------|
| `obj.Field` / `obj.Get*()` の後に `if` | 判定メソッドを obj に追加 |
| `obj.Field = obj.Field + n` | 操作メソッドを obj に追加 |
| `a.B.C.DoD()` の連鎖 | 委譲メソッドを各レベルに追加 |
| `switch v := x.(type)` | インターフェースで置換 |
| `for` + フィールド取得 + 集計変数 | 集計メソッドをコレクション側に追加 |

## 過剰適用を避ける

以下には適用しない:

- 表示・レポート目的の値取得
- DTO / 値オブジェクトからの単純な値取得
- フレームワーク・ライブラリの制約がある場合
- メソッド追加でその型が肥大化する場合（責任分割を先に検討）
