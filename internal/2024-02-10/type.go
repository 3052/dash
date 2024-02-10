package hls

// rfc-editor.org/rfc/rfc8216#section-4.3.4.1
type MediaPlaylist struct {
   Type string
   URI string
}

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4.2
type VariantStream struct {
   Bandwidth string
   Resolution string
   URI string
}

// datatracker.ietf.org/doc/html/rfc8216#section-4.3.4
type MasterPlaylist struct {
   Media []MediaPlaylist
   Stream []VariantStream
}

type Segment struct {
   Key string
   Map string
   RawIv string
   URI []string
}
