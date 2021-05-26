package printer

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
	coreConfig "github.com/monitoror/monitoror/config"
	coreModels "github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/pkg/system"
	"github.com/monitoror/monitoror/pkg/templates"
)

var startupTemplate = `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / {{ with .BuildTags }}{{ printf " %s " . | inverseColor }}{{ end }}
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  {{ .Version | green }}

{{ "https://monitoror.com" | blue }}
{{ if .DisableUI }}

┌─ {{ "DEVELOPMENT MODE" | yellow }} ──────────────────────────────┐
│ UI must be started via {{ "npm serve" | green }} from ./ui      │
│ For more details, check our development guide:  │
│ {{ "https://monitoror.com/guides/#development" | blue }}       │
└─────────────────────────────────────────────────┘
{{ end }}

{{ "ENABLED MONITORABLES" | green }}
{{ range .Monitorables }}{{ if not .IsDisabled }}
  {{- if not .ErroredVariants }}
  {{ "✓ " | green }}
  {{- else if .EnabledVariants }}
  {{ "! " | yellow }}
  {{- else }}
  {{ "x " | red }}
  {{- end }}
  {{- .MonitorableName }} {{ .StringifyEnabledVariants | grey }}

  {{- range .ErroredVariants }}
    {{- if eq .VariantName "` + string(coreModels.DefaultVariantName) + `" }}
    {{ printf "/!\\ Errored %s configuration" .VariantName | red }}
    {{- else }}
    {{ printf "/!\\ Errored %q variant configuration" .VariantName | red }}
    {{- end }}
    {{- range .Errors }}
      {{ . }}
    {{- end }}
  {{- end }}
{{- end }}{{- end }}

{{ if ne .DisabledMonitorableCount 0 -}}
{{ printf "%d more monitorables were ignored" .DisabledMonitorableCount | yellow }}
Check the documentation to know how to enabled them:
{{ printf "https://monitoror.com/%sdocumentation/" .DocumentationVersion | blue }}

{{ end }}
{{- with .NamedConfigs }}
{{ "AVAILABLE NAMED CONFIGURATIONS" | green }}
{{ range . }}
  {{ .Name }}{{ printf " -> %s" .Value | grey }}
{{- end }}

{{ end }}
─────────────────────────────────────────────────

{{ "MONITOROR IS RUNNING AT:" | green }}
{{- range .DisplayedAddresses }}
  {{ printf "http://%s:%d" . $.LookupPort | blue }}
{{- end }}

─────────────────────────────────────────────────

`

type (
	startupInfo struct {
		Version       string // From ldflags
		BuildTags     string // From ldflagsl
		LookupPort    int    // From .env
		LookupAddress string // From .env
		DisableUI     bool   // From .env
		NamedConfigs  []namedConfigInfo
		Monitorables  []monitorableInfo
	}

	namedConfigInfo struct {
		Name  string
		Value string
	}

	monitorableInfo struct {
		MonitorableName string     // From registry
		EnabledVariants []struct { // From registry
			VariantName string
		}
		ErroredVariants []struct { // From registry
			VariantName string
			Errors      []error
		}
	}
)

var parsedTemplate *template.Template

func init() {
	// Print this error when you want to debug template
	parsedTemplate, _ = templates.New("monitoror").Parse(startupTemplate)
}

func PrintStartupLog(monitororCli *cli.MonitororCli) error {
	monitororInfo := &startupInfo{
		Version:       version.Version,
		BuildTags:     version.BuildTags,
		DisableUI:     monitororCli.Store.CoreConfig.DisableUI,
		LookupPort:    monitororCli.Store.CoreConfig.Port,
		LookupAddress: monitororCli.Store.CoreConfig.Address,
	}

	// Named config Info
	for name, config := range monitororCli.Store.CoreConfig.NamedConfigs {
		monitororInfo.NamedConfigs = append(monitororInfo.NamedConfigs, namedConfigInfo{Name: string(name), Value: config})
	}
	monitororInfo.NamedConfigs = sortNamedConfigs(monitororInfo.NamedConfigs)

	// Monitorables info
	for _, mm := range monitororCli.Store.Registry.GetMonitorables() {
		monitorableInfo := monitorableInfo{
			MonitorableName: mm.Monitorable.GetDisplayName(),
		}

		for _, v := range mm.VariantsMetadata {
			if v.Enabled {
				monitorableInfo.EnabledVariants = append(monitorableInfo.EnabledVariants, struct {
					VariantName string
				}{string(v.VariantName)})
			}

			if len(v.Errors) > 0 {
				monitorableInfo.ErroredVariants = append(monitorableInfo.ErroredVariants, struct {
					VariantName string
					Errors      []error
				}{string(v.VariantName), v.Errors})
			}
		}

		monitororInfo.Monitorables = append(monitororInfo.Monitorables, monitorableInfo)
	}

	return parsedTemplate.Execute(monitororCli.Output, monitororInfo)
}

func sortNamedConfigs(namedConfigs []namedConfigInfo) []namedConfigInfo {
	sort.Slice(namedConfigs, func(i, j int) bool {
		if namedConfigs[i].Name == string(coreConfig.DefaultConfigName) {
			return true
		}
		return namedConfigs[i].Name < namedConfigs[j].Name
	})

	return namedConfigs
}

func (mi *startupInfo) DocumentationVersion() string {
	if !strings.HasSuffix(mi.Version, "-dev") {
		return ""
	}
	documentationVersion := ""
	if splittedVersion := strings.Split(mi.Version, "."); len(splittedVersion) == 3 {
		documentationVersion = fmt.Sprintf("%s.%s/", splittedVersion[0], splittedVersion[1])
	}
	return documentationVersion
}

func (mi *startupInfo) DisabledMonitorableCount() int {
	disabledMonitorableCount := 0
	for _, m := range mi.Monitorables {
		if m.IsDisabled() {
			disabledMonitorableCount++
		}
	}
	return disabledMonitorableCount
}

func (mi *startupInfo) DisplayedAddresses() []string {
	var adressess []string

	if mi.LookupAddress == "" || mi.LookupAddress == "0.0.0.0" {
		adressess = append(adressess, "localhost")
		adressess = append(adressess, system.GetNetworkIP())
	} else {
		adressess = append(adressess, mi.LookupAddress)
	}

	return adressess
}

func (mi *monitorableInfo) IsDisabled() bool {
	return len(mi.EnabledVariants) == 0 && len(mi.ErroredVariants) == 0
}

func (mi *monitorableInfo) StringifyEnabledVariants() string {
	var strVariants string
	if len(mi.EnabledVariants) == 1 && mi.EnabledVariants[0].VariantName == string(coreModels.DefaultVariantName) {
		if len(mi.ErroredVariants) > 0 {
			strVariants = fmt.Sprintf("[%s]", coreModels.DefaultVariantName)
		}
	} else {
		var strDefault string
		var variantsWithoutDefault []string

		for _, v := range mi.EnabledVariants {
			if v.VariantName == string(coreModels.DefaultVariantName) {
				strDefault = fmt.Sprintf("%s, ", v.VariantName)
			} else {
				variantsWithoutDefault = append(variantsWithoutDefault, v.VariantName)
			}
		}
		if len(variantsWithoutDefault) > 0 {
			strVariants = fmt.Sprintf("[%svariants: [%s]]", strDefault, strings.Join(variantsWithoutDefault, ", "))
		}
	}

	return strVariants
}
