package main

import (
   "log"
   "os"
   "time"
)

func main() {
   logger := log.New(os.Stdout, "", log.Ltime)
   for range 99 {
      logger.Print("hello world")
      time.Sleep(time.Second)
   }
}
