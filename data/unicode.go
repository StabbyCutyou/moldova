package data

// If you'd like to see a new range added here, please open a PR, and include a link
// to wikipedia or another resource demonstrating what characters are in the range
var PrintableRanges = [][]int{
	// Cyrillic
	{0x0400, 0x04ff},
	// Greek
	{0x0377, 0x03ff},
	// Hangul
	{0x1100, 0x11ff},
	// Chinese / Kanji
	{0x4e00, 0x4f80},
	{0x5000, 0x9fa0},
	{0x3400, 0x4db0},
	// Arabic
	{0x0600, 0x06ff},
	// Japanese Kana
	{0x30a0, 0x30f0},
	// Arabic Presentation
	{0xfb50, 0xfdff},
	// Thai
	{0x0e00, 0x0e7f},
	// Phoenician
	{0x10900, 0x1091f},
}
