// Package xd implements extended hexdump designed for ease of recognition
// of patterns in binary data, while being a superset of classic hexdump -C.
//
// Extended hexdump makes the structure of data stand out:
//   00000000  01 00 00 00 01 00 0c 00  02 00 00 00              |☺▪▪▪☺▪♀▪☻▪▪▪|
//
//   00000000  01 00 00 00 00 00 24 00  01 00 00 00 03 00 00 00  |☺▪▪▪▪▪$▪☺▪▪▪♥▪▪▪|
//   00000010  10 00 00 00 75 6e 6b 6e  6f 77 6e 20 72 65 71 75  |►▪▪▪unknown requ|
//   00000020  65 73 74 00                                       |est▪|
//
// The format makes it obvious at a quick glance that the data is little-endian,
// uses 4-byte and 2-byte integers and NUL-terminated strings.
//
// Compare to the classic hexdump -C:
//   00000000  01 00 00 00 01 00 0c 00  02 00 00 00              |............|
//
//   00000000  01 00 00 00 00 00 24 00  01 00 00 00 03 00 00 00  |......$.........|
//   00000010  10 00 00 00 75 6e 6b 6e  6f 77 6e 20 72 65 71 75  |....unknown requ|
//   00000020  65 73 74 00                                       |est.|
//
// Extended hexdump is invented by Ange Albertini to use in Corkami
// (https://github.com/angea/corkami/blob/master/src/HexII/braille/braille-ange),
// and implemented in braille dump by Justine Tunney (https://justine.lol/braille/).
//
// This implementation differs only in character used for NUL byte, as ▁ and ░ used by
// Justine and Ange are hard to count, being block characters. ▪ is perfectly suited
// for NUL.
//
package xd

import (
	"fmt"
	"io"
	"strings"
)

var charmap = []string{
	"▪", "☺", "☻", "♥", "♦", "♣", "♠", "•",
	"◘", "○", "◙", "♂", "♀", "♪", "♫", "☼",
	"►", "◄", "↕", "‼", "¶", "§", "▬", "↨",
	"↑", "↓", "→", "←", "∟", "↔", "▲", "▼",
	" ", "!", `"`, "#", "$", "%", "&", "'",
	"(", ")", "*", "+", ",", "-", ".", "/",
	"0", "1", "2", "3", "4", "5", "6", "7",
	"8", "9", ":", ";", "<", "=", ">", "?",
	"@", "A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W",
	"X", "Y", "Z", "[", `\`, "]", "^", "_",
	"`", "a", "b", "c", "d", "e", "f", "g",
	"h", "i", "j", "k", "l", "m", "n", "o",
	"p", "q", "r", "s", "t", "u", "v", "w",
	"x", "y", "z", "{", "|", "}", "~", "⌂",
	"█", "⡀", "⢀", "⣀", "⠠", "⡠", "⢠", "⣠",
	"⠄", "⡄", "⢄", "⣄", "⠤", "⡤", "⢤", "⣤",
	"⠁", "⡁", "⢁", "⣁", "⠡", "⡡", "⢡", "⣡",
	"⠅", "⡅", "⢅", "⣅", "⠥", "⡥", "⢥", "⣥",
	"⠃", "⡃", "⢃", "⣃", "⠣", "⡣", "⢣", "⣣",
	"⠇", "⡇", "⢇", "⣇", "⠧", "⡧", "⢧", "⣧",
	"⠉", "⡉", "⢉", "⣉", "⠩", "⡩", "⢩", "⣩",
	"⠍", "⡍", "⢍", "⣍", "⠭", "⡭", "⢭", "⣭",
	"⠊", "⡊", "⢊", "⣊", "⠪", "⡪", "⢪", "⣪",
	"⠎", "⡎", "⢎", "⣎", "⠮", "⡮", "⢮", "⣮",
	"⠑", "⡑", "⢑", "⣑", "⠱", "⡱", "⢱", "⣱",
	"⠕", "⡕", "⢕", "⣕", "⠵", "⡵", "⢵", "⣵",
	"⠚", "⡚", "⢚", "⣚", "⠺", "⡺", "⢺", "⣺",
	"⠞", "⡞", "⢞", "⣞", "⠾", "⡾", "⢾", "⣾",
	"⠛", "⡛", "⢛", "⣛", "⠻", "⡻", "⢻", "⣻",
	"⠟", "⡟", "⢟", "⣟", "⠿", "⡿", "⢿", "⣿",
}

func line(w *strings.Builder, data []byte, dataOffset int, lineStart int) {
	fmt.Fprintf(w, "%08x ", lineStart)
	for i := 0; i < 16; i++ {
		if i == 8 {
			w.WriteString(" ")
		}
		if lineStart+i < dataOffset || lineStart+i >= len(data)+dataOffset {
			w.WriteString("   ")
		} else {
			fmt.Fprintf(w, " %02x", data[lineStart+i-dataOffset])
		}
	}
	w.WriteString("  ")
	for i := 0; i < 16 && lineStart+i < len(data)+dataOffset; i++ {
		if lineStart+i < dataOffset {
			w.WriteString(" ")
		} else {
			if lineStart+i == dataOffset || i == 0 {
				w.WriteString("|")
			}
			w.WriteString(charmap[data[lineStart+i-dataOffset]])
		}
	}
	w.WriteString("|\n")
}

// Sprint takes the binary data and formats it. Offset is the
// offset in the binary stream, used to print address, and to
// align printed data to 16 bytes boundary.
func Sprint(data []byte, offset int) string {
	builder := strings.Builder{}

	// Every byte of data produces between 5 and 7 bytes of output
	// (fewer bytes for printable ASCII, more bytes for non-ASCII)
	builder.Grow(6 * len(data))

	for lineStart := offset / 16 * 16; lineStart < offset+len(data); lineStart += 16 {
		line(&builder, data, offset, lineStart)
	}

	return builder.String()
}

// Print is a convenience function to format data like Sprint
// and print it to stdout
func Print(data []byte, offset int) (int, error) {
	return fmt.Print(Sprint(data, offset))
}

// Print is a convenience function to format data like Sprint
// and print it to the passed writer.
func Fprint(w io.Writer, data []byte, offset int) (int, error) {
	builder := strings.Builder{}

	written := 0
	for lineStart := offset / 16 * 16; lineStart < offset+len(data); lineStart += 16 {
		builder.Grow(7 * 16) // Reserve space for the whole line once
		line(&builder, data, offset, lineStart)
		n, err := fmt.Print(builder.String())
		if err != nil {
			return written + n, err
		}
		written += n
	}
	return written, nil
}
