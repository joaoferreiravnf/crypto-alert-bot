package models

import (
	"encoding/json"
	"testing"
)

func TestFloat64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:    "Valid float string",
			input:   `"3.14159"`,
			want:    3.14159,
			wantErr: false,
		},
		{
			name:    "Valid float integer string",
			input:   `"42"`,
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "Invalid float string",
			input:   `"not-a-float"`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   `""`,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var f Float64

			err := json.Unmarshal([]byte(tt.input), &f)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && float64(f) != tt.want {
				t.Errorf("Unmarshaled value = %v, want %v", float64(f), tt.want)
			}
		})
	}
}

func TestFloat64Method(t *testing.T) {
	f := Float64(10.5)
	got := f.Float64()
	want := 10.5

	if got != want {
		t.Errorf("Float64() = %v, want %v", got, want)
	}
}
