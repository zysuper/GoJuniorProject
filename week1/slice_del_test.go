package week1

import "testing"

func TestIntSliceDel(t *testing.T) {
	var emptySlice []int
	oneEleSlice := []int{1}
	moreEleSlice := []int{1, 2, 3}

	if _, err := intSliceDel(emptySlice, 0); err == nil {
		t.Error("must out of range")
	}

	if _, err := intSliceDel(oneEleSlice, -1); err == nil {
		t.Error("must out of range")
	}

	if _, err := intSliceDel(oneEleSlice, 1); err == nil {
		t.Error("must out of range")
	}

	// head del
	r, err := intSliceDel(oneEleSlice, 0)
	if err != nil {
		t.Error("must del success")
	}

	logCompare(t, "head del", oneEleSlice, r)

	// tail del
	r2, err := intSliceDel(moreEleSlice, 2)
	if err != nil {
		t.Error("must del success")
	}

	logCompare(t, "tail del", moreEleSlice, r2)

	// mid del
	r3, err := intSliceDel(moreEleSlice, 1)
	if err != nil {
		t.Error("must del success")
	}

	logCompare(t, "mid del", moreEleSlice, r3)

	if cap(r3) != 2 {
		t.Error("没有实现缩容操作")
	}
}

func logCompare(t *testing.T, msg string, l []int, r []int) {
	t.Logf("%v      : %+v, len:%d cap:%d", msg, l, len(l), cap(l))
	t.Logf("%v after: %+v, len:%d cap:%d", msg, r, len(r), cap(r))
}
