package quiz

import "testing"

func Test_calcRate(t *testing.T) {
	type args struct {
		questionTimeSec float64
		answerTimeSec   float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "rate: 0 (8s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   8,
			},
			want: 0,
		},
		{
			name: "rate: 0 (7s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   7,
			},
			want: 0,
		},
		{
			name: "rate: 1 (6s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   6,
			},
			want: 1,
		},
		{
			name: "rate: 1 (5s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   5,
			},
			want: 1,
		},
		{
			name: "rate: 2 (4s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   4,
			},
			want: 2,
		},
		{
			name: "rate: 2 (3s)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   3,
			},
			want: 2,
		},
		{
			name: "rate: 3 (2 sec)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   2,
			},
			want: 3,
		},
		{
			name: "rate: 3 (1 sec)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   1,
			},
			want: 3,
		},
		{
			name: "rate: 3 (1 sec)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   1,
			},
			want: 3,
		},
		{
			name: "rate: 3 (0 sec)",
			args: args{
				questionTimeSec: 8,
				answerTimeSec:   0.5,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcRate(tt.args.questionTimeSec, tt.args.answerTimeSec); got != tt.want {
				t.Errorf("calcRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
