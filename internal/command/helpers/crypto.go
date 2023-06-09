package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/mundanelizard/envi/internal/lockfile"
	"io"
	"os"
	"strings"
)

func encryptCompressedEnvironment(dir, secret string) (string, error) {
	in, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer in.Close()

	data, err := io.ReadAll(in)
	if err != nil {
		return "", err
	}

	outDir := dir + ".enc"
	lock := lockfile.New(outDir)
	err = lock.Hold()
	if err != nil {
		return "", err
	}
	defer lock.Commit()

	cphr, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(data))
	iv := make([]byte, aes.BlockSize)
	cipher.NewCFBEncrypter(cphr, iv).XORKeyStream(cipherText, data)

	err = lock.Write(cipherText)
	if err != nil {
		return "", err
	}

	err = lock.Commit()
	if err != nil {
		return "", err
	}

	return outDir, nil
}

func DecryptCompressedEnvironment(dir, secret string) (string, error) {
	key := []byte(secret)
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encryptedFile, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer encryptedFile.Close()

	content, err := io.ReadAll(encryptedFile)
	if err != nil {
		return "", err
	}

	plainText := make([]byte, len(content))
	iv := make([]byte, aes.BlockSize)
	cipher.NewCFBDecrypter(aesCipher, iv).XORKeyStream(plainText, content)

	outDir := strings.ReplaceAll(dir, ".enc", "")

	decryptedFile, err := os.Create(outDir)
	if err != nil {
		return "", nil
	}
	defer decryptedFile.Close()

	_, err = decryptedFile.Write(plainText)
	if err != nil {
		return "", nil
	}

	return outDir, nil
}
