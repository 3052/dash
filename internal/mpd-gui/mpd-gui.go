package main

import (
   "41.neocities.org/dash"
   "fmt"
   "io"
   "github.com/rodrigocfd/windigo/ui"
   "github.com/rodrigocfd/windigo/win"
   "github.com/rodrigocfd/windigo/win/co"
   "net/http"
   "runtime"
   "strings"
)

func (w *window) New() {
   w.main = ui.NewWindowMain(nil)
   w.edit = ui.NewEdit(w.main, nil)
   w.button = ui.NewButton(w.main, ui.ButtonOpts().
      Position(win.POINT{X: 200}).
      Text("Text"),
   )
   w.button.On().BnClicked(func() {
      data, err := get(w.edit.Text())
      if err != nil {
         data = err.Error()
      }
      w.main.Hwnd().MessageBox(data, "", co.MB_ICONINFORMATION)
   })
}

type window struct {
   edit   ui.Edit
   button ui.Button
   main   ui.WindowMain
}

func get(address string) (string, error) {
   resp, err := http.Get(address)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }
   var mpd dash.Mpd
   err = mpd.Unmarshal(data)
   if err != nil {
      return "", err
   }
   var (
      data1 strings.Builder
      line bool
   )
   for represent := range mpd.Representation() {
      if line {
         fmt.Fprintln(&data1)
      } else {
         line = true
      }
      fmt.Fprintln(&data1, &represent)
   }
   return data1.String(), nil
}

func main() {
   runtime.LockOSThread()
   var win window
   win.New()
   win.main.RunAsMain() // ...and run
}
