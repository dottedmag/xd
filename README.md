# Extended hexdump for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/dottedmag/xd.svg)](https://pkg.go.dev/github.com/dottedmag/xd)

This package implements extended hexdump designed for ease of recognition
of patterns in binary data, while being a superset of classic `hexdump -C`.

Extended hexdump makes the structure of data stand out:

    00000000  01 00 00 00 01 00 0c 00  02 00 00 00              |☺▪▪▪☺▪♀▪☻▪▪▪|

    00000000  01 00 00 00 00 00 24 00  01 00 00 00 03 00 00 00  |☺▪▪▪▪▪$▪☺▪▪▪♥▪▪▪|
    00000010  10 00 00 00 75 6e 6b 6e  6f 77 6e 20 72 65 71 75  |►▪▪▪unknown requ|
    00000020  65 73 74 00                                       |est▪|

The format makes it obvious at a quick glance that the data is little-endian,
uses 4-byte and 2-byte integers and NUL-terminated strings.

Compare to the classic hexdump -C:

    00000000  01 00 00 00 01 00 0c 00  02 00 00 00              |............|

    00000000  01 00 00 00 00 00 24 00  01 00 00 00 03 00 00 00  |......$.........|
    00000010  10 00 00 00 75 6e 6b 6e  6f 77 6e 20 72 65 71 75  |....unknown requ|
    00000020  65 73 74 00                                       |est.|

## Usage

Print any block of data to stdout:

    xd.Print([]byte("hello world"), 0)

Or format to string:

    xd.Sprint([]byte("hello world"), 0)

Second argument is the offset in the printed stream, to properly print addresses
and to align printed data:

    xd.Sprint([]byte("first chunk"), 0)
    xd.Sprint([]byte("second chunk"), 11)
    xd.Sprint([]byte("third chunk"), 23)
    ...

## Command-line tool

There is a trivial command-line tool:

    go run github.com/dottedmag/xd/xd@master < your-binary

## Credits

Extended hexdump is invented by Ange Albertini in
[Corkami](https://github.com/angea/corkami/blob/master/src/HexII/braille/braille-ange).
and implemented in [braille dump](https://justine.lol/braille/) by Justine Tunney.

This implementation differs only in character used for NUL byte, as ▁ and ░ used
by Justine and Ange are hard to count quickly: they are block characters, so
consecutive NUL bytes in data become a continuous pattern. ▪ is perfectly suited
for NUL.

## Legal

© Mikhail Gusarov <dottedmag@dottedmag.net>.

Licensed under [Apache 2.0](LICENSE) license.
