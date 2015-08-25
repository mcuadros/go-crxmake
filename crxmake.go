package crxmake

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

var header = []byte{'C', 'r', '2', '4'}
var version uint32 = 2

// Builder creates CRX files using as source a Chrome extension folder.
type Builder struct {
	Content *bytes.Buffer
}

// NewBuilder returns a new CRX Builder
func NewBuilder() *Builder {
	return &Builder{
		Content: new(bytes.Buffer),
	}
}

// BuildZip loades the given folder and include all the files on the zip.
func (b *Builder) BuildZip(folder string) error {
	w := zip.NewWriter(b.Content)

	defer w.Close()
	return filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename, _ := filepath.Rel(folder, path)
		if info.IsDir() {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = filename
		dst, err := w.CreateHeader(header)
		if err != nil {
			return err
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}

		bytes, err := io.Copy(dst, src)
		if err != nil {
			return err
		}

		fmt.Printf("Included %q %s\n", filename, humanize.Bytes(uint64(bytes)))
		return nil
	})
}

// WriteToFile writes the generated CRX file.
func (b *Builder) WriteToFile(w io.Writer) error {
	if _, err := w.Write(header); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, version); err != nil {
		return err
	}

	key, signature, err := b.signContent()
	if err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint32(len(key))); err != nil {
		fmt.Println("key")
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint32(len(signature))); err != nil {
		return err
	}

	if _, err := w.Write(key); err != nil {
		return err
	}

	if _, err := w.Write(signature); err != nil {
		return err
	}

	if _, err := io.Copy(w, b.Content); err != nil {
		return err
	}

	return nil
}

func (b *Builder) signContent() (publicKey, signature []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return
	}

	privateKey.Precompute()
	if err = privateKey.Validate(); err != nil {
		return
	}

	publicKey, err = x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return
	}

	h := crypto.SHA1.New()
	h.Write(b.Content.Bytes())

	signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, h.Sum(nil))
	if err != nil {
		fmt.Println("foo")
	}

	return
}
