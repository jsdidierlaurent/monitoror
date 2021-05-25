package delivery

import (
	"net/http/httptest"
	"testing"

	"github.com/monitoror/monitoror/internal/pkg/monitorable/params"
	"github.com/monitoror/monitoror/internal/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type (
	Params1 struct {
		params.Default
		Field string `query:"field" validate:"required"`
	}

	Params2 struct {
		Field string `query:"field" validate:"required"`
	}

	Params3 struct {
		Field string `query:"field" validate:"required"`
	}

	Params4 struct {
		Field chan string `query:"field" validate:"required"`
	}

	Params5 map[string]string
)

func (p *Params2) Validate() []validator.Error { return nil }

func (p *Params3) Validate() []validator.Error {
	return []validator.Error{validator.NewDefaultError("Field", "boom")}
}

func (m *Params4) Validate() []validator.Error { return nil }

func (m *Params5) Validate() []validator.Error { return nil }

func TestBindAndValidateParams(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v1/xxx", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	p := &Params1{}
	err := BindAndValidateParams(ctx, p)
	assert.Error(t, err)
	assert.Equal(t, `Required "field" field is missing.`, err.Error())

	ctx.QueryParams().Add("field", "test")

	p2 := &Params2{}
	err = BindAndValidateParams(ctx, p2)
	assert.NoError(t, err)

	p3 := &Params3{}
	err = BindAndValidateParams(ctx, p3)
	assert.Error(t, err)
	assert.Equal(t, `Invalid "field" field. Must be boom.`, err.Error())

	req = httptest.NewRequest(echo.GET, "/api/v1/xxx?field=test", nil)
	res = httptest.NewRecorder()
	ctx = e.NewContext(req, res)

	p4 := &Params4{}
	err = BindAndValidateParams(ctx, p4)
	assert.Error(t, err)
	assert.Equal(t, `invalid configuration, unable to parse request parameters`, err.Error())

	p5 := &Params5{}
	assert.Panics(t, func() {
		_ = BindAndValidateParams(ctx, p5)
	})
}
