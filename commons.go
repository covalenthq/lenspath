package lenspath

import (
	"reflect"
	"testing"
)

func checkSetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, setFail bool) {
	lp, err := Create(lens)
	if err != nil {
		t.Fatalf("create: Expected no error, got %v", err)
		return
	}

	containsArr := false
	var expectedArr reflect.Value
	for _, lensv := range lens {
		if lensv == "*" {
			containsArr = true
			expectedArr = reflect.ValueOf(expectedValue)
			break
		}
	}
	index := 0

	err = lp.Setter(structure, func(value any) any {
		if containsArr {
			if index >= expectedArr.Len() {
				t.Fatalf("set: expectedValue array length mismatch")
				return nil
			}
			defer func() { index++ }()
			return expectedArr.Index(index).Interface()
		}

		return expectedValue
	})

	if err != nil && !setFail {
		t.Errorf("set: Expected no error, got %v", err)
		return
	} else if err == nil && setFail {
		t.Errorf("set: Expected error, got %v", err)
		return
	}

	index = 0

	err = lp.Getter(structure, func(value any) any {
		if !containsArr {
			comparev(t, value, expectedValue)
		} else {
			comparev(t, value, expectedArr.Index(index).Interface())
			index++
		}

		return nil
	})

	if err != nil {
		t.Errorf("og_get: Expected no error, got %v", err)
		return
	}
}

func checkGetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, createFail bool, getFail bool, assumeNil bool) {
	lp, err := Create(lens)

	switch {
	case err != nil && !createFail:
		t.Fatalf("Expected no error, got %v", err)

	case err == nil && createFail:
		t.Fatalf("Expected error, got %v", lp)

	case err != nil && createFail:
		// success
		return
	}

	lp.WithOptions(WithAssumeNil(assumeNil))

	containsArr := false
	var expectedArr reflect.Value
	for _, lensv := range lens {
		if lensv == "*" {
			containsArr = true
			expectedArr = reflect.ValueOf(expectedValue)
			break
		}
	}
	index := 0

	exec_err := lp.Getter(structure, func(value any) any {
		if !containsArr {
			comparev(t, value, expectedValue)
		} else {
			comparev(t, value, expectedArr.Index(index).Interface())
			index++
		}

		return nil
	})

	if exec_err != nil && !getFail {
		t.Fatalf("Expected no error, got %v", exec_err)
	} else if exec_err == nil && getFail {
		t.Fatalf("Expected error, got %v", exec_err)
	} else if containsArr && index != expectedArr.Len() {
		t.Fatalf("expected array size mismatch")
	}

	// success

}

func comparev(t *testing.T, value any, expectedValue any) {
	if value == nil {
		if expectedValue != nil {
			t.Fatalf("Expected %v, got %v", expectedValue, value)
		} else {
			return
		}
	}
	kind := reflect.TypeOf(value).Kind()
	switch {
	case kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map:
		if !reflect.DeepEqual(value, expectedValue) {
			t.Fatalf("Expected %v, got %v", expectedValue, value)
		}

	case value != expectedValue:
		t.Fatalf("Expected %v, got %v", expectedValue, value)
	default:
		// fine
	}
}
