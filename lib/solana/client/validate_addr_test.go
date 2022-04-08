package client

import (
	"fmt"
	"testing"

	"github.com/portto/solana-go-sdk/types"
)

func TestValidateSolanaWalletAddr(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid public key",
			args: args{
				addr: types.NewAccount().PublicKey.ToBase58(),
			},
			wantErr: false,
		},
		{
			name: "invalid public key",
			args: args{
				addr: "invalid",
			},
			wantErr: true,
		},

		{
			name: "invalid public key: too long",
			args: args{
				addr: fmt.Sprintf("%s%s", types.NewAccount().PublicKey.ToBase58(), "q"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSolanaWalletAddr(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSolanaWalletAddr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
