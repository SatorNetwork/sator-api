package gapi

import (
	"testing"
	"time"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
)

func Test_recoveryEnergyPoints(t *testing.T) {
	type args struct {
		player               repository.UnityGamePlayer
		energyFull           int32
		energyRecoveryPeriod time.Duration
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		{
			name: "recovery points: 0",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     10,
					EnergyRefilledAt: time.Now().Add(-time.Hour * 2),
				},
				energyFull:           10,
				energyRecoveryPeriod: time.Hour,
			},
			want: 0,
		},
		{
			name: "recovery points: 1",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     9,
					EnergyRefilledAt: time.Now().Add(-time.Hour * 2),
				},
				energyFull:           10,
				energyRecoveryPeriod: time.Hour * 2,
			},
			want: 1,
		},
		{
			name: "recovery points: max",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     0,
					EnergyRefilledAt: time.Now().Add(-time.Hour * 11),
				},
				energyFull:           10,
				energyRecoveryPeriod: time.Hour,
			},
			want: 10,
		},
		{
			name: "recovery points: recovery period is not reached",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     0,
					EnergyRefilledAt: time.Now().Add(-time.Minute * 11),
				},
				energyFull:           10,
				energyRecoveryPeriod: time.Hour,
			},
			want: 0,
		},

		{
			name: "recovery points: empty energy full",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     0,
					EnergyRefilledAt: time.Now().Add(-time.Hour * 11),
				},
				energyRecoveryPeriod: time.Hour,
			},
			want: 3,
		},
		{
			name: "recovery points: empty energy recovery period",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints:     0,
					EnergyRefilledAt: time.Now().Add(-time.Hour * 11),
				},
			},
			want: 2,
		},
		{
			name: "recovery points: empty energy refilled at",
			args: args{
				player: repository.UnityGamePlayer{
					EnergyPoints: 0,
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := recoveryEnergyPoints(tt.args.player, tt.args.energyFull, tt.args.energyRecoveryPeriod); got != tt.want {
				t.Errorf("recoveryEnergyPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}
