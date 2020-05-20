package mbits

// Copyright(c) Dorin Duminica. All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//   1. Redistributions of source code must retain the above copyright notice,
// 	 this list of conditions and the following disclaimer.
//
//   2. Redistributions in binary form must reproduce the above copyright notice,
// 	 this list of conditions and the following disclaimer in the documentation
// 	 and/or other materials provided with the distribution.
//
//   3. Neither the name of the copyright holder nor the names of its
// 	 contributors may be used to endorse or promote products derived from this
// 	 software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"bytes"
	"fmt"
	"reflect" // using some structs to make magic
	"unsafe"  // it's safe, trust me :)
)

const (
	KSZ_U8           = uint(unsafe.Sizeof(uint8(0)))
	KSZ_U64          = uint(unsafe.Sizeof(uint64(0)))
	KBITS_PER_BYTE   = uint(8)
	KWORD_SIZE_BYTES = uint(unsafe.Sizeof(uint(0)))
	KWORD_SIZE_BITS  = KWORD_SIZE_BYTES * KBITS_PER_BYTE
)

type BitBuffer struct {
	// internal buffer
	buff []uint
	// number of bytes "wanted" from buffer
	byte_len uint
}

// constructs a BitBuffer with given length in bytes and returns pointer to instance
func NewBitBuffer(length uint) *BitBuffer {
	return (&BitBuffer{}).SetBufferLen(length)
}

// set buffer length in bytes
// discards previous data if any
// returns pointer to self
func (m *BitBuffer) SetBufferLen(length uint) *BitBuffer {
	m.byte_len = length
	lbuff := m.byte_len/KWORD_SIZE_BYTES + 1
	if m.byte_len == 0 {
		m.byte_len = lbuff * KWORD_SIZE_BYTES
	}
	m.buff = make([]uint, lbuff)
	return m
}

// length of buffer in bits
func (m *BitBuffer) LenBits() uint {
	return m.LenBytes() * KBITS_PER_BYTE
}

// length of buffer in bytes
func (m *BitBuffer) LenBytes() uint {
	return m.byte_len
}

// returns a copy of buffer as a byte slice
func (m *BitBuffer) Bytes() []byte {
	r := make([]byte, m.LenBytes())
	copy(r, m.MutableByteSlice())
	return r
}

// returns a bool slice in which each item represents the state of a bit
func (m *BitBuffer) Bool() []bool {
	l := m.LenBits()
	r := make([]bool, l)
	for i := uint(0); i < l; i++ {
		r[i] = m.IsSet(i)
	}
	return r
}

// checks if bitIndex is out of bounds, if so, it grows the internal buffer
func (m *BitBuffer) growIfNeeded(bitIndex uint) {
	l := bitIndex/KWORD_SIZE_BITS + 1
	if uint(len(m.buff)) < l {
		r := make([]uint, l)
		copy(r, m.buff)
		m.buff = r
		m.byte_len = l * KWORD_SIZE_BYTES
	}
}

// copies the byte slice to the internal buffer
// internal buffer will be reset
// returns pointer to self
func (m *BitBuffer) LoadBuffer(buffer []byte) *BitBuffer {
	m.SetBufferLen(uint(len(buffer)))
	copy(m.MutableByteSlice(), buffer)
	return m
}

// toggle bit state at index
// returns pointer to self
func (m *BitBuffer) Toggle(bitIndex uint) *BitBuffer {
	m.growIfNeeded(bitIndex)
	m.buff[bitIndex/KWORD_SIZE_BITS] ^= 1 << (bitIndex % KWORD_SIZE_BITS)
	return m
}

// turn bit on at index
// returns pointer to self
func (m *BitBuffer) Set(bitIndex uint) *BitBuffer {
	m.growIfNeeded(bitIndex)
	m.buff[bitIndex/KWORD_SIZE_BITS] |= 1 << (bitIndex % KWORD_SIZE_BITS)
	return m
}

// set bit off
// returns pointer to self
func (m *BitBuffer) Clear(bitIndex uint) *BitBuffer {
	m.growIfNeeded(bitIndex)
	m.buff[bitIndex/KWORD_SIZE_BITS] &^= 1 << (bitIndex % KWORD_SIZE_BITS)
	return m
}

// returns pointer to self
func (m *BitBuffer) SetAll(x uint) *BitBuffer {
	// avoid calling len() in the loop
	l := len(m.buff)

	// fetch address of first uint in buff
	p := uintptr(unsafe.Pointer(&m.buff[0]))

	// loop over all uint items in buff
	for i := 0; i < l; i++ {
		// cast current pointer to a pointer of our element type
		v := (*uint)(unsafe.Pointer(p))

		// set value
		*v = x

		// move to next element
		p++
	}
	return m
}

// turns off all the bits
// returns pointer to self
func (m *BitBuffer) ClearAll() *BitBuffer {
	return m.SetAll(0)
}

// returns pointer to self
func (m *BitBuffer) SetOnAll() *BitBuffer {
	return m.SetAll(wordBitsOn)
}

// returns true if bit at index is set
func (m *BitBuffer) IsSet(bitIndex uint) bool {
	m.growIfNeeded(bitIndex)
	return m.buff[bitIndex/KWORD_SIZE_BITS]&(1<<(bitIndex%KWORD_SIZE_BITS)) != 0
}

