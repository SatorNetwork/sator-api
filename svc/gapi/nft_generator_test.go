package gapi

import (
	"testing"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/segmentio/ksuid"
)

func Test_craftNFT(t *testing.T) {
	type args struct {
		nftsToCraft []repository.UnityGameNft
	}
	tests := []struct {
		name    string
		args    args
		want    *NFTInfo
		wantErr bool
	}{
		{
			name: "should return error when nftsToCraft is empty",
			args: args{
				nftsToCraft: []repository.UnityGameNft{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error when nftsToCraft is less than 2",
			args: args{
				nftsToCraft: []repository.UnityGameNft{
					{
						ID:       "1",
						MaxLevel: 1,
						NftType:  "1",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return NFTInfo when nftsToCraft is more than 2",
			args: args{
				nftsToCraft: []repository.UnityGameNft{
					{
						ID:       "1",
						MaxLevel: 1,
						NftType:  NFTTypeCommon.String(),
					},
					{
						ID:       "2",
						MaxLevel: 2,
						NftType:  NFTTypeCommon.String(),
					},
				},
			},
			want: &NFTInfo{
				ID:       "3",
				MaxLevel: 2,
				NftType:  NFTTypeRare,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := craftNFT(tt.args.nftsToCraft)
			if (err != nil) != tt.wantErr {
				t.Errorf("craftNFT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.MaxLevel != tt.want.MaxLevel {
					t.Errorf("craftNFT() got = %v, want %v", got.MaxLevel, tt.want.MaxLevel)
				}
				if got.NftType != tt.want.NftType {
					t.Errorf("craftNFT() got = %v, want %v", got.NftType, tt.want.NftType)
				}
			} else if tt.want != nil && got == nil {
				t.Errorf("craftNFT() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateNFT(t *testing.T) {
	type args struct {
		nftPack repository.UnityGameNftPack
	}
	tests := []struct {
		name    string
		args    args
		want    *NFTInfo
		wantErr bool
	}{
		{
			name: "should return error when dropChances is empty",
			args: args{
				nftPack: repository.UnityGameNftPack{
					DropChances: []byte{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return NFTInfo when dropChances is valid",
			args: args{
				nftPack: repository.UnityGameNftPack{
					DropChances: []byte(`{"common":55.3,"rare":28.2,"super_rare":12.8,"epic":3.5,"legend":0.2}`),
				},
			},
			want: &NFTInfo{
				ID:       ksuid.New().String(),
				MaxLevel: getNFTLevelByType(NFTTypeRare),
				NftType:  NFTTypeCommon,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateNFT(tt.args.nftPack)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateNFT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.NftType != tt.want.NftType {
					t.Errorf("generateNFT() got = %v, want %v", got.NftType, tt.want.NftType)
				}
			}
		})
	}
}
