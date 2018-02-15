package xmpp

import "encoding/xml"

// CompressionFeature implements XEP-0138: Stream Compression - 10.1 Stream Feature
type CompressionFeature struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	Methods []string `xml:"method"`
}

// CompressionCompress implements XEP-0138: Stream Compression - 10.2 Protocol Namespace
type CompressionCompress struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compress"`
	Methods []string `xml:"method"`
}

// CompressionCompressed implements XEP-0138: Stream Compression - 10.2 Protocol Namespace
type CompressionCompressed struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/compress compressed"`
}

// CompressionFailure implements XEP-0138: Stream Compression - 10.2 Protocol Namespace
type CompressionFailure struct {
	XMLName           xml.Name  `xml:"http://jabber.org/protocol/compress failure"`
	SetupFailed       *xml.Name `xml:"http://jabber.org/protocol/compress setup-failed"`
	ProcessingFailed  *xml.Name `xml:"http://jabber.org/protocol/compress processing-failed"`
	UnsupportedFailed *xml.Name `xml:"http://jabber.org/protocol/compress unsupported-failed"`
	Text              *Text
	StanzaErrorGroup
}
