package utils

import (
	"fmt"
	"testing"
)

func TestHashId(t *testing.T) {
	fmt.Println(HashId(100, 256))
}
