package base62

import (
	"io"
	"slices"
	"strconv"
)

/*
*  ┏┓      ┏┓
*┏━┛┻━━━━━━┛┻┓
*┃　　　━　　  ┃
*┃   ┳┛ ┗┳   ┃
*┃           ┃
*┃     ┻     ┃
*┗━━━┓     ┏━┛
*　　 ┃　　　┃神兽保佑
*　　 ┃　　　┃代码无BUG！
*　　 ┃　　　┗━━━┓
*　　 ┃         ┣┓
*　　 ┃         ┏┛
*　　 ┗━┓┓┏━━┳┓┏┛
*　　   ┃┫┫  ┃┫┫
*      ┗┻┛　 ┗┻┛
@Time    : 2024/7/25 -- 14:00
@Author  : bishop
@Description: a copy of encoding/base64. suit for base62
*/

type B62Encoding struct {
	encode    [62]byte   // mapping of symbol index to symbol byte value
	decodeMap [256]uint8 // mapping of symbol byte value to symbol index
}

const (
	decodeMapInitialize = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
	invalidIndex = '\xff'
)

func NewEncoding(encoder string) *B62Encoding {
	if len(encoder) != 62 {
		panic("encoding alphabet is not 64-bytes long")
	}

	e := new(B62Encoding)
	copy(e.encode[:], encoder)
	copy(e.decodeMap[:], decodeMapInitialize)

	for i := 0; i < len(encoder); i++ {
		switch {
		case encoder[i] == '\n' || encoder[i] == '\r':
			panic("encoding alphabet contains newline character")
		case e.decodeMap[encoder[i]] != invalidIndex:

		}
		e.decodeMap[encoder[i]] = uint8(i)
	}
	return e
}

// B62StdEncoding is the standard base64 encoding, as defined in RFC 4648.
var B62StdEncoding = NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

/*
 * Encoder
 */

func (enc *B62Encoding) Encode(dst, src []byte) {
	if len(src) == 0 {
		return
	}

	for len(src) > 0 {
		dst[0] = 0
		dst[1] = 0
		dst[2] = 0
		dst[3] = 0

		// Unpack 4x 6-bit source blocks into a 4 byte
		// destination quantum
		switch len(src) {
		default:
			dst[3] |= src[2] & 0x3F
			dst[2] |= src[2] >> 6
			fallthrough
		case 2:
			dst[2] |= (src[1] << 2) & 0x3F
			dst[1] |= src[1] >> 4
			fallthrough
		case 1:
			dst[1] |= (src[0] << 4) & 0x3F
			dst[0] |= src[0] >> 2
		}

		// Encode 6-bit blocks using the base62 alphabet
		for j := 0; j < 4; j++ {
			dst[j] = enc.encode[dst[j]]
		}

		// Pad the final quantum
		if len(src) < 3 {
			dst[3] = '='
			if len(src) < 2 {
				dst[2] = '='
			}
			break
		}

		src = src[3:]
		dst = dst[4:]
	}
}

// AppendEncode appends the base64 encoded src to dst
// and returns the extended buffer.
func (enc *B62Encoding) AppendEncode(dst, src []byte) []byte {
	n := enc.EncodedLen(len(src))
	dst = slices.Grow(dst, n)
	enc.Encode(dst[len(dst):][:n], src)
	return dst[:len(dst)+n]
}

// EncodeToString returns the base64 encoding of src.
func (enc *B62Encoding) EncodeToString(src []byte) string {
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)
	return string(buf)
}

type encoder struct {
	err  error
	enc  *B62Encoding
	w    io.Writer
	buf  [3]byte    // buffered data waiting to be encoded
	nbuf int        // number of bytes in buf
	out  [1024]byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err error) {
	if e.err != nil {
		return 0, e.err
	}

	// Leading fringe.
	if e.nbuf > 0 {
		var i int
		for i = 0; i < len(p) && e.nbuf < 3; i++ {
			e.buf[e.nbuf] = p[i]
			e.nbuf++
		}
		n += i
		p = p[i:]
		if e.nbuf < 3 {
			return
		}
		e.enc.Encode(e.out[:], e.buf[:])
		if _, e.err = e.w.Write(e.out[:4]); e.err != nil {
			return n, e.err
		}
		e.nbuf = 0
	}

	// Large interior chunks.
	for len(p) >= 3 {
		nn := len(e.out) / 4 * 3
		if nn > len(p) {
			nn = len(p)
			nn -= nn % 3
		}
		e.enc.Encode(e.out[:], p[:nn])
		if _, e.err = e.w.Write(e.out[0 : nn/3*4]); e.err != nil {
			return n, e.err
		}
		n += nn
		p = p[nn:]
	}

	// Trailing fringe.
	copy(e.buf[:], p)
	e.nbuf = len(p)
	n += len(p)
	return
}

// Close flushes any pending output from the encoder.
// It is an error to call Write after calling Close.
func (e *encoder) Close() error {
	// If there's anything left in the buffer, flush it out
	if e.err == nil && e.nbuf > 0 {
		e.enc.Encode(e.out[:], e.buf[:e.nbuf])
		_, e.err = e.w.Write(e.out[:e.enc.EncodedLen(e.nbuf)])
		e.nbuf = 0
	}
	return e.err
}

