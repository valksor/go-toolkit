package helper_test

import (
	"testing"
	"time"
)

func TestWriteAndReadFile(t *testing.T) {
	dir := TempDir(t)
	path := CreateFile(t, dir, "test.txt", "hello world")
	content := ReadFile(t, path)
	AssertEqual(t, content, "hello world")
}

func TestAssertions(t *testing.T) {
	t.Run("AssertEqual", func(t *testing.T) {
		AssertEqual(t, 1, 1)
		AssertEqual(t, "foo", "foo")
	})

	t.Run("AssertTrue", func(t *testing.T) {
		AssertTrue(t, true, "should be true")
	})

	t.Run("AssertFalse", func(t *testing.T) {
		AssertFalse(t, false, "should be false")
	})

	t.Run("AssertNoError", func(t *testing.T) {
		AssertNoError(t, nil)
	})

	t.Run("AssertLen", func(t *testing.T) {
		AssertLen(t, []int{1, 2, 3}, 3)
	})

	t.Run("AssertContains", func(t *testing.T) {
		AssertContains(t, []string{"a", "b", "c"}, "b")
	})

	t.Run("AssertNotContains", func(t *testing.T) {
		AssertNotContains(t, []string{"a", "b", "c"}, "d")
	})
}

func TestPointerHelpers(t *testing.T) {
	b := BoolPtr(true)
	if *b != true {
		t.Error("BoolPtr failed")
	}

	s := StringPtr("test")
	if *s != "test" {
		t.Error("StringPtr failed")
	}

	i := IntPtr(42)
	if *i != 42 {
		t.Error("IntPtr failed")
	}

	f := Float64Ptr(3.14)
	if *f != 3.14 {
		t.Error("Float64Ptr failed")
	}

	tm := TimePtr(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC))
	if !tm.Equal(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)) {
		t.Error("TimePtr failed")
	}
}

func TestLogger(t *testing.T) {
	log := Logger()
	if log == nil {
		t.Error("Logger returned nil")
	}

	discard := DiscardLogger()
	if discard == nil {
		t.Error("DiscardLogger returned nil")
	}
}
