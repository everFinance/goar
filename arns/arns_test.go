package arns

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewArNS(t *testing.T) {
	type args struct {
		dreUrl   string
		arNSAddr string
		timout   time.Duration
	}
	tests := []struct {
		name string
		args args
		want *ArNS
	}{
		{
			name: "test",
			args: args{
				dreUrl:   "https://dre-3.warp.cc",
				arNSAddr: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				timout:   3 * time.Second,
			},
			want: &ArNS{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				HttpClient:  &http.Client{Timeout: 3 * time.Second},
			},
		},

		{
			name: "test",
			args: args{
				dreUrl:   "https://dre-3.warp.cc",
				arNSAddr: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				timout:   0,
			},
			want: &ArNS{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				HttpClient:  &http.Client{Timeout: 5 * time.Second},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewArNS(tt.args.dreUrl, tt.args.arNSAddr, tt.args.timout)
			if got.DreUrl != tt.want.DreUrl {
				t.Errorf("NewArNS() = %v, want %v", got.DreUrl, tt.want.DreUrl)
			}
			if got.ArNSAddress != tt.want.ArNSAddress {
				t.Errorf("NewArNS() = %v, want %v", got.ArNSAddress, tt.want.ArNSAddress)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArNS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArNS_GetArNSTxID(t *testing.T) {
	type fields struct {
		DreUrl      string
		ArNSAddress string
		Timeout     time.Duration
	}
	type args struct {
		caAddress string
		domain    string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantTxId string
		wantErr  bool
	}{
		{
			// "arseeding":{"contractTxId":"jr4P6Y_Olv3QGho0uo7p9DpvSn33mUC_XgJSKB3JDZ4","endTimestamp":1720865538,"tier":"a27dbfe4-6992-4276-91fb-5b97ae8c3ffa"}
			name: "test success",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				caAddress: "jr4P6Y_Olv3QGho0uo7p9DpvSn33mUC_XgJSKB3JDZ4",
				domain:    "@",
			},
			wantTxId: "wQk7txuMvlrlYlVozj6aeF7E9dlwar8nNtfs3iNTpbQ",
			wantErr:  false,
		},

		{
			name: "test error",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding",
			},
			wantTxId: "",
			wantErr:  true,
		},

		{
			name: "test error",
			fields: fields{
				DreUrl:      "",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arweave",
			},
			wantTxId: "",
			wantErr:  true,
		},

		{
			name: "test not exist",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding123456",
			},
			wantTxId: "",
			wantErr:  true,
		},

		// https://httpbin.org/xxx

		{
			name: "test 404",
			fields: fields{
				DreUrl:      "https://httpbin.org/xxx",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding123456",
			},
			wantTxId: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArNS(tt.fields.DreUrl, tt.fields.ArNSAddress, tt.fields.Timeout)
			gotTxId, err := a.GetArNSTxID(tt.args.caAddress, tt.args.domain)

			t.Logf("name: %v", tt.name)
			t.Logf("gotTxId: %v", gotTxId)
			t.Logf("wantErr: %v", tt.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArNSTxID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTxId != tt.wantTxId {
				t.Errorf("GetArNSTxID() gotTxId = %v, want %v", gotTxId, tt.wantTxId)
			}
		})
	}
}

func TestArNS_QueryLatestRecord(t *testing.T) {
	type fields struct {
		DreUrl      string
		ArNSAddress string
		Timeout     time.Duration
	}
	type args struct {
		domain string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantTxId string
		wantErr  bool
	}{
		{
			name: "test success",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding",
			},
			wantTxId: "wQk7txuMvlrlYlVozj6aeF7E9dlwar8nNtfs3iNTpbQ",
			wantErr:  false,
		},

		{
			name: "test init domain data success",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "web3infra",
			},
			wantTxId: "wQk7txuMvlrlYlVozj6aeF7E9dlwar8nNtfs3iNTpbQ",
			wantErr:  false,
		},
		{
			name: "test error",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "null",
			},
			wantTxId: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArNS(tt.fields.DreUrl, tt.fields.ArNSAddress, tt.fields.Timeout)
			gotTxId, err := a.QueryLatestRecord(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryLatestRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTxId != tt.wantTxId {
				t.Errorf("QueryLatestRecord() gotTxId = %v, want %v", gotTxId, tt.wantTxId)
			}
		})
	}
}

func TestArNS_QueryNameCa(t *testing.T) {
	type fields struct {
		DreUrl      string
		ArNSAddress string
		Timeout     time.Duration
	}
	type args struct {
		domain string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCaAddress string
		wantErr       bool
	}{
		{
			// "arseeding":{"contractTxId":"jr4P6Y_Olv3QGho0uo7p9DpvSn33mUC_XgJSKB3JDZ4","endTimestamp":1720865538,"tier":"a27dbfe4-6992-4276-91fb-5b97ae8c3ffa"}
			name: "test success",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding",
			},
			wantCaAddress: "jr4P6Y_Olv3QGho0uo7p9DpvSn33mUC_XgJSKB3JDZ4",
			wantErr:       false,
		},

		{
			name: "test init domain success",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "web3infra",
			},
			wantCaAddress: "Vx4bW_bh7nXMyq-Jy24s9EiCyY_BXZuToshhSqabc9o",
			wantErr:       false,
		},

		{
			name: "test error",
			fields: fields{
				DreUrl:      "",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arweave",
			},
			wantCaAddress: "",
			wantErr:       true,
		},

		{
			name: "test not exist",
			fields: fields{
				DreUrl:      "https://dre-3.warp.cc",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding123456",
			},
			wantCaAddress: "",
			wantErr:       true,
		},

		// https://httpbin.org/xxx

		{
			name: "test 404",
			fields: fields{
				DreUrl:      "https://httpbin.org/xxx",
				ArNSAddress: "bLAgYxAdX2Ry-nt6aH2ixgvJXbpsEYm28NgJgyqfs-U",
				Timeout:     10 * time.Second,
			},
			args: args{
				domain: "arseeding123456",
			},
			wantCaAddress: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("test %s", tt.name)
			t.Logf("CaAddr url %s", tt.fields.ArNSAddress)
			t.Logf("wantErr: %v", tt.wantErr)
			a := NewArNS(tt.fields.DreUrl, tt.fields.ArNSAddress, tt.fields.Timeout)
			gotCaAddress, err := a.QueryNameCa(tt.args.domain)
			t.Logf("gotCaAddress: %v", gotCaAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryNameCa() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCaAddress != tt.wantCaAddress {
				t.Errorf("QueryNameCa() gotCaAddress = %v, want %v", gotCaAddress, tt.wantCaAddress)
			}
		})
	}
}
