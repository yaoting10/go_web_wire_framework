package util_test

import (
	"goboot/pkg/util"
	"testing"
)

func TestNewUUID(t *testing.T) {
	println(util.UUID())
	println(util.UUID32())
}
