package project

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"myapp", false},
		{"my-app-2", false},
		{"a", false},
		{"", true},
		{"1app", true},
		{"My-App", true},
		{"my_app", true}, // underscores not allowed in project names (reserved for db)
		{"app space", true},
	}
	for _, tt := range tests {
		err := ValidateName(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateName(%q) err=%v wantErr=%v", tt.name, err, tt.wantErr)
		}
	}
}

func TestDocRoot(t *testing.T) {
	tests := []struct {
		t    Type
		name string
		want string
	}{
		{Laravel, "myapp", "/var/www/html/myapp/public"},
		{WordPress, "blog", "/var/www/html/blog"},
		{NextJS, "web", "/var/www/html/web"},
	}
	for _, tt := range tests {
		got := tt.t.DocRoot(tt.name)
		if got != tt.want {
			t.Errorf("%s.DocRoot(%s) = %s, want %s", tt.t, tt.name, got, tt.want)
		}
	}
}

func TestParseType(t *testing.T) {
	for _, s := range []string{"wordpress", "LARAVEL", "Node", "nextjs", "astro"} {
		if _, err := ParseType(s); err != nil {
			t.Errorf("ParseType(%q) err=%v", s, err)
		}
	}
	if _, err := ParseType("react"); err == nil {
		t.Error("ParseType(react) should fail")
	}
}
