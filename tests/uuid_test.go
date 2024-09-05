package tests

import (
	"github.com/google/uuid"
	"necron.dev/zed"
	"testing"
)

func TestUUID(t *testing.T) {
	testData := []string{
		uuid.New().String(),
		"00000000-0000-0000-0000-000000000000",
	}
	antitheses := []any{
		"foo",
		"bar",
		0,
		1,
	}
	f := zed.UUID("expected uuid value")
	testMulti[string, uuid.UUID](t, f, testData, false, uuidEq)
	testMulti[any, uuid.UUID](t, f, antitheses, true, nil)
}
