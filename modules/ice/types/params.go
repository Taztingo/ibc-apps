package types

import (
	"fmt"
)

const (
	// DefaultHostEnabled is the default value for the host param (set to true)
	DefaultHostEnabled = true
)

// NewParams creates a new parameter configuration
func NewParams(enableHost bool) Params {
	return Params{
		HostEnabled: enableHost,
	}
}

// DefaultParams is the default parameter configuration
func DefaultParams() Params {
	return NewParams(DefaultHostEnabled)
}

// Validate validates all parameters
func (p Params) Validate() error {
	return validateEnabled(p.HostEnabled)
}

func validateEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
