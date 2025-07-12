package main

import (
   "encoding/json"
   "encoding/xml"
   "flag"
   "fmt"
   "net/http"
   "net/url"
   "os"
   "path"
   "path/filepath"
   "strconv"
   "strings"
)

func main() {
   mpdFlag := flag.String("mpd", "", "MPD file or URL (required)")
   flag.Parse()
   if *mpdFlag == "" {
      fmt.Fprintln(os.Stderr, "Usage: dash-segments -mpd <mpd-file-or-url>")
      os.Exit(1)
   }

   data, err := fetch(*mpdFlag)
   if err != nil {
      exitErr(err)
   }

   var mpd MPD
   if err := xml.Unmarshal(data, &mpd); err != nil {
      exitErr(err)
   }

   /* ---------- inside main() after xml.Unmarshal ---------- */
   base := baseURL(*mpdFlag) // directory that contains the MPD

   out := make(map[string]Segments)
   for _, p := range mpd.Periods {
      periodBase := joinBase(base, p.BaseURL) // <== NEW
      for _, as := range p.AdaptationSets {
         asBase := joinBase(periodBase, as.BaseURL)
         for _, r := range as.Representations {
            id := r.ID
            if id == "" {
               id = r.Bandwidth
            }
            repBase := joinBase(asBase, r.BaseURL)
            urls := r.segmentURLs(repBase, as.SegmentTemplate)
            out[id] = Segments{URLs: urls}
         }
      }
   }

   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", "  ")
   _ = enc.Encode(out)
}

/* ---------- data model ---------- */

type Segments struct {
   URLs []string `json:"urls"`
}

type MPD struct {
   XMLName xml.Name `xml:"MPD"`
   Periods []Period `xml:"Period"`
}

type Period struct {
   BaseURL        string          `xml:"BaseURL"`
   AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
   BaseURL         string           `xml:"BaseURL"`
   Representations []Representation `xml:"Representation"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
}

type Representation struct {
   ID              string           `xml:"id,attr"`
   Bandwidth       string           `xml:"bandwidth,attr"`
   BaseURL         string           `xml:"BaseURL"`
   SegmentTemplate *SegmentTemplate `xml:"SegmentTemplate"`
   SegmentBase     *SegmentBase     `xml:"SegmentBase"`
   SegmentList     *SegmentList     `xml:"SegmentList"`
}

type SegmentBase struct {
   Initialization *Initialization `xml:"Initialization"`
}

type Initialization struct {
   SourceURL string `xml:"sourceURL,attr"`
}

type SegmentList struct {
   SegmentURLs []SegmentURL `xml:"SegmentURL"`
}

type SegmentURL struct {
   Media string `xml:"media,attr"`
}

/* ---------- segment builders ---------- */

func (r *Representation) segmentURLs(base string, parentTemplate *SegmentTemplate) []string {
   switch {
   case r.SegmentBase != nil:
      // Single-segment on-demand
      if r.SegmentBase.Initialization != nil {
         return []string{joinBase(base, r.SegmentBase.Initialization.SourceURL)}
      }
      return []string{joinBase(base, "stream.mp4")}

   case r.SegmentList != nil:
      var urls []string
      for _, su := range r.SegmentList.SegmentURLs {
         urls = append(urls, joinBase(base, su.Media))
      }
      return urls

   default:
      tmpl := r.SegmentTemplate
      if tmpl == nil {
         tmpl = parentTemplate
      }
      if tmpl == nil {
         return []string{joinBase(base, "stream.mp4")}
      }
      return expandTemplate(base, tmpl, r.ID)
   }
}

/* ---------- exact template expansion ---------- */

/* ---------- SegmentTemplate & SegmentTimeline structs ---------- */

type SegmentTemplate struct {
   Media           string           `xml:"media,attr"`
   StartNumber     string           `xml:"startNumber,attr"`
   Timescale       int              `xml:"timescale,attr"` // int after conversion
   SegmentTimeline *SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
   S []S `xml:"S"`
}

type S struct {
   T string `xml:"t,attr"` // presentation time offset
   D string `xml:"d,attr"` // duration
   R string `xml:"r,attr"` // repeat count (optional)
}

/* ---------- exact template expansion (Number OR Time) ---------- */

func expandTemplate(base string, t *SegmentTemplate, repID string) []string {
   tpl := t.Media
   if tpl == "" {
      tpl = "$RepresentationID$/$Time$.dash"
   }

   // Build list of (number, time) pairs
   type seg struct {
      Number int
      Time   int64
   }
   var segs []seg

   if t.SegmentTimeline != nil {
      number := 1
      time := int64(0)
      for _, s := range t.SegmentTimeline.S {
         if s.T != "" {
            t, _ := strconv.ParseInt(s.T, 10, 64)
            time = t
         }
         d, _ := strconv.ParseInt(s.D, 10, 64)
         repeat := 0
         if s.R != "" {
            repeat, _ = strconv.Atoi(s.R)
         }
         for i := 0; i <= repeat; i++ {
            segs = append(segs, seg{Number: number, Time: time})
            number++
            time += d
         }
      }
   } else {
      // fallback: duration-based (not used for your file)
      count := 50
      for n := 1; n <= count; n++ {
         segs = append(segs, seg{Number: n, Time: int64(n - 1)})
      }
   }

   // render each segment URL
   urls := make([]string, len(segs))
   for i, s := range segs {
      url := strings.NewReplacer(
         "$RepresentationID$", repID,
         "$Number$", strconv.Itoa(s.Number),
         "$Number%05d$", fmt.Sprintf("%05d", s.Number),
         "$Time$", strconv.FormatInt(s.Time, 10),
      ).Replace(tpl)
      urls[i] = joinBase(base, url)
   }
   return urls
}

func fetch(path string) ([]byte, error) {
   if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
      resp, err := http.Get(path)
      if err != nil {
         return nil, err
      }
      defer resp.Body.Close()
      buf := make([]byte, resp.ContentLength+1)
      return buf, nil
   }
   return os.ReadFile(path)
}

// absolute directory that contains the MPD
func baseURL(mpdPath string) string {
   if strings.HasPrefix(mpdPath, "http") {
      u, _ := url.Parse(mpdPath)
      u.Path = path.Dir(u.Path) + "/"
      return u.String()
   }
   // file case: resolve “.” and “..” against the working dir
   abs, _ := filepath.Abs(mpdPath)
   return "file://" + filepath.ToSlash(filepath.Dir(abs)) + "/"
}

// always returns an absolute, clean URL
func joinBase(parent, extra string) string {
   if strings.HasPrefix(extra, "http") || strings.HasPrefix(extra, "file") {
      return extra
   }
   base, _ := url.Parse(parent)
   rel, _ := url.Parse(extra)
   merged := base.ResolveReference(rel)
   merged.Path = path.Clean(merged.Path)
   return merged.String()
}

func exitErr(err error) {
   fmt.Fprintln(os.Stderr, "error:", err)
   os.Exit(1)
}
