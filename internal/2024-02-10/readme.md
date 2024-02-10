# HLS

HTTP Live Streaming
<https://wikipedia.org/wiki/HTTP_Live_Streaming>

HTTP Live Streaming (HLS) authoring specification
<https://developer.apple.com/documentation/http_live_streaming/http_live_streaming_hls_authoring_specification_for_apple_devices/hls_authoring_specification_for_apple_devices_appendixes>

## CBC

Why does this:

~~~
#EXT-X-KEY:METHOD=AES-128,URI="https://cbsios-vh.akamaihd.net/i/temp_hd_galle...
~~~

mean CBC?

> An encryption method of AES-128 signals that Media Segments are completely
> encrypted using the Advanced Encryption Standard (AES) [`AES_128`] with a
> 128-bit key, Cipher Block Chaining (CBC)

HTTP Live Streaming
https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4

## EXT-X-KEY

If IV is missing, then use KEY for both.

## Padding

Padding (cryptography):
<https://wikipedia.org/wiki/Padding_(cryptography)#PKCS#5_and_PKCS#7>

Cryptographic Message Syntax (CMS):
https://datatracker.ietf.org/doc/html/rfc5652

> Public-Key Cryptography Standards #7 (PKCS7) padding [RFC5652]

HTTP Live Streaming:
https://datatracker.ietf.org/doc/html/rfc8216#section-4.3.2.4

## prior art

move away from RegExp:
https://github.com/grafov/m3u8/issues/206

publish module:
https://github.com/orestonce/m3u8d/issues/28

bufio.Scanner: token too long:
https://github.com/eswarantg/m3u8reader/issues/11

move away from RegExp:
https://github.com/jamesnetherton/m3u/issues/8

publish module:
https://github.com/antoinecaputo/m3u/issues/5

should not call shell:
<https://github.com/ByteTu/Api-N_m3u8DL/issues/4>

add license:
https://github.com/mattetti/m3u8Grabber/issues/3

add go.mod:
https://github.com/ushis/m3u/issues/3

m3uparser.M3uParser add byte input:
https://github.com/pawanpaudel93/go-m3u-parser/issues/3
