# dash

- <https://f002.backblazeb2.com/file/minerals/ISO_IEC_23009-1_2022(en).pdf>
- https://dashif.org/Guidelines-TimingModel/Timing-Model.pdf
- https://dashif.org/docs/DASH-IF-IOP-v4.3.pdf

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
11. do not skip tests
12. resolve BaseURL using
   1. MPD URL
   2. all parent BaseURL
13. resolve Initialization@sourceURL using
   1. MPD URL
   2. all parent BaseURL
14. resolve SegmentTemplate@initialization using
   1. MPD URL
   2. all parent BaseURL
15. resolve SegmentTemplate@media using
   1. MPD URL
   2. all parent BaseURL
16. resolve SegmentURL@media using
   1. MPD URL
   2. all parent BaseURL
17. resolve function should return `*url.URL`
18. add method to get all Representation, group by id
19. add method to get codecs
20. add method to get height
21. add method to get width
22. add method to get mimeType
23. AdaptationSet.Role is single element not slice
24. add Representation.String method. each value on its own line
   - AdaptationSet@lang
   - Period@id
   - Representation.GetCodecs
   - Representation.GetHeight
   - Representation.GetMimeType
   - Representation.GetWidth
   - Representation@bandwidth
   - Role@value
25. add method to get ContentProtection
26. add method to get SegmentTemplate
27. add method to replace SegmentTemplate@initialization
   - `$RepresentationID$`
28. add method to replace SegmentTemplate@media
   - `$Number$`
   - `$Number%02d$`
   - `$Number%03d$`
   - `$Number%04d$`
   - `$Number%05d$`
   - `$Number%06d$`
   - `$Number%07d$`
   - `$Number%08d$`
   - `$Number%09d$`
   - `$RepresentationID$`
29. add method to replace SegmentTemplate@media
   - `$RepresentationID$`
   - `$Time$`
30. SegmentTemplate@startNumber is 1 if missing
31. SegmentTemplate@timescale is 1 if missing
32. Period@duration is MPD@mediaPresentationDuration if missing
33. add method to get `Time` values from SegmentTimeline
34. add method to get `Number` values from SegmentTimeline
35. add method to get `Number` values from
   SegmentTemplate@startNumber to SegmentTemplate@endNumber
36. add method to get `Number` values from
   Ceil(
      AsSeconds(Period@duration) /
      (SegmentTemplate@duration / SegmentTemplate@timescale)
   )
37. add a method that returns the SegmentTemplate URLs
38. use SegmentTemplate@presentationTimeOffset as inital `$Time$`
39. add method to get `time.Duration`
   time.ParseDuration(strings.ToLower(
      strings.TrimPrefix(Period@duration, "PT"),
   ))
40. add function to get start and end from range or indexRange

