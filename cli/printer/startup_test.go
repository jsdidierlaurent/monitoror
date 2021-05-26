package printer

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/monitoror/monitoror/cli"
	"github.com/monitoror/monitoror/cli/version"
	"github.com/monitoror/monitoror/config"
	"github.com/monitoror/monitoror/models"
	"github.com/monitoror/monitoror/models/mocks"
	"github.com/monitoror/monitoror/registry"
	"github.com/monitoror/monitoror/store"

	"github.com/stretchr/testify/assert"
)

func initCli(writer io.Writer) *cli.MonitororCli {
	version.Version = "1.0.0"
	version.BuildTags = ""

	return &cli.MonitororCli{
		Output: writer,
		Store: &store.Store{
			CoreConfig: &config.CoreConfig{
				Port:    3000,
				Address: "1.2.3.4",
			},
			Registry: registry.NewRegistry(),
		},
	}
}

func TestPrintMonitororStartupLog_Small(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES



─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://1.2.3.4:3000

─────────────────────────────────────────────────

`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_Full(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)
	monitororCli.Store.CoreConfig.DisableUI = true
	version.Version = "1.0.0-dev"

	monitorableMock1 := new(mocks.Monitorable)
	monitorableMock1.On("GetDisplayName").Return("Monitorable 4")
	monitorableMock1.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName})
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, nil)

	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock1)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0-dev

https://monitoror.com


┌─ DEVELOPMENT MODE ──────────────────────────────┐
│ UI must be started via npm serve from ./ui      │
│ For more details, check our development guide:  │
│ https://monitoror.com/guides/#development       │
└─────────────────────────────────────────────────┘


ENABLED MONITORABLES


1 more monitorables were ignored
Check the documentation to know how to enabled them:
https://monitoror.com/1.0/documentation/


─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://1.2.3.4:3000

─────────────────────────────────────────────────

`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithoutAddress(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)
	monitororCli.Store.CoreConfig.Address = ""

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES



─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://localhost:3000
  http://`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.True(t, strings.HasPrefix(output.String(), expected))
}

func TestPrintMonitororStartupLog_WithMonitorable(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	monitorableMock1 := new(mocks.Monitorable)
	monitorableMock1.On("GetDisplayName").Return("Monitorable 1")
	monitorableMock1.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName, "variant1", "variant2"})
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock2 := new(mocks.Monitorable)
	monitorableMock2.On("GetDisplayName").Return("Monitorable 2")
	monitorableMock2.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName})
	monitorableMock2.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock3 := new(mocks.Monitorable)
	monitorableMock3.On("GetDisplayName").Return("Monitorable 3")
	monitorableMock3.On("GetVariantsNames").Return([]models.VariantName{"variant1"})
	monitorableMock3.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)
	monitorableMock4 := new(mocks.Monitorable)
	monitorableMock4.On("GetDisplayName").Return("Monitorable 4")
	monitorableMock4.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName})
	monitorableMock4.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, nil)

	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock1)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock2)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock3)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock4)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES

  ✓ Monitorable 1 [default, variants: [variant1, variant2]]
  ✓ Monitorable 2 ` + `
  ✓ Monitorable 3 [variants: [variant1]]

1 more monitorables were ignored
Check the documentation to know how to enabled them:
https://monitoror.com/documentation/


─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://1.2.3.4:3000

─────────────────────────────────────────────────

`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithErroredMonitorable(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)

	monitorableMock1 := new(mocks.Monitorable)
	monitorableMock1.On("GetDisplayName").Return("Monitorable 1")
	monitorableMock1.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName, "variant1"})
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil).Once()
	monitorableMock1.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, []error{errors.New("error 1"), errors.New("error 2")})
	monitorableMock2 := new(mocks.Monitorable)
	monitorableMock2.On("GetDisplayName").Return("Monitorable 2")
	monitorableMock2.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName})
	monitorableMock2.On("Validate", mock.AnythingOfType("models.VariantName")).Return(false, []error{errors.New("error 1"), errors.New("error 2")})
	monitorableMock3 := new(mocks.Monitorable)
	monitorableMock3.On("GetDisplayName").Return("Monitorable 3")
	monitorableMock3.On("GetVariantsNames").Return([]models.VariantName{models.DefaultVariantName, "variant1", "variant2"})
	monitorableMock3.On("Validate", mock.AnythingOfType("models.VariantName")).Return(true, nil)

	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock1)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock2)
	monitororCli.Store.Registry.RegisterMonitorable(monitorableMock3)

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES

  ! Monitorable 1 [default]
    /!\ Errored "variant1" variant configuration
      error 1
      error 2
  x Monitorable 2 ` + `
    /!\ Errored default configuration
      error 1
      error 2
  ✓ Monitorable 3 [default, variants: [variant1, variant2]]


─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://1.2.3.4:3000

─────────────────────────────────────────────────

`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestPrintMonitororStartupLog_WithNamedConfigs(t *testing.T) {
	output := &bytes.Buffer{}
	monitororCli := initCli(output)
	monitororCli.Store.CoreConfig.Address = "1.2.3.4"
	monitororCli.Store.CoreConfig.NamedConfigs = make(map[config.ConfigName]string)
	monitororCli.Store.CoreConfig.NamedConfigs[config.DefaultConfigName] = "default"
	monitororCli.Store.CoreConfig.NamedConfigs["test2"] = "test2"
	monitororCli.Store.CoreConfig.NamedConfigs["test"] = "test"

	expected := `
    __  ___            _ __
   /  |/  /___  ____  (_) /_____  _________  _____
  / /|_/ / __ \/ __ \/ / __/ __ \/ ___/ __ \/ ___/
 / /  / / /_/ / / / / / /_/ /_/ / /  / /_/ / / ` + `
/_/  /_/\____/_/ /_/_/\__/\____/_/   \____/_/  1.0.0

https://monitoror.com


ENABLED MONITORABLES



AVAILABLE NAMED CONFIGURATIONS

  default -> default
  test -> test
  test2 -> test2


─────────────────────────────────────────────────

MONITOROR IS RUNNING AT:
  http://1.2.3.4:3000

─────────────────────────────────────────────────

`

	assert.NoError(t, PrintStartupLog(monitororCli))
	assert.Equal(t, expected, output.String())
}

func TestSortNamedConfigs(t *testing.T) {
	namedConfigs := []namedConfigInfo{
		{Name: "test3"},
		{Name: "default"},
		{Name: "test2"},
		{Name: "test1"},
	}

	namedConfigs = sortNamedConfigs(namedConfigs)

	assert.Equal(t, "default", namedConfigs[0].Name)
	assert.Equal(t, "test1", namedConfigs[1].Name)
	assert.Equal(t, "test2", namedConfigs[2].Name)
	assert.Equal(t, "test3", namedConfigs[3].Name)
}
