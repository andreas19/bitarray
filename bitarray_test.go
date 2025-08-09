package bitarray

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	size := 10
	tests := []struct {
		args []int
		want []uint8
	}{
		{[]int{}, []uint8{0, 0}},
		{[]int{0}, []uint8{1, 0}},
		{[]int{1}, []uint8{0b10, 0}},
		{[]int{7, 0}, []uint8{0b10000001, 0}},
		{[]int{8}, []uint8{0, 1}},
		{[]int{9}, []uint8{0, 0b10}},
		{[]int{9, 0}, []uint8{1, 0b10}},
		{[]int{9, 7, 1}, []uint8{0b10000010, 0b10}},
		{[]int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, []uint8{0b11111111, 0b11}},
	}
	for i, test := range tests {
		ba := New(size, test.args...)
		if got := ba.data; ba.size != size || !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: got %d and %v, want %d and %v", i, ba.size, got, size, test.want)
		}
	}
}

func TestNewPanic(t *testing.T) {
	defer func() { recover() }()
	New(0)
	t.Error("did not panic")
}

func TestMustParse(t *testing.T) {
	size := 10
	tests := []struct {
		arg  string
		want []uint8
	}{
		{"0000000000", []uint8{0, 0}},
		{"0000000001", []uint8{1, 0}},
		{"0000000010", []uint8{0b10, 0}},
		{"0010000001", []uint8{0b10000001, 0}},
		{"0100000000", []uint8{0, 1}},
		{"1000000000", []uint8{0, 0b10}},
		{"1000000001", []uint8{1, 0b10}},
		{"1010000010", []uint8{0b10000010, 0b10}},
		{"1111111111", []uint8{0b11111111, 0b11}},
		{"11 11111111", []uint8{0b11111111, 0b11}},
	}
	for i, test := range tests {
		ba := MustParse(test.arg)
		if got := ba.data; ba.size != size || !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: got %d and %v, want %d and %v", i, ba.size, got, size, test.want)
		}
	}
}

func TestMustParsePanic(t *testing.T) {
	defer func() { recover() }()
	MustParse("012")
	t.Error("did not panic")
}

func TestString(t *testing.T) {
	tests := []string{
		"0000000000",
		"0000000001",
		"0000000010",
		"0010000001",
		"0100000000",
		"1000000000",
		"1000000001",
		"1010000010",
		"1111111111",
		"0101",
		"01010101",
		"0101010101010101",
	}
	for i, test := range tests {
		if got := MustParse(test).String(); got != test {
			t.Errorf("%d: got %q, want %q", i, got, test)
		}
	}
}

func TestClone(t *testing.T) {
	tests := []string{
		"0000000000",
		"0000000001",
		"0000000010",
		"0010000001",
		"0100000000",
		"1000000000",
		"1000000001",
		"1010000010",
		"1111111111",
	}
	for i, test := range tests {
		ba := MustParse(test)
		if got := Clone(ba).String(); got != test {
			t.Errorf("%d: got %q, want %q", i, got, test)
		}
	}
}

