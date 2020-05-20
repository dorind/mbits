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

type ByteTableLookupU64 [256]uint64
type ByteTableLookupByte [256]byte

// pre-computed lookup table storing the number of ON bits for each byte
// 256 bytes(1 byte X 256 entries) of raw data
var LookupByteBitsOn = ByteTableLookupByte{
	0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	1, 2, 2, 3, 2, 3, 3, 4, 2, 3, 3, 4, 3, 4, 4, 5,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	2, 3, 3, 4, 3, 4, 4, 5, 3, 4, 4, 5, 4, 5, 5, 6,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	3, 4, 4, 5, 4, 5, 5, 6, 4, 5, 5, 6, 5, 6, 6, 7,
	4, 5, 5, 6, 5, 6, 6, 7, 5, 6, 6, 7, 6, 7, 7, 8,
}

// pre-computed lookup table used for printing bits
// 2 KiB(8 bytes X 256 entries) of raw data, quite a bit -- I know!
var LookupByteBinStr = ByteTableLookupU64{
	0x3030303030303030, 0x3030303030303031, 0x3030303030303130, 0x3030303030303131,
	0x3030303030313030, 0x3030303030313031, 0x3030303030313130, 0x3030303030313131,
	0x3030303031303030, 0x3030303031303031, 0x3030303031303130, 0x3030303031303131,
	0x3030303031313030, 0x3030303031313031, 0x3030303031313130, 0x3030303031313131,
	0x3030303130303030, 0x3030303130303031, 0x3030303130303130, 0x3030303130303131,
	0x3030303130313030, 0x3030303130313031, 0x3030303130313130, 0x3030303130313131,
	0x3030303131303030, 0x3030303131303031, 0x3030303131303130, 0x3030303131303131,
	0x3030303131313030, 0x3030303131313031, 0x3030303131313130, 0x3030303131313131,
	0x3030313030303030, 0x3030313030303031, 0x3030313030303130, 0x3030313030303131,
	0x3030313030313030, 0x3030313030313031, 0x3030313030313130, 0x3030313030313131,
	0x3030313031303030, 0x3030313031303031, 0x3030313031303130, 0x3030313031303131,
	0x3030313031313030, 0x3030313031313031, 0x3030313031313130, 0x3030313031313131,
	0x3030313130303030, 0x3030313130303031, 0x3030313130303130, 0x3030313130303131,
	0x3030313130313030, 0x3030313130313031, 0x3030313130313130, 0x3030313130313131,
	0x3030313131303030, 0x3030313131303031, 0x3030313131303130, 0x3030313131303131,
	0x3030313131313030, 0x3030313131313031, 0x3030313131313130, 0x3030313131313131,
	0x3031303030303030, 0x3031303030303031, 0x3031303030303130, 0x3031303030303131,
	0x3031303030313030, 0x3031303030313031, 0x3031303030313130, 0x3031303030313131,
	0x3031303031303030, 0x3031303031303031, 0x3031303031303130, 0x3031303031303131,
	0x3031303031313030, 0x3031303031313031, 0x3031303031313130, 0x3031303031313131,
	0x3031303130303030, 0x3031303130303031, 0x3031303130303130, 0x3031303130303131,
	0x3031303130313030, 0x3031303130313031, 0x3031303130313130, 0x3031303130313131,
	0x3031303131303030, 0x3031303131303031, 0x3031303131303130, 0x3031303131303131,
	0x3031303131313030, 0x3031303131313031, 0x3031303131313130, 0x3031303131313131,
	0x3031313030303030, 0x3031313030303031, 0x3031313030303130, 0x3031313030303131,
	0x3031313030313030, 0x3031313030313031, 0x3031313030313130, 0x3031313030313131,
	0x3031313031303030, 0x3031313031303031, 0x3031313031303130, 0x3031313031303131,
	0x3031313031313030, 0x3031313031313031, 0x3031313031313130, 0x3031313031313131,
	0x3031313130303030, 0x3031313130303031, 0x3031313130303130, 0x3031313130303131,
	0x3031313130313030, 0x3031313130313031, 0x3031313130313130, 0x3031313130313131,
	0x3031313131303030, 0x3031313131303031, 0x3031313131303130, 0x3031313131303131,
	0x3031313131313030, 0x3031313131313031, 0x3031313131313130, 0x3031313131313131,
	0x3130303030303030, 0x3130303030303031, 0x3130303030303130, 0x3130303030303131,
	0x3130303030313030, 0x3130303030313031, 0x3130303030313130, 0x3130303030313131,
	0x3130303031303030, 0x3130303031303031, 0x3130303031303130, 0x3130303031303131,
	0x3130303031313030, 0x3130303031313031, 0x3130303031313130, 0x3130303031313131,
	0x3130303130303030, 0x3130303130303031, 0x3130303130303130, 0x3130303130303131,
	0x3130303130313030, 0x3130303130313031, 0x3130303130313130, 0x3130303130313131,
	0x3130303131303030, 0x3130303131303031, 0x3130303131303130, 0x3130303131303131,
	0x3130303131313030, 0x3130303131313031, 0x3130303131313130, 0x3130303131313131,
	0x3130313030303030, 0x3130313030303031, 0x3130313030303130, 0x3130313030303131,
	0x3130313030313030, 0x3130313030313031, 0x3130313030313130, 0x3130313030313131,
	0x3130313031303030, 0x3130313031303031, 0x3130313031303130, 0x3130313031303131,
	0x3130313031313030, 0x3130313031313031, 0x3130313031313130, 0x3130313031313131,
	0x3130313130303030, 0x3130313130303031, 0x3130313130303130, 0x3130313130303131,
	0x3130313130313030, 0x3130313130313031, 0x3130313130313130, 0x3130313130313131,
	0x3130313131303030, 0x3130313131303031, 0x3130313131303130, 0x3130313131303131,
	0x3130313131313030, 0x3130313131313031, 0x3130313131313130, 0x3130313131313131,
	0x3131303030303030, 0x3131303030303031, 0x3131303030303130, 0x3131303030303131,
	0x3131303030313030, 0x3131303030313031, 0x3131303030313130, 0x3131303030313131,
	0x3131303031303030, 0x3131303031303031, 0x3131303031303130, 0x3131303031303131,
	0x3131303031313030, 0x3131303031313031, 0x3131303031313130, 0x3131303031313131,
	0x3131303130303030, 0x3131303130303031, 0x3131303130303130, 0x3131303130303131,
	0x3131303130313030, 0x3131303130313031, 0x3131303130313130, 0x3131303130313131,
	0x3131303131303030, 0x3131303131303031, 0x3131303131303130, 0x3131303131303131,
	0x3131303131313030, 0x3131303131313031, 0x3131303131313130, 0x3131303131313131,
	0x3131313030303030, 0x3131313030303031, 0x3131313030303130, 0x3131313030303131,
	0x3131313030313030, 0x3131313030313031, 0x3131313030313130, 0x3131313030313131,
	0x3131313031303030, 0x3131313031303031, 0x3131313031303130, 0x3131313031303131,
	0x3131313031313030, 0x3131313031313031, 0x3131313031313130, 0x3131313031313131,
	0x3131313130303030, 0x3131313130303031, 0x3131313130303130, 0x3131313130303131,
	0x3131313130313030, 0x3131313130313031, 0x3131313130313130, 0x3131313130313131,
	0x3131313131303030, 0x3131313131303031, 0x3131313131303130, 0x3131313131303131,
	0x3131313131313030, 0x3131313131313031, 0x3131313131313130, 0x3131313131313131,
}
