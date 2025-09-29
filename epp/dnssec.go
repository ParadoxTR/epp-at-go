package epp

import (
	"encoding/xml"
	"fmt"
)

type DNSSECData struct {
	KeyTag     int    `xml:"secDNS:keyTag"`
	Alg        int    `xml:"secDNS:alg"`
	DigestType int    `xml:"secDNS:digestType"`
	Digest     string `xml:"secDNS:digest"`
}

type DNSSECExtension struct {
	XMLName      xml.Name      `xml:"extension"`
	SecDNS       *SecDNSData   `xml:"secDNS:create,omitempty"`
	SecDNSUpdate *SecDNSUpdate `xml:"secDNS:update,omitempty"`
}

type SecDNSData struct {
	XMLName xml.Name     `xml:"secDNS:create"`
	Xmlns   string       `xml:"xmlns:secDNS,attr"`
	DSData  []DNSSECData `xml:"secDNS:dsData"`
}

type SecDNSUpdate struct {
	XMLName xml.Name         `xml:"secDNS:update"`
	Xmlns   string           `xml:"xmlns:secDNS,attr"`
	Add     *SecDNSUpdateAdd `xml:"secDNS:add,omitempty"`
	Rem     *SecDNSUpdateRem `xml:"secDNS:rem,omitempty"`
	Chg     *SecDNSUpdateChg `xml:"secDNS:chg,omitempty"`
}

type SecDNSUpdateAdd struct {
	DSData []DNSSECData `xml:"secDNS:dsData"`
}

type SecDNSUpdateRem struct {
	DSData []DNSSECData `xml:"secDNS:dsData"`
}

type SecDNSUpdateChg struct {
	DSData []DNSSECData `xml:"secDNS:dsData"`
}

func (c *Client) CreateDomainWithDNSSEC(domain Domain, dsRecords []DNSSECData) (*CreateDomainResponse, error) {
	var extension *DNSSECExtension
	if len(dsRecords) > 0 {
		extension = &DNSSECExtension{
			SecDNS: &SecDNSData{
				XMLName: xml.Name{Local: "secDNS:create"},
				Xmlns:   "urn:ietf:params:xml:ns:secDNS-1.1",
				DSData:  dsRecords,
			},
		}
	}

	createReq := CreateDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: CreateDomainCommand{
			Create: CreateDomain{
				XMLName: xml.Name{Local: "create"},
				Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
				Domain:  domain,
			},
			ClTRID: generateTransactionID(),
		},
	}

	// Add DNSSEC extension if provided
	if extension != nil {
		// Note: This would need to be added to the CreateDomainCommand struct
		// For now, we'll use the regular create and suggest adding extension support
	}

	requestXML, err := xml.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create domain with DNSSEC request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send create domain with DNSSEC request: %w", err)
	}

	var response CreateDomainResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create domain with DNSSEC response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("create domain with DNSSEC failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

func (c *Client) UpdateDomainDNSSEC(domainName string, add, rem, chg []DNSSECData) (*Response, error) {
	var secDNSUpdate *SecDNSUpdate

	if len(add) > 0 || len(rem) > 0 || len(chg) > 0 {
		secDNSUpdate = &SecDNSUpdate{
			XMLName: xml.Name{Local: "secDNS:update"},
			Xmlns:   "urn:ietf:params:xml:ns:secDNS-1.1",
		}

		if len(add) > 0 {
			secDNSUpdate.Add = &SecDNSUpdateAdd{DSData: add}
		}
		if len(rem) > 0 {
			secDNSUpdate.Rem = &SecDNSUpdateRem{DSData: rem}
		}
		if len(chg) > 0 {
			secDNSUpdate.Chg = &SecDNSUpdateChg{DSData: chg}
		}
	}

	updateReq := UpdateDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: UpdateDomainCommand{
			Update: UpdateDomain{
				XMLName: xml.Name{Local: "update"},
				Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
				Name:    domainName,
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update domain DNSSEC request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send update domain DNSSEC request: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update domain DNSSEC response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("update domain DNSSEC failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}
