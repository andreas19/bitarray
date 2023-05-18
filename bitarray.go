// Package bitarray implements a [BitArray].
//
// The least significant bit is at index 0. In the string representation used by
// [Parse], [MustParse], and [BitArray.String] it is the rightmost digit.
package bitarray

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"math/bits"
	"strings"
)

const bitsN = 8

// BitArray type.
type BitArray struct {
	size int
	data []uint8
}

// New creates a new BitArray with size bits and the bits at the given indexes
// set to 1. Panics if size <= 0 or one of the indexes is out of range.
func New(size int, idx ...int) *BitArray {
	if size <= 0 {
		panic("size must be > 0")
	}
	n, r := size/bitsN, size%bitsN
	if r > 0 {
		n++
	}
	ba := BitArray{size, make([]uint8, n)}
	for _, i := range idx {
		ba.Set(i)
	}
	return &ba
}

// Parse creates a new BitArray by parsing the given string. Space characters are ignored.
// Returns an error if one of the characters in the string is not space, 0, or 1.
func Parse(s string) (*BitArray, error) {
	rs := []rune(strings.ReplaceAll(s, " ", ""))
	rsLen := len(rs)
	ba := New(len(rs))
	for i, c := range rs {
		if c == '1' {
			ba.set(rsLen - 1 - i)
		} else if rs[i] != '0' {
			return nil, fmt.Errorf("unknown character: %c", rs[i])
		}
	}
	return ba, nil
}

// MustParse creates a new BitArray by parsing the given string. Space characters are
// ignored. Panics if one of the characters in the string is not space, 0, or 1.
func MustParse(s string) *BitArray {
	ba, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return ba
}

// Clone clones the BitArray.
func Clone(ba *BitArray) *BitArray {
	sl := make([]uint8, len(ba.data))
	copy(sl, ba.data)
	return &BitArray{ba.size, sl}
}

// Clear sets all bits to 0.
func (ba *BitArray) Clear() {
	for i := range ba.data {
		ba.data[i] = 0
	}
}

// SetAll sets all bits to 1.
func (ba *BitArray) SetAll() {
	if x := ba.size % bitsN; x == 0 {
		ba.data[len(ba.data)-1] = math.MaxUint8
	} else {
		ba.data[len(ba.data)-1] = (2 << (x - 1)) - 1
	}
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] = math.MaxUint8
	}
}

// Set sets the bit at index idx to 1.
func (ba *BitArray) Set(idx int) {
	ba.checkIdx(idx)
	ba.set(idx)
}

func (ba *BitArray) set(idx int) {
	n, i := idx/bitsN, idx%bitsN
	ba.data[n] |= 1 << i
}

// Unset sets the bit at index idx to 0.
func (ba *BitArray) Unset(idx int) {
	ba.checkIdx(idx)
	ba.unset(idx)
}

func (ba *BitArray) unset(idx int) {
	n, i := idx/bitsN, idx%bitsN
	ba.data[n] &^= 1 << i
}

// Get reports whether the bit at index idx is set.
func (ba *BitArray) Get(idx int) bool {
	ba.checkIdx(idx)
	return ba.get(idx)
}

func (ba *BitArray) get(idx int) bool {
	n, i := idx/bitsN, idx%bitsN
	return ba.data[n]&(1<<i) != 0
}

// And sets ba = ba & other (bitwise AND).
func (ba *BitArray) And(other *BitArray) {
	ba.checkSize(other)
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] &= other.data[i]
	}
	if ba.size%bitsN == 0 {
		i := len(ba.data) - 1
		ba.data[i] &= other.data[i]
	} else {
		for i := (ba.size / bitsN) * bitsN; i < ba.size; i++ {
			if !other.get(i) {
				ba.unset(i)
			}
		}
	}
}

// Or sets ba = ba | other (bitwise OR).
func (ba *BitArray) Or(other *BitArray) {
	ba.checkSize(other)
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] |= other.data[i]
	}
	if ba.size%bitsN == 0 {
		i := len(ba.data) - 1
		ba.data[i] |= other.data[i]
	} else {
		for i := (ba.size / bitsN) * bitsN; i < ba.size; i++ {
			if other.get(i) {
				ba.set(i)
			}
		}
	}
}

// Xor sets ba = ba ^ other (bitwise XOR).
func (ba *BitArray) Xor(other *BitArray) {
	ba.checkSize(other)
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] ^= other.data[i]
	}
	if ba.size%bitsN == 0 {
		i := len(ba.data) - 1
		ba.data[i] ^= other.data[i]
	} else {
		for i := (ba.size / bitsN) * bitsN; i < ba.size; i++ {
			b1 := ba.get(i)
			b2 := other.get(i)
			b := (b1 || b2) && !(b1 && b2)
			if b && !b1 {
				ba.set(i)
			} else if !b && b1 {
				ba.unset(i)
			}
		}
	}
}

// AndNot sets ba = ba &^ other (bit clear).
func (ba *BitArray) AndNot(other *BitArray) {
	ba.checkSize(other)
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] &^= other.data[i]
	}
	if ba.size%bitsN == 0 {
		i := len(ba.data) - 1
		ba.data[i] &^= other.data[i]
	} else {
		for i := (ba.size / bitsN) * bitsN; i < ba.size; i++ {
			if ba.get(i) && other.get(i) {
				ba.unset(i)
			}
		}
	}
}

