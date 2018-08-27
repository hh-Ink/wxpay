package wxpay

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"sync"
)

var (
	syncAesMutex sync.Mutex
	commonAeskey []byte
)

func SetAesKey(key string) (err error) {
	syncAesMutex.Lock()
	defer syncAesMutex.Unlock()
	b := []byte(key)
	if len(b) == 16 || len(b) == 24 || len(b) == 32 {
		commonAeskey = b
		return nil
	}
	return errors.New(fmt.Sprintf("key size is not 16 or 24 or 32, but %d", len(b)))
}

func AesECBEncrypt(plaintext []byte) (ciphertext []byte, err error) {
	plaintext = PKCS5Padding(plaintext, aes.BlockSize)
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, errors.New("plaintext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(commonAeskey)
	if err != nil {
		return nil, err
	}
	ciphertext = make([]byte, len(plaintext))
	NewECBEncrypter(block).CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func AesECBDecrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(commonAeskey)
	if err != nil {
		return nil, err
	}
	NewECBDecrypter(block).CryptBlocks(ciphertext, ciphertext)
	plaintext = PKCS5UnPadding(ciphertext)
	return plaintext, nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

type ecbDecrypter ecb

// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
// mode, using the given Block.
func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
