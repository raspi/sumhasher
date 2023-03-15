package sumhasher

import (
	"encoding/binary"
	"hash"
	"math"
)

// Check implementation
var _ hash.Hash = &Hasher{}

// Hasher is simple fast non-cryptographic hasher.
// It simply takes a byte, increments it by one and sums all bytes found.
//
// final hash = total sum and count of bytes read is decreased from max(uint64)
// Endianness is BigEndian
type Hasher struct {
	sum       uint64 // All bytes summed together
	readBytes uint64 // Bytes hashed total
}

// New returns new Hasher
func New() (h *Hasher) {
	h = &Hasher{}
	h.Reset()

	return h
}

// Write updates hash
func (h *Hasher) Write(buffer []byte) (n int, err error) {
	if buffer == nil {
		return n, err
	}

	add := uint64(0)

	for _, b := range buffer {
		add += uint64(b) + 1
		n++
	}

	h.sum += add

	// Increase read bytes
	h.readBytes += uint64(n)

	return n, nil
}

// Sum returns final hash result
func (h *Hasher) Sum(buffer []byte) (result []byte) {
	result = make([]byte, h.Size())

	if buffer != nil {
		_, err := h.Write(buffer)
		if err != nil {
			panic(err)
		}
	}

	// Calculate hash
	res := uint64(math.MaxUint64)
	res -= h.sum
	res -= h.readBytes

	binary.BigEndian.PutUint64(result, res)
	return result
}

// Reset resets hasher's state
func (h *Hasher) Reset() {
	h.sum = 0
	h.readBytes = 0
}

// Size returns bytes required for uint64
func (h *Hasher) Size() int {
	return 8
}

// BlockSize returns recommended size for hashing buffer
func (h *Hasher) BlockSize() int {
	return 65535
}