func TestClear(t *testing.T) {
	s := "1010101010"
	want := "0000000000"
	ba := MustParse(s)
	ba.Clear()
	if got := ba.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetAll(t *testing.T) {
	tests := []struct {
		size int
		want string
	}{
		{4, "1111"},
		{8, "11111111"},
		{10, "1111111111"},
		{16, "1111111111111111"},
	}
	for i, test := range tests {
		ba := New(test.size)
		ba.SetAll()
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestSet(t *testing.T) {
	want := "0100000010"
	ba := New(10, 1)
	ba.Set(1)
	ba.Set(8)
	if got := ba.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetIdxNeg(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Set(-1)
	t.Error("did not panic")
}

func TestSetIdxSize(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Set(4)
	t.Error("did not panic")
}

func TestUnset(t *testing.T) {
	want := "1000000001"
	ba := New(10, 9, 0, 7)
	ba.Unset(7)
	ba.Unset(8)
	if got := ba.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestUnsetIdxNeg(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Unset(-1)
	t.Error("did not panic")
}

func TestUnsetIdxSize(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Unset(4)
	t.Error("did not panic")
}

func TestToggle(t *testing.T) {
	want := "1000000010"
	ba := New(10, 9, 0)
	a := ba.Toggle(0) // false
	b := ba.Toggle(1) // true
	if got := ba.String(); got != want || a || !b {
		t.Errorf("got %q, want %q (%v, %v)", got, want, a, b)
	}
}

func TestGet(t *testing.T) {
	s := "0100110101"
	tests := []struct {
		idx  int
		want bool
	}{
		{0, true},
		{1, false},
		{2, true},
		{3, false},
		{4, true},
		{5, true},
		{6, false},
		{7, false},
		{8, true},
		{9, false},
	}
	ba := MustParse(s)
	for i, test := range tests {
		if got := ba.Get(test.idx); got != test.want {
			t.Errorf("%d: got %t, want %t", i, got, test.want)
		}
	}
}

func TestGetIdxNeg(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Get(-1)
	t.Error("did not panic")
}

func TestGetIdxSize(t *testing.T) {
	defer func() { recover() }()
	ba := New(4)
	ba.Get(4)
	t.Error("did not panic")
}

func TestAnd(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   string
	}{
		{"0000", "1111", "0000"},
		{"1111", "0000", "0000"},
		{"0101", "0100", "0100"},
		{"00000000", "11111111", "00000000"},
		{"11111111", "00000000", "00000000"},
		{"01010101", "01000110", "01000100"},
		{"0000000000", "1111111111", "0000000000"},
		{"1111111111", "0000000000", "0000000000"},
		{"0101010101", "0100000110", "0100000100"},
		{"0000000000000000", "1111111111111111", "0000000000000000"},
		{"1111111111111111", "0000000000000000", "0000000000000000"},
		{"0101010101010101", "0100000000000110", "0100000000000100"},
	}
	for i, test := range tests {
		ba := MustParse(test.s1)
		ba.And(MustParse(test.s2))
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestAndDiffSize(t *testing.T) {
	defer func() { recover() }()
	New(4).And(New(5))
	t.Error("did not panic")
}

func TestOr(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   string
	}{
		{"0000", "1111", "1111"},
		{"1111", "0000", "1111"},
		{"0101", "0100", "0101"},
		{"00000000", "11111111", "11111111"},
		{"11111111", "00000000", "11111111"},
		{"01010101", "01000110", "01010111"},
		{"0000000000", "1111111111", "1111111111"},
		{"1111111111", "0000000000", "1111111111"},
		{"0101010101", "0100000110", "0101010111"},
		{"0000000000000000", "1111111111111111", "1111111111111111"},
		{"1111111111111111", "0000000000000000", "1111111111111111"},
		{"0101010101010101", "0100000000000110", "0101010101010111"},
	}
	for i, test := range tests {
		ba := MustParse(test.s1)
		ba.Or(MustParse(test.s2))
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestOrDiffSize(t *testing.T) {
	defer func() { recover() }()
	New(4).Or(New(5))
	t.Error("did not panic")
}

func TestXor(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   string
	}{
		{"0000", "1111", "1111"},
		{"1111", "0000", "1111"},
		{"0101", "0100", "0001"},
		{"00000000", "11111111", "11111111"},
		{"11111111", "00000000", "11111111"},
		{"01010101", "01000110", "00010011"},
		{"0000000000", "1111111111", "1111111111"},
		{"1111111111", "0000000000", "1111111111"},
		{"0101010101", "0100000110", "0001010011"},
		{"0000000000000000", "1111111111111111", "1111111111111111"},
		{"1111111111111111", "0000000000000000", "1111111111111111"},
		{"0101010101010101", "0100000000000110", "0001010101010011"},
	}
	for i, test := range tests {
		ba := MustParse(test.s1)
		ba.Xor(MustParse(test.s2))
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestXorDiffSize(t *testing.T) {
	defer func() { recover() }()
	New(4).Xor(New(5))
	t.Error("did not panic")
}

func TestAndNot(t *testing.T) {
	tests := []struct {
		s1, s2 string
		want   string
	}{
		{"0000", "1111", "0000"},
		{"1111", "0000", "1111"},
		{"0101", "0100", "0001"},
		{"00000000", "11111111", "00000000"},
		{"11111111", "00000000", "11111111"},
		{"01010101", "01000110", "00010001"},
		{"0000000000", "1111111111", "0000000000"},
		{"1111111111", "0000000000", "1111111111"},
		{"0101010101", "0100000110", "0001010001"},
		{"0000000000000000", "1111111111111111", "0000000000000000"},
		{"1111111111111111", "0000000000000000", "1111111111111111"},
		{"0101010101010101", "0100000000000110", "0001010101010001"},
	}
	for i, test := range tests {
		ba := MustParse(test.s1)
		ba.AndNot(MustParse(test.s2))
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestAndNotDiffSize(t *testing.T) {
	defer func() { recover() }()
	New(4).AndNot(New(5))
	t.Error("did not panic")
}

func TestNot(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{"0000", "1111"},
		{"1111", "0000"},
		{"0101", "1010"},
		{"00000000", "11111111"},
		{"11111111", "00000000"},
		{"01010101", "10101010"},
		{"0000000000", "1111111111"},
		{"1111111111", "0000000000"},
		{"0101010101", "1010101010"},
		{"0000000000000000", "1111111111111111"},
		{"1111111111111111", "0000000000000000"},
		{"0101010101010101", "1010101010101010"},
	}
	for i, test := range tests {
		ba := MustParse(test.s)
		ba.Not()
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		s, want string
	}{
		{"0", "0"},
		{"1", "1"},
		{"10", "01"},
		{"110", "011"},
		{"01010101", "10101010"},
		{"010101011", "110101010"},
		{"0101010101010101", "1010101010101010"},
		{"01010101010101011", "11010101010101010"},
	}
	for i, test := range tests {
		ba := MustParse(test.s)
		ba.Reverse()
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"0001", 0, "0001"},
		{"0001", 1, "0010"},
		{"0001", 3, "1000"},
		{"0001", 4, "0001"},
		{"0001", 5, "0010"},
		{"0001", -1, "1000"},
		{"0001", -4, "0001"},
		{"00000001", 7, "10000000"},
		{"00000001", 8, "00000001"},
		{"00000001", 9, "00000010"},
		{"00000001", -1, "10000000"},
		{"0100000001", 0, "0100000001"},
		{"0100000001", 1, "1000000010"},
		{"0100000001", 2, "0000000101"},
		{"0100000001", 10, "0100000001"},
		{"0100000001", 11, "1000000010"},
		{"0100000001", -1, "1010000000"},
		{"0100000001", -2, "0101000000"},
		{"0100000001", -11, "1010000000"},
		{"0100000001", -12, "0101000000"},
	}
	for i, test := range tests {
		ba := MustParse(test.s)
		ba.Rotate(test.n)
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestShift(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"0001", 0, "0001"},
		{"0001", 1, "0010"},
		{"0001", 3, "1000"},
		{"0001", 4, "0000"},
		{"0001", 5, "0000"},
		{"1000", -1, "0100"},
		{"1000", -4, "0000"},
		{"00000001", 7, "10000000"},
		{"00000001", 8, "00000000"},
		{"00000001", 9, "00000000"},
		{"10000000", -1, "01000000"},
		{"0100000001", 0, "0100000001"},
		{"0100000001", 1, "1000000010"},
		{"0100000001", 2, "0000000100"},
		{"0100000001", 10, "0000000000"},
		{"0100000001", 11, "0000000000"},
		{"0100000001", -1, "0010000000"},
		{"0100000001", -2, "0001000000"},
		{"0100000001", -11, "0000000000"},
		{"0100000001", -12, "0000000000"},
	}
	for i, test := range tests {
		ba := MustParse(test.s)
		ba.Shift(test.n)
		if got := ba.String(); got != test.want {
			t.Errorf("%d: got %q, want %q", i, got, test.want)
		}
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		s1, s2 string
	}{
		{"0101", "010"},
		{"0101", "1010"},
		{"01010101", "10101010"},
		{"0101010101", "1010101010"},
	}
	for _, test := range tests {
		ba1a := MustParse(test.s1)
		ba1b := MustParse(test.s1)
		ba2a := MustParse(test.s2)
		ba2b := MustParse(test.s2)
		if !ba1a.Equal(ba1b) {
			t.Errorf("%v = %v: got false, want true", ba1a, ba1b)
		}
		if !ba2a.Equal(ba2b) {
			t.Errorf("%v = %v: got false, want true", ba2a, ba2b)
		}
		if ba1a.Equal(ba2a) {
			t.Errorf("%v = %v: got true, want false", ba1a, ba2a)
		}
		if ba2a.Equal(ba1a) {
			t.Errorf("%v = %v: got true, want false", ba2a, ba1a)
		}
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"0000", 0},
		{"1111", 4},
		{"0101", 2},
		{"00000000", 0},
		{"11111111", 8},
		{"01010101", 4},
		{"0000000000", 0},
		{"1111111111", 10},
		{"0101010101", 5},
		{"0000000000000000", 0},
		{"1111111111111111", 16},
		{"0101010101010101", 8},
	}
	for i, test := range tests {
		if got := MustParse(test.s).Count(); got != test.want {
			t.Errorf("%d: got %d, want %d", i, got, test.want)
		}
	}
}

func TestSize(t *testing.T) {
	want := 4
	tests := []*BitArray{New(4), MustParse("1010")}
	for i, test := range tests {
		if got := test.Size(); got != want {
			t.Errorf("%d: got %d, want %d", i, got, want)
		}
	}
}

func TestLeadingTrailingZeros(t *testing.T) {
	tests := []struct {
		s              string
		want_l, want_t int
	}{
		{"1001", 0, 0},
		{"0110", 1, 1},
		{"0000", 4, 4},
		{"10000001", 0, 0},
		{"01000010", 1, 1},
		{"00000000", 8, 8},
		{"1000000001", 0, 0},
		{"0100000010", 1, 1},
		{"0000000000", 10, 10},
		{"1000000000000001", 0, 0},
		{"0100000000000010", 1, 1},
		{"0000000000000000", 16, 16},
	}
	for i, test := range tests {
		ba := MustParse(test.s)
		got_l := ba.LeadingZeros()
		got_t := ba.TrailingZeros()
		if got_l != test.want_l || got_t != test.want_t {
			t.Errorf("%d: got %d and %d, want %d and %d", i, got_l, test.want_l, got_t, test.want_t)
		}
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		s1, s2 string
	}{
		{"0101", "101"},
		{"0101", "1010"},
		{"0101", "10101"},
		{"01010101", "1010"},
		{"01010101", "10101010"},
		{"0101010101", "1010"},
		{"0101010101", "101010"},
		{"0101010101", "10101010"},
	}
	for i, test := range tests {
		ba := Concat(MustParse(test.s1), MustParse(test.s2))
		want := test.s1 + test.s2
		if got := ba.String(); got != want {
			t.Errorf("%d: got %q, want %q", i, got, want)
		}
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	tests := []string{
		"0101", "01010101", "0101010101", "0101010101010101",
	}
	for i, test := range tests {
		ba1 := MustParse(test)
		buf, err := ba1.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		ba2 := new(BitArray)
		err = ba2.UnmarshalBinary(buf)
		if err != nil {
			t.Fatal(err)
		}
		if got := ba2.String(); got != test {
			t.Errorf("%d: got %q, want %q", i, got, test)
		}
	}
}
