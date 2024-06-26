package slices

func index_func[S ~[]E, E any](s S, f func(int) bool) int {
   for i := range s {
      if f(i) {
         return i
      }
   }
   return -1
}
