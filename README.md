# jwt-with-grpc-performance-tests

Micro-benchmark on GRPC + JWT signing and verification on Node.js and Go.

## Results

- Node.js: v22.1.0
- Go: v1.22.1
- Hardware: Apple M3 Pro with 32GB RAM

## Single core

Running the servers:

```bash
node node-grpc-server.js
GOMAXPROCS=1 go run go-grpc-server.go
```

### HS256 - 1 connection, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50051

  Summary:
    Count:	374255
    Total:	10.00 s
    Slowest:	22.55 ms
    Fastest:	0.09 ms
    Average:	1.16 ms
    Requests/sec:	37425.74

  Response time histogram:
    0.095  [1]      |
    2.340  [365833] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    4.585  [8116]   |∎
    6.830  [151]    |
    9.075  [53]     |
    11.320 [1]      |
    13.565 [0]      |
    15.810 [0]      |
    18.055 [0]      |
    20.301 [1]      |
    22.546 [49]     |

  Latency distribution:
    10 % in 0.88 ms
    25 % in 0.96 ms
    50 % in 1.03 ms
    75 % in 1.20 ms
    90 % in 1.72 ms
    95 % in 2.04 ms
    99 % in 2.75 ms

  Status code distribution:
    [OK]            374205 responses
    [Unavailable]   50 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57379->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50052

  Summary:
    Count:	556070
    Total:	10.00 s
    Slowest:	9.92 ms
    Fastest:	0.04 ms
    Average:	0.69 ms
    Requests/sec:	55606.10

  Response time histogram:
    0.040 [1]      |
    1.028 [475103] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    2.016 [77999]  |∎∎∎∎∎∎∎
    3.005 [2359]   |
    3.993 [155]    |
    4.981 [173]    |
    5.970 [107]    |
    6.958 [49]     |
    7.946 [49]     |
    8.935 [32]     |
    9.923 [15]     |

  Latency distribution:
    10 % in 0.36 ms
    25 % in 0.45 ms
    50 % in 0.59 ms
    75 % in 0.86 ms
    90 % in 1.13 ms
    95 % in 1.30 ms
    99 % in 1.79 ms

  Status code distribution:
    [Canceled]      3 responses
    [Unavailable]   25 responses
    [OK]            556042 responses

  Error distribution:
    [25]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:55599->127.0.0.1:50052: use of closed network connection
    [3]    rpc error: code = Canceled desc = grpc: the client connection is closing
  ```
</details>

Go is ~1.48x faster than Node.js.   
Go uses less memory than Node.js (18MB peak vs 105MB peak) albeit Node.js has a higher base memory (~40MB).

### HS256 - 5 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ghz --insecure -z 10s --connections 5 --concurrency 250 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50051

  Summary:
    Count:	346427
    Total:	10.00 s
    Slowest:	35.55 ms
    Fastest:	1.01 ms
    Average:	6.87 ms
    Requests/sec:	34641.56

  Response time histogram:
    1.012  [1]      |
    4.466  [6710]   |∎
    7.920  [278366] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    11.373 [59290]  |∎∎∎∎∎∎∎∎∎
    14.827 [1495]   |
    18.281 [25]     |
    21.735 [130]    |
    25.189 [34]     |
    28.643 [69]     |
    32.097 [33]     |
    35.551 [40]     |

  Latency distribution:
    10 % in 5.53 ms
    25 % in 6.03 ms
    50 % in 6.75 ms
    75 % in 7.50 ms
    90 % in 8.53 ms
    95 % in 9.32 ms
    99 % in 10.78 ms

  Status code distribution:
    [OK]            346193 responses
    [Unavailable]   234 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57385->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57381->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57383->127.0.0.1:50051: use of closed network connection
    [49]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57382->127.0.0.1:50051: use of closed network connection
    [35]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57384->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ghz --insecure -z 10s --connections 5 --concurrency 250 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50052

  Summary:
    Count:	503190
    Total:	10.00 s
    Slowest:	19.02 ms
    Fastest:	0.05 ms
    Average:	4.18 ms
    Requests/sec:	50320.55

  Response time histogram:
    0.054  [1]      |
    1.950  [65374]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    3.847  [169445] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    5.744  [161000] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    7.640  [79643]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    9.537  [23651]  |∎∎∎∎∎∎
    11.434 [3411]   |∎
    13.331 [145]    |
    15.227 [68]     |
    17.124 [151]    |
    19.021 [51]     |

  Latency distribution:
    10 % in 1.71 ms
    25 % in 2.70 ms
    50 % in 4.02 ms
    75 % in 5.47 ms
    90 % in 6.84 ms
    95 % in 7.75 ms
    99 % in 9.32 ms

  Status code distribution:
    [OK]            502940 responses
    [Unavailable]   250 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54509->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54507->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54510->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54506->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54508->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.45x faster than Node.js.   
Go uses less memory than Node.js (18MB peak vs 105MB peak) albeit Node.js has a higher base memory (~40MB).

### ES256 - 1 connection, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50051

  Summary:
    Count:	114994
    Total:	10.00 s
    Slowest:	32.32 ms
    Fastest:	0.21 ms
    Average:	4.26 ms
    Requests/sec:	11215.61

  Response time histogram:
    0.212  [1]      |
    3.423  [3069]   |∎
    6.633  [106959] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    9.843  [1933]   |∎
    13.054 [93]     |
    16.264 [5]      |
    19.474 [0]      |
    22.684 [0]      |
    25.895 [0]      |
    29.105 [0]      |
    32.315 [50]     |

  Latency distribution:
    10 % in 3.88 ms
    25 % in 4.00 ms
    50 % in 4.11 ms
    75 % in 4.26 ms
    90 % in 4.90 ms
    95 % in 5.65 ms
    99 % in 7.13 ms

  Status code distribution:
    [OK]            112110 responses
    [Unavailable]   46 responses

  Error distribution:
    [46]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54673->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50052

  Summary:
    Count:	134062
    Total:	10.00 s
    Slowest:	11.94 ms
    Fastest:	0.14 ms
    Average:	3.56 ms
    Requests/sec:	13405.87

  Response time histogram:
    0.138  [1]     |
    1.318  [2520]  |∎∎
    2.498  [12122] |∎∎∎∎∎∎∎∎
    3.678  [60907] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    4.858  [52793] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    6.038  [4281]  |∎∎∎
    7.218  [1174]  |∎
    8.397  [94]    |
    9.577  [25]    |
    10.757 [65]    |
    11.937 [30]    |

  Latency distribution:
    10 % in 2.40 ms
    25 % in 3.34 ms
    50 % in 3.57 ms
    75 % in 3.93 ms
    90 % in 4.39 ms
    95 % in 4.72 ms
    99 % in 6.06 ms

  Status code distribution:
    [Unavailable]   50 responses
    [OK]            134012 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54678->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.17x faster than Node.js.   
Go uses less memory than Node.js (18MB peak vs 105MB peak) albeit Node.js has a higher base memory (~40MB).

### ES256 - 5 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 5 --concurrency 250 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50051

  Summary:
    Count:	115057
    Total:	10.00 s
    Slowest:	47.15 ms
    Fastest:	8.02 ms
    Average:	21.35 ms
    Requests/sec:	11505.75

  Response time histogram:
    8.016  [1]     |
    11.930 [82]    |
    15.843 [397]   |
    19.756 [8664]  |∎∎∎∎
    23.670 [98817] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    27.583 [3586]  |∎
    31.497 [1972]  |∎
    35.410 [1144]  |
    39.324 [62]    |
    43.237 [47]    |
    47.151 [35]    |

  Latency distribution:
    10 % in 20.12 ms
    25 % in 20.71 ms
    50 % in 21.16 ms
    75 % in 21.84 ms
    90 % in 22.46 ms
    95 % in 24.33 ms
    99 % in 31.69 ms

  Status code distribution:
    [OK]            114807 responses
    [Unavailable]   250 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57411->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57409->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57408->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57410->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57407->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 5 --concurrency 250 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50052

  Summary:
    Count:	127735
    Total:	10.00 s
    Slowest:	51.44 ms
    Fastest:	0.18 ms
    Average:	18.92 ms
    Requests/sec:	12773.04

  Response time histogram:
    0.180  [1]     |
    5.305  [2852]  |∎∎∎
    10.431 [13259] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    15.556 [21732] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    20.682 [34743] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    25.807 [36409] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    30.933 [14790] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    36.058 [3298]  |∎∎∎∎
    41.184 [242]   |
    46.310 [61]    |
    51.435 [98]    |

  Latency distribution:
    10 % in 9.17 ms
    25 % in 14.06 ms
    50 % in 19.43 ms
    75 % in 23.21 ms
    90 % in 27.21 ms
    95 % in 29.81 ms
    99 % in 33.71 ms

  Status code distribution:
    [OK]            127485 responses
    [Unavailable]   250 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54775->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54773->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54772->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54771->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:54774->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.11x faster than Node.js.   
Go uses less memory than Node.js (18MB peak vs 105MB peak) albeit Node.js has a higher base memory (~40MB).

## Multi-core with 3 cores

Running the servers:

```bash
NODEMAXPROCS=3 node node-grpc-server-cluster.js
GOMAXPROCS=3 go run go-grpc-server.go
```

### HS256 - 3 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 3 --concurrency 150 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50051

  Summary:
    Count:	545116
    Total:	10.00 s
    Slowest:	21.76 ms
    Fastest:	0.07 ms
    Average:	2.27 ms
    Requests/sec:	54510.34

  Response time histogram:
    0.073  [1]      |
    2.242  [306396] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    4.411  [214045] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    6.579  [21935]  |∎∎∎
    8.748  [2313]   |
    10.917 [214]    |
    13.086 [23]     |
    15.255 [23]     |
    17.424 [52]     |
    19.593 [51]     |
    21.762 [47]     |

  Latency distribution:
    10 % in 1.03 ms
    25 % in 1.54 ms
    50 % in 2.11 ms
    75 % in 2.76 ms
    90 % in 3.68 ms
    95 % in 4.32 ms
    99 % in 5.87 ms

  Status code distribution:
    [OK]            545100 responses
    [Unavailable]   16 responses

  Error distribution:
    [13]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57439->127.0.0.1:50051: use of closed network connection
    [3]    rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57438->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 3 --concurrency 150 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50052

  Summary:
    Count:	686222
    Total:	10.00 s
    Slowest:	16.86 ms
    Fastest:	0.03 ms
    Average:	1.10 ms
    Requests/sec:	68614.26

  Response time histogram:
    0.032  [1]      |
    1.715  [567391] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    3.398  [101329] |∎∎∎∎∎∎∎
    5.080  [12062]  |∎
    6.763  [3349]   |
    8.446  [1053]   |
    10.128 [568]    |
    11.811 [217]    |
    13.494 [112]    |
    15.176 [46]     |
    16.859 [35]     |

  Latency distribution:
    10 % in 0.22 ms
    25 % in 0.45 ms
    50 % in 0.89 ms
    75 % in 1.45 ms
    90 % in 2.12 ms
    95 % in 2.71 ms
    99 % in 4.65 ms

  Status code distribution:
    [Unavailable]   59 responses
    [OK]            686163 responses

  Error distribution:
    [32]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:55261->127.0.0.1:50052: use of closed network connection
    [3]    rpc error: code = Unavailable desc = transport is closing
    [23]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:55259->127.0.0.1:50052: use of closed network connection
    [1]    rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:55260->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.25x faster than Node.js.   
Go uses less memory than Node.js (20MB peak vs 3 * 105MB peak) albeit Node.js has a higher base memory (3 * ~40MB).

### HS256 - 15 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 15 --concurrency 750 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50051

  Summary:
    Count:	517687
    Total:	10.00 s
    Slowest:	63.92 ms
    Fastest:	0.08 ms
    Average:	13.07 ms
    Requests/sec:	51751.83

  Response time histogram:
    0.082  [1]      |
    6.466  [58756]  |∎∎∎∎∎∎∎∎∎∎∎∎
    12.850 [199108] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    19.234 [200387] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    25.618 [47952]  |∎∎∎∎∎∎∎∎∎∎
    32.002 [5655]   |∎
    38.386 [2418]   |
    44.770 [1016]   |
    51.154 [945]    |
    57.538 [759]    |
    63.922 [272]    |

  Latency distribution:
    10 % in 6.06 ms
    25 % in 9.22 ms
    50 % in 12.87 ms
    75 % in 16.19 ms
    90 % in 19.72 ms
    95 % in 22.13 ms
    99 % in 32.31 ms

  Status code distribution:
    [OK]            517269 responses
    [Unavailable]   415 responses
    [Canceled]      3 responses

  Error distribution:
    [8]    rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57453->127.0.0.1:50051: use of closed network connection
    [40]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57457->127.0.0.1:50051: use of closed network connection
    [34]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57463->127.0.0.1:50051: use of closed network connection
    [39]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57452->127.0.0.1:50051: use of closed network connection
    [39]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57462->127.0.0.1:50051: use of closed network connection
    [10]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57455->127.0.0.1:50051: use of closed network connection
    [3]    rpc error: code = Canceled desc = grpc: the client connection is closing
    [5]    rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57454->127.0.0.1:50051: use of closed network connection
    [18]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57458->127.0.0.1:50051: use of closed network connection
    [16]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57460->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57461->127.0.0.1:50051: use of closed network connection
    [38]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57466->127.0.0.1:50051: use of closed network connection
    [27]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57464->127.0.0.1:50051: use of closed network connection
    [44]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57456->127.0.0.1:50051: use of closed network connection
    [42]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57465->127.0.0.1:50051: use of closed network connection
    [5]    rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57459->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 15 --concurrency 750 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"HS256"}' 0.0.0.0:50052

  Summary:
    Count:	683605
    Total:	10.00 s
    Slowest:	70.89 ms
    Fastest:	0.03 ms
    Average:	5.24 ms
    Requests/sec:	68354.53

  Response time histogram:
    0.033  [1]      |
    7.118  [505622] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    14.204 [150190] |∎∎∎∎∎∎∎∎∎∎∎∎
    21.289 [16992]  |∎
    28.374 [5708]   |
    35.460 [2582]   |
    42.545 [1405]   |
    49.631 [601]    |
    56.716 [304]    |
    63.801 [113]    |
    70.887 [84]     |

  Latency distribution:
    10 % in 0.56 ms
    25 % in 1.79 ms
    50 % in 4.13 ms
    75 % in 7.28 ms
    90 % in 10.51 ms
    95 % in 13.24 ms
    99 % in 25.60 ms

  Status code distribution:
    [OK]            683602 responses
    [Unavailable]   3 responses

  Error distribution:
    [2]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57155->127.0.0.1:50052: use of closed network connection
    [1]   rpc error: code = Unavailable desc = transport is closing
  ```
