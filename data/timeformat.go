package data

// TimeFormats is a lookup map of a few common, defined formatting strings for ease
// of use.
var TimeFormats = map[string]string{
	// simple is a bare-bones, MySQL insert friendly time format
	"simple": "2006-01-02 15:04:05",
	// SimpleTimeWithZoneFormat is the same as SimpleTimeFormat, but with the Timezone set
	"simpletz": "2006-01-02 15:04:05 -0700",
}