// Not sets ba = ^ba.
func (ba *BitArray) Not() {
	for i := 0; i < len(ba.data)-1; i++ {
		ba.data[i] = ^ba.data[i]
	}
	if ba.size%bitsN == 0 {
		i := len(ba.data) - 1
		ba.data[i] = ^ba.data[i]
	} else {
		for i := (ba.size / bitsN) * bitsN; i < ba.size; i++ {
			if ba.get(i) {
				ba.unset(i)
			} else {
				ba.set(i)
			}
		}
	}
}

// Rotate rotates the bit array by |n| bits. If n > 0 to the left, if n < 0 to the right.
func (ba *BitArray) Rotate(n int) {
	n = n % ba.size
	var aux *BitArray
	if n > 0 {
		aux = Slice(ba, ba.size-n, ba.size)
	} else if n < 0 {
		aux = Slice(ba, 0, -n)
	}
	start, end := ba.moveBits(n)
	for i := start; i < end; i++ {
		if aux.get(i - start) {
			ba.set(i)
		} else {
			ba.unset(i)
		}
	}
}

// Shift shifts the bit array by |n| bits. If n > 0 to the left, if n < 0 to the right.
func (ba *BitArray) Shift(n int) {
	if n > 0 && n >= ba.size || n < 0 && -n >= ba.size {
		ba.Clear()
		return
	}
	start, end := ba.moveBits(n)
	for i := start; i < end; i++ {
		ba.unset(i)
	}
}

func (ba *BitArray) moveBits(n int) (int, int) {
	if n == 0 {
		return 0, 0
	}
	var start, end int
	if n > 0 {
		for i := ba.size - 1; i >= n; i-- {
			if ba.get(i - n) {
				ba.set(i)
			} else {
				ba.unset(i)
			}
		}
		end = n
	} else {
		n = -n
		for i := 0; i < ba.size-n; i++ {
			if ba.get(i + n) {
				ba.set(i)
			} else {
				ba.unset(i)
			}
		}
		start = ba.size - n
		end = ba.size
	}
	return start, end
}

// Equal reports whether the two bit arrays are equal.
func (ba *BitArray) Equal(other *BitArray) bool {
	if ba.size != other.size {
		return false
	}
	for i := range ba.data {
		if ba.data[i] != other.data[i] {
			return false
		}
	}
	return true
}

// Count returns the number of set bits.
func (ba *BitArray) Count() int {
	cnt := 0
	for _, x := range ba.data {
		cnt += bits.OnesCount8(x)
	}
	return cnt
}

// LeadingZeros returns the number of leading unset bits.
func (ba *BitArray) LeadingZeros() int {
	cnt := 0
	for i := len(ba.data) - 1; i >= 0; i-- {
		if ba.data[i] == 0 {
			cnt += bitsN
		} else {
			cnt += bits.LeadingZeros8(ba.data[i])
			break
		}
	}
	if x := ba.size % bitsN; x != 0 {
		cnt -= bitsN - x
	}
	return cnt
}

// TrailingZeros returns the number of trailing unset bits.
func (ba *BitArray) TrailingZeros() int {
	cnt := 0
	for _, x := range ba.data {
		if x == 0 {
			cnt += bitsN
		} else {
			cnt += bits.TrailingZeros8(x)
			break
		}
	}
	if cnt > ba.size {
		cnt = ba.size
	}
	return cnt
}

// Size returns the size of the bit array.
func (ba *BitArray) Size() int {
	return ba.size
}

// String returns a string representation of the bit array.
func (ba *BitArray) String() string {
	baLen := len(ba.data)
	n := ba.size % bitsN
	if n == 0 {
		n = bitsN
	}
	s := fmt.Sprintf("%0*b", n, ba.data[baLen-1])
	if baLen > 1 {
		for i := baLen - 2; i >= 0; i-- {
			s += fmt.Sprintf("%0*b", bitsN, ba.data[i])
		}
	}
	return s
}

func (ba *BitArray) checkIdx(idx int) {
	if idx < 0 || idx >= ba.size {
		panic("index out of range")
	}
}

func (ba *BitArray) checkSize(other *BitArray) {
	if ba.size != other.size {
		panic("bit array sizes must be equal")
	}
}

// MarshalBinary implements the [encoding/BinaryMarshaler] interface.
func (ba *BitArray) MarshalBinary() ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	err := enc.Encode(ba.size)
	if err != nil {
		return nil, err
	}
	err = enc.Encode(ba.data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// UnmarshalBinary implements the [encoding/BinaryUnmarshaler] interface.
func (ba *BitArray) UnmarshalBinary(data []byte) error {
	b := bytes.NewReader(data)
	dec := gob.NewDecoder(b)
	err := dec.Decode(&ba.size)
	if err != nil {
		return err
	}
	err = dec.Decode(&ba.data)
	if err != nil {
		return err
	}
	return nil
}

// Slice returns a new BitArray with the bits from ba at indexes [start, end).
func Slice(ba *BitArray, start, end int) *BitArray {
	ba.checkIdx(start)
	ba.checkIdx(end - 1)
	result := New(end - start)
	idx := 0
	for i := start; i < end; i++ {
		if ba.get(i) {
			result.set(idx)
		}
		idx++
	}
	return result
}

// Concat returns a new BitArray with the bits from ba1 and ba2 concatenated.
func Concat(ba1, ba2 *BitArray) *BitArray {
	ba := New(ba1.size + ba2.size)
	n := ba2.size / bitsN
	for i := 0; i < n; i++ {
		ba.data[i] = ba2.data[i]
	}
	for i := n * bitsN; i < ba2.size; i++ {
		if ba2.get(i) {
			ba.set(i)
		}
	}
	i2 := 0
	for i := ba2.size; i < ba.size; i++ {
		if ba1.get(i2) {
			ba.set(i)
		}
		i2++
	}
	return ba
}
