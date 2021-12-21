package envelope

import (
	"crypto/rand"
	"crypto/rsa"

	internal_aes "github.com/SatorNetwork/sator-api/internal/encryption/aes"
	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
)

type Envelope struct {
	Ciphertext     []byte `json:"ciphertext"`
	CipheredAESKey []byte `json:"ciphered_aes_key"`
}

type Encryptor struct {
	publicKey *rsa.PublicKey
}

func NewEncryptor(publicKey *rsa.PublicKey) *Encryptor {
	return &Encryptor{
		publicKey: publicKey,
	}
}

func (e *Encryptor) Encrypt(plaintext []byte) (*Envelope, error) {
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, err
	}

	ciphertext, err := internal_aes.Encrypt(aesKey, plaintext)
	if err != nil {
		return nil, err
	}
	cipheredAESKey, err := internal_rsa.EncryptWithPublicKey(aesKey, e.publicKey)
	if err != nil {
		return nil, err
	}

	return &Envelope{
		Ciphertext:     ciphertext,
		CipheredAESKey: cipheredAESKey,
	}, nil
}

type Decryptor struct {
	privateKey *rsa.PrivateKey
}

func NewDecryptor(privateKey *rsa.PrivateKey) *Decryptor {
	return &Decryptor{
		privateKey: privateKey,
	}
}

func (d *Decryptor) Decrypt(e *Envelope) ([]byte, error) {
	aesKey, err := internal_rsa.DecryptWithPrivateKey(e.CipheredAESKey, d.privateKey)
	if err != nil {
		return nil, err
	}

	decrypted, err := internal_aes.Decrypt(aesKey, e.Ciphertext)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
