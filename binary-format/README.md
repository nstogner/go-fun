# Parsing a Custom Binary Format

This project implements a parser of a theoretical "mps7" binary format.

## Usage

There is a CLI utility provided that reads MPS7 format from stdin and writes JSON to stdout.

```sh
cat txnlog.dat | go run ./main.go > txnlog.json
```

