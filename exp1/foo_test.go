package foo

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mock_foo "gomock_study/exp1/mock"
)

func TestSUT(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	m := mock_foo.NewMockFoo(ctl)
	bar := m.
		EXPECT().
		Bar(gomock.Eq(99)).
		DoAndReturn(func(_ int) int {
			time.Sleep(1 * time.Second)
			return 101
		}).
		AnyTimes()
	m.EXPECT().Bar1(gomock.Any()).After(bar)

	// Does not make any assertions. Returns 103 when Bar is invoked with 101.
	m.
		EXPECT().
		Bar(gomock.Eq(101)).
		Return(103).
		AnyTimes()

	SUT(m)
	type args struct {
		f Foo
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				f: m,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SUT(tt.args.f)
		})
	}
}
