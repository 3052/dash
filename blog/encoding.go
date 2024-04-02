package encoding

type Namer interface {
   Show() (string, bool)
   Season() (string, bool)
   Episode() (string, bool)
   Title() (string, bool)
   Year() (string, bool)
}

// Show - Season - Episode - Title - Year
func Name(n Namer) string {
   var b []byte
   if v, ok := n.Show(); ok {
      b = append(b, v...)
   }
   if v, ok := n.Season(); ok {
      if b != nil {
         b = append(b, " - "...)
      }
      b = append(b, v...)
   }
   if v, ok := n.Episode(); ok {
      if b != nil {
         b = append(b, " - "...)
      }
      b = append(b, v...)
   }
   if v, ok := n.Title(); ok {
      if b != nil {
         b = append(b, " - "...)
      }
      b = append(b, v...)
   }
   if v, ok := n.Year(); ok {
      if b != nil {
         b = append(b, " - "...)
      }
      b = append(b, v...)
   }
   return string(b)
}
