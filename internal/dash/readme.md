# dash

Go language, I need a package for DASH (MPD) files

1. standard library only
2. use a separate file for each type
3. only send new or updated files
4. package is named "dash"
5. package will include a parse method, byte slice input
6. BaseURL is a single element not a slice
7. include test Go file in same package
8. test will read all ".mpd" files in the "testdata" folder. user will provide
9. support these elements and attributes
   - MPD
      - @mediaPresentationDuration
      - BaseURL
      - Period
         - @duration
         - @id
         - BaseURL
         - AdaptationSet
            - @codecs
            - @height
            - @lang
            - @mimeType
            - @width
            - ContentProtection
               - @schemeIdUri
               - cenc:pssh
            - Role
               - @value
            - SegmentTemplate
               - @duration
               - @endNumber
               - @initialization
               - @media
               - @presentationTimeOffset
               - @startNumber
               - @timescale
               - SegmentTimeline
                  - S
                     - @d
                     - @r
            - Representation
               - @bandwidth
               - @codecs
               - @height
               - @id
               - @mimeType
               - @width
               - BaseURL
               - SegmentTemplate
               - ContentProtection
               - SegmentBase
                  - @indexRange
                  - Initialization
                     - @range
               - SegmentList
                  - Initialization
                     - @sourceURL
                  - SegmentURL
                     - @media
10. add navigation
   1. from AdaptationSet to Period
   2. from Initialization to SegmentList
   3. from Period to MPD
   4. from Representation to AdaptationSet
   5. from SegmentList to Representation
   6. from SegmentTemplate to AdaptationSet
   7. from SegmentTemplate to Representation
   8. from SegmentURL to SegmentList
