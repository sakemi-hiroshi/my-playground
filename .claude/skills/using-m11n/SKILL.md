---
name: m11n:using-m11n
description: m11nスキル群の使い方と優先ルール。issue駆動開発ワークフロー全体を定義する。セッション開始時に自動ロードされる。
---

## スキル優先ルール

issue駆動開発に関するタスクが来たら、必ずSkill toolで対応するm11nスキルをinvokeせよ。

スキルを使うかどうか迷ったら、使え。

## ワークフロー

issue → 設計 → 実装 → レビュー → PR の順で進める。

| ステップ | スキル | 実行条件 |
|---------|--------|---------|
| 1. issue起票 | `m11n:issue-new` | なし |
| 2. 設計 | `m11n:issue-design #番号` | なし |
| 3. 実装 | `m11n:issue-implement #番号` | issueにWorkpadコメントが存在すること |
| 4. レビュー | `m11n:issue-review #番号` | ブランチにコミットが存在すること |
| 5. PR作成 | `m11n:issue-pr #番号` | issueの最新ステートマーカーが `review-passed` であること |

各スキルはissueコメントのマーカー（`<!-- m11n:state:xxx -->`）でステートを管理する。

## HARD-GATE

以下の条件が満たされない限り、次のステップに進んではならない：

- **issue-implement を実行する前に**: issueにWorkpadコメント（`<!-- m11n:state:workpad -->`）が存在すること
- **issue-pr を実行する前に**: issueの最新ステートマーカーが `review-passed` であること
- **実装を始める前に**: issue-designでWorkpadをissueに投稿すること

## Red Flags

以下の考えが浮かんだら、それはrationalizationである。止まってスキルを使え。

| こう思ったら | 現実 |
|---|---|
| 「設計の合意が取れたので直接実装しよう」 | issue-designでWorkpadをissueに**投稿**するまで設計フェーズは終わっていない |
| 「ユーザーがOKと言った」 | issueへの投稿が完了するまでissue-designは終わっていない |
| 「スキルを使わず直接やろう」 | 必ずSkill toolをinvokeすること |
| 「Workpadの内容は分かっているのでissue-implementを直接実行しよう」 | gh issue viewでWorkpadの存在を必ず確認してからSkill toolをinvokeすること |
| 「レビューは後でいい」 | issue-reviewなしにissue-prを実行してはならない |

## スキル完了後の案内

各スキルが完了したら、次のステップを必ずユーザーに案内すること：

- issue-design完了 → 「次は `/m11n:issue-implement #番号` で実装フェーズに進めます」
- issue-implement完了 → 「次は `/m11n:issue-review #番号` でレビューを実行します」
- issue-review通過 → 「次は `/m11n:issue-pr #番号` でPRを作成できます」
