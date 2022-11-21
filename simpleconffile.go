//
//
// simpleconffile package
//
// This is not the intended to be a high class way of storing off
// a config file, let alone encrypt the contents.  However, for simple
// tools, it is a great lift.
//
// Allows one to quickly write a struct into a conf file and read back. 
//
// Also if wanted - encrypt strings so a someone cannot read them in the
// exported conf file
//
//
package simpleconffile


import (
	"os"
	"fmt"
        "io"
        "crypto/aes"
        "crypto/rand"
        "crypto/cipher"
        "encoding/base64"
        "encoding/json"
        "bytes"
        "github.com/seldonsmule/logmsg"

)

//
// internal structure of all our stuff
//

type SimpleConfFile struct {

  ByteEncryptKey []byte // holds the key for the crypto work
  Filename string // name of the file to write/read 

}

// New()
//
// Helper function to create and setup the object
//

func New(keystring string, confname string) *SimpleConfFile {

  s := new(SimpleConfFile)

  s.Filename = confname

  //switch len(s.GetEncryptKey()) {
  switch len(keystring) {

    case 16:
      fallthrough
    case 24:
      fallthrough
    case 32:
    

    default:
      logmsg.Print(logmsg.Error, "Key must be 16, 24, or 32 bytes long->", len(keystring) )
      return nil

  }

  s.SetEncryptKey(keystring)

  return (s)
}

//
// GetEncryptKey()
//
// Returns the key used to encrypt/decrypt data
//

func (pS *SimpleConfFile) GetEncryptKey() []byte{
  return(pS.ByteEncryptKey)
}

//
// SaveConf
//
// Saves off a struct (interface) to a predefined configuration filename
//
//   v - Of type interface{}.  Pass in a struct
//

func (pS *SimpleConfFile) SaveConf(v interface{}) bool{

  j, err := json.Marshal(v)

  if(err != nil){
    logmsg.Print(logmsg.Error, "json.Marshal failed: ", err)
    return false
  }

  //fmt.Println(string(j))

  writeFile, err := os.Create(pS.Filename)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to write config: ", err)
     return false
  }
  defer writeFile.Close()

  writeFile.Write(j)
  //os.Stdout.Write(j)
  writeFile.Close()

  return true
}

//
// ReadConf
//
// Reads in a struct (interface) from a predefined configuration filename
//
//   v - Of type interface{}.  Pass in a struct
//

func (pS *SimpleConfFile) ReadConf(v interface{}) bool{

  file, err := os.Open(pS.Filename) // For read access.

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to config config: ", err," ",  pS.Filename)
     return false
  }

  defer file.Close()

  data := make([]byte, 500)

  count, err := file.Read(data)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to read config: ", err, count)
     return false
  }

  err = json.NewDecoder(bytes.NewReader(data)).Decode(&v)

  if err != nil {
     logmsg.Print(logmsg.Error,"Unable to decode config: ", err)
     return false
  }


  return true
}

//
// Dump()
//
// For debugging purposes - dumps the contents of SimpleConfFile
//

func (pS *SimpleConfFile) Dump(){

  fmt.Printf("SimpleConfFile.ByteEncryptKey [%v]\n", pS.ByteEncryptKey)
  fmt.Printf("SimpleConfFile.ByteEncryptKey [%s]\n", string(pS.ByteEncryptKey))
  fmt.Printf("SimpleConfFile.Filename [%s]\n", pS.Filename)

}

//
// SetEncryptKey() - Stores off a key for the crypto tools
//
//  input - 16, 24 or 32 bit long string
//

func (pS *SimpleConfFile) SetEncryptKey(input string) {

  pS.ByteEncryptKey = []byte(input)

}

//
// EncryptString() - Encrypts a string
//
//    text - string to encrypt
//
// addapted from https://gist.github.com/manishtpatel/8222606
//
// encrypt string to base64 crypto using AES

func (pS *SimpleConfFile) EncryptString(text string) string {

  plaintext := []byte(text)

  block, err := aes.NewCipher(pS.GetEncryptKey())
  if err != nil {
    fmt.Println("aes.NewCipher err: ", err)
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  ciphertext := make([]byte, aes.BlockSize+len(plaintext))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    panic(err)
  }

  stream := cipher.NewCFBEncrypter(block, iv)
  stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

  // convert to base64
  return base64.URLEncoding.EncodeToString(ciphertext)

}

//
// DecryptString() - Decrypts a string
//
//    text - string to decrypt
//
// addapted from https://gist.github.com/manishtpatel/8222606
//
// encrypt string to base64 crypto using AES
//
func (pS *SimpleConfFile) DecryptString(cryptoText string) string {

  ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

  block, err := aes.NewCipher(pS.GetEncryptKey())
  if err != nil {
    panic(err)
  }

  // The IV needs to be unique, but not secure. Therefore it's common to
  // include it at the beginning of the ciphertext.
  if len(ciphertext) < aes.BlockSize {
    panic("ciphertext too short")
  }
  iv := ciphertext[:aes.BlockSize]
  ciphertext = ciphertext[aes.BlockSize:]

  stream := cipher.NewCFBDecrypter(block, iv)

  // XORKeyStream can work in-place if the two arguments are the same.
  stream.XORKeyStream(ciphertext, ciphertext)

  return fmt.Sprintf("%s", ciphertext)

}

