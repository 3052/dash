# DASH

> Some things are hidden for a reason, and if you see them, you'll be changed
> forever, but I **wanted** to be changed forever.
>
> [Miranda July](//youtube.com/watch?v=7dMGWporaFE&t=142s)

- <https://f002.backblazeb2.com/file/minerals/ISO_IEC_23009-1_2022(en).pdf>
- <https://wikipedia.org/wiki/Dynamic_Adaptive_Streaming_over_HTTP>
- https://dashif.org/Guidelines-TimingModel/Timing-Model.pdf
- https://dashif.org/docs/DASH-IF-IOP-v4.3.pdf

Go language, I need a package for DASH (MPD) files

1. standard library only
2. use a separate file for each type
3. only send new or updated files
4. package is named "dash"
5. package will include a parse method, byte slice input
6. BaseURL is a single element not a slice
7. support these elements and attributes
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
8. add navigation
   1. from AdaptationSet to Period
   2. from Initialization to SegmentList
   3. from Period to MPD
   4. from Representation to AdaptationSet
   5. from SegmentList to Representation
   6. from SegmentTemplate to AdaptationSet
   7. from SegmentTemplate to Representation
   8. from SegmentURL to SegmentList
9. resolve BaseURL using
   1. MPD URL
   2. all parent BaseURL
10. resolve Initialization@sourceURL using
   1. MPD URL
   2. all parent BaseURL
11. resolve SegmentTemplate@initialization using
   1. MPD URL
   2. all parent BaseURL
12. resolve SegmentTemplate@media using
   1. MPD URL
   2. all parent BaseURL
13. resolve SegmentURL@media using
   1. MPD URL
   2. all parent BaseURL
14. resolve function should return `*url.URL`
15. add method to get all Representation, group by id
16. add method to get codecs
17. add method to get height
18. add method to get width
19. add method to get mimeType
20. AdaptationSet.Role is single element not slice
21. add method to get SegmentTemplate
22. add method to replace SegmentTemplate@initialization
   - `$RepresentationID$`
23. add method to replace SegmentTemplate@media
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
24. add method to replace SegmentTemplate@media
   - `$RepresentationID$`
   - `$Time$`
25. SegmentTemplate@startNumber is 1 if missing
26. SegmentTemplate@timescale is 1 if missing
27. Period@duration is MPD@mediaPresentationDuration if missing
28. add method to get `Time` values from SegmentTimeline
29. add method to get `Number` values from SegmentTimeline
30. add method to get `Number` values from
   SegmentTemplate@startNumber to SegmentTemplate@endNumber
31. add method to get `Number` values from
   Ceil(
      AsSeconds(Period@duration) /
      (SegmentTemplate@duration / SegmentTemplate@timescale)
   )
32. add a method that returns the SegmentTemplate URLs
33. use SegmentTemplate@presentationTimeOffset as initial `$Time$`
34. add method to get `time.Duration`
    ```
    time.ParseDuration(strings.ToLower(
       strings.TrimPrefix(Period@duration, "PT"),
    ))
    ```
35. add method to get unique ContentProtection from
   - AdaptationSet
   - Representation
36. include test file in same package
   - test will read all ".mpd" files in the "testdata" folder
   - user will provide test files
   - for each file, get the slice of replaced `SegmentTemplate@media` URLs
   - only get one slice of URLs per mimeType
   - print slice length, and first and last URLs
37. add Representation.String method. each value on its own line
   1. Representation@bandwidth
   2. Representation.GetWidth
   3. Representation.GetHeight
   4. Representation.GetCodecs
   5. Representation.GetMimeType
   6. AdaptationSet@lang
   7. Role@value
   8. Period@id
   9. Representation@id

## contact

<dl>
   <dt>email</dt>
   <dd>27@riseup.net</dd>
   <dt>Discord username</dt>
   <dd>10308</dd>
</dl>

