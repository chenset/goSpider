package bloomFilter

import (
	"github.com/willf/bloom"
	"fmt"
)

func Filter() {
	filter := bloom.New(20000, 5) // load of 20, 5 keys
	filter.Add([]byte("Love"))

	if filter.Test([]byte("Love")) {

		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}

}
