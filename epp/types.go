package epp

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"
)

func generateTransactionID() string {
	return fmt.Sprintf("epp-go-%d-%d", time.Now().Unix(), rand.Int63n(10000))
}

type Response struct {
	XMLName xml.Name `xml:"epp"`
	Result  Result   `xml:"response>result"`
	TrID    TrID     `xml:"response>trID"`
}

type Result struct {
	Code string `xml:"code,attr"`
	Msg  string `xml:"msg"`
}

type TrID struct {
	ClTRID string `xml:"clTRID"` // Client transaction ID
	SvTRID string `xml:"svTRID"` // Server transaction ID
}

type LoginRequest struct {
	XMLName xml.Name     `xml:"epp"`
	Xmlns   string       `xml:"xmlns,attr"`
	Command LoginCommand `xml:"command"`
}

type LoginCommand struct {
	Login  Login  `xml:"login"`
	ClTRID string `xml:"clTRID"`
}

type Login struct {
	ClID    string        `xml:"clID"`    // Client identifier (username)
	Pw      string        `xml:"pw"`      // Password
	Options LoginOptions  `xml:"options"` // Protocol options
	Svcs    LoginServices `xml:"svcs"`    // Available services
}

type LoginOptions struct {
	Version string `xml:"version"` // EPP protocol version
	Lang    string `xml:"lang"`    // Language code
}

type LoginServices struct {
	ObjURI       []string               `xml:"objURI"`                 // Object URIs
	SvcExtension *LoginServiceExtension `xml:"svcExtension,omitempty"` // Service extensions
}

type LoginServiceExtension struct {
	ExtURI []string `xml:"extURI"` // Extension URIs
}

type LogoutRequest struct {
	XMLName xml.Name      `xml:"epp"`
	Xmlns   string        `xml:"xmlns,attr"`
	Command LogoutCommand `xml:"command"`
}

type LogoutCommand struct {
	Logout struct{} `xml:"logout"`
	ClTRID string   `xml:"clTRID"`
}

type HelloRequest struct {
	XMLName xml.Name `xml:"epp"`
	Xmlns   string   `xml:"xmlns,attr"`
	Hello   struct{} `xml:"hello"`
}

type Contact struct {
	ID         string            `xml:"contact:id,omitempty"`
	PostalInfo ContactPostalInfo `xml:"contact:postalInfo"`
	Voice      string            `xml:"contact:voice,omitempty"`
	Fax        string            `xml:"contact:fax,omitempty"`
	Email      string            `xml:"contact:email"`
	AuthInfo   ContactAuthInfo   `xml:"contact:authInfo,omitempty"`
	Status     []ContactStatus   `xml:"contact:status,omitempty"`
	Disclose   *ContactDisclose  `xml:"contact:disclose,omitempty"`
	Type       string            `xml:"-"` // Austrian EPP extension: privateperson, organisation, role
}

type ContactAuthInfo struct {
	Pw string `xml:"contact:pw"`
}

type ContactStatus struct {
	Status string `xml:"s,attr"`
	Text   string `xml:",chardata"`
}

type ContactDisclose struct {
	Flag  int    `xml:"flag,attr"`
	Voice string `xml:"voice,omitempty"`
	Fax   string `xml:"fax,omitempty"`
	Email string `xml:"email,omitempty"`
}

type Domain struct {
	Name        string          `xml:"name"`
	Registrant  string          `xml:"registrant,omitempty"`
	Contacts    []DomainContact `xml:"contact,omitempty"`
	Nameservers []string        `xml:"ns>hostObj,omitempty"`
	AuthInfo    string          `xml:"authInfo>pw,omitempty"`
	Status      []DomainStatus  `xml:"status,omitempty"`
	Period      *Period         `xml:"period,omitempty"`
}

type DomainContact struct {
	Type string `xml:"type,attr"` // Contact type: admin, tech, billing
	ID   string `xml:",chardata"`
}

type DomainStatus struct {
	Status string `xml:"s,attr"`
	Text   string `xml:",chardata"`
}

type Period struct {
	Unit  string `xml:"unit,attr"` // Time unit: y (years), m (months)
	Value int    `xml:",chardata"`
}

type CheckDomainRequest struct {
	XMLName xml.Name           `xml:"epp"`
	Xmlns   string             `xml:"xmlns,attr"`
	Command CheckDomainCommand `xml:"command"`
}

type CheckDomainCommand struct {
	Check  CheckDomain `xml:"check"`
	ClTRID string      `xml:"clTRID"`
}

type CheckDomain struct {
	XMLName     xml.Name    `xml:"check"`
	DomainCheck DomainCheck `xml:"domain:check"`
}

type DomainCheck struct {
	XMLName xml.Name `xml:"domain:check"`
	Xmlns   string   `xml:"xmlns:domain,attr"`
	Names   []string `xml:"domain:name"`
}

type CheckDomainResponse struct {
	XMLName xml.Name                `xml:"epp"`
	Result  Result                  `xml:"response>result"`
	ResData CheckDomainResponseData `xml:"response>resData"`
	TrID    TrID                    `xml:"response>trID"`
}

type CheckDomainResponseData struct {
	ChkData CheckDomainData `xml:"chkData"`
}

type CheckDomainData struct {
	XMLName xml.Name              `xml:"chkData"`
	Xmlns   string                `xml:"xmlns,attr"`
	Names   []CheckDomainNameData `xml:"cd"`
}

type CheckDomainNameData struct {
	Name   CheckDomainName `xml:"name"`
	Reason string          `xml:"reason,omitempty"`
}

type CheckDomainName struct {
	Name      string `xml:",chardata"`
	Available string `xml:"avail,attr"` // "1" for available, "0" for unavailable
}
