# Documentation

This directory contains supplementary documentation for `es-typed-go`.

The top-level `README*` files are intentionally focused on library users: installation, quick start, core concepts, and API usage. More detailed reference material and contributor-facing guidance live under `docs/`.

## For library users

### Guides
- [Search Guide](search-guide.md) — Search API selection notes and query building patterns
- [検索ガイド](search-guide.ja.md) — 検索 API の使い分けとクエリ構築パターン（日本語）
- [搜索指南](search-guide.zh-CN.md) — 搜索 API 选择指引与查询构建模式（简体中文）
- [Migration Guide (v2)](migration-v2.md) — Breaking changes and migration steps for the v2 architecture

### Reference
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
- internal tooling guidance
- migration guides for breaking changes

If a document is mainly useful for contributors or maintainers, it belongs in `docs/`, not in the top-level README.
