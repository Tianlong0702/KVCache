package test_util

import (
	"testing"
)

// Verify two string variables are equel
func VerifyStringValue(t *testing.T, expectedValue, actualValue string) {
	if expectedValue != actualValue {
		t.Fatalf("Unable to match. expectedValue: %s, actualValue: %s", expectedValue, actualValue)
	}
}

// Verify two bool variables are equel
func VerifyBoolValue(t *testing.T, expectedValue, actualValue bool) {
	if expectedValue != actualValue {
		t.Fatalf("Unable to match. expectedValue: %v, actualValue: %v", expectedValue, actualValue)
	}
}

// verify is the val is nil
func IsNil(t *testing.T, val interface{}) {
	if val != nil {
		t.Fatal("value is not nil ", val)
	}
}
