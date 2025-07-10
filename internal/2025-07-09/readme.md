# dash

## MPD.BaseURL

~~~go
strings.HasPrefix(MPD.BaseURL, "http")
!strings.HasPrefix(MPD.BaseURL, "http")
~~~

## MPD.Period

~~~go
len(MPD.Period) == 1
len(MPD.Period) >= 2
~~~

## Period.duration

~~~go
Period.duration != nil
Period.duration == nil
~~~

## Representation

~~~go
Representation.SegmentBase != nil
Representation.SegmentList != nil
Representation.SegmentTemplate != nil
~~~

## SegmentTemplate.SegmentTimeline

~~~go
SegmentTemplate.SegmentTimeline != nil
SegmentTemplate.SegmentTimeline == nil
~~~

## SegmentTemplate.duration

~~~go
SegmentTemplate.duration == 0
SegmentTemplate.duration >= 1
~~~

## SegmentTemplate.endNumber

~~~go
SegmentTemplate.endNumber == 0
SegmentTemplate.endNumber >= 1
~~~

## SegmentTemplate.media

~~~
strings.Contains(SegmentTemplate.media, "$Time$")
strings.Contains(SegmentTemplate.media, "$Number$")
strings.Contains(SegmentTemplate.media, "$Number%0")
~~~

## SegmentTemplate.startNumber

~~~go
SegmentTemplate.startNumber != nil
SegmentTemplate.startNumber == nil
~~~

## SegmentTemplate.timescale

~~~go
SegmentTemplate.timescale != nil
SegmentTemplate.timescale == nil
~~~
