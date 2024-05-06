package verify

import (
	"testing"

	"github.com/ovechkin-dm/mockio/tests/common"

	. "github.com/ovechkin-dm/mockio/mock"
)

type iface interface {
	Foo(a int) int
}

type iface_mock_arg interface {
	MockAsArg(m iface)
}

func TestVerifyMockAsArg(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	m2 := Mock[iface_mock_arg]()

	m2.MockAsArg(m)

	// This works:
	// Verify(m2, Once()).MockAsArg(Exact[iface](m))
	// This doesn't:
	Verify(m2, Once()).MockAsArg(Any[iface]())

	r.AssertNoError()
}

func TestVerifySimple(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	Verify(m, Once()).Foo(10)
	r.AssertNoError()
}

func TestVerifyAny(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	Verify(m, Once()).Foo(Any[int]())
	r.AssertNoError()
}

func TestVerifyMultipleAny(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Times(2)).Foo(Any[int]())
	r.AssertNoError()
}

func TestVerifyNever(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Never()).Foo(13)
	r.AssertNoError()
}

func TestVerifyNeverFails(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Never()).Foo(10)
	r.AssertError()
}

func TestNoMoreInteractionsFails(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	VerifyNoMoreInteractions(m)
	r.AssertError()
}

func TestNoMoreInteractionsSuccess(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(Any[int]())).ThenReturn(10)
	m.Foo(10)
	Verify(m, Once()).Foo(10)
	VerifyNoMoreInteractions(m)
	r.AssertNoError()
}

func TestNoMoreInteractionsComplexFail(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(10)).ThenReturn(10)
	WhenSingle(m.Foo(11)).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, Once()).Foo(10)
	VerifyNoMoreInteractions(m)
	r.AssertError()
}

func TestNoMoreInteractionsComplexSuccess(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(10)).ThenReturn(10)
	WhenSingle(m.Foo(11)).ThenReturn(10)
	m.Foo(10)
	m.Foo(11)
	Verify(m, AtLeastOnce()).Foo(AnyInt())
	Verify(m, Once()).Foo(11)
	VerifyNoMoreInteractions(m)
	r.AssertNoError()
}

func TestVerifyInsideReturnerPass(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(AnyInt())).ThenReturn(11).Verify(Once())
	m.Foo(10)
	r.TriggerCleanup()
	r.AssertNoError()
}

func TestVerifyInsideReturnerNoMoreInteractionsFail(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(AnyInt())).ThenReturn(11).Verify(Once())
	VerifyNoMoreInteractions(m)
	r.AssertError()
}

func TestVerifyInsideReturnerFail(t *testing.T) {
	r := common.NewMockReporter(t)
	SetUp(r)
	m := Mock[iface]()
	WhenSingle(m.Foo(AnyInt())).ThenReturn(11).Verify(Once())
	r.TriggerCleanup()
	r.AssertError()
}
