package encrypt

import "testing"

func TestEncrypt(t *testing.T) {
	type args struct {
		src string
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{src: "2020-12-15 10:00", key: "rtGrm478RY8KoiugzrdduHLB"},
			want:    "kfzsZsv4sKvhORwz9STRdKTvrVZ9z6Dd",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt3Des(tt.args.src, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		enStr string
		key   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test2",
			args:    args{enStr: "kfzsZsv4sKvhORwz9STRdKTvrVZ9z6Dd", key: "rtGrm478RY8KoiugzrdduHLB"},
			want:    "2020-12-15 10:00",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt3Des(tt.args.enStr, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
