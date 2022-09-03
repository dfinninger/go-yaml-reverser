package reverser

import (
	"bufio"
	"strings"
)

const (
	DOCUMENT_START = "---"
	DOCUMENT_END   = "..."
)

// Document represents a Yaml document within a string
// Since this library is only focused on reversing a Yaml document we do not care
// about the content, just the delimiters. Since some versions of Yaml have optional
// delimiters (and handling a "bare" document), we track the delimiters so that we
// can put them back when we are finished.
type Document struct {
	prelude        string
	startDelimiter string
	content        string
	endDelimiter   string
}

func (d Document) String() string {
	parts := []string{d.prelude, d.startDelimiter, d.content, d.endDelimiter}
	var filteredParts []string
	for _, p := range parts {
		if p == "" {
			continue
		}
		filteredParts = append(filteredParts, p)
	}
	return strings.Join(filteredParts, "\n")
}

// DocumentStream is the content of a single file's Yaml documents
type DocumentStream []Document

func (ds DocumentStream) String() string {
	var parts []string
	for _, doc := range ds {
		s := doc.String()
		if s == "" {
			continue
		}
		parts = append(parts, doc.String())
	}
	return strings.Join(parts, "\n")
}

func (ds DocumentStream) Reverse() DocumentStream {
	var reversed DocumentStream
	l := len(ds)
	for i, _ := range ds {
		reversed = append(reversed, ds[l-1-i])
	}
	return reversed
}

func reverseStream(scanner *bufio.Scanner) (DocumentStream, error) {
	scanner.Split(bufio.ScanLines)

	docStream := DocumentStream{}
	doc := Document{}
	for scanner.Scan() {
		// Document has not officially "started", so we fill the prelude
		if doc.startDelimiter == "" {
			if scanner.Text() != DOCUMENT_START {
				doc.prelude += scanner.Text()
			} else {
				doc.startDelimiter = DOCUMENT_START
			}
			// Document has "started" so fill the content
		} else {
			switch scanner.Text() {
			case DOCUMENT_END:
				// Document has "finished", so append it to the stream and start a new document
				doc.endDelimiter = DOCUMENT_END
				docStream = append(docStream, doc)
				doc = Document{}
			case DOCUMENT_START:
				// Document has "finished" by starting a new document
				docStream = append(docStream, doc)
				doc = Document{}
				doc.startDelimiter = DOCUMENT_START
			default:
				doc.content += scanner.Text()
			}
		}
	}

	// EOF, append what we have to the stream
	docStream = append(docStream, doc)

	return docStream.Reverse(), nil
}