// returns the count of on and off bits
func (m *BitBuffer) CountBits() (on uint, off uint) {
	// got any ideas as to how we can speed this up?

	// counting ON bits -- fast -- is more challening than it seems
	// use pre-computed lookup table
	lookup := LookupByteBitsOn

	// how many iterations can we perform given the buffer size in bytes
	// and size of a pointer in bytes
	l := m.byte_len / KWORD_SIZE_BYTES

	// pointer to buffer
	p := uintptr(unsafe.Pointer(&m.buff[0]))

	// check pointer size only once, and then unroll the madness
	switch KWORD_SIZE_BYTES {
	case 2:
		{
			// 16-bit
			for i := uint(0); i < l; i++ {
				v := (*[2]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				p++
			}
		}
	case 4:
		{
			// 32-bit
			for i := uint(0); i < l; i++ {
				v := (*[4]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
				p++
			}
		}
	case 8:
		{
			// 64-bit
			for i := uint(0); i < l; i++ {
				v := (*[8]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
				on += uint(lookup[(*v)[4]])
				on += uint(lookup[(*v)[5]])
				on += uint(lookup[(*v)[6]])
				on += uint(lookup[(*v)[7]])
				p++
			}
		}
	}

	// check if there are still unprocessed bytes left
	rem := m.byte_len % KWORD_SIZE_BYTES
	if rem > 0 {
		// alright, count these ones
		// increment pointer
		p++

		switch rem {
		case 1:
			{
				v := (*byte)(unsafe.Pointer(p))
				on += uint(lookup[*v])
			}
		case 2:
			{
				v := (*[2]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
			}
		case 3:
			{
				v := (*[3]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
			}
		case 4:
			{
				v := (*[4]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
			}
		case 5:
			{
				v := (*[5]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
				on += uint(lookup[(*v)[4]])
			}
		case 6:
			{
				v := (*[6]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
				on += uint(lookup[(*v)[4]])
				on += uint(lookup[(*v)[5]])
			}
		case 7:
			{
				v := (*[7]byte)(unsafe.Pointer(p))
				on += uint(lookup[(*v)[0]])
				on += uint(lookup[(*v)[1]])
				on += uint(lookup[(*v)[2]])
				on += uint(lookup[(*v)[3]])
				on += uint(lookup[(*v)[4]])
				on += uint(lookup[(*v)[5]])
				on += uint(lookup[(*v)[6]])
			}
		default:
			panic("unexpected word size")
		}
	}

	// these can be computed last
	off = (m.byte_len * KBITS_PER_BYTE) - on

	return
}

// returns the number of on bits
// calls CountBits() internally
func (m *BitBuffer) CountBitsOn() uint {
	on, _ := m.CountBits()

	return on
}

// returns the number of off bits
// calls CountBits() internally
func (m *BitBuffer) CountBitsOff() uint {
	_, off := m.CountBits()

	return off
}

// compares (this) buffer with (other) buffer
func (m *BitBuffer) CmpWith(other *BitBuffer) int {
	return bytes.Compare(m.MutableByteSlice(), other.MutableByteSlice())
}

// copy internal buffer(and state) to (other)
func (m *BitBuffer) CopyTo(other *BitBuffer) {
	other.SetBufferLen(m.byte_len)
	copy(other.buff, m.buff)
}

// copy internal buffer(and state) from (other)
// returns pointer to self
func (m *BitBuffer) CopyFrom(other *BitBuffer) *BitBuffer {
	other.CopyTo(m)
	return m
}

// returns a clone (this)
func (m *BitBuffer) Clone() *BitBuffer {
	r := &BitBuffer{}
	m.CopyTo(r)
	return r
}

// returns a mutable byte slice of internal buffer
// use with care!
func (m *BitBuffer) MutableByteSlice() []byte {
	len_bytes := m.LenBytes()

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&m.buff[0])),
		Len:  int(len_bytes),
	}))
}

// returns a string of 1's and 0's representing the state of the bits
func (m *BitBuffer) String() string {
	// got any ideas as to how we can speed this up?

	// use pre-computed lookup table
	lookup := LookupByteBinStr

	// init locals
	len_bytes := m.LenBytes()
	len_bits := m.LenBits()

	// byte slice
	r := make([]byte, len_bits)

	// pointer to buffer
	p := uintptr(unsafe.Pointer(&m.buff[0]))

	// pointer to return buffer
	rp := uintptr(unsafe.Pointer(&r[0]))

	// itrate over all bytes
	for i := uint(0); i < len_bytes; i++ {
		// pointer to current byte
		v := (*byte)(unsafe.Pointer(p))

		// pointer to current return buffer as a 64-bit unsigned int
		dr := (*uint64)(unsafe.Pointer(rp))

		// deref and set value
		*dr = lookup[*v]

		// increment pointer to return buffer
		rp += uintptr(KSZ_U64)

		// increment pointer to byte buffer
		p += uintptr(KSZ_U8)
	}

	// "convert" the []byte slice into a string with magic
	// avoids copying r's data into a string
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&r[0])),
		Len:  int(len_bits),
	}))
}

var wordBitsOn uint = 0

func init() {
	x := uint(0xff)
	switch KWORD_SIZE_BITS {
	case 16:
		x = uint(0xffff)
	case 32:
		x = uint(0xffffffff)
	case 64:
		x = uint(0xffffffffffffffff)
	default:
		panic(fmt.Sprintf("unexpected word size: %v", KWORD_SIZE_BITS))
	}
	wordBitsOn = x
}
