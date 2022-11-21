package main


import (
	"fmt"
	"reflect"
	"strings"
        "github.com/seldonsmule/simpleconffile"

)

type MyConf struct {

  Email string
  Passwd string
  TokenLocation string
  Url string

}

func (c *MyConf) Dump(){
  fmt.Println("Email = ", c.Email)
  fmt.Println("Passwd = ", c.Passwd)
  fmt.Println("TokenLocation = ", c.TokenLocation)
  fmt.Println("Url = ", c.Url)

}

const EMAIL = "simple@is.com"
const PASSWD = "supercryptic123"
const TOKENLOCATION = "$HOME/tmp"
const URL = "https://somewhere.overtherainbow.com"

func main() {

  fmt.Println("Example use of SimpleConfFile")

  simple := simpleconffile.New("1234567890123456", "my.conf")

  conf := new(MyConf)

  conf.Email = simple.EncryptString(EMAIL)
  conf.Passwd = simple.EncryptString(PASSWD)
  conf.TokenLocation = TOKENLOCATION
  conf.Url = URL

  //simple.Dump()

  if(!simple.SaveConf(conf)){
    fmt.Println("FAIL - simple.SaveConf()")
    return
  }
  
  fmt.Println("PASS - simple.SaveConf()")

  readconf := new(MyConf)

  if(!simple.ReadConf(readconf)){
    fmt.Println("FAIL - simple.ReadConf()")
    return
  }

  fmt.Println("PASS - simple.ReadConf()")

  if(!reflect.DeepEqual(conf, readconf)) {
    fmt.Println("FAIL - two structs are not equal")
    fmt.Println("Dump 1st conf")
    conf.Dump()
    fmt.Println("Dump 2st readconf")
    readconf.Dump()
    return
  }

  fmt.Println("PASS - Configs are the same")


  if(strings.Compare(simple.DecryptString(readconf.Email), EMAIL) != 0){
    fmt.Println("FAIL - simple.DecryptString(EMAIL)")
    return
  }

  fmt.Println("PASS - simple.DecryptString(EMAIL)")

  if(strings.Compare(simple.DecryptString(readconf.Passwd), PASSWD) != 0){
    fmt.Println("FAIL - simple.DecryptString(PASSED)")
    return
  }

  fmt.Println("PASS - simple.DecryptString(PASSED)")


}
