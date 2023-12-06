package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

var (
	secretKey string = "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
)

func encrypt(plaintext string) string {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}
	//nonce: like IV(initial value)
	nonce := make([]byte, gcm.NonceSize())
	tempres, err := rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	fmt.Printf("tempres %+v\n", tempres)

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return string(ciphertext)
}
func decrypt(ciphertext string) string {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}
func main() {
	ciphertext1 := encrypt("This is some sensitive information")
	fmt.Printf("Encrypted ciphertext 1: %x \n", ciphertext1)

	plaintext1 := decrypt(ciphertext1)
	fmt.Printf("Decrypted plaintext 1: %s \n", plaintext1)

	ciphertext2 := encrypt("Hello")
	fmt.Printf("Encrypted ciphertext 2: %x \n", ciphertext2)

	plaintext2 := decrypt(ciphertext2)
	fmt.Printf("Decrypted plaintext 2: %s \n", plaintext2)
}
