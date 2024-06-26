package slices

import (
   "fmt"
   "strings"
   "testing"
)

func TestSlices(t *testing.T) {
   vs := []string{"Sunday", "Monday", "Tuesday"}
   i := index_func(vs, func(j int) bool {
      return strings.HasPrefix(vs[j], "Mon")
   })
   fmt.Println(i == 1)
}
