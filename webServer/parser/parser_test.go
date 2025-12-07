package parser

import "testing"

func TestParseAndClassify_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantN       int
		wantParity  Parity
		wantSign    Sign
		expectError bool
	}{
		{"simple even positive", "42", 42, Even, NonNegative, false},
		{"simple odd negative", "-3", -3, Odd, Negative, false},
		{"zero", "0", 0, Even, NonNegative, false},
		{"minus zero", "-0", 0, Even, NonNegative, false}, // strconv treats "-0" as 0
		{"leading/trailing spaces", "  7  ", 7, Odd, NonNegative, false},
		{"plus sign", "+8", 8, Even, NonNegative, false},
		{"invalid string", "abc", 0, Even, NonNegative, true},
		{"empty", "", 0, Even, NonNegative, true},
		{"large overflow", "9223372036854775808", 0, Even, NonNegative, true}, // too big for int on 64-bit
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseAndClassify(tc.input)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil (input=%q)", tc.input)
				}
				// мы не проверяем текст ошибки здесь — достаточно, что она есть
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v (input=%q)", err, tc.input)
			}
			assertNoEqual(t, "N", got.N, tc.wantN)
			assertNoEqual(t, "Parity", got.Parity, tc.wantParity)
			assertNoEqual(t, "Sign", got.Sign, tc.wantSign)
		})
	}
}

func assertNoEqual[T comparable](t *testing.T, name string, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v want %v", name, got, want)
	}
}
