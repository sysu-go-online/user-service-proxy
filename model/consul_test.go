package model

import "testing"

func TestGetValueWithKey(t *testing.T) {
	type args struct {
		key string
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValueWithKey(tt.args.key, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValueWithKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetValueWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
