//go:generate mockery --name Monitorable

package models

type Monitorable interface {
	//GetDisplayName return monitorable name display in console
	GetDisplayName() string

	//GetVariantsNames return variant list extract from config
	GetVariantsNames() []VariantName

	//Validate test if config variant is valid
	// return false if empty and error if config have an error (ex: wrong url format)
	Validate(variantName VariantName) (bool, []error)

	//Enable monitorable variant (add route to echo and enable tile for config verify / hydrate)
	Enable(variantName VariantName)
}
