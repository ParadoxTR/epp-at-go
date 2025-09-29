package epp

import (
	"encoding/xml"
	"fmt"
)

type WithdrawRequest struct {
	XMLName xml.Name        `xml:"epp"`
	Xmlns   string          `xml:"xmlns,attr"`
	Command WithdrawCommand `xml:"command"`
}

type WithdrawCommand struct {
	Withdraw  WithdrawDomainData `xml:"withdraw"`
	Extension *WithdrawExtension `xml:"extension,omitempty"`
	ClTRID    string             `xml:"clTRID"`
}

type WithdrawDomainData struct {
	XMLName xml.Name       `xml:"withdraw"`
	Domain  WithdrawDomain `xml:"domain:withdraw"`
}

type WithdrawDomain struct {
	XMLName xml.Name `xml:"domain:withdraw"`
	Xmlns   string   `xml:"xmlns:domain,attr"`
	Name    string   `xml:"domain:name"`
}

type WithdrawExtension struct {
	XMLName xml.Name             `xml:"extension"`
	AtExt   *AtWithdrawExtension `xml:"at-ext-epp:withdraw,omitempty"`
}

type AtWithdrawExtension struct {
	XMLName xml.Name `xml:"at-ext-epp:withdraw"`
	Xmlns   string   `xml:"xmlns:at-ext-epp,attr"`
	Domain  string   `xml:"at-ext-epp:domain"`
}

type WithdrawResponse struct {
	XMLName xml.Name `xml:"epp"`
	Result  Result   `xml:"response>result"`
	TrID    TrID     `xml:"response>trID"`
}

func (c *Client) WithdrawDomainProper(domainName string) (*WithdrawResponse, error) {
	withdrawReq := WithdrawRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: WithdrawCommand{
			Withdraw: WithdrawDomainData{
				XMLName: xml.Name{Local: "withdraw"},
				Domain: WithdrawDomain{
					XMLName: xml.Name{Local: "domain:withdraw"},
					Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
					Name:    domainName,
				},
			},
			Extension: &WithdrawExtension{
				AtExt: &AtWithdrawExtension{
					XMLName: xml.Name{Local: "at-ext-epp:withdraw"},
					Xmlns:   "http://www.nic.at/xsd/at-ext-epp-1.0",
					Domain:  domainName,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(withdrawReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal withdraw request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send withdraw request: %w", err)
	}

	var response WithdrawResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal withdraw response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("withdraw failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}
