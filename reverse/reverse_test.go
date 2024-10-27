package main

import "testing"

func TestReverse(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, test := range tests {
		if got := Reverse(test.in); got != test.want {
			t.Errorf("Reverse(%q) = %q, want %q", test.in, got, test.want)
		}
	}
}

func TestReverse2(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}

	for _, test := range tests {
		// what the
		if got := Reverse2(test.in); got != test.want {
			t.Errorf("Reverse2(%q) = %q, want %q", test.in, got, test.want)
		}
	}
}

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse("Hello, world")
	}
}

func BenchmarkReverse2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reverse2("Hello, world")
	}
}
