package main

import (
   "154.pages.dev/sofia"
   "fmt"
   "os"
   "strconv"
)

var tos = []int64{0, 2000, 4000, 6000, 8000, 10000, 12000, 14000, 16000, 18000,
20000, 22000, 24000, 26000, 28000, 30000, 32000, 34000, 36000, 38000, 40000,
42000, 44000, 46000, 48000, 50000, 52000, 54000, 56000, 58000, 60000, 62000,
64000, 66000, 68000, 70000, 72000, 74000, 76000, 78000, 80000, 82000, 84000,
86000, 88000, 90000, 92000, 94000, 96000, 98000, 100000, 102000, 104000,
106000, 108000, 110000, 112000, 114000, 116000, 118000, 120000, 122000, 124000,
126000, 128000, 130000, 132000, 134000, 136000, 138000, 140000, 142000, 144000,
146000, 148000, 150000, 152000, 154000, 156000, 158000, 160000, 162000, 164000,
166000, 168000, 170000, 172000, 174000, 176000, 178000, 180000, 182000, 184000,
186000, 188000, 190000, 192000, 194000, 196000, 198000, 200000, 202000, 204000,
206000, 208000, 210000, 212000, 214000, 216000, 218000, 220000, 222000, 224000,
226000, 228000, 230000, 232000, 234000, 236000, 238000, 240000, 242000, 244000,
246000, 248000, 250000, 252000, 254000, 256000, 258000, 260000, 262000, 264000,
266000, 268000, 270000, 272000, 274000, 276000, 278000, 280000, 282000, 284000,
286000, 288000, 290000, 292000, 294000, 296000, 298000, 300000, 302000, 304000,
306000, 308000, 310000, 312000, 314000, 316000, 318000, 320000, 322000, 324000,
326000, 328000, 330000, 332000, 334000, 336000, 338000, 340000, 342000, 344000,
346000, 348000, 350000, 352000, 354000, 356000, 358000, 360000, 362000, 364000,
366000, 368000, 370000, 372000, 374000, 376000, 378000, 380000, 382000, 384000,
386000, 388000, 390000, 392000, 394000, 396000, 398000, 400000, 402000, 404000,
406000, 408000, 410000, 412000, 414000, 416000, 418000, 420000, 422000, 424000,
426000, 428000, 430000, 432000, 434000, 436000, 438000, 440000, 442000, 444000,
446000, 448000, 450000, 452000, 454000, 456000, 458000, 460000, 462000, 464000,
466000, 468000, 470000, 472000, 474000, 476000, 478000, 480000, 482000, 484000,
486000, 488000, 490000, 492000, 494000, 496000, 498000, 500000, 502000, 504000,
506000, 508000, 510000, 512000, 514000, 516000, 518000, 520000, 522000, 524000,
526000, 528000, 530000, 532000, 534000, 536000, 538000, 540000, 542000, 544000,
546000, 548000, 550000, 552000, 554000, 556000, 558000, 560000, 562000, 564000,
566000, 568000, 570000, 572000, 574000, 576000, 578000, 580000, 582000, 584000,
586000, 588000, 590000, 592000, 594000, 596000, 598000, 600000, 602000, 604000,
606000, 608000, 610000, 612000, 614000, 616000, 618000, 620000, 622000, 624000,
626000, 628000, 630000, 632000, 634000, 636000, 638000, 640000, 642000, 644000,
646000, 648000, 650000, 652000, 654000, 656000, 658000, 660000, 662000, 664000,
666000, 668000, 670000, 672000, 674000, 676000, 678000, 680000, 682000, 684000,
686000, 688000, 690000, 692000, 694000, 696000, 698000, 700000, 702000, 704000,
706000, 708000, 710000, 712000, 714000, 716000, 718000, 720000, 722000, 724000,
726000, 728000, 730000, 732000, 734000, 736000, 738000, 740000, 742000, 744000,
746000, 748000, 750000, 752000, 754000, 756000, 758000, 760000, 762000, 764000,
766000, 768000, 770000, 772000, 774000, 776000, 778000, 780000, 782000, 784000,
786000, 788000, 790000, 792000, 794000, 796000, 798000, 800000, 802000, 804000,
806000, 808000, 810000, 812000, 814000, 816000, 818000, 820000, 822000, 824000,
826000, 828000, 830000, 832000, 834000, 836000, 838000, 840000, 842000, 844000,
846000, 848000, 850000, 852000, 854000, 856000, 858000, 860000, 862000, 864000,
866000, 868000, 870000, 872000, 874000, 876000, 878000, 880000, 882000, 884000,
886000, 888000, 890000, 892000, 894000, 896000, 898000, 900000, 902000, 904000,
906000, 908000, 910000, 912000, 914000, 916000, 918000, 920000, 922000, 924000,
926000, 928000, 930000, 932000, 934000, 936000, 938000, 940000, 942000, 944000,
946000, 948000, 950000, 952000, 954000, 956000, 958000, 960000, 962000, 964000,
966000, 968000, 970000, 972000, 974000, 976000, 978000, 980000, 982000, 984000,
986000, 988000, 990000, 992000, 994000, 996000, 998000, 1000000, 1002000,
1004000, 1006000, 1008000, 1010000, 1012000, 1014000, 1016000, 1018000,
1020000, 1022000, 1024000, 1026000, 1028000, 1030000, 1032000, 1034000,
1036000, 1038000, 1040000, 1042000, 1044000, 1046000, 1048000, 1050000,
1052000, 1054000, 1056000, 1058000, 1060000, 1062000, 1064000, 1066000,
1068000, 1070000, 1072000, 1074000, 1076000, 1078000, 1080000, 1082000,
1084000, 1086000, 1088000, 1090000, 1092000, 1094000, 1096000, 1098000,
1100000, 1102000, 1104000, 1106000, 1108000, 1110000, 1112000, 1114000,
1116000, 1118000, 1120000, 1122000, 1124000, 1126000, 1128000, 1130000,
1132000, 1134000, 1136000, 1138000, 1140000, 1142000, 1144000, 1146000,
1148000, 1150000, 1152000, 1154000, 1156000, 1158000, 1160000, 1162000,
1164000, 1166000, 1168000, 1170000, 1172000, 1174000, 1176000, 1178000,
1180000, 1182000, 1184000, 1186000, 1188000, 1190000, 1192000, 1194000,
1196000, 1198000, 1200000, 1202000, 1204000, 1206000, 1208000, 1210000,
1212000, 1214000, 1216000, 1218000, 1220000, 1222000, 1224000, 1226000,
1228000, 1230000, 1232000, 1234000, 1236000, 1238000, 1240000, 1242000,
1244000, 1246000, 1248000, 1250000, 1252000, 1254000, 1256000, 1258000,
1260000, 1262000, 1264000, 1266000, 1268000, 1270000, 1272000, 1274000,
1276000, 1278000, 1280000, 1282000, 1284000, 1286000, 1288000, 1290000,
1292000, 1294000, 1296000, 1298000, 1300000, 1302000, 1304000, 1306000,
1308000, 1310000, 1312000, 1314000, 1316000, 1318000, 1320000, 1322000,
1324000, 1326000, 1328000, 1330000, 1332000, 1334000, 1336000, 1338000,
1340000, 1342000, 1344000, 1346000, 1348000, 1350000, 1352000, 1354000,
1356000, 1358000, 1360000, 1362000, 1364000, 1366000, 1368000, 1370000,
1372000, 1374000, 1376000, 1378000, 1380000, 1382000, 1384000, 1386000,
1388000, 1390000, 1392000, 1394000, 1396000, 1398000, 1400000, 1402000,
1404000, 1406000, 1408000, 1410000, 1412000, 1414000, 1416000, 1418000,
1420000, 1422000, 1424000, 1426000, 1428000, 1430000, 1432000, 1434000,
1436000, 1438000, 1440000, 1442000, 1444000, 1446000, 1448000, 1450000,
1452000, 1454000, 1456000, 1458000, 1460000, 1462000, 1464000, 1466000,
1468000, 1470000, 1472000, 1474000, 1476000, 1478000, 1480000, 1482000,
1484000, 1486000, 1488000, 1490000, 1492000, 1494000, 1496000, 1498000,
1500000, 1502000, 1504000, 1506000, 1508000, 1510000, 1512000, 1514000,
1516000, 1518000, 1520000, 1522000, 1524000, 1526000, 1528000, 1530000,
1532000, 1534000, 1536000, 1538000, 1540000, 1542000, 1544000, 1546000,
1548000, 1550000, 1552000, 1554000, 1556000, 1558000, 1560000, 1562000,
1564000, 1566000, 1568000, 1570000, 1572000, 1574000, 1576000, 1578000,
1580000, 1582000, 1584000, 1586000, 1588000, 1590000, 1592000, 1594000,
1596000, 1598000, 1600000, 1602000, 1604000, 1606000, 1608000, 1610000,
1612000, 1614000, 1616000, 1618000, 1620000, 1622000, 1624000, 1626000,
1628000, 1630000, 1632000, 1634000, 1636000, 1638000, 1640000, 1642000,
1644000, 1646000, 1648000, 1650000, 1652000, 1654000, 1656000, 1658000,
1660000, 1662000, 1664000, 1666000, 1668000, 1670000, 1672000, 1674000,
1676000, 1678000, 1680000, 1682000, 1684000, 1686000, 1688000, 1690000,
1692000, 1694000, 1696000, 1698000, 1700000, 1702000, 1704000, 1706000,
1708000, 1710000, 1712000, 1714000, 1716000, 1718000, 1720000, 1722000,
1724000, 1726000, 1728000, 1730000, 1732000, 1734000, 1736000, 1738000,
1740000, 1742000, 1744000, 1746000, 1748000, 1750000, 1752000, 1754000,
1756000, 1758000, 1760000, 1762000, 1764000, 1766000, 1768000, 1770000,
1772000, 1774000, 1776000, 1778000, 1780000, 1782000, 1784000, 1786000,
1788000, 1790000, 1792000, 1794000, 1796000, 1798000, 1800000, 1802000,
1804000, 1806000, 1808000, 1810000, 1812000, 1814000, 1816000, 1818000,
1820000, 1822000, 1824000, 1826000, 1828000, 1830000, 1832000, 1834000,
1836000, 1838000, 1840000, 1842000, 1844000, 1846000, 1848000, 1850000,
1852000, 1854000, 1856000, 1858000, 1860000, 1862000, 1864000, 1866000,
1868000, 1870000, 1872000, 1874000, 1876000, 1878000, 1880000, 1882000,
1884000, 1886000, 1888000, 1890000, 1892000, 1894000, 1896000, 1898000,
1900000, 1902000, 1904000, 1906000, 1908000, 1910000, 1912000, 1914000,
1916000, 1918000, 1920000, 1922000, 1924000, 1926000, 1928000, 1930000,
1932000, 1934000, 1936000, 1938000, 1940000, 1942000, 1944000, 1946000,
1948000, 1950000, 1952000, 1954000, 1956000, 1958000, 1960000, 1962000,
1964000, 1966000, 1968000, 1970000, 1972000, 1974000, 1976000, 1978000,
1980000, 1982000, 1984000, 1986000, 1988000, 1990000, 1992000, 1994000,
1996000, 1998000, 2000000, 2002000, 2004000, 2006000, 2008000, 2010000,
2012000, 2014000, 2016000, 2018000, 2020000, 2022000, 2024000, 2026000,
2028000, 2030000, 2032000, 2034000, 2036000, 2038000, 2040000, 2042000,
2044000, 2046000, 2048000, 2050000, 2052000, 2054000, 2056000, 2058000,
2060000, 2062000, 2064000, 2066000, 2068000, 2070000, 2072000, 2074000,
2076000, 2078000, 2080000, 2082000, 2084000, 2086000, 2088000, 2090000,
2092000, 2094000, 2096000, 2098000, 2100000, 2102000, 2104000, 2106000,
2108000, 2110000, 2112000, 2114000, 2116000, 2118000, 2120000, 2122000,
2124000, 2126000, 2128000, 2130000, 2132000, 2134000, 2136000, 2138000,
2140000, 2142000, 2144000, 2146000, 2148000, 2150000, 2152000, 2154000,
2156000, 2158000, 2160000, 2162000, 2164000, 2166000, 2168000, 2170000,
2172000, 2174000, 2176000, 2178000, 2180000, 2182000, 2184000, 2186000,
2188000, 2190000, 2192000, 2194000, 2196000, 2198000, 2200000, 2202000,
2204000, 2206000, 2208000, 2210000, 2212000, 2214000, 2216000, 2218000,
2220000, 2222000, 2224000, 2226000, 2228000, 2230000, 2232000, 2234000,
2236000, 2238000, 2240000, 2242000, 2244000, 2246000, 2248000, 2250000,
2252000, 2254000, 2256000, 2258000, 2260000, 2262000, 2264000, 2266000,
2268000, 2270000, 2272000, 2274000, 2276000, 2278000, 2280000, 2282000,
2284000, 2286000, 2288000, 2290000, 2292000, 2294000, 2296000, 2298000,
2300000, 2302000, 2304000, 2306000, 2308000, 2310000, 2312000, 2314000,
2316000, 2318000, 2320000, 2322000, 2324000, 2326000, 2328000, 2330000,
2332000, 2334000, 2336000, 2338000, 2340000, 2342000, 2344000, 2346000,
2348000, 2350000, 2352000, 2354000, 2356000, 2358000, 2360000, 2362000,
2364000, 2366000, 2368000, 2370000, 2372000, 2374000, 2376000, 2378000,
2380000, 2382000, 2384000, 2386000, 2388000, 2390000, 2392000, 2394000,
2396000, 2398000, 2400000, 2402000, 2404000, 2406000, 2408000, 2410000,
2412000, 2414000, 2416000, 2418000, 2420000, 2422000, 2424000, 2426000,
2428000, 2430000, 2432000, 2434000, 2436000, 2438000, 2440000, 2442000,
2444000, 2446000, 2448000, 2450000, 2452000, 2454000, 2456000, 2458000,
2460000, 2462000, 2464000, 2466000, 2468000, 2470000, 2472000, 2474000,
2476000, 2478000, 2480000, 2482000, 2484000, 2486000, 2488000, 2490000,
2492000, 2494000, 2496000, 2498000, 2500000, 2502000, 2504000, 2506000,
2508000, 2510000, 2512000, 2514000, 2516000, 2518000, 2520000, 2522000,
2524000, 2526000, 2528000, 2530000, 2532000, 2534000, 2536000, 2538000,
2540000, 2542000, 2544000, 2546000, 2548000, 2550000, 2552000, 2554000,
2556000, 2558000, 2560000, 2562000, 2564000, 2566000, 2568000, 2570000,
2572000, 2574000, 2576000, 2578000, 2580000, 2582000, 2584000, 2586000,
2588000, 2590000, 2592000, 2594000, 2596000, 2598000, 2600000, 2602000,
2604000, 2606000, 2608000, 2610000, 2612000, 2614000, 2616000, 2618000,
2620000, 2622000, 2624000, 2626000, 2628000, 2630000, 2632000, 2634000,
2636000, 2638000, 2640000, 2642000, 2644000, 2646000, 2648000, 2650000,
2652000, 2654000, 2656000, 2658000, 2660000, 2662000, 2664000, 2666000,
2668000, 2670000, 2672000, 2674000, 2676000, 2678000, 2680000, 2682000,
2684000, 2686000, 2688000, 2690000, 2692000, 2694000, 2696000, 2698000,
2700000, 2702000, 2704000, 2706000, 2708000, 2710000, 2712000, 2714000,
2716000, 2718000, 2720000, 2722000, 2724000, 2726000, 2728000, 2730000,
2732000, 2734000, 2736000, 2738000, 2740000, 2742000, 2744000, 2746000,
2748000, 2750000, 2752000, 2754000, 2756000, 2758000, 2760000, 2762000,
2764000, 2766000, 2768000, 2770000, 2772000, 2774000, 2776000, 2778000,
2780000, 2782000, 2784000, 2786000, 2788000, 2790000, 2792000, 2794000,
2796000, 2798000, 2800000, 2802000, 2804000, 2806000, 2808000, 2810000,
2812000, 2814000, 2816000, 2818000, 2820000, 2822000, 2824000, 2826000,
2828000, 2830000, 2832000, 2834000, 2836000, 2838000, 2840000, 2842000,
2844000, 2846000, 2848000, 2850000, 2852000, 2854000, 2856000, 2858000,
2860000, 2862000, 2864000, 2866000, 2868000, 2870000, 2872000, 2874000,
2876000, 2878000, 2880000, 2882000, 2884000, 2886000, 2888000, 2890000,
2892000, 2894000, 2896000, 2898000, 2900000, 2902000, 2904000, 2906000,
2908000, 2910000, 2912000, 2914000, 2916000, 2918000, 2920000, 2922000,
2924000, 2926000, 2928000, 2930000, 2932000, 2934000, 2936000, 2938000,
2940000, 2942000, 2944000, 2946000, 2948000, 2950000, 2952000, 2954000,
2956000, 2958000, 2960000, 2962000, 2964000, 2966000, 2968000, 2970000,
2972000, 2974000, 2976000, 2978000, 2980000, 2982000, 2984000, 2986000,
2988000, 2990000, 2992000, 2994000, 2996000, 2998000, 3000000, 3002000,
3004000, 3006000, 3008000, 3010000, 3012000, 3014000, 3016000, 3018000,
3020000, 3022000, 3024000, 3026000, 3028000, 3030000, 3032000, 3034000,
3036000, 3038000, 3040000, 3042000, 3044000, 3046000, 3048000, 3050000,
3052000, 3054000, 3056000, 3058000, 3060000, 3062000, 3064000, 3066000,
3068000, 3070000, 3072000, 3074000, 3076000, 3078000, 3080000, 3082000,
3084000, 3086000, 3088000, 3090000, 3092000, 3094000, 3096000, 3098000,
3100000, 3102000, 3104000, 3106000, 3108000, 3110000, 3112000, 3114000,
3116000, 3118000, 3120000, 3122000, 3124000, 3126000, 3128000, 3130000,
3132000, 3134000, 3136000, 3138000, 3140000, 3142000, 3144000, 3146000,
3148000, 3150000, 3152000, 3154000, 3156000, 3158000, 3160000, 3162000,
3164000, 3166000, 3168000, 3170000, 3172000, 3174000, 3176000, 3178000,
3180000, 3182000, 3184000, 3186000, 3188000, 3190000, 3192000, 3194000,
3196000, 3198000, 3200000, 3202000, 3204000, 3206000, 3208000, 3210000,
3212000, 3214000, 3216000, 3218000, 3220000, 3222000, 3224000, 3226000,
3228000, 3230000, 3232000, 3234000, 3236000, 3238000, 3240000, 3242000,
3244000, 3246000, 3248000, 3250000, 3252000, 3254000, 3256000, 3258000,
3260000, 3262000, 3264000, 3266000, 3268000, 3270000, 3272000, 3274000,
3276000, 3278000, 3280000, 3282000, 3284000, 3286000, 3288000, 3290000,
3292000, 3294000, 3296000, 3298000, 3300000, 3302000, 3304000, 3306000,
3308000, 3310000, 3312000, 3314000, 3316000, 3318000, 3320000, 3322000,
3324000, 3326000, 3328000, 3330000, 3332000, 3334000, 3336000, 3338000,
3340000, 3342000, 3344000, 3346000, 3348000, 3350000, 3352000, 3354000,
3356000, 3358000, 3360000, 3362000, 3364000, 3366000, 3368000, 3370000,
3372000, 3374000, 3376000, 3378000, 3380000, 3382000, 3384000, 3386000,
3388000, 3390000, 3392000, 3394000, 3396000, 3398000, 3400000, 3402000,
3404000, 3406000, 3408000, 3410000, 3412000, 3414000, 3416000, 3418000,
3420000, 3422000, 3424000, 3426000, 3428000, 3430000, 3432000, 3434000,
3436000, 3438000, 3440000, 3442000, 3444000, 3446000, 3448000, 3450000,
3452000, 3454000, 3456000, 3458000, 3460000, 3462000, 3464000, 3466000,
3468000, 3470000, 3472000, 3474000, 3476000, 3478000, 3480000, 3482000,
3484000, 3486000, 3488000, 3490000, 3492000, 3494000, 3496000, 3498000,
3500000, 3502000, 3504000, 3506000, 3508000, 3510000, 3512000, 3514000,
3516000, 3518000, 3520000, 3522000, 3524000, 3526000, 3528000, 3530000,
3532000, 3534000, 3536000, 3538000, 3540000, 3542000, 3544000, 3546000,
3548000, 3550000, 3552000, 3554000, 3556000, 3558000, 3560000, 3562000,
3564000, 3566000, 3568000, 3570000, 3572000, 3574000, 3576000, 3578000,
3580000, 3582000, 3584000, 3586000, 3588000, 3590000, 3592000, 3594000,
3596000, 3598000, 3600000, 3602000, 3604000, 3606000, 3608000, 3610000,
3612000, 3614000, 3616000, 3618000, 3620000, 3622000, 3624000, 3626000,
3628000, 3630000, 3632000, 3634000, 3636000, 3638000, 3640000, 3642000,
3644000, 3646000, 3648000, 3650000, 3652000, 3654000, 3656000, 3658000,
3660000, 3662000, 3664000, 3666000, 3668000, 3670000, 3672000, 3674000,
3676000, 3678000, 3680000, 3682000, 3684000, 3686000, 3688000, 3690000,
3692000, 3694000, 3696000, 3698000, 3700000, 3702000, 3704000, 3706000,
3708000, 3710000, 3712000, 3714000, 3716000, 3718000, 3720000, 3722000,
3724000, 3726000, 3728000, 3730000, 3732000, 3734000, 3736000, 3738000,
3740000, 3742000, 3744000, 3746000, 3748000, 3750000, 3752000, 3754000,
3756000, 3758000, 3760000, 3762000, 3764000, 3766000, 3768000, 3770000,
3772000, 3774000, 3776000, 3778000, 3780000, 3782000, 3784000, 3786000,
3788000, 3790000, 3792000, 3794000, 3796000, 3798000, 3800000, 3802000,
3804000, 3806000, 3808000, 3810000, 3812000, 3814000, 3816000, 3818000,
3820000, 3822000, 3824000, 3826000, 3828000, 3830000, 3832000, 3834000,
3836000, 3838000, 3840000, 3842000, 3844000, 3846000, 3848000, 3850000,
3852000, 3854000, 3856000, 3858000, 3860000, 3862000, 3864000, 3866000,
3868000, 3870000, 3872000, 3874000, 3876000, 3878000, 3880000, 3882000,
3884000, 3886000, 3888000, 3890000, 3892000, 3894000, 3896000, 3898000,
3900000, 3902000, 3904000, 3906000, 3908000, 3910000, 3912000, 3914000,
3916000, 3918000, 3920000, 3922000, 3924000, 3926000, 3928000, 3930000,
3932000, 3934000, 3936000, 3938000, 3940000, 3942000, 3944000, 3946000,
3948000, 3950000, 3952000, 3954000, 3956000, 3958000, 3960000, 3962000,
3964000, 3966000, 3968000, 3970000, 3972000, 3974000, 3976000, 3978000,
3980000, 3982000, 3984000, 3986000, 3988000, 3990000, 3992000, 3994000,
3996000, 3998000, 4000000, 4002000, 4004000, 4006000, 4008000, 4010000,
4012000, 4014000, 4016000, 4018000, 4020000, 4022000, 4024000, 4026000,
4028000, 4030000, 4032000, 4034000, 4036000, 4038000, 4040000, 4042000,
4044000, 4046000, 4048000, 4050000, 4052000, 4054000, 4056000, 4058000,
4060000, 4062000, 4064000, 4066000, 4068000, 4070000, 4072000, 4074000,
4076000, 4078000, 4080000, 4082000, 4084000, 4086000, 4088000, 4090000,
4092000, 4094000, 4096000, 4098000, 4100000, 4102000, 4104000, 4106000,
4108000, 4110000, 4112000, 4114000, 4116000, 4118000, 4120000, 4122000,
4124000, 4126000, 4128000, 4130000, 4132000, 4134000, 4136000, 4138000,
4140000, 4142000, 4144000, 4146000, 4148000, 4150000, 4152000, 4154000,
4156000, 4158000, 4160000, 4162000, 4164000, 4166000, 4168000, 4170000,
4172000, 4174000, 4176000, 4178000, 4180000, 4182000, 4184000, 4186000,
4188000, 4190000, 4192000, 4194000, 4196000, 4198000, 4200000, 4202000,
4204000, 4206000, 4208000, 4210000, 4212000, 4214000, 4216000, 4218000,
4220000, 4222000, 4224000, 4226000, 4228000, 4230000, 4232000, 4234000,
4236000, 4238000, 4240000, 4242000, 4244000, 4246000, 4248000, 4250000,
4252000, 4254000, 4256000, 4258000, 4260000, 4262000, 4264000, 4266000,
4268000, 4270000, 4272000, 4274000, 4276000, 4278000, 4280000, 4282000,
4284000, 4286000, 4288000, 4290000, 4292000, 4294000, 4296000, 4298000,
4300000, 4302000, 4304000, 4306000, 4308000, 4310000, 4312000, 4314000,
4316000, 4318000, 4320000, 4322000, 4324000, 4326000, 4328000, 4330000,
4332000, 4334000, 4336000, 4338000, 4340000, 4342000, 4344000, 4346000,
4348000, 4350000, 4352000, 4354000, 4356000, 4358000, 4360000, 4362000,
4364000, 4366000, 4368000, 4370000, 4372000, 4374000, 4376000, 4378000,
4380000, 4382000, 4384000, 4386000, 4388000, 4390000, 4392000, 4394000,
4396000, 4398000, 4400000, 4402000, 4404000, 4406000, 4408000, 4410000,
4412000, 4414000, 4416000, 4418000, 4420000, 4422000, 4424000, 4426000,
4428000, 4430000, 4432000, 4434000, 4436000, 4438000, 4440000, 4442000,
4444000, 4446000, 4448000, 4450000, 4452000, 4454000, 4456000, 4458000,
4460000, 4462000, 4464000, 4466000, 4468000, 4470000, 4472000, 4474000,
4476000, 4478000, 4480000, 4482000, 4484000, 4486000, 4488000, 4490000,
4492000, 4494000, 4496000, 4498000, 4500000, 4502000, 4504000, 4506000,
4508000, 4510000, 4512000, 4514000, 4516000, 4518000, 4520000, 4522000,
4524000, 4526000, 4528000, 4530000, 4532000, 4534000, 4536000, 4538000,
4540000, 4542000, 4544000, 4546000, 4548000, 4550000, 4552000, 4554000,
4556000, 4558000, 4560000, 4562000, 4564000, 4566000, 4568000, 4570000,
4572000, 4574000, 4576000, 4578000, 4580000, 4582000, 4584000, 4586000,
4588000, 4590000, 4592000, 4594000, 4596000, 4598000, 4600000, 4602000,
4604000, 4606000, 4608000, 4610000, 4612000, 4614000, 4616000, 4618000,
4620000, 4622000, 4624000, 4626000, 4628000, 4630000, 4632000, 4634000,
4636000, 4638000, 4640000, 4642000, 4644000, 4646000, 4648000, 4650000,
4652000, 4654000, 4656000, 4658000, 4660000, 4662000, 4664000, 4666000,
4668000, 4670000, 4672000, 4674000, 4676000, 4678000, 4680000, 4682000,
4684000, 4686000, 4688000, 4690000, 4692000, 4694000, 4696000, 4698000,
4700000, 4702000, 4704000, 4706000, 4708000, 4710000, 4712000, 4714000,
4716000, 4718000, 4720000, 4722000, 4724000, 4726000, 4728000, 4730000,
4732000, 4734000, 4736000, 4738000, 4740000, 4742000, 4744000, 4746000,
4748000, 4750000, 4752000, 4754000, 4756000, 4758000, 4760000, 4762000,
4764000, 4766000, 4768000, 4770000, 4772000, 4774000, 4776000, 4778000,
4780000, 4782000, 4784000, 4786000, 4788000, 4790000, 4792000, 4794000,
4796000, 4798000, 4800000, 4802000, 4804000, 4806000, 4808000, 4810000,
4812000, 4814000, 4816000, 4818000, 4820000, 4822000, 4824000, 4826000,
4828000, 4830000, 4832000, 4834000, 4836000, 4838000, 4840000, 4842000,
4844000, 4846000, 4848000, 4850000, 4852000, 4854000, 4856000, 4858000,
4860000, 4862000, 4864000, 4866000, 4868000, 4870000, 4872000, 4874000,
4876000, 4878000, 4880000, 4882000, 4884000, 4886000, 4888000, 4890000,
4892000, 4894000, 4896000, 4898000, 4900000, 4902000, 4904000, 4906000,
4908000, 4910000, 4912000, 4914000, 4916000, 4918000, 4920000, 4922000,
4924000, 4926000, 4928000, 4930000, 4932000, 4934000, 4936000, 4938000,
4940000, 4942000, 4944000, 4946000, 4948000, 4950000, 4952000, 4954000,
4956000, 4958000, 4960000, 4962000, 4964000, 4966000, 4968000, 4970000,
4972000, 4974000, 4976000, 4978000, 4980000, 4982000, 4984000, 4986000,
4988000, 4990000, 4992000, 4994000, 4996000, 4998000, 5000000, 5002000,
5004000, 5006000, 5008000, 5010000, 5012000, 5014000, 5016000, 5018000,
5020000, 5022000, 5024000, 5026000, 5028000, 5030000, 5032000, 5034000,
5036000, 5038000, 5040000, 5042000, 5044000, 5046000, 5048000, 5050000,
5052000, 5054000, 5056000, 5058000, 5060000, 5062000, 5064000, 5066000,
5068000, 5070000, 5072000, 5074000, 5076000, 5078000, 5080000, 5082000,
5084000, 5086000, 5088000, 5090000, 5092000, 5094000, 5096000, 5098000,
5100000, 5102000, 5104000, 5106000, 5108000, 5110000, 5112000, 5114000,
5116000, 5118000, 5120000, 5122000, 5124000, 5126000, 5128000, 5130000,
5132000, 5134000, 5136000, 5138000, 5140000, 5142000, 5144000, 5146000,
5148000, 5150000, 5152000, 5154000, 5156000, 5158000, 5160000, 5162000,
5164000, 5166000, 5168000, 5170000, 5172000, 5174000, 5176000, 5178000,
5180000, 5182000, 5184000, 5186000, 5188000, 5190000, 5192000, 5194000,
5196000, 5198000, 5200000, 5202000, 5204000, 5206000, 5208000, 5210000,
5212000, 5214000, 5216000, 5218000, 5220000, 5222000, 5224000, 5226000,
5228000, 5230000, 5232000, 5234000, 5236000, 5238000, 5240000, 5242000,
5244000, 5246000, 5248000, 5250000, 5252000, 5254000, 5256000, 5258000,
5260000, 5262000, 5264000, 5266000, 5268000, 5270000, 5272000, 5274000,
5276000, 5278000, 5280000, 5282000, 5284000, 5286000, 5288000, 5290000,
5292000, 5294000, 5296000, 5298000, 5300000, 5302000, 5304000, 5306000,
5308000, 5310000, 5312000, 5314000, 5316000, 5318000, 5320000, 5322000,
5324000, 5326000, 5328000, 5330000, 5332000, 5334000, 5336000, 5338000,
5340000, 5342000, 5344000, 5346000, 5348000, 5350000, 5352000, 5354000,
5356000, 5358000, 5360000, 5362000, 5364000, 5366000, 5368000, 5370000,
5372000, 5374000, 5376000, 5378000, 5380000, 5382000, 5384000, 5386000,
5388000, 5390000, 5392000, 5394000, 5396000, 5398000, 5400000, 5402000,
5404000, 5406000, 5408000, 5410000, 5412000, 5414000, 5416000, 5418000,
5420000, 5422000, 5424000, 5426000, 5428000, 5430000, 5432000, 5434000,
5436000, 5438000, 5440000, 5442000, 5444000, 5446000, 5448000, 5450000,
5452000, 5454000, 5456000, 5458000, 5460000, 5462000, 5464000, 5466000,
5468000, 5470000, 5472000, 5474000, 5476000, 5478000, 5480000, 5482000,
5484000, 5486000, 5488000, 5490000, 5492000, 5494000, 5496000, 5498000,
5500000, 5502000, 5504000, 5506000, 5508000, 5510000, 5512000, 5514000,
5516000, 5518000, 5520000, 5522000, 5524000, 5526000, 5528000, 5530000,
5532000, 5534000, 5536000, 5538000, 5540000, 5542000, 5544000, 5546000,
5548000, 5550000, 5552000, 5554000, 5556000, 5558000, 5560000, 5562000,
5564000, 5566000, 5568000, 5570000, 5572000, 5574000, 5576000, 5578000,
5580000, 5582000, 5584000, 5586000, 5588000, 5590000, 5592000, 5594000,
5596000, 5598000, 5600000, 5602000, 5604000, 5606000, 5608000, 5610000,
5612000, 5614000, 5616000, 5618000, 5620000, 5622000, 5624000, 5626000,
5628000, 5630000, 5632000, 5634000, 5636000, 5638000, 5640000, 5642000,
5644000, 5646000, 5648000, 5650000, 5652000, 5654000, 5656000, 5658000,
5660000, 5662000}

func main() {
   create, err := os.Create("textstream_eng.ttml")
   if err != nil {
      panic(err)
   }
   defer create.Close()
   for _, to := range tos {
      fmt.Println(to)
      name := func() string {
         b := []byte("dash/drm_playlist.af343964d7-textstream_eng=1000-")
         b = strconv.AppendInt(b, to, 10)
         b = append(b, ".dash"...)
         return string(b)
      }()
      open, err := os.Open(name)
      if err != nil {
         panic(err)
      }
      func() {
         defer open.Close()
         var file sofia.File
         err := file.Decode(open)
         if err != nil {
            panic(err)
         }
         for _, data := range file.MediaData.Data {
            create.Write(data)
            create.WriteString("\n")
         }
      }()
   }
}
