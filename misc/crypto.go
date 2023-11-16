package misc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func Md5Hash(source string) string {
	hash := md5.New()
	return fmt.Sprintf("%x", hash.Sum([]byte(source)))
}

/*CBC加密*/

// 使用PKCS7进行填充，IOS也是7
func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) (result []byte, err error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length > unpadding {
		result = origData[:(length - unpadding)]
	} else {
		err = errors.New("解密失败，密文数据错误或密码错误")
	}
	return
}

// aes加密，填充秘钥key的16，24，32位分别对应AES-128, AES-192, or AES-256.
func aesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = pKCS7Padding(rawData, blockSize)
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))
	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func aesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData, _ = pKCS7UnPadding(encryptData)
	return encryptData, nil
}

func AesCBCEncrypt(rawData, key []byte) (string, error) {
	data, err := aesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func AesCBCDecrypt(rawData string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := aesCBCDecrypt(data, key)
	if err != nil {
		return "", err
	}
	return string(dnData), nil
}

// 对字符串进行SHA1哈希
func SHA1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// CryptoJS解密(适配CryptoJS, CBC+pKCS7)
func CryptoJSAESDecrypt(rawData string, key []byte) (result string, err error) {
	var pErr error
	data, _ := base64.StdEncoding.DecodeString(rawData)
	if block, newErr := aes.NewCipher(key); newErr != nil {
		err = newErr
	} else {
		blockSize := block.BlockSize()
		if len(key)%blockSize == 0 {
			blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
			decryptData := make([]byte, len(data))
			blockMode.CryptBlocks(decryptData, data)
			if len(decryptData) > 0 {
				if decryptData, pErr = pKCS7UnPadding(decryptData); pErr == nil {
					result = string(decryptData)
				} else {
					err = pErr
				}
			} else {
				err = errors.New("解密失败")
			}
		} else {
			err = errors.New("解密失败, 数据长度检查失败")
		}
	}
	return
}

// CryptoJS加密(适配CryptoJS, CBC+pKCS7)
func CryptoJSAESEncrypt(rawData string, key []byte) (result string, err error) {
	data := []byte(rawData)
	if block, newErr := aes.NewCipher(key); newErr != nil {
		err = newErr
	} else {
		blockSize := block.BlockSize()
		if len(key)%blockSize == 0 {
			blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
			data = pKCS7Padding(data, blockSize)
			encryptData := make([]byte, len(data))
			blockMode.CryptBlocks(encryptData, data)
			if len(encryptData) > 0 {
				result = base64.StdEncoding.EncodeToString(encryptData)
			} else {
				err = errors.New("加密失败")
			}
		} else {
			err = errors.New("解密失败, 数据长度检查失败")
		}
	}
	return
}
