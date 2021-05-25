//go:generate mockery --name Registry|TileEnabler|GeneratorEnabler|TileMetadataExplorer|VariantMetadataExplorer

package registry

import (
	"fmt"

	"github.com/monitoror/monitoror/api/config/models"
	"github.com/monitoror/monitoror/api/config/versions"
	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	coreModels "github.com/monitoror/monitoror/models"
)

type (
	// Registry is used to register Tile and Generator in config for verify / hydrate
	Registry interface {
		RegisterMonitorable(monitorable coreModels.Monitorable)
		RegisterTile(tileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) TileEnabler
		RegisterGenerator(generatedTileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) GeneratorEnabler
		GetMonitorables() []*MonitorableMetadata
	}
	// TileEnabler is returned to monitorable after register to enable monitorable tile with this variant if she is "valid"
	TileEnabler interface {
		Enable(variantName coreModels.VariantName, paramsValidator params.Validator, routePath string)
	}
	// GeneratorEnabler is returned to monitorable after register to enable monitorable generator with this variant if she is "valid"
	GeneratorEnabler interface {
		Enable(variantName coreModels.VariantName, generatorParamsValidator params.Validator, tileGeneratorFunction models.TileGeneratorFunction)
	}

	// TileMetadataExplorer is used in verify. Matching tileMetadata and generatorMetadata.
	TileMetadataExplorer interface {
		GetMinimalVersion() versions.RawVersion
		GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool)
		GetVariantsNames() []coreModels.VariantName
	}
	// VariantMetadataExplorer is used in verify. Matching tileVariantMetadata and generatorVariantMetadata.
	VariantMetadataExplorer interface {
		IsEnabled() bool
		GetValidator() params.Validator
	}
)

type (
	MetadataRegistry struct {
		MonitorableMetadata []*MonitorableMetadata
		TileMetadata        map[coreModels.TileType]*tileMetadata
		GeneratorMetadata   map[coreModels.TileType]*generatorMetadata
	}

	MonitorableMetadata struct {
		Monitorable      coreModels.Monitorable
		VariantsMetadata []*MonitorableVariantMetadata
	}

	MonitorableVariantMetadata struct {
		VariantName coreModels.VariantName

		Enabled bool
		Errors  []error
	}

	tileMetadata struct {
		// TileType
		TileType coreModels.TileType
		// MinimalVersion is the version that makes the tile available
		MinimalVersion versions.RawVersion
		// VariantsMetadata list all registered variants (can be available or not)
		VariantsMetadata map[coreModels.VariantName]*tileVariantMetadata
	}

	tileVariantMetadata struct {
		// Enabled define if variant is enabled
		Enabled bool

		// VariantName
		VariantName coreModels.VariantName

		// RoutePath path of the api endpoint for this tile. Used by hydrate
		RoutePath *string
		// ParamsValidator is used to validate given params
		ParamsValidator params.Validator
	}

	generatorMetadata struct {
		// TileType
		TileType coreModels.TileType
		// GeneratedTileType
		GeneratedTileType coreModels.TileType
		// MinimalVersion is the version that makes the tile available
		MinimalVersion versions.RawVersion
		// VariantsMetadata list all registered variants (can be available or not)
		VariantsMetadata map[coreModels.VariantName]*generatorVariantMetadata
	}

	generatorVariantMetadata struct {
		// Enabled define if variant is enabled
		Enabled bool

		// VariantName
		VariantName coreModels.VariantName

		// GeneratorFunction function used to generate tile config
		GeneratorFunction models.TileGeneratorFunction
		// GeneratorParamsValidator is used to validate given params for generator
		GeneratorParamsValidator params.Validator
	}
)

func NewRegistry() *MetadataRegistry {
	return &MetadataRegistry{
		MonitorableMetadata: []*MonitorableMetadata{},
		TileMetadata:        make(map[coreModels.TileType]*tileMetadata),
		GeneratorMetadata:   make(map[coreModels.TileType]*generatorMetadata),
	}
}

// REGISTRY
// ----------------------------------------
func (r *MetadataRegistry) RegisterMonitorable(monitorable coreModels.Monitorable) {
	monitorableMetadata := &MonitorableMetadata{
		Monitorable: monitorable,
	}

	for _, variantName := range monitorable.GetVariantsNames() {
		isValid, errors := monitorable.Validate(variantName)

		monitorableMetadata.VariantsMetadata = append(monitorableMetadata.VariantsMetadata, &MonitorableVariantMetadata{
			VariantName: variantName,
			Enabled:     isValid,
			Errors:      errors,
		})
	}

	r.MonitorableMetadata = append(r.MonitorableMetadata, monitorableMetadata)
}

