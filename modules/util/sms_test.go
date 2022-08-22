package util

import (
	"reflect"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

func Test_createCLient(t *testing.T) {
	type args struct {
		region    string
		keyID     string
		keySecret string
	}
	tests := []struct {
		name    string
		args    args
		want    *dysmsapi.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createCLient(tt.args.region, tt.args.keyID, tt.args.keySecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCLient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCLient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSendSms(t *testing.T) {
	type args struct {
		phone     string
		region    string
		keyID     string
		keySecret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "t1",
			args: args{
				phone:     "18810975701",
				region:    "cn-hangzhou",
				keyID:     "aFmL",
				keySecret: "9RBqn8mN2uFW",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendSms(tt.args.phone, tt.args.region, tt.args.keyID, tt.args.keySecret); (err != nil) != tt.wantErr {
				t.Errorf("SendSms() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name     string
		wantCode string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCode := GenerateCode(); gotCode != tt.wantCode {
				t.Errorf("GenerateCode() = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestCodeIsEq(t *testing.T) {
	type args struct {
		phone string
		code  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CodeIsEq(tt.args.phone, tt.args.code); got != tt.want {
				t.Errorf("CodeIsEq() = %v, want %v", got, tt.want)
			}
		})
	}
}
