// Package crypto: telefon/e-posta gibi kişisel iletişim verilerini
// DB'ye yazmadan önce şifrelemek, okurken çözmek için AES-256-GCM sarmalayıcısı.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type Encryptor struct {
	key []byte // 32 byte olmalı (AES-256)
}

// NewEncryptor: hex encode edilmiş 32 byte'lık anahtarı env'den alır.
// Anahtar üretmek için: `openssl rand -hex 32`
func NewEncryptor(hexKey string) (*Encryptor, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("ENCRYPTION_KEY geçerli bir hex string değil")
	}
	if len(key) != 32 {
		return nil, errors.New("ENCRYPTION_KEY 32 byte (64 hex karakter) olmalı")
	}
	return &Encryptor{key: key}, nil
}

// Encrypt: düz metni şifreler, (cipherText, nonce) döner. İkisi de ayrı ayrı DB'ye yazılır.
func (e *Encryptor) Encrypt(plain string) (cipherText []byte, nonce []byte, err error) {
	if plain == "" {
		return nil, nil, nil
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	cipherText = gcm.Seal(nil, nonce, []byte(plain), nil)
	return cipherText, nonce, nil
}

// Decrypt: şifreli veri + nonce'tan düz metni geri döndürür.
// Bu fonksiyon YALNIZCA yetkili, audit-log'lanan servis çağrılarında kullanılmalı
// (ör. mail worker'ın gönderim için adresi çözmesi, admin'in açık rızayla iletişim
// bilgisini görüntülemesi). Sonuç asla API response'unda ham olarak dönülmemeli.
func (e *Encryptor) Decrypt(cipherText, nonce []byte) (string, error) {
	if len(cipherText) == 0 {
		return "", nil
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
