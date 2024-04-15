package util

import (
	"github.com/google/uuid"
	"goboot/pkg/value"
	"strings"
)

func UUID() string {
	return value.Must(uuid.NewUUID()).String()
}

func UUID32() string {
	uid := UUID()
	return strings.ReplaceAll(uid, "-", "")
}
