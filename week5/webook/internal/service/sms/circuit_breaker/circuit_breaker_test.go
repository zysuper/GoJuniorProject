package cb

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCircuitBreakerService_Do(t *testing.T) {
	tests := []struct {
		name        string
		state       int32
		cnt         int32
		threshold   int32
		recoverTime time.Duration
		doWhat      func(args ...any) error
		wantErr     error
		wantState   int32
		wantCnt     int32
	}{
		{
			name:        "Open 状态,一直很正常",
			state:       Open,
			cnt:         1,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return nil
			},
			wantState: Open,
			wantCnt:   int32(1),
			wantErr:   nil,
		},
		{
			name:        "Open 状态,开局不顺正常",
			state:       Open,
			cnt:         0,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return errors.New("呜呜呜，干翻了")
			},
			wantState: Open,
			wantCnt:   int32(1),
			wantErr:   errors.New("呜呜呜，干翻了"),
		},
		{
			name:        "Open 状态逆转了",
			state:       Open,
			cnt:         2,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return errors.New("呜呜呜，干翻了")
			},
			wantState: Close,
			wantCnt:   int32(0),
			wantErr:   errors.New("呜呜呜，干翻了"),
		},
		{
			name:        "HalfOpen 到 Open 状态",
			state:       HalfOpen,
			cnt:         2,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return nil
			},
			wantState: Open,
			wantCnt:   int32(0),
			wantErr:   nil,
		},
		{
			name:        "HalfOpen 态在努力努力",
			state:       HalfOpen,
			cnt:         1,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return nil
			},
			wantState: HalfOpen,
			wantCnt:   2,
			wantErr:   nil,
		},
		{
			name:        "HalfOpen 态中道崩图",
			state:       HalfOpen,
			cnt:         2,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(args ...any) error {
				return errors.New("呜呜呜，干翻了")
			},
			wantState: Close,
			wantCnt:   0,
			wantErr:   errors.New("呜呜呜，干翻了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &CircuitBreakerService{
				state:       tt.state,
				cnt:         tt.cnt,
				threshold:   tt.threshold,
				recoverTime: tt.recoverTime,
				doWhat:      tt.doWhat,
			}
			err := svc.Do()
			assert.Equal(t, tt.wantState, svc.state)
			assert.Equal(t, tt.wantCnt, svc.cnt)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCircuitBreakerService_Do_Close2Half(t *testing.T) {
	tests := []struct {
		name        string
		state       int32
		cnt         int32
		threshold   int32
		recoverTime time.Duration
		doWhat      func(...any) error
		wantErr     error
		wantState   int32
		wantCnt     int32
	}{
		{name: "Open 状态逆转了, 过一会儿变 half 了,",
			state:       Open,
			cnt:         2,
			threshold:   3,
			recoverTime: time.Second,
			doWhat: func(...any) error {
				return errors.New("呜呜呜，干翻了")
			},
			wantState: HalfOpen,
			wantCnt:   int32(0),
			wantErr:   errors.New("呜呜呜，干翻了"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &CircuitBreakerService{
				state:       tt.state,
				cnt:         tt.cnt,
				threshold:   tt.threshold,
				recoverTime: tt.recoverTime,
				doWhat:      tt.doWhat,
			}
			err := svc.Do()
			// 睡一会儿，等到 half 状态.
			time.Sleep(time.Second * 2)
			assert.Equal(t, tt.wantState, svc.state)
			assert.Equal(t, tt.wantCnt, svc.cnt)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
