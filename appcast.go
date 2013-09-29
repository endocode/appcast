package appcast

import (
	"encoding/xml"
	"github.com/c9s/rss"
	"io/ioutil"
	"net/http"
)

// XXX: better solution? use lower-case, because we need to encode it with lowercase
type Appcast struct {
	XMLName      xml.Name `xml:"rss"`
	XmlNSSparkle string   `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle sparkle,attr"`
	XmlNSDC      string   `xml:"http://purl.org/dc/elements/1.1 dc,attr"`
	Channel      Channel  `xml:"channel"`
	rss.RSS
	/*
		<rss version="2.0"
			xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle"
			xmlns:dc="http://purl.org/dc/elements/1.1/">
	*/
}

func (self *Appcast) MarshalIndent() ([]byte, error) {
	content, err := xml.MarshalIndent(self, "", "  ")
	if err != nil {
		return nil, err
	}
	return content, nil
}

/*
Write appcast XML content to file.
*/
func (self *Appcast) WriteFile(path string) error {
	self.Version = "2.0"
	self.XmlNSSparkle = "http://www.andymatuschak.org/xml-namespaces/sparkle"
	self.XmlNSDC = "http://purl.org/dc/elements/1.1/"
	content, err := self.MarshalIndent()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func New() *Appcast {
	appcast := Appcast{}
	appcast.Version = "2.0"
	appcast.XmlNSSparkle = "http://www.andymatuschak.org/xml-namespaces/sparkle"
	appcast.XmlNSDC = "http://purl.org/dc/elements/1.1/"
	return &appcast
}

/*
Parse appcast XML content from bytes
*/
func ParseContent(text []byte) (*Appcast, error) {
	var appcast = New()
	err := xml.Unmarshal(text, appcast)
	if err != nil {
		return nil, err
	}
	return appcast, nil
}

func ReadFile(file string) (*Appcast, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ParseContent(text)
}

func ReadUrl(url string) (*Appcast, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	text, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return ParseContent(text)
}
