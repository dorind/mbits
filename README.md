**m**emory **bit** buffer

use it as you see fit, I'm interested in speeding up String(), SetAll() and CountBits() functions

should be relatively easy to use, check `mbits_test.go` for some examples of usage

### get

```shell
$ go get github.com/dorind/mbits
```

### example

```go
package main

import (
	"fmt"

	"github.com/dorind/mbits"
)

func main() {
	buff_odd := mbits.NewBitBuffer(0)
	buff_even := mbits.NewBitBuffer(0)

	for i := uint(0); i < 128; i++ {
		if i&1 == 1 {
			buff_odd.Set(i)
		} else {
			buff_even.Set(i)
		}
	}

	fmt.Println("O", buff_odd.String())
	fmt.Println("E", buff_even.String())
}
```