</details>

Go is ~1.32x faster than Node.js.   
Go uses less memory than Node.js (20MB peak vs 3 * 105MB peak) albeit Node.js has a higher base memory (3 * ~40MB).

### ES256 - 3 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 3 --concurrency 150 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50051

  Summary:
    Count:	298483
    Total:	10.00 s
    Slowest:	25.01 ms
    Fastest:	0.17 ms
    Average:	4.69 ms
    Requests/sec:	29847.25

  Response time histogram:
    0.169  [1]      |
    2.653  [8086]   |∎
    5.137  [220883] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    7.622  [63349]  |∎∎∎∎∎∎∎∎∎∎∎
    10.106 [5521]   |∎
    12.590 [339]    |
    15.075 [89]     |
    17.559 [2]      |
    20.043 [11]     |
    22.527 [21]     |
    25.012 [31]     |

  Latency distribution:
    10 % in 3.94 ms
    25 % in 4.18 ms
    50 % in 4.44 ms
    75 % in 5.07 ms
    90 % in 6.06 ms
    95 % in 6.76 ms
    99 % in 8.32 ms

  Status code distribution:
    [OK]            298333 responses
    [Unavailable]   150 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57472->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57473->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57474->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 3 --concurrency 150 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50052

  Summary:
    Count:	322989
    Total:	10.00 s
    Slowest:	45.25 ms
    Fastest:	0.12 ms
    Average:	4.30 ms
    Requests/sec:	32298.41

  Response time histogram:
    0.118  [1]      |
    4.631  [200524] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    9.144  [113865] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    13.657 [7716]   |∎∎
    18.170 [521]    |
    22.683 [106]    |
    27.196 [26]     |
    31.709 [15]     |
    36.222 [42]     |
    40.736 [12]     |
    45.249 [13]     |

  Latency distribution:
    10 % in 1.73 ms
    25 % in 2.79 ms
    50 % in 4.02 ms
    75 % in 5.49 ms
    90 % in 7.14 ms
    95 % in 8.25 ms
    99 % in 10.48 ms

  Status code distribution:
    [OK]            322841 responses
    [Canceled]      1 responses
    [Unavailable]   147 responses

  Error distribution:
    [1]    rpc error: code = Canceled desc = grpc: the client connection is closing
    [49]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57178->127.0.0.1:50052: use of closed network connection
    [48]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57179->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57180->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.08x faster than Node.js.   
