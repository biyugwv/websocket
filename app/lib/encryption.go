package lib

import(
    "websocket/app/log"
    "io"
    "fmt"
    "crypto/md5"
    "crypto/aes"
    "encoding/hex"
    "crypto/cipher"
    "strings"
    "bytes"

)


func Md5(str string) string {
        m := md5.New()
        _, err := io.WriteString(m, str)
        if err != nil {
                return  ""
        }
        arr := m.Sum(nil)
        return fmt.Sprintf("%x", arr)
}

func AESEncodeStr(src, key string) string {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        panic(err)
        log.Write("error","key error1" + fmt.Sprintf("%s",err))
    }
    if src == "" {
        log.Write("error","plain content empty")
    }
    ivspec := []byte("0000000000000000")
    ecb := cipher.NewCBCEncrypter(block, ivspec)
    content := []byte(src)
    content = PKCS5Padding(content, block.BlockSize())
    crypted := make([]byte, len(content))
    ecb.CryptBlocks(crypted, content)
    return hex.EncodeToString(crypted)
}

func AESDecodeStr(crypt, key string) string {
    crypted, err := hex.DecodeString(strings.ToLower(crypt))
    if err != nil || len(crypted) == 0 {
        log.Write("error","plain content empty")
    }
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        log.Write("error","key error1" + fmt.Sprintf("%s",err))
    }
    ivspec := []byte("0000000000000000")
    ecb := cipher.NewCBCDecrypter(block, ivspec)
    decrypted := make([]byte, len(crypted))
    ecb.CryptBlocks(decrypted, crypted)

    return string(PKCS5Trimming(decrypted))
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
    padding := encrypt[len(encrypt)-1]
    return encrypt[:len(encrypt)-int(padding)]
}
