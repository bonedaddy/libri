package enc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	cerrors "github.com/drausin/libri/libri/common/errors"
	"github.com/drausin/libri/libri/librarian/api"
)

// ErrUnexpectedCiphertextSize indicates when the ciphertext size does not match the expected value.
var ErrUnexpectedCiphertextSize = errors.New("unexpected ciphertext size")

// ErrUnexpectedCiphertextMAC indicates when the ciphertext MAC does not match the expected value.
var ErrUnexpectedCiphertextMAC = errors.New("unexpected ciphertext MAC")

// ErrUnexpectedUncompressedSize indicates when the uncompressed size does not match the expected
// value.
var ErrUnexpectedUncompressedSize = errors.New("unexpected uncompressed size")

// ErrUnexpectedUncompressedMAC indicates when the uncompressed MAC does not match the expected
// value.
var ErrUnexpectedUncompressedMAC = errors.New("unexpected uncompressed MAC")

// MAC wraps a hash function to return a message authentication code (MAC) and the total number
// of bytes it has digested.
type MAC interface {
	io.Writer

	// Sum returns the MAC after writing the given bytes.
	Sum(in []byte) []byte

	// Reset resets the MAC.
	Reset()

	// MessageSize returns the total number of digested bytes.
	MessageSize() uint64
}

type sizeHMAC struct {
	inner hash.Hash
	size  uint64
}

// NewHMAC returns a MAC internally using an HMAC-256 with a a given key.
func NewHMAC(hmacKey []byte) MAC {
	return &sizeHMAC{
		inner: hmac.New(sha256.New, hmacKey),
	}
}

func (h *sizeHMAC) Write(p []byte) (int, error) {
	n, err := h.inner.Write(p)
	h.size += uint64(n)
	return n, err
}

func (h *sizeHMAC) Sum(in []byte) []byte {
	return h.inner.Sum(in)
}

func (h *sizeHMAC) Reset() {
	h.inner.Reset()
}

func (h *sizeHMAC) MessageSize() uint64 {
	return h.size
}

// HMAC returns the HMAC sum for the given input bytes and HMAC-256 key.
func HMAC(p []byte, hmacKey []byte) []byte {
	macer := NewHMAC(hmacKey)
	_, err := macer.Write(p)
	cerrors.MaybePanic(err) // should never happen b/c sha256.Write always returns nil error
	return macer.Sum(nil)
}

// CheckMACs checks that the ciphertext and uncompressed MACs are consistent with the *api.Metadata.
func CheckMACs(ciphertextMAC, uncompressedMAC MAC, md *api.EntryMetadata) error {
	if err := api.ValidateEntryMetadata(md); err != nil {
		return err
	}
	if md.CiphertextSize != ciphertextMAC.MessageSize() {
		return ErrUnexpectedCiphertextSize
	}
	if !bytes.Equal(md.CiphertextMac, ciphertextMAC.Sum(nil)) {
		return ErrUnexpectedCiphertextMAC
	}
	if md.UncompressedSize != uncompressedMAC.MessageSize() {
		return ErrUnexpectedUncompressedSize
	}
	if !bytes.Equal(md.UncompressedMac, uncompressedMAC.Sum(nil)) {
		return ErrUnexpectedUncompressedMAC
	}
	return nil
}
