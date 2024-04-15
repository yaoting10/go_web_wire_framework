package io_test

import (
	"fmt"
	"goboot/pkg/io"
	"strings"
	"testing"
)

func TestWalkAllFiles(t *testing.T) {
	fs := io.WalkAllFiles("/Users/sam/Downloads/videos")
	fmt.Println(strings.Join(fs, "\n"))
}
