package extraction

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/h2non/filetype"

	"github.com/pkg/errors"
)

type ReadOnlyMemoryFile struct {
	source          string
	data            []byte
	currentPosition int64
}

func ReadFile(source string) (*ReadOnlyMemoryFile, error) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read input file")
	}
	return &ReadOnlyMemoryFile{
		source: source,
		data:   data,
	}, nil
}

func (m *ReadOnlyMemoryFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.readAt(p, off)
}

func (m *ReadOnlyMemoryFile) readAt(p []byte, offset int64) (int, error) {
	bytesCopied := copy(p, m.data[m.currentPosition+offset:])
	if bytesCopied < len(p) {
		return bytesCopied, io.EOF
	}
	return bytesCopied, nil
}

func (m *ReadOnlyMemoryFile) Read(p []byte) (n int, err error) {
	read, err := m.readAt(p, 0)
	m.currentPosition = m.currentPosition + int64(read)
	return read, err
}

func (m *ReadOnlyMemoryFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekCurrent:
		m.currentPosition = m.currentPosition + offset
	case io.SeekStart:
		m.currentPosition = offset
	case io.SeekEnd:
		m.currentPosition = int64(len(m.data)) + offset
	}
	return m.currentPosition, nil
}

// Reset resets the seek position to the start of the file.
// If the seek fails, there will be a panic.
func (m *ReadOnlyMemoryFile) MustReset() {
	_, err := m.Seek(0, io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("seek failed on a memory file. This should never happen! %v", err))
	}
}

func (m *ReadOnlyMemoryFile) IsImage() (bool, error) {
	header := make([]byte, 261)
	_, err := m.ReadAt(header, 0)
	if err != nil {
		return false, errors.Wrap(err, "failed to read file header")
	}
	return filetype.IsImage(header), nil
}
