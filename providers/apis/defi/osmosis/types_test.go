package osmosis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis"
)

func TestCreateURL(t *testing.T) {
	type args struct {
		baseURL    string
		poolID     uint64
		baseAsset  string
		quoteAsset string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				baseURL:    "http://localhost",
				poolID:     1,
				baseAsset:  "base",
				quoteAsset: "quote",
			},
			want: "http://localhost/osmosis/poolmanager/v2/pools/1/prices?base_asset_denom=base&quote_asset_denom" +
				"=quote",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := osmosis.CreateURL(tt.args.baseURL, tt.args.poolID, tt.args.baseAsset, tt.args.quoteAsset)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