func (r *MetadataRegistry) RegisterTile(tileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) TileEnabler {
	tileMetadata := &tileMetadata{
		TileType:         tileType,
		MinimalVersion:   minimalVersion,
		VariantsMetadata: make(map[coreModels.VariantName]*tileVariantMetadata),
	}

	// Register Variant with Enabled False
	for _, variantName := range variantNames {
		tileMetadata.VariantsMetadata[variantName] = &tileVariantMetadata{
			VariantName: variantName,
		}
	}

	r.TileMetadata[tileType] = tileMetadata

	return tileMetadata
}

func (r *MetadataRegistry) RegisterGenerator(generatedTileType coreModels.TileType, minimalVersion versions.RawVersion, variantNames []coreModels.VariantName) GeneratorEnabler {
	// Boxing tile type into generator
	tileType := coreModels.NewGeneratorTileType(generatedTileType)

	generatorMetadata := &generatorMetadata{
		TileType:          tileType,
		GeneratedTileType: generatedTileType,
		MinimalVersion:    minimalVersion,
		VariantsMetadata:  make(map[coreModels.VariantName]*generatorVariantMetadata),
	}

	// Register Variant with Enabled False
	for _, variantName := range variantNames {
		generatorMetadata.VariantsMetadata[variantName] = &generatorVariantMetadata{
			VariantName: variantName,
		}
	}

	r.GeneratorMetadata[tileType] = generatorMetadata

	return generatorMetadata
}

func (r *MetadataRegistry) GetMonitorables() []*MonitorableMetadata {
	return r.MonitorableMetadata
}

// ----------------------------------------

// TILE METADATA
// ----------------------------------------
func (tm *tileMetadata) Enable(variantName coreModels.VariantName, paramsValidator params.Validator, routePath string) {
	variantMetadata, exists := tm.VariantsMetadata[variantName]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variantName: %s for tile: %s. register it before.", variantName, tm.TileType))
	}

	variantMetadata.Enabled = true
	variantMetadata.ParamsValidator = paramsValidator
	variantMetadata.RoutePath = &routePath
}

func (tm *tileMetadata) GetMinimalVersion() versions.RawVersion {
	return tm.MinimalVersion
}

func (tm *tileMetadata) GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool) {
	v, exists := tm.VariantsMetadata[variantName]
	return v, exists
}

func (tm *tileMetadata) GetVariantsNames() []coreModels.VariantName {
	var result []coreModels.VariantName
	for variantName := range tm.VariantsMetadata {
		result = append(result, variantName)
	}
	return result
}

// ----------------------------------------

// GENERATOR METADATA
// ----------------------------------------
func (gm *generatorMetadata) Enable(variantName coreModels.VariantName, generatorParamsValidator params.Validator, tileGeneratorFunction models.TileGeneratorFunction) {
	variantMetadata, exists := gm.VariantsMetadata[variantName]
	if !exists {
		panic(fmt.Sprintf("unable to enable unknown variantName: %s for tile: %s. register it before.", variantName, gm.TileType))
	}

	variantMetadata.Enabled = true
	variantMetadata.GeneratorParamsValidator = generatorParamsValidator
	variantMetadata.GeneratorFunction = tileGeneratorFunction
}

func (gm *generatorMetadata) GetMinimalVersion() versions.RawVersion {
	return gm.MinimalVersion
}

func (gm *generatorMetadata) GetVariant(variantName coreModels.VariantName) (VariantMetadataExplorer, bool) {
	v, exists := gm.VariantsMetadata[variantName]
	return v, exists
}

func (gm *generatorMetadata) GetVariantsNames() []coreModels.VariantName {
	var result []coreModels.VariantName
	for variantName := range gm.VariantsMetadata {
		result = append(result, variantName)
	}
	return result
}

// ----------------------------------------

// TILE VARIANT METADATA
// ----------------------------------------
func (tvm *tileVariantMetadata) IsEnabled() bool {
	return tvm.Enabled
}
func (tvm *tileVariantMetadata) GetValidator() params.Validator {
	return tvm.ParamsValidator
}

// ----------------------------------------

// GENERATOR VARIANT METADATA
// ----------------------------------------
func (gvm *generatorVariantMetadata) IsEnabled() bool {
	return gvm.Enabled
}
func (gvm *generatorVariantMetadata) GetValidator() params.Validator {
	return gvm.GeneratorParamsValidator
}

// ----------------------------------------
