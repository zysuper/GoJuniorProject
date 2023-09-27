package week1

import "testing"

func TestSliceDel(t *testing.T) {
	var emptySlice []int
	oneEleSlice := []int{1}
	moreEleSlice := []int{1, 2, 3}

	t.Log(">> test []int....")
	testAnySlice(t, emptySlice, oneEleSlice, moreEleSlice)

	var s0 []string
	s1 := []string{"a"}
	sm := []string{"a", "b", "c"}

	t.Log(">> test []string....")
	testAnySlice(t, s0, s1, sm)
}

func testAnySlice[T any](t *testing.T, emptySlice []T, oneEleSlice []T, moreEleSlice []T) {
	if _, err := sliceDel(emptySlice, 0); err == nil {
		t.Error("must out of range")
	}

	if _, err := sliceDel(oneEleSlice, -1); err == nil {
		t.Error("must out of range")
	}

	if _, err := sliceDel(oneEleSlice, 1); err == nil {
		t.Error("must out of range")
	}

	// head del
	r, err := sliceDel(oneEleSlice, 0)
	if err != nil {
		t.Error("must del success")
	}

	logCompareGeneric(t, "head del", oneEleSlice, r)

	// tail del
	r2, err := sliceDel(moreEleSlice, 2)
	if err != nil {
		t.Error("must del success")
	}

	logCompareGeneric(t, "tail del", moreEleSlice, r2)

	// mid del
	r3, err := sliceDel(moreEleSlice, 1)
	if err != nil {
		t.Error("must del success")
	}

	logCompareGeneric(t, "mid del", moreEleSlice, r3)

	if cap(r3) != 2 {
		t.Error("没有实现缩容操作")
	}
}

func logCompareGeneric[T any](t *testing.T, msg string, l []T, r []T) {
	t.Logf("%v      : %+v, len:%d cap:%d", msg, l, len(l), cap(l))
	t.Logf("%v after: %+v, len:%d cap:%d", msg, r, len(r), cap(r))
}
