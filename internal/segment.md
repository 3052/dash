# segment

~~~
Aud *CENC a0 | 128 Kbps | ec-3 | en-US | 2CH | 2 Segments | ~01h32m55s
Aud *CENC a1 | 63 Kbps | mp4a.40.5 | en-US | 2CH | 2 Segments | ~01h32m55s
Aud *CENC a2 | 258 Kbps | ec-3 | en-US | 6CH | 1 Segment | ~01h32m51s
Aud *CENC a3 | 66 Kbps | mp4a.40.5 | en-US | 2CH | 1 Segment | ~01h32m51s
~~~

whats the deal with the "2 Segments" options?

## Period id="0" duration="PT4.004S" start="PT0S"

~~~
curl -o 0.mp4 https://fly.prd.media.h264.io/0cca0350-b0d9-47bb-98d5-d5d81a73fee5/a/1_0dc23f/a1.mp4
~~~

result:

~~~
[ftyp] size=36
 - majorBrand: mp41
 - minorVersion: 0
 - compatibleBrand: iso8
 - compatibleBrand: isom
 - compatibleBrand: mp41
 - compatibleBrand: dash
 - compatibleBrand: cmfc
[moov] size=571
  [mvhd] size=108 version=0 flags=000000
   - timeScale: 24000
   - duration: 0
   - creation time: 2023-11-14T19:14:10Z
   - modification time: 2023-11-14T19:14:10Z
  [trak] size=415
    [tkhd] size=92 version=0 flags=000007
     - trackID: 1
     - duration: 0
     - creation time: 2023-11-14T19:14:10Z
     - modification time: 2023-11-14T19:14:10Z
    [mdia] size=315
      [mdhd] size=32 version=0 flags=000000
       - timeScale: 24000
       - creation time: 2023-11-14T19:14:10Z
       - modification time: 2023-11-14T19:14:10Z
       - language: eng
      [hdlr] size=45 version=0 flags=000000
       - handlerType: soun
       - handlerName: "SoundHandler"
      [minf] size=230
        [smhd] size=16 version=0 flags=000000
        [dinf] size=36
          [dref] size=28 version=0 flags=000000
            [url ] size=12
             - location: ""
        [stbl] size=170
          [stsd] size=94 version=0 flags=000000
            [mp4a] size=78
              [esds] size=42 version=0 flags=000000
                Descriptor "tag=3 ES" size=2+28
                  Descriptor "tag=4 DecoderConfig" size=2+20
                   - BufferSizeDB: 0
                   - MaxBitrate: 67363
                   - AvgBitrate: 62688
                    Descriptor "tag=5 DecoderSpecificInfo" size=2+5
                     - DecConfig (5B): 131056e598
                  Descriptor "tag=6 SLConfig" size=2+1
          [stts] size=16 version=0 flags=000000
          [stsc] size=16 version=0 flags=000000
          [stsz] size=20 version=0 flags=000000
          [stco] size=16 version=0 flags=000000
  [mvex] size=40
    [trex] size=32 version=0 flags=000000
     - trackID: 1
     - defaultSampleDescriptionIndex: 1
     - defaultSampleDuration: 1024
     - defaultSampleSize: 0
     - defaultSampleFlags: 00000000 (isLeading=0 dependsOn=0 isDependedOn=0 hasRedundancy=0 padding=0 isNonSync=false degradationPriority=0)
[sidx] size=56 version=0 flags=000000
 - referenceID: 1
 - timeScale: 24000
 - earliestPresentationTime: 0
 - firstOffset: 0
[moof] size=472
  [mfhd] size=16 version=0 flags=000000
   - sequenceNumber: 1
  [traf] size=448
    [tfhd] size=28 version=0 flags=02002a
     - trackID: 1
     - defaultBaseIsMoof: true
     - sampleDescriptionIndex: 1
     - defaultSampleDuration: 1024
     - defaultSampleFlags: 00000000 (isLeading=0 dependsOn=0 isDependedOn=0 hasRedundancy=0 padding=0 isNonSync=false degradationPriority=0)
    [tfdt] size=16 version=0 flags=000000
     - baseMediaDecodeTime: 0
    [trun] size=396 version=0 flags=000201
     - sampleCount: 94
[mdat] size=31436
[moof] size=108
  [mfhd] size=16 version=0 flags=000000
   - sequenceNumber: 2
  [traf] size=84
    [tfhd] size=28 version=0 flags=02002a
     - trackID: 1
     - defaultBaseIsMoof: true
     - sampleDescriptionIndex: 1
     - defaultSampleDuration: 1024
     - defaultSampleFlags: 00000000 (isLeading=0 dependsOn=0 isDependedOn=0 hasRedundancy=0 padding=0 isNonSync=false degradationPriority=0)
    [tfdt] size=16 version=0 flags=000000
     - baseMediaDecodeTime: 96256
    [trun] size=32 version=0 flags=000201
     - sampleCount: 3
[mdat] size=1011
~~~

## Period id="1" duration="PT17M0.853166666S" start="PT4.004S"

~~~
curl -o 1.mp4 https://fly.prd.media.h264.io/59da086b-1d1e-48fa-b318-782408318b54/a/0_fa3c08/a1.mp4
~~~