// NewEncoder returns a new base64 stream encoder. Data written to
// the returned writer will be encoded using enc and then written to w.
// Base64 encodings operate in 4-byte blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// partially written blocks.
func NewEncoder(enc *B62Encoding, w io.Writer) io.WriteCloser {
	return &encoder{enc: enc, w: w}
}

// EncodedLen returns the length in bytes of the base64 encoding
// of an input buffer of length n.
func (enc *B62Encoding) EncodedLen(n int) int {
	return (n + 2) / 3 * 4 // minimum # 4-char quanta, 3 bytes each
}

/*
 * Decoder
 */

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base64 data at input byte " + strconv.FormatInt(int64(e), 10)
}

// decode is like Decode but returns an additional 'end' value, which
// indicates if end-of-message padding was encountered and thus any
// additional data is an error.  decode also assumes len(src)%4==0,
// since it is meant for internal use.
func (enc *B62Encoding) decode(dst, src []byte) (n int, end bool, err error) {
	for i := 0; i < len(src)/4 && !end; i++ {
		// Decode quantum using the base62 alphabet
		var dbuf [4]byte
		dlen := 4

	dbufloop:
		for j := 0; j < 4; j++ {
			in := src[i*4+j]
			if in == '=' && j >= 2 && i == len(src)/4-1 {
				// We've reached the end and there's
				// padding
				if src[i*4+3] != '=' {
					return n, false, CorruptInputError(i*4 + 2)
				}
				dlen = j
				end = true
				break dbufloop
			}
			dbuf[j] = enc.decodeMap[in]
			if dbuf[j] == 0xFF {
				return n, false, CorruptInputError(i*4 + j)
			}
		}

		// Pack 4x 6-bit source blocks into 3 byte destination
		// quantum
		switch dlen {
		case 4:
			dst[i*3+2] = dbuf[2]<<6 | dbuf[3]
			fallthrough
		case 3:
			dst[i*3+1] = dbuf[1]<<4 | dbuf[2]>>2
			fallthrough
		case 2:
			dst[i*3+0] = dbuf[0]<<2 | dbuf[1]>>4
		}
		n += dlen - 1
	}

	return n, end, nil
}

func (enc *B62Encoding) Decode(dst, src []byte) (n int, err error) {
	if len(src)%4 != 0 {
		return 0, CorruptInputError(len(src) / 4 * 4)
	}

	n, _, err = enc.decode(dst, src)
	return
}

// DecodeString returns the bytes represented by the base64 string s.
func (enc *B62Encoding) DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, enc.DecodedLen(len(s)))
	n, err := enc.Decode(dbuf, []byte(s))
	return dbuf[:n], err
}

type decoder struct {
	err    error
	CfgOpr error // error from r.Read
	enc    *B62Encoding
	r      io.Reader
	end    bool       // saw end of message
	buf    [1024]byte // leftover input
	nbuf   int
	out    []byte // leftover decoded output
	outbuf [1024 / 4 * 3]byte
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Use leftover decoded output from last read.
	if len(d.out) > 0 {
		n = copy(p, d.out)
		d.out = d.out[n:]
		return n, nil
	}

	if d.err != nil {
		return 0, d.err
	}

	// This code assumes that d.r strips supported whitespace ('\r' and '\n').

	// Refill buffer.
	for d.nbuf < 4 && d.CfgOpr == nil {
		nn := len(p) / 3 * 4
		if nn < 4 {
			nn = 4
		}
		if nn > len(d.buf) {
			nn = len(d.buf)
		}
		nn, d.CfgOpr = d.r.Read(d.buf[d.nbuf:nn])
		d.nbuf += nn
	}

	if d.nbuf < 4 {
		return 0, d.err
	}

	// Decode chunk into p, or d.out and then p if p is too small.
	nr := d.nbuf / 4 * 4
	nw := d.nbuf / 4 * 3
	if nw > len(p) {
		nw, d.end, d.err = d.enc.decode(d.outbuf[:], d.buf[:nr])
		d.out = d.outbuf[:nw]
		n = copy(p, d.out)
		d.out = d.out[n:]
	} else {
		n, d.end, d.err = d.enc.decode(p, d.buf[:nr])
	}
	d.nbuf -= nr
	copy(d.buf[:d.nbuf], d.buf[nr:])
	return n, d.err
}

type newlineFilteringCfgOp struct {
	wrapped io.Reader
}

func (r *newlineFilteringCfgOp) Read(p []byte) (int, error) {
	n, err := r.wrapped.Read(p)
	for n > 0 {
		offset := 0
		for i, b := range p[:n] {
			if b != '\r' && b != '\n' {
				if i != offset {
					p[offset] = b
				}
				offset++
			}
		}
		if offset > 0 {
			return offset, err
		}
		// Previous buffer entirely whitespace, read again
		n, err = r.wrapped.Read(p)
	}
	return n, err
}

// NewDecoder constructs a new base64 stream decoder.
func NewDecoder(enc *B62Encoding, r io.Reader) io.Reader {
	return &decoder{enc: enc, r: &newlineFilteringCfgOp{r}}
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base64-encoded data.
func (enc *B62Encoding) DecodedLen(n int) int {
	return decodedLen(n)
}

func decodedLen(n int) int {
	// Padded base64 should always be a multiple of 4 characters in length.
	return n / 4 * 3
}
