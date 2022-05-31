package gapi

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

func Test_calculateElectricityCost(t *testing.T) {
	type args struct {
		conf    configer
		nftType string
		result  int32
		rewards float64
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
				conf:    &configerMock{intVal: 40, strVal: GameResultWin.String()},
				nftType: NFTTypeCommon.String(),
				result:  int32(GameResultWin),
				rewards: 100,
			},
			want:    40,
			wantErr: false,
		},
		{
			name: "lose",
			args: args{
				conf:    &configerMock{intVal: 40, strVal: GameResultLose.String()},
				nftType: NFTTypeCommon.String(),
				result:  int32(GameResultLose),
				rewards: 100,
			},
			want:    40,
			wantErr: false,
		},
		{
			name: "always",
			args: args{
				conf:    &configerMock{intVal: 40, strVal: GameResultWin.String()},
				nftType: NFTTypeCommon.String(),
				result:  int32(GameResultWin),
				rewards: 100,
			},
			want:    40,
			wantErr: false,
		},
		{
			name: "float",
			args: args{
				conf:    &configerMock{intVal: 32, strVal: GameResultWin.String()},
				nftType: NFTTypeLegend.String(),
				result:  int32(GameResultWin),
				rewards: 1450,
			},
			want:    464,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateElectricityCost(tt.args.conf, tt.args.nftType, tt.args.result, tt.args.rewards)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateElectricityCost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateElectricityCost() = %v, want %v", got, tt.want)
			}
		})
	}
}

type configerMock struct {
	boolVal  bool
	strVal   string
	f64Val   float64
	intVal   int
	int32Val int32
	jsonVal  []byte
	durVal   time.Duration
}

func (c *configerMock) GetBool(ctx context.Context, key string) (bool, error) {
	return c.boolVal, nil
}

func (c *configerMock) GetString(ctx context.Context, key string) (string, error) {
	return c.strVal, nil
}

func (c *configerMock) GetFloat64(ctx context.Context, key string) (float64, error) {
	return c.f64Val, nil
}

func (c *configerMock) GetInt(ctx context.Context, key string) (int, error) {
	return c.intVal, nil
}

func (c *configerMock) GetInt32(ctx context.Context, key string) (int32, error) {
	return c.int32Val, nil
}

func (c *configerMock) GetJSON(ctx context.Context, key string, result interface{}) error {
	return json.Unmarshal(c.jsonVal, result)
}

func (c *configerMock) GetDurration(ctx context.Context, key string) (time.Duration, error) {
	return c.durVal, nil
}
