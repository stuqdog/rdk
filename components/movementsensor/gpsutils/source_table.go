// Package gpsutils implements necessary functions to set and return
// NTRIP information here.
package gpsutils

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// id numbers of the different fields returned in the standard
	// Stream response from the ntrip client, numbered 1-18.
	// Information on each field is explained int the comments
	// of the Stream struct.
	mp            = 1
	id            = 2
	format        = 3
	formatDetails = 4
	carrierField  = 5
	navsystem     = 6
	network       = 7
	country       = 8
	latitude      = 9
	longitude     = 10
	nmeaBit       = 11
	solution      = 12
	generator     = 13
	compression   = 14
	auth          = 15
	feeBit        = 16
	bitRateField  = 17
	misc          = 18
	floatbitsize  = 32
	streamSize    = 200
)

// Sourcetable struct contains the stream.
type Sourcetable struct {
	Streams []Stream
}

// Stream contrains a stream record in sourcetable.
// https://software.rtcm-ntrip.org/wiki/STR
type Stream struct {
	MP             string   // Datastream mountpoint
	Identifier     string   // Source identifier (most time nearest city)
	Format         string   // Data format of generic type (https://software.rtcm-ntrip.org/wiki/STR#DataFormats)
	FormatDetails  string   // Specifics of data format (https://software.rtcm-ntrip.org/wiki/STR#DataFormats)
	Carrier        int      // Phase information about GNSS correction (https://software.rtcm-ntrip.org/wiki/STR#Carrier)
	NavSystem      []string // Multiple navigation system (https://software.rtcm-ntrip.org/wiki/STR#NavigationSystem)
	Network        string   // Network record in sourcetable (https://software.rtcm-ntrip.org/wiki/NET)
	Country        string   // ISO 3166 country code (https://en.wikipedia.org/wiki/ISO_3166-1)
	Latitude       float32  // Position, Latitude in degree
	Longitude      float32  // Position, Longitude in degree
	Nmea           bool     // Caster requires NMEA input (1) or not (0)
	Solution       int      // Generated by single base (0) or network (1)
	Generator      string   // Generating soft- or hardware
	Compression    string   // Compression algorithm
	Authentication string   // Access protection for data streams None (N), Basic (B) or Digest (D)
	Fee            bool     // User fee for data access: yes (Y) or no (N)
	BitRate        int      // Datarate in bits per second
	Misc           string   // Miscellaneous information
}

// parseStream parses a line from the sourcetable.
func parseStream(line string) (Stream, error) {
	fields := strings.Split(line, ";")

	// Standard stream contains 19 fields.
	// They are enumerated by their constants at the top of the file
	if len(fields) < 19 {
		return Stream{}, fmt.Errorf("missing fields at stream line: %s", line)
	}

	if fields[carrierField] == "" {
		fields[carrierField] = "0"
	}
	carrier, err := strconv.Atoi(fields[carrierField])
	if err != nil {
		return Stream{}, fmt.Errorf("cannot parse the streams carrier in line: %s", line)
	}

	satSystems := strings.Split(fields[navsystem], "+")

	lat, err := strconv.ParseFloat(fields[latitude], floatbitsize)
	if err != nil {
		return Stream{}, fmt.Errorf("cannot parse the streams latitude in line: %s", line)
	}
	lon, err := strconv.ParseFloat(fields[longitude], floatbitsize)
	if err != nil {
		return Stream{}, fmt.Errorf("cannot parse the streams longitude in line: %s", line)
	}

	nmea, err := strconv.ParseBool(fields[nmeaBit])
	if err != nil {
		return Stream{}, fmt.Errorf("cannot parse the streams nmea in line: %s", line)
	}

	sol, err := strconv.Atoi(fields[solution])
	if err != nil {
		return Stream{}, fmt.Errorf("cannot parse the streams solution in line: %s", line)
	}

	fee := false
	if fields[feeBit] == "Y" {
		fee = true
	}

	bitrate, err := strconv.Atoi(fields[bitRateField])
	if err != nil {
		bitrate = 0
	}

	return Stream{
		MP: fields[mp], Identifier: fields[id], Format: fields[format], FormatDetails: fields[formatDetails],
		Carrier: carrier, NavSystem: satSystems, Network: fields[network], Country: fields[country],
		Latitude: float32(lat), Longitude: float32(lon), Nmea: nmea, Solution: sol, Generator: fields[generator],
		Compression: fields[compression], Authentication: fields[auth], Fee: fee, BitRate: bitrate, Misc: fields[misc],
	}, nil
}

// HasStream checks if the sourcetable contains the given mountpoint in it's stream.
func (st *Sourcetable) HasStream(mountpoint string) (Stream, bool) {
	for _, str := range st.Streams {
		if str.MP == mountpoint {
			return str, true
		}
	}

	return Stream{}, false
}
