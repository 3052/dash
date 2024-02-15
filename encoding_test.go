package encoding

import (
   "fmt"
   "testing"
)

func TestPercent(t *testing.T) {
   fmt.Println(Percent(1234) / 10000)
}
