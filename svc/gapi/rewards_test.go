package gapi

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func Test_calculateUserRewardsForGame(t *testing.T) {
	type args struct {
		conf       configer
		nftType    string
		complexity int32
		result     int32
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "win",
			args: args{
				conf:       &calculateUserRewardsForGameConfigerMock{intVal: 10, f64Val: 4.8},
				nftType:    NFTTypeCommon.String(),
				complexity: GameLevelEasy,
				result:     int32(GameResultWin),
			},
			want:    48,
			wantErr: false,
		},
		{
			name: "lose",
			args: args{
				conf:       &calculateUserRewardsForGameConfigerMock{intVal: 5, f64Val: 4.8},
				nftType:    NFTTypeCommon.String(),
				complexity: GameLevelEasy,
				result:     int32(GameResultLose),
			},
			want:    24,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := calculateUserRewardsForGame(tt.args.conf, tt.args.nftType, tt.args.complexity, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateUserRewardsForGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateUserRewardsForGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

type calculateUserRewardsForGameConfigerMock struct {
	boolVal  bool
	strVal   string
	f64Val   float64
	intVal   int
	int32Val int32
	jsonVal  []byte
	durVal   time.Duration
}

func (c *calculateUserRewardsForGameConfigerMock) GetBool(ctx context.Context, key string) (bool, error) {
	return c.boolVal, nil
}

func (c *calculateUserRewardsForGameConfigerMock) GetString(ctx context.Context, key string) (string, error) {
	return c.strVal, nil
}

func (c *calculateUserRewardsForGameConfigerMock) GetFloat64(ctx context.Context, key string) (float64, error) {
	return c.f64Val, nil
}

func (c *calculateUserRewardsForGameConfigerMock) GetInt(ctx context.Context, key string) (int, error) {
	return c.intVal, nil
}

func (c *calculateUserRewardsForGameConfigerMock) GetInt32(ctx context.Context, key string) (int32, error) {
	return c.int32Val, nil
}

func (c *calculateUserRewardsForGameConfigerMock) GetJSON(ctx context.Context, key string, result interface{}) error {
	return json.Unmarshal(c.jsonVal, result)
}

func (c *calculateUserRewardsForGameConfigerMock) GetDurration(ctx context.Context, key string) (time.Duration, error) {
	return c.durVal, nil
}
