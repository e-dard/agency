package agency

import (
	"bytes"
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/benbjohnson/agency/data"
)

// eof is an internal EOF marker.
var eof = errors.New("eof")

// Scanner is a user agent tokenizer.
type Scanner struct {
	c      rune
	buf    []byte
	buflen int
	idx    int
	size   int
	prevstart int
	browsers []*data.Browser
}

// NewScanner creates a new user agent scanner.
func NewScanner() *Scanner {
	return &Scanner{
		browsers: make([]*data.Browser, data.MaxRank),
	}
}

// Scan scans a user agent string for device information.
func (s *Scanner) ScanBytes(b []byte) (*UserAgent, error) {
	var ua = new(UserAgent)
	s.buf = b
	s.buflen = len(b)
	s.reset()

	// Iterate over each word in the byte slice.
	for {
		unigram, bigram, err := s.readNgrams()
		if err == eof {
			break
		} else if err != nil {
			return nil, err
		}

		if ua.Browser.Type == "" {
			s.matchBrowser(unigram, bigram)
		}
	}

	// Find browser by rank level.
	for _, browser := range s.browsers {
		if browser != nil {
			ua.Browser.Type = browser.Type
			ua.Browser.Name = browser.Name
			break
		}
	}

	return ua, nil
}

// Scan scans a user agent string for device information.
func (s *Scanner) Scan(str string) (*UserAgent, error) {
	return s.ScanBytes([]byte(str))
}

// read retrieves the next rune from the string.
func (s *Scanner) read() error {
	if s.idx >= s.buflen {
		return eof
	}

	// Read a single byte and then determine if utf8 decoding is needed.
	b := s.buf[s.idx]
	if b < utf8.RuneSelf {
		s.c = rune(b)
		s.size = 1
	} else {
		s.c, s.size = utf8.DecodeRune(s.buf[s.idx:])
	}
	s.idx += s.size
	return nil
}

// unread moves back one rune. Only works once.
func (s *Scanner) unread() {
	s.idx -= s.size
	s.size = 0
}

// readWord reads a word and previous bigram from the string.
func (s *Scanner) readNgrams() ([]byte, []byte, error) {
	var index int
	start := s.idx
	for {
		if err := s.read(); err == eof {
			break
		}

		// Only read in letters, numbers and some punctuation.
		if unicode.IsLetter(s.c) || unicode.IsDigit(s.c) || s.c == '-' || s.c == '.' {
			index++
		} else if index == 0 {
			// This section skips over initial non-word characters.
			start = s.idx
		} else {
			s.unread()
			break
		}
	}

	// If nothing was read then it's EOF.
	if s.idx == start {
		return nil, nil, eof
	}

	unigram := s.buf[start:s.idx]
	bigram := s.buf[s.prevstart:s.idx]
	s.prevstart = start
	return unigram, bigram, nil
}

// match checks a unigram and bigram against the list of browser tokens.
func (s *Scanner) matchBrowser(unigram []byte, bigram []byte) {
	for _, browser := range data.Browsers {
		if bytes.Equal(unigram, browser.Token) || bytes.Equal(bigram, browser.Token) {
			s.browsers[browser.Rank] = browser
		}
	}
}

// reset re-initializes the state of the scanner.
func (s *Scanner) reset() {
	s.idx = 0
	s.size = 0
	s.prevstart = 0

	for i, _ := range s.browsers {
		s.browsers[i] = nil
	}
}

// Scan extracts properties from a user agent byte slice.
func ScanBytes(b []byte) (*UserAgent, error) {
	return NewScanner().ScanBytes(b)
}

// ScanString extracts properties from a user agent string.
func Scan(str string) (*UserAgent, error) {
	return NewScanner().Scan(str)
}