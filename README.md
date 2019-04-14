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

```
opendirindexer https://example.com/
```

Will output to stdout all links to files