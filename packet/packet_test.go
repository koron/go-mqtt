package packet

import "testing"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func compareBytes(t *testing.T, actual, expected []byte) {
	la, lb := len(actual), len(expected)
	n := min(la, lb)
	for i := 0; i < n; i++ {
		if actual[i] != expected[i] {
			t.Errorf("actual[%d] (=%02x) != expected[%d] (=%02x)", i, actual[i], i, expected[i])
			return
		}
	}
	if la != lb {
		if la > lb {
			t.Errorf("len(actual)=%d > len(expected)=%d actual[%d]=%02x",
			la, lb, lb, actual[lb])
		} else {
			t.Errorf("len(actual)=%d < len(expected)=%d expected[%d]=%02x",
			la, lb, la, expected[la])
		}
	}
}
