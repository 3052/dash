package dash

import (
   "os"
   "testing"
)

func Test_Write(t *testing.T) {
   write(os.Stdout)
}
