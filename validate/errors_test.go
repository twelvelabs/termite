package validate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_panicToErr_WhenPanicNotString(t *testing.T) {
	err := panicToErr("key", "rules", errors.New("boom"))
	assert.ErrorContains(t, err, "boom")
}
