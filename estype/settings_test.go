package estype_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/tomtwinkle/es-typed-go/estype"
)

type settingsProviderStub struct {
	settings estype.Settings
}

func (s settingsProviderStub) Settings() estype.Settings {
	return s.settings
}

type settingsAndMappingProviderStub struct {
	settings estype.Settings
	mapping  estype.Mapping
}

func (p settingsAndMappingProviderStub) Settings() estype.Settings {
	return p.settings
}

func (p settingsAndMappingProviderStub) Mapping() estype.Mapping {
	return p.mapping
}

func TestSettings_ZeroValue(t *testing.T) {
	t.Parallel()

	var settings estype.Settings

	assert.Assert(t, settings.NumberOfShards == nil)
	assert.Assert(t, settings.NumberOfReplicas == nil)
	assert.Assert(t, settings.RefreshInterval == nil)
}

func TestSettings_CanHoldAllSupportedFields(t *testing.T) {
	t.Parallel()

	settings := estype.Settings{
		NumberOfShards:   new(int(3)),
		NumberOfReplicas: new(int(1)),
		RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
	}

	assert.Assert(t, settings.NumberOfShards != nil)
	assert.Equal(t, 3, *settings.NumberOfShards)

	assert.Assert(t, settings.NumberOfReplicas != nil)
	assert.Equal(t, 1, *settings.NumberOfReplicas)

	assert.Assert(t, settings.RefreshInterval != nil)
	assert.Equal(t, estype.RefreshIntervalDefault, *settings.RefreshInterval)
}

func TestSettings_RefreshIntervalDisable(t *testing.T) {
	t.Parallel()

	settings := estype.Settings{
		RefreshInterval: new(estype.RefreshInterval(estype.RefreshIntervalDisable)),
	}

	assert.Assert(t, settings.RefreshInterval != nil)
	assert.Equal(t, estype.RefreshIntervalDisable, *settings.RefreshInterval)
	assert.Equal(t, "-1", settings.RefreshInterval.String())
}

func TestESSettingsProvider_ReturnsConfiguredSettings(t *testing.T) {
	t.Parallel()

	provider := settingsProviderStub{
		settings: estype.Settings{
			NumberOfShards:   new(int(5)),
			NumberOfReplicas: new(int(0)),
			RefreshInterval:  new(estype.RefreshInterval(estype.RefreshIntervalDefault)),
		},
	}

	got := provider.Settings()

	assert.Assert(t, got.NumberOfShards != nil)
	assert.Equal(t, 5, *got.NumberOfShards)

	assert.Assert(t, got.NumberOfReplicas != nil)
	assert.Equal(t, 0, *got.NumberOfReplicas)

	assert.Assert(t, got.RefreshInterval != nil)
	assert.Equal(t, estype.RefreshIntervalDefault, *got.RefreshInterval)
}

func TestESSettingsProvider_InterfaceSatisfaction(t *testing.T) {
	t.Parallel()

	var provider estype.SettingsProvider = settingsProviderStub{}
	assert.Assert(t, provider != nil)
}

func TestESConfig_RequiresSettingsAndMapping(t *testing.T) {
	t.Parallel()

	var provider estype.ESConfig = settingsAndMappingProviderStub{}
	assert.Assert(t, provider != nil)

	gotSettings := provider.Settings()
	assert.Assert(t, gotSettings.NumberOfShards == nil)
	assert.Assert(t, gotSettings.NumberOfReplicas == nil)
	assert.Assert(t, gotSettings.RefreshInterval == nil)

	gotMapping := provider.Mapping()
	assert.Assert(t, len(gotMapping.Fields) == 0)
}
