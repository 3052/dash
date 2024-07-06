# internal

~~~
PS D:\Desktop\etc> paramount.exe -b esJvFlqdrcS_kFHnpxSuYp449E7tTexD -w
15:19:54 INFO GET URL="https://link.theplatform.com/s/dJ5BDC/media/guid/2198311517/esJvFlqdrcS_kFHnpxSuYp449E7tTexD?assetTypes=DASH_CENC_PRECON&formats=MPEG-DASH"
15:19:55 INFO GET URL="https://www.paramountplus.com/apps-api/v2.0/androidphone/video/cid/esJvFlqdrcS_kFHnpxSuYp449E7tTexD.json?at=ABAAAAAAAAAAAAAAAAAAAAAA%2Bc9AiUS4F1JOX0w0O1uzQw%2F%2BqdAuVgybB1FK7aqonjY%3D"

PS D:\Desktop\etc> paramount.exe -b esJvFlqdrcS_kFHnpxSuYp449E7tTexD -i 5
15:19:56 INFO GET URL=https://vod-gcs-cedexis.cbsaavideo.com/intl_vms/2024/03/26/2323316803935/2709030_cenc_precon_dash/stream.mpd
15:19:57 INFO GET URL="https://www.paramountplus.com/apps-api/v3.1/androidphone/irdeto-control/anonymous-session-token.json?at=ABAAAAAAAAAAAAAAAAAAAAAA%2Bc9AiUS4F1JOX0w0O1uzQw%2F%2BqdAuVgybB1FK7aqonjY%3D&contentId=esJvFlqdrcS_kFHnpxSuYp449E7tTexD"
15:19:57 INFO GET URL=https://vod-gcs-cedexis.cbsaavideo.com/intl_vms/2024/03/26/2323316803935/2709030_cenc_precon_dash/CriminalMinds_1701_HQ_R1_ba27fc94-9d27-4854-8289-442ba83be9c3_R1_1088_Corrected_2709014_4500/init.m4v
15:19:57 INFO POST URL="https://cbsi.live.ott.irdeto.com/widevine/getlicense?CrmId=cbsi&AccountId=cbsi&SubContentType=Default&contentId=esJvFlqdrcS_kFHnpxSuYp449E7tTexD"
15:19:57 INFO CDM id=1fde0154d72a4f45912b34f0ce0777eb key=9e3e3d2cb89469998730adf788b98a1d
panic: HTTP/1.1 404 Not Found
Content-Length: 340
Accept-Ranges: bytes
Access-Control-Allow-Origin: *
Access-Control-Expose-Headers: X-CDN
Age: 0
Cache-Control: max-age=0, no-cache
Connection: keep-alive
Content-Type: video/mp4
Date: Sat, 06 Jul 2024 20:19:55 GMT
X-Cache: HIT, MISS
X-Cache-Hits: 1, 0
X-Cdn: Fastly
X-Failover-Status: 404
X-Restarts: 1
X-Served-By: cache-iad-kcgs7200116-IAD, cache-bur-kbur8200072-BUR
X-Shield: Error 404
X-Timer: S1720297196.841623,VS0,VE78

<?xml version='1.0' encoding='UTF-8'?><Error><Code>NoSuchKey</Code><Message>The specified key does not exist.</Message><Details>No such object: allaccess/intl_vms/2024/03/26/2323316803935/2709030_cenc_precon_dash/CriminalMinds_1701_HQ_R1_ba27fc94-9d27-4854-8289-442ba83be9c3_R1_1088_Corrected_2709014_4500/seg_73776703.m4s</Details></Error>
~~~

