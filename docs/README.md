# Documentation

This directory contains supplementary documentation for `es-typed-go`.

The top-level `README*` files are intentionally focused on library users: installation, quick start, core concepts, and API usage. More detailed reference material and contributor-facing guidance live under `docs/`.

## For library users

### Reference
- [Search Guide](search-guide.md) — Detailed guide to typed search helpers, counts, pagination, and raw search usage
- [検索ガイド](search-guide.ja.md) — 型付き検索 helper、件数取得、ページネーション、raw search の詳細ガイド（日本語）
- [搜索指南](search-guide.zh-CN.md) — typed search helper、计数、分页与 raw search 用法详解（简体中文）
- [Property Reference](property-reference.md) — Complete list of supported property builders and their functional options
- [プロパティリファレンス](property-reference.ja.md) — 全プロパティ型とオプションの一覧（日本語）

### Main entry points
- [`../README.md`](../README.md) — English user guide
- [`../README.ja.md`](../README.ja.md) — 日本語ユーザーガイド
- [`../README.zh-CN.md`](../README.zh-CN.md) — 中文用户指南

## For contributors

### Contribution guides
- [Contributing Guide](contributing.md) — Development workflow, repository rules, testing, generators, internal tools, and PR checklist
- [コントリビューションガイド](contributing.ja.md) — 開発フロー、テスト、内部ツール、保守ルール
- [贡献指南](contributing.zh-CN.md) — 开发流程、测试、代码生成与仓库规则

## Documentation policy

Use the following separation when adding or editing documentation:

### `README*`
Keep these focused on library users:
- what the library solves
- installation
- quick start
- public API usage
- links to detailed references

### `docs/`
Place detailed or contributor-oriented material here:
- reference documentation
- development workflow
- testing instructions
- generator maintenance notes
- repository conventions
- internal tooling guidance

If a document is mainly useful for contributors or maintainers, it belongs in `docs/`, not in the top-level README.

## Suggested future additions

If the documentation set grows, this directory is a good place for:
- migration guides for breaking changes
- generator design notes
- integration testing setup details
- architecture or API design notes
- release process documentation