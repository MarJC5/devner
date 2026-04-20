package network

import "testing"

func TestNeedsHostsEntry(t *testing.T) {
	tests := []struct {
		domain string
		want   bool
	}{
		{"myapp.localhost", false},
		{"MYAPP.LOCALHOST", false},
		{"myapp.test", true},
		{"myapp.dev", true},
		{"foo.bar.localhost", false},
	}
	for _, tt := range tests {
		got := NeedsHostsEntry(tt.domain)
		if got != tt.want {
			t.Errorf("NeedsHostsEntry(%q) = %v, want %v", tt.domain, got, tt.want)
		}
	}
}
