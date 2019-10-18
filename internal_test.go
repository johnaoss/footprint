package footprint

import (
	"os"
	"testing"
)

// Test_parsePerms checks if we can properly parse the linux permission bits
func Test_parsePerms(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    os.FileMode
		wantErr bool
	}{
		{
			name: "0777",
			arg:  "drwxrwxrwx",
			want: 0777,
		},
		{
			name: "0000",
			arg:  "l---------",
			want: 0,
		},
		{
			name: "0070",
			arg:  "----rwx---",
			want: 0070,
		},
		{
			name:    "Invalid Char",
			arg:     "drwxrwxrwc",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePerms([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePerms() error = %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePerms() = %v, want %v", got, tt.want)
			}
		})
	}
}
