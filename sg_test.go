package main

import "testing"

func Test_isPublicIP(t *testing.T) {
	type args struct {
		ipaddr string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "is public ip",
			args:    args{ipaddr: "1.1.1.1"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "is private ip",
			args:    args{ipaddr: "192.168.1.1"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "is cidr public ip",
			args:    args{ipaddr: "1.1.1.1/0"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "is cidr private ip",
			args:    args{ipaddr: "192.168.1.1/0"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "is cidr invalid ip",
			args:    args{ipaddr: "192.168.1/0"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isPublicIP(tt.args.ipaddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("isPublicIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isPublicIP() got = %v, want %v", got, tt.want)
			}
		})
	}
}
