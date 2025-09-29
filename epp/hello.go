package epp

import (
	"encoding/xml"
	"fmt"
)

type HelloResponse struct {
	XMLName  xml.Name `xml:"epp"`
	Greeting Greeting `xml:"greeting"`
}

type Greeting struct {
	SvID    string  `xml:"svID"`
	SvDate  string  `xml:"svDate"`
	SvcMenu SvcMenu `xml:"svcMenu"`
	DCP     DCP     `xml:"dcp"`
}

type SvcMenu struct {
	Version []string `xml:"version"`
	Lang    []string `xml:"lang"`
	ObjURI  []string `xml:"objURI"`
	SvcExt  []string `xml:"svcExtURI,omitempty"`
}

type DCP struct {
	Access    Access      `xml:"access"`
	Statement []Statement `xml:"statement"`
}

type Access struct {
	All string `xml:"all,omitempty"`
}

type Statement struct {
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
	Retention Retention `xml:"retention"`
}

type Purpose struct {
	Admin   string `xml:"admin,omitempty"`
	Contact string `xml:"contact,omitempty"`
	Prov    string `xml:"prov,omitempty"`
}

type Recipient struct {
	Ours   string `xml:"ours,omitempty"`
	Public string `xml:"public,omitempty"`
}

type Retention struct {
	Stated string `xml:"stated,omitempty"`
}

func (c *Client) Hello() (*HelloResponse, error) {
	helloReq := HelloRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Hello:   struct{}{},
	}

	requestXML, err := xml.Marshal(helloReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal hello request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send hello request: %w", err)
	}

	var response HelloResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hello response: %w", err)
	}

	return &response, nil
}
