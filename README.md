# bencode
Bencode implementation in Go

## Installation
```bash
$ go get github.com/fuchsi/bencode
```

## Usage
### Encode
bencode.Encode takes a map[string]interface{} as argument and returns a []byte.  
Example:
```go
package main

import (
	"fmt"
	
	"github.com/fuchsi/bencode"
)

func main() {
	dict := make(map[string]interface{})
	dict["some string key"] = "string value"
	dict["some int key"] = 3735928559
	
	fmt.Printf("bencode encoded dict: %s\n", bencode.Encode(dict))
}
```
### Decode
bencode.Decode takes an io.Reader as argument and returns (map[string]interface{}, error)  
Example:
```go
package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/fuchsi/bencode"
)

func main() {
	file, err := os.Open(os.Args[1]) 
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	dict, err := bencode.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("some string key: %s\n", dict["some string key"].(string))
	fmt.Printf("some int key: %d\n", dict["some int key"].(int64))
}
```