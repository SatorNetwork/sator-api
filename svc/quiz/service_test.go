package quiz

import (
	"testing"
)

func Test_calcPrize(t *testing.T) {
	type args struct {
		prizePool      float64
		pts            int
		rate           int
		totalWinners   int
		totalQuestions int
		totalPts       int
		totalRate      int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "pts = 0, single winner",
			args: args{
				prizePool:      250,
				pts:            0,
				rate:           0,
				totalWinners:   1,
				totalQuestions: 10,
				totalPts:       0,
				totalRate:      0,
			},
			want: 250,
		},
		{
			name: "pts = 9, single winner",
			args: args{
				prizePool:      250,
				pts:            9,
				rate:           0,
				totalWinners:   1,
				totalQuestions: 10,
				totalPts:       9,
				totalRate:      0,
			},
			want: 250,
		},
		{
			name: "pts = 10, two winners",
			args: args{
				prizePool:      250,
				pts:            10,
				rate:           0,
				totalWinners:   2,
				totalQuestions: 10,
				totalPts:       10,
				totalRate:      0,
			},
			want: 166.67,
		},
		{
			name: "pts = 0, two winners",
			args: args{
				prizePool:      250,
				pts:            0,
				rate:           0,
				totalWinners:   2,
				totalQuestions: 10,
				totalPts:       10,
				totalRate:      0,
			},
			want: 83.33,
		},
		{
			name: "pts = 10, winners = 5",
			args: args{
				prizePool:      250,
				pts:            10,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       50,
				totalRate:      0,
			},
			want: 50,
		},
		{
			name: "pts = 0, winners = 5",
			args: args{
				prizePool:      250,
				pts:            0,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       65,
				totalRate:      0,
			},
			want: 21.74,
		},
		{
			name: "pts = 5, winners = 5",
			args: args{
				prizePool:      250,
				pts:            5,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       65,
				totalRate:      0,
			},
			want: 32.61,
		},
		{
			name: "pts = 10, winners = 5",
			args: args{
				prizePool:      250,
				pts:            10,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       65,
				totalRate:      0,
			},
			want: 43.48,
		},
		{
			name: "pts = 20, winners = 5",
			args: args{
				prizePool:      250,
				pts:            20,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       65,
				totalRate:      0,
			},
			want: 65.22,
		},
		{
			name: "pts = 30, winners = 5",
			args: args{
				prizePool:      250,
				pts:            30,
				rate:           0,
				totalWinners:   5,
				totalQuestions: 10,
				totalPts:       65,
				totalRate:      0,
			},
			want: 86.96,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcPrize(
				tt.args.prizePool,
				tt.args.totalWinners,
				tt.args.totalQuestions,
				tt.args.totalPts,
				tt.args.totalRate,
				tt.args.pts,
				tt.args.rate,
			); toFixed(got, 2) != toFixed(tt.want, 2) {
				t.Errorf("calcPrize() = %v, want %v", toFixed(got, 2), toFixed(tt.want, 2))
			}
			// else {
			// 	t.Errorf("calcPrize() = %v, want %v", math.Round(got*100)/100, math.Round(tt.want*100)/100)
			// 	t.Errorf("calcPrize() = %v, want %v", toFixed(got, 2), toFixed(tt.want, 2))
			// }
		})
	}
}
