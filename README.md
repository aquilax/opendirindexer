# opendirindexer

Open directory indexer command line tool

## Installation

```
go install github.com/aquilax/opendirindexer
```

## Usage

```
Usage of opendirindexer:
opendirindexer [OPTIONS] URL
  -debug
    	Enable debugging
  -ingoreRobots
    	Ignores robots.txt restrictions
  -userAgent string
    	set user agent (default "opendirindexer/1.0")
```

Example:

Given the following structure:

```bash
$ tree
.
├── 4.txt
├── test1
│   ├── 1.txt
│   └── test1.1
│       └── 1.txt
└── test2
    └── 3.txt

3 directories, 4 files
```
The output will be:

```bash
$ opendirindexer http://localhost:8000
http://localhost:8000/4.txt
http://localhost:8000/test1/1.txt
http://localhost:8000/test2/3.txt
http://localhost:8000/test1/test1.1/1.txt
```

And the server log looks like:
```bash
$ python -m http.server 8000
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...
127.0.0.1 - - [14/Apr/2019 08:43:25] "GET / HTTP/1.1" 200 -
127.0.0.1 - - [14/Apr/2019 08:43:25] "GET /test1/ HTTP/1.1" 200 -
127.0.0.1 - - [14/Apr/2019 08:43:25] "GET /test2/ HTTP/1.1" 200 -
127.0.0.1 - - [14/Apr/2019 08:43:25] "GET /test1/test1.1/ HTTP/1.1" 200 -
```
