package car

import (
	"fmt"
	cid "github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-varint"
	"io"
)

func ReadVersion(r io.Reader) (uint64, error) {
	cbor.RegisterCborType(CarHeader{})
	header, err := ReadHeader(r)
	if err != nil {
		return 0, err
	}
	return header.Version, nil
}

type CarHeader struct {
	Roots   []cid.Cid
	Version uint64
}

func ReadHeader(r io.Reader) (*CarHeader, error) {
	hb, err := LdRead(r, false)
	if err != nil {
		return nil, err
	}

	var ch CarHeader
	if err := cbor.DecodeInto(hb, &ch); err != nil {
		return nil, fmt.Errorf("invalid header: %v", err)
	}

	return &ch, nil
}

func LdRead(r io.Reader, zeroLenAsEOF bool) ([]byte, error) {
	l, err := varint.ReadUvarint(ToByteReader(r))
	if err != nil {
		// If the length of bytes read is non-zero when the error is EOF then signal an unclean EOF.
		if l > 0 && err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, err
	} else if l == 0 && zeroLenAsEOF {
		return nil, io.EOF
	}

	buf := make([]byte, l)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func ToByteReader(r io.Reader) io.ByteReader {
	if br, ok := r.(io.ByteReader); ok {
		return br
	}
	return &readerPlusByte{Reader: r}
}

type readerPlusByte struct {
	io.Reader

	byteBuf [1]byte // escapes via io.Reader.Read; preallocate
}

func (rb *readerPlusByte) ReadByte() (byte, error) {
	_, err := io.ReadFull(rb, rb.byteBuf[:])
	return rb.byteBuf[0], err
}
