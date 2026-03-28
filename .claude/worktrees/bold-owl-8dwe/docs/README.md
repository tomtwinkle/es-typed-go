# Documentation

This directory contains supplementary documentation for `es-typed-go`.

The top-level `README*` files are intentionally focused on library users: installation, quick start, core concepts, and API usage. More detailed reference material and contributor-facing guidance live under `docs/`.

## For library users

### Reference
- [Search Guide](search-guide.md) — Concise search API selection notes
- [検索ガイド](search-guide.ja.md) — 検索 API の使い分けを簡潔に整理したメモ（日本語）
- [搜索指南](search-guide.zh-CN.md) — 搜索 API 选择指引的简要说明（简体中文）
- [Property Reference](property-reference.md) — Complete list of supported property builders and their functional options
- [プロパティリファレンス](property-reference.ja.md) — 全プロパティ型とオプションの一覧（日本語）

### Main entry points
- [`../README.md`](../README.md) — English user guide
- [`../README.ja.md`](../README.ja.md) — 日本語ユーザーガイド
- [`../README.zh-CN.md`](../README.zh-CN.md) — 中文用户指南

## For contributors

### Contribution guides
- [Contributing Guide](contributing.md) — Development prerequisites, local Elasticsearch setup, build/test commands, and validation steps

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
- testing instructions
- environment setup details
- reference documentation
- internal tooling guidance

If a document is mainly useful for contributors or maintainers, it belongs in `docs/`, not in the top-level README.

## Suggested future additions

If the documentation set grows, this directory is a good place for:
- migration guides for breaking changes
- generator design notes
- integration testing setup details
- architecture or API design notes
- release process documentation