package mahi

import (
	"io"
	"os"
)

type FileBlobStorage interface {
	Upload(bucket string, b *FileBlob) error
}

// FileBlob is holds the properties needed for the blob of a file.
type FileBlob struct {
	ID        string
	Data      io.ReadCloser
	MIMEValue string
	Size      int64

	// TempFileName this is used to determine if we need to delete the temp file after using the FileBlob
	TempFileName string
}

// Bytes transforms the data of the FileBlob into a byte array
func (b *FileBlob) Bytes(p []byte) (int, error) {
	n, err := b.Data.Read(p)
	if err != nil && err != io.EOF {
		return 0, err
	}

	if int64(n) != b.Size {
		return 0, io.ErrShortWrite
	}

	return n, nil
}

// Close closes the io.ReadCloser and deletes the temp file if one was used
func (b *FileBlob) Close() error {
	if err := b.Data.Close(); err != nil {
		return err
	}

	if b.TempFileName != "" {
		return os.RemoveAll(b.TempFileName)
	}

	return nil
}
