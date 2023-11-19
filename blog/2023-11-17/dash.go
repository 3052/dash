package dash

type AdaptationSet struct {
   ContentProtection []ContentProtection
   Representation []Representation
}

type ContentProtection struct {
   PSSH string `xml:"pssh"`
}

type Representation struct {
   AdaptationSet *AdaptationSet
   ContentProtection []ContentProtection
}

// github.com/grpc/grpc-go/blob/master/examples/helloworld/helloworld/helloworld.pb.go
func (r Representation) GetContentProtection() []ContentProtection {
   return r.s.ContentProtection
}
