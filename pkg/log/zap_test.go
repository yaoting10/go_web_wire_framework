package log_test

import (
	"fmt"
	"goboot/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefLogger(t *testing.T) {
	go func() {
		err := recover()
		fmt.Printf("need an error: %v", err)
		assert.True(t, err != nil)
		log.Default.Debug("this is a debug log")
		log.Default.Info("this is an info log")
		log.Default.Warn("this is an warning log")
		log.Default.Error("this is an error log")
		log.Default.Fatal("this is a fatal log")
	}()
}
