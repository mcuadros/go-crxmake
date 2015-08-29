package crxmake

import (
	"archive/zip"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
)

var (
	header         = []byte{'C', 'r', '2', '4'}
	version uint32 = 2
)

const keyFilename = "key.pem"

// Builder creates CRX files using as source a Chrome extension folder.
type Builder struct {
	Content    *bytes.Buffer
	PrivateKey *rsa.PrivateKey
}

// NewBuilder returns a new CRX Builder
func NewBuilder() *Builder {
	return &Builder{
		Content: new(bytes.Buffer),
	}
}

func (b *Builder) LoadKeyFile(pemFile string) error {
	buf, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(buf)
	if block == nil {
		return errors.New("key not found")
	}

	r, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	b.PrivateKey = r
	return nil
}

// BuildZip loades the given folder and include all the files on the zip.
func (b *Builder) BuildZip(folder string) error {
	w := zip.NewWriter(b.Content)
	defer w.Close()

	if err := b.generateKeyIfNeeded(); err != nil {
		return err
	}

	keyFile, err := w.Create(keyFilename)
	if err != nil {
		return err
	}

	size, err := b.saveKeyFile(keyFile)
	if err != nil {
		return err
	}

	fmt.Printf("Included %q %s\n", keyFilename, humanize.Bytes(uint64(size)))

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

func (b *Builder) generateKeyIfNeeded() error {
	if b.PrivateKey != nil {
		return nil
	}

	var err error
	b.PrivateKey, err = rsa.GenerateKey(rand.Reader, 1024)

	return err
}

func (b *Builder) saveKeyFile(file io.Writer) (int, error) {
	bytes := x509.MarshalPKCS1PrivateKey(b.PrivateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytes,
	}

	return file.Write(pem.EncodeToMemory(block))
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
	b.PrivateKey.Precompute()
	if err = b.PrivateKey.Validate(); err != nil {
		return
	}

	publicKey, err = x509.MarshalPKIXPublicKey(&b.PrivateKey.PublicKey)
	if err != nil {
		return
	}

	h := crypto.SHA1.New()
	h.Write(b.Content.Bytes())

	signature, err = rsa.SignPKCS1v15(rand.Reader, b.PrivateKey, crypto.SHA1, h.Sum(nil))
	if err != nil {
		return
	}

	return
}
