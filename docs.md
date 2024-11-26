# Overview

package `dash`

## Index

- [Types](#types)
  - [type AdaptationSet](#type-adaptationset)
    - [func (a *AdaptationSet) GetPeriod() *Period](#func-adaptationset-getperiod)
  - [type BaseUrl](#type-baseurl)
    - [func (b *BaseUrl) UnmarshalText(text []byte) error](#func-baseurl-unmarshaltext)
  - [type ContentProtection](#type-contentprotection)
  - [type Duration](#type-duration)
    - [func (d *Duration) UnmarshalText(text []byte) error](#func-duration-unmarshaltext)
  - [type Mpd](#type-mpd)
  - [type Period](#type-period)
  - [type Pssh](#type-pssh)
    - [func (p *Pssh) UnmarshalText(src []byte) error](#func-pssh-unmarshaltext)
  - [type Range](#type-range)
    - [func (r *Range) MarshalText() ([]byte, error)](#func-range-marshaltext)
    - [func (r *Range) UnmarshalText(text []byte) error](#func-range-unmarshaltext)
  - [type Representation](#type-representation)
    - [func Unmarshal(text []byte, base *url.URL) ([]Representation, error)](#func-unmarshal)
    - [func (r *Representation) Ext() (string, bool)](#func-representation-ext)
    - [func (r *Representation) GetAdaptationSet() *AdaptationSet](#func-representation-getadaptationset)
    - [func (r *Representation) GetBaseUrl() (*BaseUrl, bool)](#func-representation-getbaseurl)
    - [func (r *Representation) GetMimeType() string](#func-representation-getmimetype)
    - [func (r *Representation) Initialization() (string, bool)](#func-representation-initialization)
    - [func (r *Representation) Media() []string](#func-representation-media)
    - [func (r *Representation) String() string](#func-representation-string)
    - [func (r *Representation) Widevine() (Pssh, bool)](#func-representation-widevine)
  - [type SegmentBase](#type-segmentbase)
  - [type SegmentTemplate](#type-segmenttemplate)
- [Source files](#source-files)

## Types

### type [AdaptationSet](./dash.go#L14)

```go
type AdaptationSet struct {
  Codecs            string `xml:"codecs,attr"`
  ContentProtection []ContentProtection
  Height            uint64 `xml:"height,attr"`
  Lang              string `xml:"lang,attr"`
  MaxHeight         int    `xml:"maxHeight,attr"`
  MaxWidth          int    `xml:"maxWidth,attr"`
  MimeType          string `xml:"mimeType,attr"`
  Representation    []Representation
  Role              *struct {
    Value string `xml:"value,attr"`
  }
  SegmentTemplate *SegmentTemplate
  Width           uint64 `xml:"width,attr"`
  // contains filtered or unexported fields
}
```

### func (*AdaptationSet) [GetPeriod](./dash.go#L31)

```go
func (a *AdaptationSet) GetPeriod() *Period
```

### type [BaseUrl](./dash.go#L35)

```go
type BaseUrl struct {
  Url *url.URL
}
```

### func (*BaseUrl) [UnmarshalText](./dash.go#L39)

```go
func (b *BaseUrl) UnmarshalText(text []byte) error
```

### type [ContentProtection](./dash.go#L44)

```go
type ContentProtection struct {
  Pssh        Pssh   `xml:"pssh"`
  SchemeIdUri string `xml:"schemeIdUri,attr"`
}
```

### type [Duration](./dash.go#L60)

```go
type Duration struct {
  Duration time.Duration
}
```

### func (*Duration) [UnmarshalText](./dash.go#L49)

```go
func (d *Duration) UnmarshalText(text []byte) error
```

### type [Mpd](./dash.go#L64)

```go
type Mpd struct {
  BaseUrl                   *BaseUrl  `xml:"BaseURL"`
  MediaPresentationDuration *Duration `xml:"mediaPresentationDuration,attr"`
  Period                    []Period
}
```

### type [Period](./dash.go#L70)

```go
type Period struct {
  AdaptationSet []AdaptationSet
  BaseUrl       *BaseUrl  `xml:"BaseURL"`
  Duration      *Duration `xml:"duration,attr"`
  Id            string    `xml:"id,attr"`
  // contains filtered or unexported fields
}
```

### type [Pssh](./dash.go#L85)

```go
type Pssh []byte
```

### func (*Pssh) [UnmarshalText](./dash.go#L87)

```go
func (p *Pssh) UnmarshalText(src []byte) error
```

### type [Range](./dash.go#L105)

```go
type Range struct {
  Start uint64
  End   uint64
}
```

SegmentIndexBox uses:
unsigned int(32) subsegment_duration;
but range values can exceed 32 bits

### func (*Range) [MarshalText](./dash.go#L96)

```go
func (r *Range) MarshalText() ([]byte, error)
```

### func (*Range) [UnmarshalText](./dash.go#L110)

```go
func (r *Range) UnmarshalText(text []byte) error
```

### type [Representation](./dash.go#L296)

```go
type Representation struct {
  Bandwidth         uint64   `xml:"bandwidth,attr"`
  BaseUrl           *BaseUrl `xml:"BaseURL"`
  Codecs            string   `xml:"codecs,attr"`
  ContentProtection []ContentProtection
  Height            uint64 `xml:"height,attr"`
  Id                string `xml:"id,attr"`
  MimeType          string `xml:"mimeType,attr"`
  SegmentBase       *SegmentBase
  SegmentTemplate   *SegmentTemplate
  Width             uint64 `xml:"width,attr"`
  // contains filtered or unexported fields
}
```

### func [Unmarshal](./dash.go#L263)

```go
func Unmarshal(text []byte, base *url.URL) ([]Representation, error)
```

### func (*Representation) [Ext](./dash.go#L335)

```go
func (r *Representation) Ext() (string, bool)
```

### func (*Representation) [GetAdaptationSet](./dash.go#L320)

```go
func (r *Representation) GetAdaptationSet() *AdaptationSet
```

### func (*Representation) [GetBaseUrl](./dash.go#L198)

```go
func (r *Representation) GetBaseUrl() (*BaseUrl, bool)
```

### func (*Representation) [GetMimeType](./dash.go#L125)

```go
func (r *Representation) GetMimeType() string
```

### func (*Representation) [Initialization](./dash.go#L164)

```go
func (r *Representation) Initialization() (string, bool)
```

### func (*Representation) [Media](./dash.go#L132)

```go
func (r *Representation) Media() []string
```

### func (*Representation) [String](./dash.go#L222)

```go
func (r *Representation) String() string
```

### func (*Representation) [Widevine](./dash.go#L187)

```go
func (r *Representation) Widevine() (Pssh, bool)
```

### type [SegmentBase](./dash.go#L347)

```go
type SegmentBase struct {
  Initialization struct {
    Range Range `xml:"range,attr"`
  }
  IndexRange Range `xml:"indexRange,attr"`
}
```

### type [SegmentTemplate](./dash.go#L406)

```go
type SegmentTemplate struct {
  Duration               uint64 `xml:"duration,attr"`
  Initialization         string `xml:"initialization,attr"`
  Media                  string `xml:"media,attr"`
  PresentationTimeOffset uint   `xml:"presentationTimeOffset,attr"`
  SegmentTimeline        *struct {
    S []struct {
      D uint `xml:"d,attr"` // duration
      R uint `xml:"r,attr"` // repeat
    }
  }
  Timescale   uint64 `xml:"timescale,attr"`
  StartNumber *uint  `xml:"startNumber,attr"`
}
```

## Source files

[dash.go](./dash.go)
