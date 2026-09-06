package test

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/charlieparkes/go-testcmp"
	"github.com/stretchr/testify/require"
)

func AssertStructFieldsEqual(t *testing.T, got, want any, ignore ...string) {
	t.Helper()

	publicFieldNames := func(v any) ([]string, error) {
		typ := reflect.TypeOf(v)
		for typ != nil && typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}
		if typ == nil || typ.Kind() != reflect.Struct {
			return nil, fmt.Errorf("got %T, want struct", v)
		}

		names := make([]string, 0, typ.NumField())
		for field := range typ.Fields() {
			if field.IsExported() && !slices.Contains(ignore, field.Name) {
				names = append(names, field.Name)
			}
		}
		return names, nil
	}

	gotNames, err := publicFieldNames(got)
	require.NoError(t, err)

	wantNames, err := publicFieldNames(want)
	require.NoError(t, err)

	slices.Sort(gotNames)
	slices.Sort(wantNames)
	testcmp.Compare(t, gotNames, wantNames)
	require.Equal(t, wantNames, gotNames)
}
