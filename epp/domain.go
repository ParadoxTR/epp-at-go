package epp

import (
	"encoding/xml"
	"fmt"

	"github.com/ParadoxTR/epp-at-go/internal/errors"
	"github.com/ParadoxTR/epp-at-go/internal/validator"
)

func (c *Client) CheckDomain(domains []string) (*CheckDomainResponse, error) {
	if len(domains) == 0 {
		return nil, fmt.Errorf("at least one domain name is required")
	}

	for _, domain := range domains {
		if err := validator.ValidateDomainName(domain); err != nil {
			return nil, fmt.Errorf("invalid domain name '%s': %w", domain, err)
		}
	}

	checkReq := CheckDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: CheckDomainCommand{
			Check: CheckDomain{
				XMLName: xml.Name{Local: "check"},
				DomainCheck: DomainCheck{
					XMLName: xml.Name{Local: "domain:check"},
					Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
					Names:   domains,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(checkReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal domain check request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send domain check request: %w", err)
	}

	var response CheckDomainResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal domain check response: %w", err)
	}

	if !errors.IsSuccessCode(response.Result.Code) {
		return nil, errors.NewEPPError(response.Result.Code, response.Result.Msg, "domain check operation failed")
	}

	return &response, nil
}

type CreateDomainRequest struct {
	XMLName xml.Name            `xml:"epp"`
	Xmlns   string              `xml:"xmlns,attr"`
	Command CreateDomainCommand `xml:"command"`
}

type CreateDomainCommand struct {
	Create    CreateDomain     `xml:"create"`
	Extension *DNSSECExtension `xml:"extension,omitempty"`
	ClTRID    string           `xml:"clTRID"`
}

type CreateDomain struct {
	XMLName xml.Name           `xml:"create"`
	Domain  CreateDomainDetail `xml:"domain:create"`
}

type CreateDomainDetail struct {
	XMLName     xml.Name                 `xml:"domain:create"`
	Xmlns       string                   `xml:"xmlns:domain,attr"`
	Name        string                   `xml:"domain:name"`
	Period      *CreateDomainPeriod      `xml:"domain:period,omitempty"`
	Nameservers *CreateDomainNameservers `xml:"domain:ns,omitempty"`
	Registrant  string                   `xml:"domain:registrant,omitempty"`
	Contacts    []DomainContact          `xml:"domain:contact,omitempty"`
	AuthInfo    *CreateDomainAuthInfo    `xml:"domain:authInfo,omitempty"`
}

type CreateDomainPeriod struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type CreateDomainNameservers struct {
	HostAttrs []CreateDomainHostAttr `xml:"domain:hostAttr"`
}

type CreateDomainHostAttr struct {
	HostName string                 `xml:"domain:hostName"`
	HostAddr []CreateDomainHostAddr `xml:"domain:hostAddr,omitempty"`
}

type CreateDomainHostAddr struct {
	IP   string `xml:"ip,attr"`
	Addr string `xml:",chardata"`
}

type CreateDomainAuthInfo struct {
	Pw string `xml:"domain:pw"`
}

type CreateDomainResponse struct {
	XMLName xml.Name                 `xml:"epp"`
	Result  Result                   `xml:"response>result"`
	ResData CreateDomainResponseData `xml:"response>resData"`
	TrID    TrID                     `xml:"response>trID"`
}

type CreateDomainResponseData struct {
	CreData CreateDomainData `xml:"creData"`
}

type CreateDomainData struct {
	XMLName xml.Name `xml:"creData"`
	Xmlns   string   `xml:"xmlns,attr"`
	Name    string   `xml:"name"`
	CrDate  string   `xml:"crDate"`
	ExDate  string   `xml:"exDate"`
}

func (c *Client) CreateDomain(domain Domain) (*CreateDomainResponse, error) {
	// Convert nameservers to proper structure (NIC.at requires hostAttr format)
	var nameservers *CreateDomainNameservers
	if len(domain.Nameservers) > 0 {
		var hostAttrs []CreateDomainHostAttr
		for _, ns := range domain.Nameservers {
			hostAttrs = append(hostAttrs, CreateDomainHostAttr{
				HostName: ns,
			})
		}
		nameservers = &CreateDomainNameservers{
			HostAttrs: hostAttrs,
		}
	}

	// Convert authInfo to proper structure
	var authInfo *CreateDomainAuthInfo
	if domain.AuthInfo != "" {
		authInfo = &CreateDomainAuthInfo{Pw: domain.AuthInfo}
	}

	createReq := CreateDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: CreateDomainCommand{
			Create: CreateDomain{
				XMLName: xml.Name{Local: "create"},
				Domain: CreateDomainDetail{
					XMLName:     xml.Name{Local: "domain:create"},
					Xmlns:       "urn:ietf:params:xml:ns:domain-1.0",
					Name:        domain.Name,
					Nameservers: nameservers,
					Registrant:  domain.Registrant,
					Contacts:    domain.Contacts,
					AuthInfo:    authInfo,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create domain request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send create domain request: %w", err)
	}

	var response CreateDomainResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create domain response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("create domain failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type InfoDomainRequest struct {
	XMLName xml.Name          `xml:"epp"`
	Xmlns   string            `xml:"xmlns,attr"`
	Command InfoDomainCommand `xml:"command"`
}

type InfoDomainCommand struct {
	Info   InfoDomain `xml:"info"`
	ClTRID string     `xml:"clTRID"`
}

type InfoDomain struct {
	XMLName xml.Name `xml:"info"`
	Xmlns   string   `xml:"xmlns,attr"`
	Name    string   `xml:"name"`
}

type InfoDomainResponse struct {
	XMLName xml.Name               `xml:"epp"`
	Result  Result                 `xml:"response>result"`
	ResData InfoDomainResponseData `xml:"response>resData"`
	TrID    TrID                   `xml:"response>trID"`
}

type InfoDomainResponseData struct {
	InfData InfoDomainData `xml:"infData"`
}

type InfoDomainData struct {
	XMLName     xml.Name        `xml:"infData"`
	Xmlns       string          `xml:"xmlns,attr"`
	Name        string          `xml:"name"`
	ROID        string          `xml:"roid"`
	Status      []DomainStatus  `xml:"status"`
	Registrant  string          `xml:"registrant"`
	Contacts    []DomainContact `xml:"contact"`
	Nameservers []string        `xml:"ns>hostObj"`
	ClID        string          `xml:"clID"`
	CrID        string          `xml:"crID"`
	CrDate      string          `xml:"crDate"`
	ExDate      string          `xml:"exDate"`
	AuthInfo    string          `xml:"authInfo>pw"`
}

func (c *Client) InfoDomain(domainName string) (*InfoDomainResponse, error) {
	infoReq := InfoDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: InfoDomainCommand{
			Info: InfoDomain{
				XMLName: xml.Name{Local: "info"},
				Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
				Name:    domainName,
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(infoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal info domain request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send info domain request: %w", err)
	}

	var response InfoDomainResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info domain response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("info domain failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type UpdateDomainRequest struct {
	XMLName xml.Name            `xml:"epp"`
	Xmlns   string              `xml:"xmlns,attr"`
	Command UpdateDomainCommand `xml:"command"`
}

type UpdateDomainCommand struct {
	Update    UpdateDomain     `xml:"update"`
	Extension *DNSSECExtension `xml:"extension,omitempty"`
	ClTRID    string           `xml:"clTRID"`
}

type UpdateDomain struct {
	XMLName xml.Name           `xml:"update"`
	Domain  UpdateDomainDetail `xml:"domain:update"`
}

type UpdateDomainDetail struct {
	XMLName xml.Name         `xml:"domain:update"`
	Xmlns   string           `xml:"xmlns:domain,attr"`
	Name    string           `xml:"domain:name"`
	Add     *DomainUpdateAdd `xml:"domain:add,omitempty"`
	Rem     *DomainUpdateRem `xml:"domain:rem,omitempty"`
	Chg     *DomainUpdateChg `xml:"domain:chg,omitempty"`
}

type DomainUpdateAdd struct {
	Nameservers []string        `xml:"domain:ns>domain:hostObj,omitempty"`
	Contacts    []DomainContact `xml:"domain:contact,omitempty"`
	Status      []DomainStatus  `xml:"domain:status,omitempty"`
}

type DomainUpdateRem struct {
	Nameservers []string        `xml:"domain:ns>domain:hostObj,omitempty"`
	Contacts    []DomainContact `xml:"domain:contact,omitempty"`
	Status      []DomainStatus  `xml:"domain:status,omitempty"`
}

type DomainUpdateChg struct {
	Registrant string `xml:"domain:registrant,omitempty"`
	AuthInfo   string `xml:"domain:authInfo>domain:pw,omitempty"`
}

func (c *Client) UpdateDomain(domainName string, add *DomainUpdateAdd, rem *DomainUpdateRem, chg *DomainUpdateChg) (*Response, error) {
	updateReq := UpdateDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: UpdateDomainCommand{
			Update: UpdateDomain{
				XMLName: xml.Name{Local: "update"},
				Domain: UpdateDomainDetail{
					XMLName: xml.Name{Local: "domain:update"},
					Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
					Name:    domainName,
					Add:     add,
					Rem:     rem,
					Chg:     chg,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update domain request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send update domain request: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update domain response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("update domain failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type DeleteDomainRequest struct {
	XMLName xml.Name            `xml:"epp"`
	Xmlns   string              `xml:"xmlns,attr"`
	Command DeleteDomainCommand `xml:"command"`
}

type DeleteDomainCommand struct {
	Delete DeleteDomain `xml:"delete"`
	ClTRID string       `xml:"clTRID"`
}

type DeleteDomain struct {
	XMLName xml.Name `xml:"delete"`
	Xmlns   string   `xml:"xmlns,attr"`
	Name    string   `xml:"name"`
}

func (c *Client) DeleteDomain(domainName string) (*Response, error) {
	deleteReq := DeleteDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: DeleteDomainCommand{
			Delete: DeleteDomain{
				XMLName: xml.Name{Local: "delete"},
				Xmlns:   "urn:ietf:params:xml:ns:domain-1.0",
				Name:    domainName,
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(deleteReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delete domain request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send delete domain request: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delete domain response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("delete domain failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type TransferDomainRequest struct {
	XMLName xml.Name              `xml:"epp"`
	Xmlns   string                `xml:"xmlns,attr"`
	Command TransferDomainCommand `xml:"command"`
}

type TransferDomainCommand struct {
	Transfer TransferDomain `xml:"transfer"`
	ClTRID   string         `xml:"clTRID"`
}

type TransferDomain struct {
	Op     string               `xml:"op,attr"`
	Domain TransferDomainDetail `xml:"domain:transfer"`
}

type TransferDomainDetail struct {
	XMLName  xml.Name           `xml:"domain:transfer"`
	Xmlns    string             `xml:"xmlns:domain,attr"`
	Name     string             `xml:"domain:name"`
	AuthInfo *TransferAuthInfo  `xml:"domain:authInfo,omitempty"`
}

type TransferAuthInfo struct {
	Pw string `xml:"domain:pw"`
}

type TransferDomainResponse struct {
	XMLName   xml.Name                   `xml:"epp"`
	Result    Result                     `xml:"response>result"`
	ResData   TransferDomainResponseData `xml:"response>resData"`
	Extension *TransferExtension         `xml:"response>extension"`
	TrID      TrID                       `xml:"response>trID"`
}

type TransferExtension struct {
	KeyDate string `xml:"keydate"`
}

type TransferDomainResponseData struct {
	TrnData TransferDomainData `xml:"trnData"`
}

type TransferDomainData struct {
	XMLName  xml.Name `xml:"trnData"`
	Xmlns    string   `xml:"xmlns,attr"`
	Name     string   `xml:"name"`
	TrStatus string   `xml:"trStatus"`
	ReID     string   `xml:"reID"`
	ReDate   string   `xml:"reDate"`
	AcID     string   `xml:"acID"`
	AcDate   string   `xml:"acDate"`
}

func (c *Client) TransferRequestDomain(domainName, authInfo string) (*TransferDomainResponse, error) {
	return c.transferDomain(domainName, "request", authInfo)
}

func (c *Client) TransferQueryDomain(domainName string) (*TransferDomainResponse, error) {
	return c.transferDomain(domainName, "query", "")
}

func (c *Client) TransferCancelDomain(domainName string) (*TransferDomainResponse, error) {
	return c.transferDomain(domainName, "cancel", "")
}

func (c *Client) transferDomain(domainName, operation, authInfo string) (*TransferDomainResponse, error) {
	var authInfoStruct *TransferAuthInfo
	if authInfo != "" {
		authInfoStruct = &TransferAuthInfo{Pw: authInfo}
	}

	transferReq := TransferDomainRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: TransferDomainCommand{
			Transfer: TransferDomain{
				Op: operation,
				Domain: TransferDomainDetail{
					XMLName:  xml.Name{Local: "domain:transfer"},
					Xmlns:    "urn:ietf:params:xml:ns:domain-1.0",
					Name:     domainName,
					AuthInfo: authInfoStruct,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(transferReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer domain request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send transfer domain request: %w", err)
	}

	var response TransferDomainResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfer domain response: %w", err)
	}

	if response.Result.Code != "1000" && response.Result.Code != "1001" {
		return nil, fmt.Errorf("transfer domain failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

func (c *Client) WithdrawDomain(domainName string) (*Response, error) {

	chg := &DomainUpdateChg{}
	add := &DomainUpdateAdd{
		Status: []DomainStatus{{Status: "clientHold"}},
	}

	return c.UpdateDomain(domainName, add, nil, chg)
}
