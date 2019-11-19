package regist

import (
	"testing"
	"net/http"
)

func TestPlatformRegister_Regist(t *testing.T) {
	type fields struct {
		PlatformType string
		ReqMethod    string
		ReqUrl       string
	}
	type args struct {
		param *RegParam
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{name: PlatformTypeTBK, fields: struct {
			PlatformType string
			ReqMethod    string
			ReqUrl       string
		}{PlatformType: PlatformTypeTBK, ReqMethod: http.MethodPost, ReqUrl: "http://8553d.com/cn/register"},
			args: args{param: &RegParam{
				Account:      "itceshi004",
				Password:     "ffffff",
				RePassword:   "ffffff",
				RealName:     "程胜",
				Mobile:       "13111254456",
				WxNo:         "dfidf",
				WithdrawPass: "ffffff",
				Question:     "您所在的城市",
				Answer:       "北京",
			},
			},
			want: false,
			want1: "该名称已经被使用,请更换",
		},
		{name: "test2", fields: struct {
			PlatformType string
			ReqMethod    string
			ReqUrl       string
		}{PlatformType: PlatformTypeBoss, ReqMethod: http.MethodPost, ReqUrl: "https://y09n.com/signup#"},
			args: args{param: &RegParam{
				Account:      "itceshi008",
				Password:     "ffffff",
				RePassword:   "ffffff",
				WithdrawPass: "1234",
				RealName:     "程胜",
				Mobile:       "13117254456",
				WxNo:         "dfidf",
			},
			},
			want: true,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PlatformRegister{
				PlatformType: tt.fields.PlatformType,
				ReqMethod:    tt.fields.ReqMethod,
				ReqUrl:       tt.fields.ReqUrl,
			}
			got, got1 := a.Regist(tt.args.param)
			if got != tt.want {
				t.Errorf("PlatformRegister.Regist() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlatformRegister.Regist() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPlatformRegister_tbk(t *testing.T) {
	type fields struct {
		PlatformType string
		ReqMethod    string
		ReqUrl       string
	}
	type args struct {
		param *RegParam
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PlatformRegister{
				PlatformType: tt.fields.PlatformType,
				ReqMethod:    tt.fields.ReqMethod,
				ReqUrl:       tt.fields.ReqUrl,
			}
			got, got1 := a.tbk(tt.args.param)
			if got != tt.want {
				t.Errorf("PlatformRegister.tbk() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PlatformRegister.tbk() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
