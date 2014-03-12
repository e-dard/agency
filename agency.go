package agency

import (
	"bytes"
	"encoding/csv"
	"strconv"
)

// maxRank is the highest rank that can be used by the data files.
const maxRank = 5

type browser struct {
	rank  int
	typ   string
	name  string
	token []byte
}

var browsers []*browser

func init() {
	data, _ := browser_csv()
	records, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		panic("parse error (browser.csv): " + err.Error())
	}
	for _, record := range records {
		rank, _ := strconv.Atoi(record[0])
		browsers = append(browsers, &browser{rank, record[1], record[2], []byte(record[3])})
	}
}

type device struct {
	rank  int
	typ   string
	token []byte
}

var devices []*device

func init() {
	data, _ := device_csv()
	records, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		panic("parse error (device.csv): " + err.Error())
	}
	for _, record := range records {
		rank, _ := strconv.Atoi(record[0])
		devices = append(devices, &device{rank, record[1], []byte(record[2])})
	}
}

type mobile struct {
	token []byte
}

var mobiles []*mobile

func init() {
	data, _ := mobile_csv()
	records, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		panic("parse error (mobile.csv): " + err.Error())
	}
	for _, record := range records {
		mobiles = append(mobiles, &mobile{[]byte(record[0])})
	}
}

type os struct {
	rank    int
	name    string
	version string
	token   []byte
}

var oses []*os

func init() {
	data, _ := os_csv()
	records, err := csv.NewReader(bytes.NewBuffer(data)).ReadAll()
	if err != nil {
		panic("parse error (os.csv): " + err.Error())
	}
	for _, record := range records {
		rank, _ := strconv.Atoi(record[0])
		oses = append(oses, &os{rank, record[1], record[2], []byte(record[3])})
	}
}
