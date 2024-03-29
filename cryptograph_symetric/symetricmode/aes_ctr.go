package symetricmode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Error is %+v\n", err)
		log.Fatal("ERROR:", err)
	}
}

func encrypt(keyByte []byte, nonce []byte, plainText string) string {

	plainTextByte := []byte(plainText)
	cipherTextByte := make([]byte, len(plainText))

	// GET CIPHER BLOCK USING KEY
	block, err := aes.NewCipher(keyByte)
	checkErr(err)

	// GET CTR
	ctr := cipher.NewCTR(block, nonce)

	// ENCRYPT DATA
	ctr.XORKeyStream(cipherTextByte, plainTextByte)

	// RETURN HEX
	cipherText := hex.EncodeToString(cipherTextByte)
	return cipherText
}

func decrypt(keyByte []byte, nonce []byte, cipherText string) string {

	cipherTextByte, _ := hex.DecodeString(cipherText)
	plainTextByte := make([]byte, len(cipherTextByte))

	// CHECK cipherTextByte
	// CBC mode always works in whole blocks.
	if len(cipherTextByte)%aes.BlockSize != 0 {
		panic("cipherTextByte is not a multiple of the block size")
	}

	// GET CIPHER BLOCK USING KEY
	block, err := aes.NewCipher(keyByte)
	checkErr(err)

	// GET CTR
	ctr := cipher.NewCTR(block, nonce)

	// DECRYPT DATA
	ctr.XORKeyStream(plainTextByte, cipherTextByte)

	// RETURN STRING
	plainText := string(plainTextByte[:])
	return plainText
}

func AES_CTR() {
	// DATA
	// IN CBC Must be Block Size of AES (Multiple of 16)
	plainText := "This is AES-256 CTR (32 Bytes)!!"
	if len(plainText)%aes.BlockSize != 0 {
		panic("Plaintext is not a multiple of the block size")
	}
	fmt.Printf("\nOriginal Text:           %s\n\n", plainText)

	// KEY
	keyText := "myverystrongpasswordo32bitlength"
	keyByte := []byte(keyText)
	fmt.Printf("The 32-byte Key:         %s\n", keyText)

	// CREATE A NONCE
	// For this example I'm not including in the cipherText
	nonce := make([]byte, aes.BlockSize)
	// Populates the nonce with a cryptographically secure random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	fmt.Printf("The Nonce:               %x\n\n", nonce)

	// ENCRYPT
	cipherText := encrypt(keyByte, nonce, plainText)
	fmt.Printf("Encrypted Text:          %s\n", cipherText)

	// DECRYPT
	plainText = decrypt(keyByte, nonce, cipherText)
	fmt.Printf("Decrypted Text:          %s\n\n", plainText)

}
