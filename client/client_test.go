package client

import (
	"context"
	"reflect"
	"testing"

	"github.com/aftab-hussain-93/crypto-price-finder-microservice/types"
)

func TestClient_FindPrice(t *testing.T) {
	type fields struct {
		url string
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.FindPriceResponse
		wantErr bool
	}{
		{
			name: "testing positive case",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				ctx: context.Background(),
				key: "ETH",
			},
			want: &types.FindPriceResponse{
				Price:  2_000,
				Ticker: "ETH",
			},
		},
		{
			name: "testing negative case",
			fields: fields{
				url: "http://localhost:3000",
			},
			args: args{
				ctx: context.Background(),
				key: "ET",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				url: tt.fields.url,
			}
			got, err := c.FindPrice(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FindPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.FindPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
