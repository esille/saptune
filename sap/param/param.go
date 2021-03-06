package param

/*
Parameters are tunable switches on the system that are tuned in very specific
ways.

Each tunable parameter is able to inspect the system to determine the current
value, calculate a new value according to system environment and/or user input,
and be able to apply the new value upon request.

There's only one way to tune a parameter, however a parameter can be referred
to by one or more SAP notes.
*/

import (
	"math"
)

/*
Parameter is a tunable parameter, usually calculated based on user input and/or automatically.
Parameter is immutable. Internal state changes can only be made to copies.
*/
type Parameter interface {
	Inspect() (Parameter, error)                             // Read the parameter values from current system.
	Optimise(additionalInput interface{}) (Parameter, error) // Calculate new values based on internal states, and return a copy of new states.
	Apply() error                                            // Apply the parameter value without further calculation.
}

// MaxI64 returns the maximum among the input values.
// If there isn't any input value, return 0.
func MaxI64(values ...int64) int64 {
	var ret int64
	for _, value := range values {
		if ret < value {
			ret = value
		}
	}
	return ret
}

// MaxI returns the maximum among the input values.
// If there isn't any input value, return 0.
func MaxI(values ...int) int {
	var ret int
	for _, value := range values {
		if ret < value {
			ret = value
		}
	}
	return ret
}

// MaxU64 returns the maximum among the input values.
// If there isn't any input value, return 0.
func MaxU64(values ...uint64) uint64 {
	var ret uint64
	for _, value := range values {
		if ret < value {
			ret = value
		}
	}
	return ret
}

// MinU64 returns the minimum among the input values.
// If there isn't any input value, return 0.
func MinU64(values ...uint64) uint64 {
	if len(values) == 0 {
		return 0
	}
	var ret uint64 = math.MaxUint64
	for _, value := range values {
		if ret > value {
			ret = value
		}
	}
	return ret
}