Go uses less memory than Node.js (20MB peak vs 3 * 105MB peak) albeit Node.js has a higher base memory (3 * ~40MB).

### ES256 - 15 connections, 50 concurrent calls per connection

<details>
  <summary>Node.js</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 15 --concurrency 750 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50051

  Summary:
    Count:	302979
    Total:	10.00 s
    Slowest:	75.19 ms
    Fastest:	2.03 ms
    Average:	23.89 ms
    Requests/sec:	30297.08

  Response time histogram:
    2.032  [1]      |
    9.348  [419]    |
    16.664 [7737]   |∎∎
    23.980 [171025] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    31.296 [109646] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    38.612 [12296]  |∎∎∎
    45.928 [823]    |
    53.245 [237]    |
    60.561 [76]     |
    67.877 [0]      |
    75.193 [1]      |

  Latency distribution:
    10 % in 19.63 ms
    25 % in 21.93 ms
    50 % in 23.37 ms
    75 % in 25.70 ms
    90 % in 28.97 ms
    95 % in 30.97 ms
    99 % in 35.81 ms

  Status code distribution:
    [OK]            302261 responses
    [Unavailable]   718 responses

  Error distribution:
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57486->127.0.0.1:50051: use of closed network connection
    [20]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57491->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57493->127.0.0.1:50051: use of closed network connection
    [48]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57499->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57492->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57496->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57500->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57489->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57487->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57488->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57490->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57497->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57495->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57494->127.0.0.1:50051: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57498->127.0.0.1:50051: use of closed network connection
  ```
</details>

<details>
  <summary>Go</summary>

  ```
  ❯ ghz --insecure -z 10s --connections 15 --concurrency 750 --proto ./jwt.proto --call jwt.services.JwtService.GenerateJwt -d '{"algorithm":"ES256"}' 0.0.0.0:50052

  Summary:
    Count:	323963
    Total:	10.00 s
    Slowest:	95.93 ms
    Fastest:	0.14 ms
    Average:	21.14 ms
    Requests/sec:	32392.82

  Response time histogram:
    0.136  [1]      |
    9.715  [41777]  |∎∎∎∎∎∎∎∎∎∎∎∎∎
    19.294 [96084]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    28.873 [125620] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    38.452 [44818]  |∎∎∎∎∎∎∎∎∎∎∎∎∎∎
    48.031 [10183]  |∎∎∎
    57.610 [2966]   |∎
    67.189 [1217]   |
    76.768 [562]    |
    86.347 [51]     |
    95.926 [25]     |

  Latency distribution:
    10 % in 8.48 ms
    25 % in 14.19 ms
    50 % in 20.87 ms
    75 % in 26.81 ms
    90 % in 33.19 ms
    95 % in 37.91 ms
    99 % in 53.03 ms

  Status code distribution:
    [Unavailable]   659 responses
    [OK]            323304 responses

  Error distribution:
    [34]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57247->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57251->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57252->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57256->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57255->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57257->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57260->127.0.0.1:50052: use of closed network connection
    [25]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57248->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57250->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57254->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57258->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57261->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57249->127.0.0.1:50052: use of closed network connection
    [50]   rpc error: code = Unavailable desc = error reading from server: read tcp 127.0.0.1:57259->127.0.0.1:50052: use of closed network connection
  ```
</details>

Go is ~1.07x faster than Node.js.   
Go uses less memory than Node.js (20MB peak vs 3 * 105MB peak) albeit Node.js has a higher base memory (3 * ~40MB).
