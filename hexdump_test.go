package xd

import (
	"bytes"
	"fmt"
	"testing"
)

func ExamplePrint_simple() {
	Print([]byte("Hello world!\n"), 0)
	// Output: 00000000  48 65 6c 6c 6f 20 77 6f  72 6c 64 21 0a           |Hello world!◙|
}

func ExamplePrint_offset() {
	// If the printed data forms binary stream, you can supply offset
	// to align printed chunks of the stream properly.

	fmt.Println("[0]")
	Print([]byte("First line\n"), 0)
	fmt.Println("[1]")
	Print([]byte("Second line\n"), 11)
	fmt.Println("[2]")
	Print([]byte("Third line\n"), 23)
	fmt.Println("[3]")
	Print([]byte("END\n"), 34)

	// Output:
	// [0]
	// 00000000  46 69 72 73 74 20 6c 69  6e 65 0a                 |First line◙|
	// [1]
	// 00000000                                    53 65 63 6f 6e             |Secon|
	// 00000010  64 20 6c 69 6e 65 0a                              |d line◙|
	// [2]
	// 00000010                       54  68 69 72 64 20 6c 69 6e         |Third lin|
	// 00000020  65 0a                                             |e◙|
	// [3]
	// 00000020        45 4e 44 0a                                   |END◙|
}

func TestEmpty(t *testing.T) {
	if Sprint(nil, 0) != "" {
		t.Fail()
	}
}

func TestSimple(t *testing.T) {
	expe := "00000000  41                                                |A|\n"
	act := Sprint([]byte{'A'}, 0)
	if act != expe {
		t.Errorf("Expected\n%q, got\n%q", expe, act)
	}
}

func TestAlphabet(t *testing.T) {
	expe := "00000000  00 01 02 03 04 05 06 07  08 09 0a 0b 0c 0d 0e 0f  |▪☺☻♥♦♣♠•◘○◙♂♀♪♫☼|\n" +
		"00000010  10 11 12 13 14 15 16 17  18 19 1a 1b 1c 1d 1e 1f  |►◄↕‼¶§▬↨↑↓→←∟↔▲▼|\n" +
		`00000020  20 21 22 23 24 25 26 27  28 29 2a 2b 2c 2d 2e 2f  | !"#$%&'()*+,-./|` + "\n" +
		"00000030  30 31 32 33 34 35 36 37  38 39 3a 3b 3c 3d 3e 3f  |0123456789:;<=>?|\n" +
		"00000040  40 41 42 43 44 45 46 47  48 49 4a 4b 4c 4d 4e 4f  |@ABCDEFGHIJKLMNO|\n" +
		`00000050  50 51 52 53 54 55 56 57  58 59 5a 5b 5c 5d 5e 5f  |PQRSTUVWXYZ[\]^_|` + "\n" +
		"00000060  60 61 62 63 64 65 66 67  68 69 6a 6b 6c 6d 6e 6f  |`abcdefghijklmno|\n" +
		"00000070  70 71 72 73 74 75 76 77  78 79 7a 7b 7c 7d 7e 7f  |pqrstuvwxyz{|}~⌂|\n" +
		"00000080  80 81 82 83 84 85 86 87  88 89 8a 8b 8c 8d 8e 8f  |█⡀⢀⣀⠠⡠⢠⣠⠄⡄⢄⣄⠤⡤⢤⣤|\n" +
		"00000090  90 91 92 93 94 95 96 97  98 99 9a 9b 9c 9d 9e 9f  |⠁⡁⢁⣁⠡⡡⢡⣡⠅⡅⢅⣅⠥⡥⢥⣥|\n" +
		"000000a0  a0 a1 a2 a3 a4 a5 a6 a7  a8 a9 aa ab ac ad ae af  |⠃⡃⢃⣃⠣⡣⢣⣣⠇⡇⢇⣇⠧⡧⢧⣧|\n" +
		"000000b0  b0 b1 b2 b3 b4 b5 b6 b7  b8 b9 ba bb bc bd be bf  |⠉⡉⢉⣉⠩⡩⢩⣩⠍⡍⢍⣍⠭⡭⢭⣭|\n" +
		"000000c0  c0 c1 c2 c3 c4 c5 c6 c7  c8 c9 ca cb cc cd ce cf  |⠊⡊⢊⣊⠪⡪⢪⣪⠎⡎⢎⣎⠮⡮⢮⣮|\n" +
		"000000d0  d0 d1 d2 d3 d4 d5 d6 d7  d8 d9 da db dc dd de df  |⠑⡑⢑⣑⠱⡱⢱⣱⠕⡕⢕⣕⠵⡵⢵⣵|\n" +
		"000000e0  e0 e1 e2 e3 e4 e5 e6 e7  e8 e9 ea eb ec ed ee ef  |⠚⡚⢚⣚⠺⡺⢺⣺⠞⡞⢞⣞⠾⡾⢾⣾|\n" +
		"000000f0  f0 f1 f2 f3 f4 f5 f6 f7  f8 f9 fa fb fc fd fe ff  |⠛⡛⢛⣛⠻⡻⢻⣻⠟⡟⢟⣟⠿⡿⢿⣿|\n"
	var data []byte
	for i := 0; i < 256; i++ {
		data = append(data, byte(i))
	}
	act := Sprint(data, 0)
	if act != expe {
		t.Errorf("Expected\n%q, got\n%q", expe, act)
	}
}

func TestPartialLine(t *testing.T) {
	exp := "00000000     31 32 33 34 35 36 37  38 39 41 42 43 44 45      |123456789ABCDE|\n"
	act := Sprint([]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E'}, 1)
	if act != exp {
		t.Errorf("Expected\n%q, got\n%q", exp, act)
	}
}

func TestFprint(t *testing.T) {
	const expected = "00000020  30 31 32 33 34 35 36 37  38 39 41 42 43 44 45 46  |0123456789ABCDEF|\n"

	var buf bytes.Buffer
	Fprint(&buf, []byte("0123456789ABCDEF"), 32)
	if buf.String() != expected {
		t.Errorf("Expected\n%q, got\n%q", expected, buf.String())
	}
}
