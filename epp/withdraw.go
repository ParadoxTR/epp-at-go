package epp

import (
	"encoding/xml"
	"fmt"
)

type WithdrawRequest struct {
	XMLName   xml.Name             `xml:"epp"`
	Xmlns     string               `xml:"xmlns,attr"`
	Extension WithdrawEPPExtension `xml:"extension"`
}

type WithdrawEPPExtension struct {
	Command WithdrawCommand `xml:"command"`
}

type WithdrawCommand struct {
	XMLName  xml.Name           `xml:"command"`
	Xmlns    string             `xml:"xmlns,attr"`
	Withdraw WithdrawDomainData `xml:"withdraw"`
	ClTRID   string             `xml:"clTRID"`
}

type WithdrawDomainData struct {
	XMLName xml.Name       `xml:"withdraw"`
	Domain  WithdrawDomain `xml:"domain:withdraw"`
}

type WithdrawDomain struct {
	XMLName    xml.Name            `xml:"domain:withdraw"`
	Xmlns      string              `xml:"xmlns:domain,attr"`
	Name       string              `xml:"domain:name"`
	ZoneDelete *WithdrawZoneDelete `xml:"domain:zd,omitempty"`
}

type WithdrawZoneDelete struct {
	Value int `xml:"value,attr"`
}

type WithdrawResponse = Response

func (c *Client) WithdrawDomainProper(domainName string) (*WithdrawResponse, error) {
	return c.withdrawDomain(domainName, nil)
}

func (c *Client) WithdrawDomainWithZoneDelete(domainName string, zoneDelete bool) (*Response, error) {
	value := 0
	if zoneDelete {
		value = 1
	}
	return c.withdrawDomain(domainName, &value)
}

func (c *Client) withdrawDomain(domainName string, zoneDelete *int) (*Response, error) {
	var zd *WithdrawZoneDelete
	if zoneDelete != nil {
		zd = &WithdrawZoneDelete{Value: *zoneDelete}
	}

	withdrawReq := WithdrawRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Extension: WithdrawEPPExtension{
			Command: WithdrawCommand{
				XMLName: xml.Name{Local: "command"},
				Xmlns:   "http://www.nic.at/xsd/at-ext-epp-1.0",
				Withdraw: WithdrawDomainData{
					XMLName: xml.Name{Local: "withdraw"},
					Domain: WithdrawDomain{
						XMLName:    xml.Name{Local: "domain:withdraw"},
						Xmlns:      "http://www.nic.at/xsd/at-ext-domain-1.0",
						Name:       domainName,
						ZoneDelete: zd,
					},
				},
				ClTRID: generateTransactionID(),
			},
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

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal withdraw response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("withdraw failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}
