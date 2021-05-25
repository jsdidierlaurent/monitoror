package validate

import (
	"regexp"

	"github.com/go-playground/validator"

	pkgValidator "github.com/monitoror/monitoror/internal/pkg/validator"
)

// -------------------------------------------------------
// This file is an helper / wrapper of github.com/go-playground/validator/v10
// Its purpose is to validate monitorable struct using official validator tags
// Like:
//	type Test struct {
//    Field  string `validate:"required"`
//    Number int `validate:"required,gte=200,lte=400"`
//  }
//
// Each error raised by the validator is wrapped to be used by the cli, the APIs or when verifying the config
// I also added custom validators for monitoror needs
// -------------------------------------------------------

const (
	regexTag    = "regex"    // Valid when regex compile
	httpTag     = "http"     // Valid when string starts with http:// or https://
	notEmptyTag = "notempty" // Valid when slice is not empty (like gt=0 but with custom message / expected)

	HTTPRegex = `^https?://`
)

var (
	validateTagMapping = map[string]pkgValidator.ErrorID{
		"required":  pkgValidator.ErrorRequired,
		"eq":        pkgValidator.ErrorEq,
		"ne":        pkgValidator.ErrorNE,
		"oneof":     pkgValidator.ErrorOneOf,
		"gte":       pkgValidator.ErrorGTE,
		"gt":        pkgValidator.ErrorGT,
		"lte":       pkgValidator.ErrorLTE,
		"lt":        pkgValidator.ErrorLT,
		"url":       pkgValidator.ErrorURL,
		notEmptyTag: pkgValidator.ErrorNotEmpty,
		httpTag:     pkgValidator.ErrorHTTP,
		regexTag:    pkgValidator.ErrorRegex,
	}
)

// use a single instance of Struct, it caches struct info
var validate *validator.Validate
var httpRegex *regexp.Regexp

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation(notEmptyTag, validateNotEmpty)
	_ = validate.RegisterValidation(httpTag, validateHTTP)
	_ = validate.RegisterValidation(regexTag, validateRegex)

	httpRegex = regexp.MustCompile(HTTPRegex)
}

func Struct(s interface{}) []pkgValidator.Error {
	var errors []pkgValidator.Error

	if err := validate.Struct(s); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			panic("unsupported echo bind parameters type. must be a struct.")
		}

		// range over all validate validateError to bind then into ValidatorError
		for _, err := range err.(validator.ValidationErrors) {
			id, exists := validateTagMapping[err.Tag()]
			if !exists {
				panic("unsupported validate tag. use a tag listed in validateTagMapping instead.")
			}

			e := validateError{
				errorID:   id,
				fieldName: err.Field(),
				tagParam:  err.Param(),
			}

			errors = append(errors, &e)
		}
	}

	return errors
}

// validateRegex implements validator.Func
func validateRegex(fl validator.FieldLevel) bool {
	_, err := regexp.Compile(fl.Field().String())
	return err == nil
}

// validateHTTP implements validator.Func
func validateHTTP(fl validator.FieldLevel) bool {
	return httpRegex.MatchString(fl.Field().String())
}

// validateHTTP implements validator.Func
func validateNotEmpty(fl validator.FieldLevel) bool {
	return fl.Field().Len() != 0
}
