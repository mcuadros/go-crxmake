package main

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
)

var header = []byte{'C', 'r', '2', '4'}
var version uint32 = 2

func main() {
	// Create a buffer to write our archive to.
	b := NewBuilder()
	err := b.buildZip("examples/app/")
	fmt.Println(err)

	file, err := os.Create("foo.crx")
	defer file.Close()
	fmt.Println(err)

	err = b.write(file)
	fmt.Println(err)
}

type Builder struct {
	Content *bytes.Buffer
}

func NewBuilder() *Builder {
	return &Builder{
		Content: new(bytes.Buffer),
	}
}

func (b *Builder) buildZip(folder string) error {
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

		fmt.Println("Added", filename, bytes)
		return nil
	})
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

func (b *Builder) write(w io.Writer) error {
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
