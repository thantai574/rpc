package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		service string
		env     string
	}
	tests := []struct {
		name    string
		args    args
		want    *Logger
		wantErr bool
	}{
		{args: args{
			service: "service-user",
			env:     "production",
		}},
		{args: args{
			service: "service-user",
			env:     "dev",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLogger(tt.args.service, tt.args.env)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Info("test")
		})
	}
}
