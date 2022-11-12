package dibf

import (
	"encoding/binary"
	"testing"
)

func TestBasic(t *testing.T) {
	f := New(1000, 4, 1000/20, nil)
	n1 := []byte("Utkarsh")
	n2 := []byte("ABCXYZ")
	f.Add(n1)

	n1b := f.Contains(n1)
	n2b := f.Contains(n2)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
}

func TestBasicWithRemove(t *testing.T) {
	f := New(1000, 4, 1000/20, nil)
	n1 := []byte("Utkarsh")
	n2 := []byte("ABCXYZ")
	f.Add(n1)

	n1b := f.Contains(n1)
	n2b := f.Contains(n2)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}

	f.Add(n2)
	n2b = f.Contains(n2)
	if !n2b {
		t.Errorf("%v should be in.", n2)
	}

	f.Delete(n1)
	n1b = f.Contains(n1)
	if n1b {
		t.Errorf("%v should not be in.", n1)
	}
}

func TestNewWithLowNumbers(t *testing.T) {
	f := New(0, 0, 0, nil)
	if f.k != 1 {
		t.Errorf("%v should be 1", f.k)
	}
	if f.m != 1 {
		t.Errorf("%v should be 1", f.m)
	}
}

func testEstimated(n uint, maxFp float64, t *testing.T) {
	f := NewWithEstimates(n, maxFp, 20, nil)
	m := f.m
	k := f.k
	fpRate := f.CurrentFalsePositiveRate()
	if fpRate > 1.5*maxFp {
		t.Errorf("False positive rate too high: n: %v; m: %v; k: %v; maxFp: %f; fpRate: %f, fpRate/maxFp: %f", n, m, k, maxFp, fpRate, fpRate/maxFp)
	}
}

func TestEstimated1000_0001(t *testing.T) {
	testEstimated(1000, 0.000100, t)
}

func TestEstimated10000_0001(t *testing.T) {
	testEstimated(10000, 0.000100, t)
}
func TestEstimated100000_0001(t *testing.T) {
	testEstimated(100000, 0.000100, t)
}

func TestEstimated1000_001(t *testing.T) {
	testEstimated(1000, 0.001000, t)
}

func TestEstimated10000_001(t *testing.T) {
	testEstimated(10000, 0.001000, t)
}

func TestEstimated100000_001(t *testing.T) {
	testEstimated(100000, 0.001000, t)
}

func TestEstimated1000_01(t *testing.T) {
	testEstimated(1000, 0.010000, t)
}

func TestEstimated10000_01(t *testing.T) {
	testEstimated(10000, 0.010000, t)
}

func TestEstimated100000_01(t *testing.T) {
	testEstimated(100000, 0.010000, t)
}

func TestCap(t *testing.T) {
	f := New(1000, 4, 1000/20, nil)
	if f.Cap() != f.m {
		t.Error("not accessing Cap() correctly")
	}
}

func TestK(t *testing.T) {
	f := New(1000, 4, 1000/20, nil)
	if f.K() != f.k {
		t.Error("not accessing K() correctly")
	}
}

func BenchmarkEstimated(b *testing.B) {
	for n := uint(100000); n <= 100000; n *= 10 {
		for fp := 0.1; fp >= 0.0001; fp /= 10.0 {
			f := NewWithEstimates(n, fp, 20, nil)
			f.CurrentFalsePositiveRate()
		}
	}
}

func BenchmarkSeparateTestAndAdd(b *testing.B) {
	f := NewWithEstimates(uint(b.N), 0.0001, 20, nil)
	key := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		f.Contains(key)
		f.Add(key)
	}
}

func TestApproximatedSize(t *testing.T) {
	f := NewWithEstimates(1000, 0.001, 20, nil)
	f.Add([]byte("ABC"))
	f.Add([]byte("1234"))
	f.Add([]byte("XYZ"))
	f.Add([]byte("BBMMA1231341"))
	size := f.ApproximateCount()
	if size != 4 {
		t.Errorf("%d should equal 4.", size)
	}
}

func TestFPP(t *testing.T) {
	f := NewWithEstimates(1000, 0.001, 20, nil)
	for i := uint32(0); i < 1000; i++ {
		n := make([]byte, 4)
		binary.BigEndian.PutUint32(n, i)
		f.Add(n)
	}
	count := 0

	for i := uint32(0); i < 1000; i++ {
		n := make([]byte, 4)
		binary.BigEndian.PutUint32(n, i+1000)
		if f.Contains(n) {
			count += 1
		}
	}
	if float64(count)/1000.0 > 0.001 {
		t.Errorf("Excessive fpp: %d", count)
	}
}
