package symetricmode

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func checkErrCBC(err error) {
	if err != nil {
		fmt.Printf("Error is %+v\n", err)
	}
}

func encryptCBC(keyByte []byte, nonce []byte, plaintext string) string {
	plainTextByte := []byte(plaintext)
	cipherTextByte := make([]byte, len(plainTextByte))

	block, err := aes.NewCipher(keyByte)
	checkErrCBC(err)

	cbc := cipher.NewCBCEncrypter(block, nonce)

	//ENCRYPT DATA
	cbc.CryptBlocks(cipherTextByte, plainTextByte)

	//RETURN
	cipherText := hex.EncodeToString(cipherTextByte)
	return cipherText
}
func decryptCBC(keyByte []byte, nonce []byte, cipherText string) string {
	cipherTextByte, _ := hex.DecodeString(cipherText)
	plainTextByte := make([]byte, len(cipherTextByte))

	//CHECK cipherTextByte
	//CBC mode always work in whole blocks.
	if len(cipherTextByte)%aes.BlockSize != 0 {
		panic("cipherTextByte is not a multiple of the block size")
	}

	// GET CIPHER BLOCK USING KEY
	block, err := aes.NewCipher(keyByte)
	checkErrCBC(err)

	//GET CBC DECRYPTER
	cbc := cipher.NewCBCDecrypter(block, nonce)

	//DECRYPT DATA
	cbc.CryptBlocks(plainTextByte, cipherTextByte)
	plainText := string(plainTextByte)
	return plainText
}
func AES_CBC() {
	//DATA
	//IN CBC Must be Block Size
	plainText := "This is AES-256 CBC (32 Bytes)!!"
	if len(plainText)%aes.BlockSize != 0 {
		panic("plainText is not a multiple of the block size")
	}

	fmt.Printf("\nOriginal Text:           %s\n\n", plainText)
	//KEY
	keyText := "myverystrongpasswordo32bitlength"
	keyByte := []byte(keyText)
	//fmt.Printf("The 32-byte Key:           %s\n", keyText)
	fmt.Printf("The 32-byte Key:         %s\n", keyText)
	//CREATE A NONCE
	//For this example I'm not including in the cipherText
	nonce := make([]byte, aes.BlockSize)
	//Populate the nonce with a cryptographically secure random sequnce
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	//fmt.Printf("The Nonce:           %x\n\n", nonce)
	fmt.Printf("The Nonce:               %x\n\n", nonce)
	cipherText := encryptCBC(keyByte, nonce, plainText)
	//fmt.Printf("Encrypted Text:           %s\n\n", cipherText)
	fmt.Printf("Encrypted Text:          %s\n", cipherText)

	//DECRYPT
	plainText = decryptCBC(keyByte, nonce, cipherText)
	fmt.Printf("Decrypted Text:          %s\n\n", plainText)
}
