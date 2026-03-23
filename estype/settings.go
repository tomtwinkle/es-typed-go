package estype

// Settings represents Elasticsearch index settings in a library-owned,
// type-safe form.
//
// This type is intentionally minimal to start with. It can be expanded
// incrementally as the library adds first-class support for more index
// creation settings without forcing application model packages to depend on
// Elasticsearch typed client request types directly.
type Settings struct {
	// NumberOfShards controls how many primary shards the index has.
	NumberOfShards *int

	// NumberOfReplicas controls how many replica shards each primary shard has.
	NumberOfReplicas *int

	// RefreshInterval controls how often the index is refreshed.
	RefreshInterval *RefreshInterval
}

// ESConfig is implemented by types that provide both model-owned Elasticsearch
// index settings and field mappings.
//
// This is the unified configuration contract used by high-level index-creation
// helpers so callers can pass a single model/config value that fully describes
// the index.
type ESConfig interface {
	Settings() Settings
	Mapping() Mapping
}

// SettingsProvider is implemented by types that describe their Elasticsearch
// index settings in estype-owned types.
//
// This pairs with MappingProvider so that application model types can define
// both their index settings and mapping without importing Elasticsearch typed
// client request types.
//
// Example:
//
//	type Product struct {
//		ID string `json:"id"`
//	}
//
//	func (Product) Settings() estype.Settings {
//		return estype.Settings{
//			NumberOfShards:   new(int(1)),
//			NumberOfReplicas: new(int(0)),
//			RefreshInterval:  new(RefreshInterval(RefreshIntervalDefault)),
//		}
//	}
type SettingsProvider interface {
	Settings() Settings
}
