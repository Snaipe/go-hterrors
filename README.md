# go-hterrors

[![GoDoc](https://godoc.org/snai.pe/go-hterrors?status.svg)](https://godoc.org/snai.pe/go-hterrors)  

```
go get snai.pe/go-hterrors
```

A Go convenience library for handling non-2xx HTTP error codes.

The library does a best-effort parsing of the response body based on common MIME-types to produce
a reasonable Go error.
