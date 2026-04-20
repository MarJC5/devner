package database

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"myapp", false},
		{"my_app_2", false},
		{"a", false},
		{"", true},
		{"1app", true},
		{"My-App", true},
		{"my app", true},
		{"'; DROP TABLE x;--", true},
	}
	for _, tt := range tests {
		err := ValidateName(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateName(%q) err=%v wantErr=%v", tt.name, err, tt.wantErr)
		}
	}
}

func TestGeneratePassword(t *testing.T) {
	a := GeneratePassword()
	b := GeneratePassword()
	if a == b {
		t.Error("passwords should be unique")
	}
	if len(a) != 32 {
		t.Errorf("password length = %d, want 32", len(a))
	}
}
