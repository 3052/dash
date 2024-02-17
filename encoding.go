package encoding

import "strconv"

type Namer interface {
   Owner() (string, bool)
   Show() (string, bool)
   Season() (string, bool)
   Episode() (string, bool)
   Title() (string, bool)
   Year() (string, bool)
}

// Owner - Show Season Episode - Title - Year
func Name(n Namer) string {
   var b []byte
   if v, ok := n.Owner(); ok {
      b = append(b, v...)
   }
   if v, ok := n.Show(); ok {
      if b != nil {
         b = append(b, " - "...)
      }
      b = append(b, v...)
   }
   if v, ok := n.Season(); ok {
      if b != nil {
         b = append(b, ' ')
      }
      b = append(b, v...)
   }
   if v, ok := n.Episode(); ok {
      if b != nil {
         b = append(b, ' ')
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
   clean(b)
   return string(b)
}

func clean(path []byte) {
   m := map[byte]bool{
      '"': true,
      '*': true,
      '/': true,
      ':': true,
      '<': true,
      '>': true,
      '?': true,
      '\\': true,
      '|': true,
   }
   for k, v := range path {
      if m[v] {
         path[k] = '-'
      }
   }
}
type Cardinal float64

type Rate float64

type Size float64

func (c Cardinal) String() string {
   units := []unit_measure{
      {1, ""},
      {1e-3, " thousand"},
      {1e-6, " million"},
      {1e-9, " billion"},
      {1e-12, " trillion"},
   }
   return scale(float64(c), units)
}

func (r Rate) String() string {
   units := []unit_measure{
      {1, " byte/s"},
      {1e-3, " kilobyte/s"},
      {1e-6, " megabyte/s"},
      {1e-9, " gigabyte/s"},
      {1e-12, " terabyte/s"},
   }
   return scale(float64(r), units)
}

func (s Size) String() string {
   units := []unit_measure{
      {1, " byte"},
      {1e-3, " kilobyte"},
      {1e-6, " megabyte"},
      {1e-9, " gigabyte"},
      {1e-12, " terabyte"},
   }
   return scale(float64(s), units)
}

func scale(value float64, units []unit_measure) string {
   var unit unit_measure
   for _, unit = range units {
      if unit.factor * value < 1000 {
         break
      }
   }
   return label(value, unit)
}

type unit_measure struct {
   factor float64
   name string
}

func label(value float64, unit unit_measure) string {
   var prec int
   if unit.factor != 1 {
      prec = 2
      value *= unit.factor
   }
   return strconv.FormatFloat(value, 'f', prec, 64) + unit.name
}

type Percent float64

func (p Percent) String() string {
   unit := unit_measure{100, " %"}
   return label(float64(p), unit)
}
