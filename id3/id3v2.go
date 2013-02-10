package id3

import (
	"bytes"
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type ID3v2Header struct {
	Version        int
	MinorVersion   int
	Unsynchronized bool
	Extended       bool
	Experimental   bool
	Footer         bool
	Size           uint32
}

type ID3v2Frame struct {
	Id   string
	Data []byte
}

type ID3v2FrameParser struct {
	HeaderLen  int
	IdLen      int
	SizeLen    int
	SizeParser func([]byte) uint32
}

func ParseSynchSafe(bin []byte) uint32 {
	var out uint32 = 0
	var mask uint32 = 0x7f000000
	var in uint32
	l := len(bin)
	binbuf := bytes.NewBuffer(bin[l-4 : l])
	binary.Read(binbuf, binary.BigEndian, &in)
	for mask != 0 {
		out >>= 1
		out |= in & mask
		mask >>= 8
	}
	if l > 4 {
		out |= uint32(bin[l-5]) << 28
	}
	return out
}

func ParseID3v23FrameSize(buf []byte) uint32 {
	var out uint32
	bufr := bytes.NewBuffer(buf)
	binary.Read(bufr, binary.BigEndian, &out)
	return out
}

func parseUTF16(buf []byte) []uint16 {
	var bo binary.ByteOrder
	var bom bool
	var br *bytes.Buffer
	var buf16 []uint16

	if buf[0] == 0xFE && buf[1] == 0xFF {
		bo = binary.BigEndian
		bom = true
	} else if buf[0] == 0xFF && buf[1] == 0xFE {
		bo = binary.LittleEndian
		bom = true
	} else {
		bo = binary.LittleEndian
		bom = false
	}

	if (bom) {
		br = bytes.NewBuffer(buf[2:])
		buf16 = make([]uint16, len(buf)/2-1)
	} else {
		br = bytes.NewBuffer(buf)
		buf16 = make([]uint16, len(buf)/2)
	}

	binary.Read(br, bo, &buf16)
	return buf16
}

func parseUTF16String(buf []byte) string {
	return string(utf16.Decode(parseUTF16(buf)))
}

func ParseID3v2String(buf []byte) string {
	switch buf[0] {
	case 0:
		return string(buf[1:])
	case 1:
		return parseUTF16String(buf[1:])
	case 2:
		return parseUTF16String(buf[1:])
	case 3:
		return string(buf[1:])
	}
	return ""
}

func ParseID3v2Header(reader io.Reader) (*ID3v2Header, error) {
	buf := make([]byte, 10)
	n, err := reader.Read(buf)
	if err != nil || n != 10 {
		return nil, ErrRead
	}
	h := new(ID3v2Header)
	h.Version = int(buf[3])
	h.MinorVersion = int(buf[4])
	h.Unsynchronized = buf[5]&(1<<7) != 0
	h.Extended = buf[5]&(1<<6) != 0
	h.Experimental = buf[5]&(1<<5) != 0
	h.Footer = buf[5]&(1<<4) != 0
	h.Size = ParseSynchSafe(buf[6:10])
	return h, nil
}

func NewID3v2FrameParser(version int) *ID3v2FrameParser {
	p := new(ID3v2FrameParser)
	switch version {
	case 2:
		p.HeaderLen = 6
		p.IdLen = 3
		p.SizeLen = 3
		p.SizeParser = ParseSynchSafe
	case 3:
		p.HeaderLen = 10
		p.IdLen = 4
		p.SizeLen = 4
		p.SizeParser = ParseID3v23FrameSize
	case 4:
		p.HeaderLen = 10
		p.IdLen = 4
		p.SizeLen = 4
		p.SizeParser = ParseSynchSafe
	}
	return p
}

func (p *ID3v2FrameParser) ReadFrame(reader io.Reader) (*ID3v2Frame, error) {
	hbuf := make([]byte, p.HeaderLen)
	n, err := reader.Read(hbuf)
	if err != nil || n != p.HeaderLen {
		return nil, ErrRead
	}
	id := string(hbuf[0:p.IdLen])
	if id[0] == '\u0000' {
		return nil, ErrEOF
	}
	size := p.SizeParser(hbuf[p.IdLen:p.IdLen+p.SizeLen])
	data := make([]byte, size)
	n, err = reader.Read(data)
	if err != nil || n != int(size) {
		return nil, ErrRead
	}
	f := new(ID3v2Frame)
	f.Id = id
	f.Data = data
	return f, nil
}
