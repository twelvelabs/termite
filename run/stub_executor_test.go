package run

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStubExecutor_StubbingMethods(t *testing.T) {
	executor := NewStubExecutor()
	executor.RegisterStub(
		MatchString("/bin/date"),
		StringResponse("Sun Nov 13 22:00:00 CST 2022"),
	)
	executor.RegisterStub(
		MatchString("/bin/date"),
		StringResponse("Sun Nov 13 22:05:00 CST 2022"),
	)

	cmd := NewClient().Command("/bin/echo")
	err := executor.Run(cmd)
	assert.ErrorContains(t, err, "no registered stubs matching: /bin/echo")

	cmd = NewClient().Command("/bin/date")
	buf, err := executor.Output(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "Sun Nov 13 22:00:00 CST 2022", string(buf))

	buf, err = executor.Output(cmd)
	assert.NoError(t, err)
	assert.Equal(t, "Sun Nov 13 22:05:00 CST 2022", string(buf))

	buf, err = executor.Output(cmd)
	assert.ErrorContains(t, err, "wanted 3 of only 2 stubs matching: /bin/date")
	assert.Nil(t, buf)
}

func TestStubExecutor_VerifyWhenNoStubs(t *testing.T) {
	mt := &mockTest{}
	executor := NewStubExecutor()

	executor.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubExecutor_VerifyWhenAllStubsMatched(t *testing.T) {
	mt := &mockTest{}
	executor := NewStubExecutor()
	executor.RegisterStub(
		MatchString("/bin/echo"),
		StringResponse(""),
	)

	cmd := NewClient().Command("/bin/echo")
	err := executor.Run(cmd)
	assert.NoError(t, err)

	executor.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, false, mt.ErrorfCalled)
}

func TestStubExecutor_VerifyWhenUnmatchedStubs(t *testing.T) {
	mt := &mockTest{}
	executor := NewStubExecutor()
	executor.RegisterStub(
		MatchString("/bin/echo"),
		StringResponse(""),
	)

	executor.VerifyStubs(mt)
	assert.Equal(t, true, mt.HelperCalled)
	assert.Equal(t, true, mt.ErrorfCalled)
	assert.Equal(t, "found 1 unmatched stub(s)", mt.Msg)
}

type mockTest struct {
	Msg          string
	HelperCalled bool
	ErrorfCalled bool
}

func (mt *mockTest) Helper() {
	mt.HelperCalled = true
}
func (mt *mockTest) Errorf(line string, args ...interface{}) {
	mt.ErrorfCalled = true
	mt.Msg = fmt.Sprintf(line, args...)
}
