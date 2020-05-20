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
	"testing"
)

func TestBitBufferNew(t *testing.T) {
	nbytes := uint(5)
	b := NewBitBuffer(nbytes)

	if b.LenBytes() != nbytes {
		t.Fatalf("Expected %v bytes, found %v", nbytes, b.LenBytes())
	}

	if b.LenBits() != nbytes*KBITS_PER_BYTE {
		t.Fatalf("Expected %v bits, found %v", nbytes*KBITS_PER_BYTE, b.LenBits())
	}
}

func TestBitBufferBool(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(0).
		Set(5).
		Set(9)

	bool_slice := b.Bool()
	if !bool_slice[0] || !bool_slice[5] || !bool_slice[9] {
		t.Fatal("Conversion error")
	}
}

func TestBitBufferLoadBuffer(t *testing.T) {
	buff := []byte{0xcc, 0x00, 0xff, 0xff, 0xee, 0xee}
	b := NewBitBuffer(0)
	b.LoadBuffer(buff)
	byte_slice := b.MutableByteSlice()

	if bytes.Compare(buff, byte_slice) != 0 {
		t.Log("buff", buff)
		t.Log("bbuf", byte_slice)
		t.Fatal("Byte slice doesn't match")
	}
}

func TestBitBufferToggle(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(1).
		Set(5).
		Set(9).
		Toggle(5)

	if b.IsSet(5) {
		t.Fatal("Toggle error")
	}
}

func TestBitBufferSet(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(14)

	if !b.IsSet(14) {
		t.Fatal("Bit set error")
	}
}

func TestBitBufferClear(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(1).
		Set(64).
		Set(128).
		Clear(64)

	if b.IsSet(64) {
		t.Fatal("Clear bit error")
	}
}

func TestBitBufferClearAll(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(7).
		Set(13).
		Set(256).
		ClearAll()
	if b.CountBitsOn() != 0 {
		t.Fatalf("ClearAll error, %v bits on", b.CountBitsOn())
	}
}

func TestBitBufferSetOnAll(t *testing.T) {
	b := NewBitBuffer(10)
	b.SetOnAll()

	if b.CountBitsOff() != 0 {
		t.Fatalf("SetOnAll error, %v bits off", b.CountBitsOff())
	}
}

func TestBitBufferCountBits(t *testing.T) {
	b := NewBitBuffer(0)
	for i := uint(0); i < 256; i += 2 {
		b.Set(i)
	}

	on, off := b.CountBits()
	if on != 128 || off != 128 {
		t.Fatalf("CountsBits fail, expected 128 on and off, actual on: %v, off: %v", on, off)
	}
}

func TestBitBufferCmpWith(t *testing.T) {
	left := NewBitBuffer(5)
	left.LoadBuffer([]byte{0xdd, 0x00, 0x44, 0x11, 0xdd})
	right := NewBitBuffer(1)

	rcmp := left.CmpWith(right)
	if rcmp != 1 {
		t.Fatalf("CmpWith fail, expected 1, found %v", rcmp)
	}

	right.CopyFrom(left)

	rcmp = left.CmpWith(right)
	if rcmp != 0 {
		t.Fatalf("CmpWith fail, expected 0, found %v", rcmp)
	}

	left.SetBufferLen(1)
	rcmp = left.CmpWith(right)
	if rcmp != -1 {
		t.Fatalf("CmpWith fail, expected -1, found %v", rcmp)
	}
}

func TestBitBufferString(t *testing.T) {
	b := NewBitBuffer(0)
	b.Set(0).
		Set(63).
		Set(127)
	s := b.String()
	expected := "10000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
	if s != expected {
		t.Fatalf("String fail,\nexpected: %v\nfound: %v", expected, s)
	}
}

func BenchmarkBitBufferNew256(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		_ = NewBitBuffer(256)
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew256String(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(256)
		_ = b.String()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew256SetOnAll(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(256)
		b.SetOnAll()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew256CountBits(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(256)
		_, _ = b.CountBits()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew1024(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		_ = NewBitBuffer(1024)
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew1024String(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(1024)
		_ = b.String()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew1024SetOnAll(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(1024)
		b.SetOnAll()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew1024CountBits(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(1024)
		_, _ = b.CountBits()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew4096(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		_ = NewBitBuffer(4096)
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew4096String(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(4096)
		_ = b.String()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew4096SetOnAll(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(4096)
		b.SetOnAll()
	}
	t.StopTimer()
}

func BenchmarkBitBufferNew4096CountBits(t *testing.B) {
	t.StartTimer()
	for i := 1; i < t.N; i++ {
		b := NewBitBuffer(4096)
		_, _ = b.CountBits()
	}
	t.StopTimer()
}
