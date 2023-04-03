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

// Fprint is a convenience function to format data like Sprint
// and print it to the passed writer.
func Fprint(w io.Writer, data []byte, offset int) (int, error) {
	builder := strings.Builder{}

	written := 0
	for lineStart := offset / 16 * 16; lineStart < offset+len(data); lineStart += 16 {
		builder.Grow(7 * 16) // Reserve space for the whole line once
		line(&builder, data, offset, lineStart)
		n, err := fmt.Fprint(w, builder.String())
		if err != nil {
			return written + n, err
		}
		written += n
	}
	return written, nil
}
