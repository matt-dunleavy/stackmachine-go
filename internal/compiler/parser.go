package compiler

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// Parser parses source code for the stack machine
type Parser struct {
	reader *bufio.Reader
	lineNo int
}

// NewParser creates a new parser from a reader
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
		lineNo: 1,
	}
}

// GetLineNo returns the current line number
func (p *Parser) GetLineNo() int {
	return p.lineNo
}

// UpdateLineNo updates the line number if a newline is encountered
func (p *Parser) UpdateLineNo(c rune) rune {
	if c == '\n' {
		p.lineNo++
	}
	return c
}

// GetChar reads a character, updating the line number if needed
func (p *Parser) GetChar() (rune, error) {
	r, _, err := p.reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return p.UpdateLineNo(r), nil
}

// UngetChar pushes a character back to the reader
func (p *Parser) UngetChar(r rune) error {
	if r == '\n' {
		p.lineNo--
	}
	return p.reader.UnreadRune()
}

// SkipWhitespace skips whitespace characters
func (p *Parser) SkipWhitespace() error {
	for {
		r, err := p.GetChar()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			return p.UngetChar(r)
		}
	}
}

// NextToken returns the next token in the source
func (p *Parser) NextToken() (string, error) {
	if err := p.SkipWhitespace(); err != nil && err != io.EOF {
		return "", err
	}

	var buf bytes.Buffer
	for {
		r, err := p.GetChar()
		if err == io.EOF {
			// Return what we have so far
			return buf.String(), nil
		}
		if err != nil {
			return "", err
		}
		if unicode.IsSpace(r) {
			break
		}
		buf.WriteRune(r)
	}

	return buf.String(), nil
}

// SkipLine skips to the end of the current line
func (p *Parser) SkipLine() error {
	for {
		r, err := p.GetChar()
		if err != nil {
			return err
		}
		if r == '\n' {
			return nil
		}
	}
}
