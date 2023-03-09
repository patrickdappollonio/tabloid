package tabloid

import (
	"testing"
	"time"
)

func Test_isready(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				args: []interface{}{
					"1/1",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "more than 1 argument",
			args: args{
				args: []interface{}{
					"1/1",
					"2/2",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not a string",
			args: args{
				args: []interface{}{
					1,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not in the form of <current>/<total>",
			args: args{
				args: []interface{}{
					"1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not ready",
			args: args{
				args: []interface{}{
					"0/1",
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isready(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("isready() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEqual(t, got, tt.want, "isready() = %v, want %v", got, tt.want)
		})
	}
}

func Test_hasrestarts(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "basic no restarts",
			args: args{
				args: []interface{}{
					"0",
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "basic with restarts",
			args: args{
				args: []interface{}{
					"1",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "basic with restarts and time",
			args: args{
				args: []interface{}{
					"1 (5s ago)",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "more than 1 argument",
			args: args{
				args: []interface{}{
					"1",
					"2",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not a string",
			args: args{
				args: []interface{}{
					1,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hasrestarts(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("hasrestarts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, got, tt.want, "hasrestarts() = %v, want %v", got, tt.want)
		})
	}
}

func Test_parseDurations(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		ret1    time.Duration
		ret2    time.Duration
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				args: []interface{}{
					"1h",
					"2h",
				},
			},
			ret1:    time.Hour,
			ret2:    2 * time.Hour,
			wantErr: false,
		},
		{
			name: "using days",
			args: args{
				args: []interface{}{
					"1d",
					"2d",
				},
			},
			ret1:    24 * time.Hour,
			ret2:    2 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name: "using weeks",
			args: args{
				args: []interface{}{
					"1w",
					"2w",
				},
			},
			ret1:    7 * 24 * time.Hour,
			ret2:    2 * 7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name: "single argument",
			args: args{
				args: []interface{}{
					"1h",
				},
			},
			wantErr: true,
		},
		{
			name: "more than 2 arguments",
			args: args{
				args: []interface{}{
					"1h",
					"2h",
					"3h",
				},
			},
			wantErr: true,
		},
		{
			name: "not a string",
			args: args{
				args: []interface{}{
					1,
					"2h",
				},
			},
			wantErr: true,
		},
		{
			name: "not a string 2nd place",
			args: args{
				args: []interface{}{
					"1h",
					1,
				},
			},
			wantErr: true,
		},
		{
			name: "not a valid duration",
			args: args{
				args: []interface{}{
					"1",
					"1h",
				},
			},
			wantErr: true,
		},
		{
			name: "not a valid duration 2nd place",
			args: args{
				args: []interface{}{
					"1h",
					"1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseDurations(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDurations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, got, tt.ret1, "parseDurations() got = %v, want %v", got, tt.ret1)
			assertEqual(t, got1, tt.ret2, "parseDurations() got1 = %v, want %v", got1, tt.ret2)
		})
	}
}

func TestDurations(t *testing.T) {
	cases := []struct {
		name         string
		d1           string
		d2           string
		isOlder      bool
		isOlderEqual bool
		isNewer      bool
		isNewerEqual bool
		isEqualDur   bool
	}{
		{
			name:         "d1 is older",
			d1:           "3h",
			d2:           "1h",
			isOlder:      true,
			isOlderEqual: true,
			isNewer:      false,
			isNewerEqual: false,
			isEqualDur:   false,
		},
		{
			name:         "d1 is newer",
			d1:           "1h",
			d2:           "3h",
			isOlder:      false,
			isOlderEqual: false,
			isNewer:      true,
			isNewerEqual: true,
			isEqualDur:   false,
		},
		{
			name:         "d1 is equal to d2",
			d1:           "1h",
			d2:           "1h",
			isOlder:      false,
			isOlderEqual: true,
			isNewer:      false,
			isNewerEqual: true,
			isEqualDur:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotIsOlder, err := olderThan(c.d1, c.d2)
			if err != nil {
				t.Errorf("olderThan() error = %v", err)
				return
			}

			gotIsOlderEqual, err := olderThanEq(c.d1, c.d2)
			if err != nil {
				t.Errorf("olderThanEq() error = %v", err)
				return
			}

			gotIsNewer, err := newerThan(c.d1, c.d2)
			if err != nil {
				t.Errorf("newerThan() error = %v", err)
				return
			}

			gotIsNewerEqual, err := newerThanEq(c.d1, c.d2)
			if err != nil {
				t.Errorf("newerThanEq() error = %v", err)
				return
			}

			gotIsEqDuration, err := eqduration(c.d1, c.d2)
			if err != nil {
				t.Errorf("eqduration() error = %v", err)
				return
			}

			assertEqual(t, gotIsOlder, c.isOlder, "olderThan() = %t, want %t", gotIsOlder, c.isOlder)
			assertEqual(t, gotIsOlderEqual, c.isOlderEqual, "olderThanEq() = %t, want %t", gotIsOlderEqual, c.isOlderEqual)
			assertEqual(t, gotIsNewer, c.isNewer, "newerThan() = %t, want %t", gotIsNewer, c.isNewer)
			assertEqual(t, gotIsNewerEqual, c.isNewerEqual, "newerThanEq() = %t, want %t", gotIsNewerEqual, c.isNewerEqual)
			assertEqual(t, gotIsEqDuration, c.isEqualDur, "eqduration() = %t, want %t", gotIsEqDuration, c.isEqualDur)
		})
	}
}

func Test_compareFirstArgumentToString(t *testing.T) {
	type args struct {
		str  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "equal",
			args: args{
				str: "test",
				args: []interface{}{
					"test",
				},
			},
			want: true,
		},
		{
			name: "not equal",
			args: args{
				str: "test",
				args: []interface{}{
					"test2",
				},
			},
			want: false,
		},
		{
			name: "more than 1 argument",
			args: args{
				str: "test",
				args: []interface{}{
					"test",
					"test2",
				},
			},
			wantErr: true,
		},
		{
			name: "less than 1 argument",
			args: args{
				str:  "test",
				args: []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "not a string",
			args: args{
				str: "test",
				args: []interface{}{
					1,
				},
			},
			wantErr: true,
		},
		{
			name: "empty string",
			args: args{
				str: "",
				args: []interface{}{
					"",
				},
			},
			want: true,
		},
		{
			name: "nil value",
			args: args{
				str: "test",
				args: []interface{}{
					nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compareFirstArgumentToString("testfn", tt.args.str, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("compareFirstArgumentToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("compareFirstArgumentToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
