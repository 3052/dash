package dash

import (
   "io"
   "net/http"
)

func write(w io.Writer) error {
   res, err := http.Get("http://example.com")
   if err != nil {
      return err
   }
   defer res.Body.Close()
   // Init
   io.CopyN(w, res.Body, 10)
   // Segments
   io.WriteString(w, "[HELLO WORLD]")
   io.Copy(w, res.Body)
   return nil
}
