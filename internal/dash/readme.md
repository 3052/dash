# dash

Go language, I need a package for DASH (MPD) files

1. standard library only
2. use a separate file for each type
3. package is named "dash"
4. package will include a parse method, byte slice input
5. BaseURL is a single element not a slice
6. include test Go file in same package
7. test will read all ".mpd" files in the "testdata" folder. user will provide
8. support these elements and attributes

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
