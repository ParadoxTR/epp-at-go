package epp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	ierr "github.com/ParadoxTR/epp-at-go/internal/errors"
)

type CreateContactRequest struct {
	XMLName xml.Name             `xml:"epp"`
	Xmlns   string               `xml:"xmlns,attr"`
	Command CreateContactCommand `xml:"command"`
}

type CreateContactCommand struct {
	Create    CreateContact     `xml:"create"`
	Extension *CommandExtension `xml:"extension,omitempty"`
	ClTRID    string            `xml:"clTRID"`
}

type CreateContact struct {
	XMLName       xml.Name      `xml:"create"`
	ContactCreate ContactCreate `xml:"contact:create"`
}

type ContactCreate struct {
	Disclose   *ContactDisclose  `xml:"contact:disclose,omitempty"`
	XMLName    xml.Name          `xml:"contact:create"`
	Xmlns      string            `xml:"xmlns:contact,attr"`
	ID         string            `xml:"contact:id"`
	Voice      string            `xml:"contact:voice,omitempty"`
	Fax        string            `xml:"contact:fax,omitempty"`
	Email      string            `xml:"contact:email"`
	AuthInfo   ContactAuthInfo   `xml:"contact:authInfo"`
	PostalInfo ContactPostalInfo `xml:"contact:postalInfo"`
}

type CommandExtension struct {
	AtExt   *AtContactExtension `xml:"at-ext-contact:create,omitempty"`
	XMLName xml.Name            `xml:"extension"`
}

type AtContactExtension struct {
	XMLName xml.Name `xml:"at-ext-contact:create"`
	Xmlns   string   `xml:"xmlns:at-ext-contact,attr"`
	Type    string   `xml:"at-ext-contact:type"`
}

type CreateContactResponse struct {
	XMLName   xml.Name                  `xml:"epp"`
	Result    Result                    `xml:"response>result"`
	ResData   CreateContactResponseData `xml:"response>resData"`
	Extension *ResponseExtension        `xml:"response>extension,omitempty"`
	TrID      TrID                      `xml:"response>trID"`
}

type ResponseExtension struct {
	Conditions *Conditions `xml:"conditions,omitempty"`
	XMLName    xml.Name    `xml:"extension"`
}

type Conditions struct {
	XMLName   xml.Name    `xml:"conditions"`
	Xmlns     string      `xml:"xmlns,attr"`
	Condition []Condition `xml:"condition"`
}

type Condition struct {
	Msg     string `xml:"msg"`
	Details string `xml:"details"`
}

type CreateContactResponseData struct {
	CreData CreateContactData `xml:"creData"`
}

type CreateContactData struct {
	XMLName xml.Name `xml:"creData"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      string   `xml:"id"`
	CrDate  string   `xml:"crDate"`
}

func (c *Client) CreateContact(contact *Contact) (*CreateContactResponse, error) {
	var extension *CommandExtension
	if contact.Type != "" {
		extension = &CommandExtension{
			AtExt: &AtContactExtension{
				XMLName: xml.Name{Local: "at-ext-contact:create"},
				Xmlns:   "http://www.nic.at/xsd/at-ext-contact-1.0",
				Type:    contact.Type,
			},
		}
	}

	createReq := CreateContactRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: CreateContactCommand{
			Create: CreateContact{
				XMLName: xml.Name{Local: "create"},
				ContactCreate: ContactCreate{
					XMLName:    xml.Name{Local: "contact:create"},
					Xmlns:      "urn:ietf:params:xml:ns:contact-1.0",
					ID:         contact.ID,
					PostalInfo: contact.PostalInfo,
					Voice:      contact.Voice,
					Fax:        contact.Fax,
					Email:      contact.Email,
					AuthInfo:   ContactAuthInfo{Pw: ""}, // Empty password for Austrian EPP
					Disclose:   contact.Disclose,
				},
			},
			Extension: extension,
			ClTRID:    generateTransactionID(),
		},
	}

	// Normalize street address lines to meet NIC.AT constraints.
	// NIC.AT commonly requires a maximum of 3 street lines and ~35 characters per line.
	normalizeContactPostalInfo(&createReq.Command.Create.ContactCreate.PostalInfo)

	requestXML, err := xml.Marshal(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create contact request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send create contact request: %w", err)
	}

	var response CreateContactResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create contact response: %w", err)
	}

	if response.Result.Code != "1000" {
		errorMsg := fmt.Sprintf("create contact failed: %s - %s", response.Result.Code, response.Result.Msg)

		if response.Extension != nil && response.Extension.Conditions != nil {
			for _, condition := range response.Extension.Conditions.Condition {
				errorMsg += fmt.Sprintf("\nCondition: %s - %s", condition.Msg, condition.Details)
			}
		}

		return nil, errors.New(errorMsg)
	}

	return &response, nil
}

type InfoContactRequest struct {
	XMLName xml.Name           `xml:"epp"`
	Xmlns   string             `xml:"xmlns,attr"`
	Command InfoContactCommand `xml:"command"`
}

type InfoContactCommand struct {
	Info   InfoContact `xml:"info"`
	ClTRID string      `xml:"clTRID"`
}

type InfoContact struct {
	XMLName     xml.Name    `xml:"info"`
	ContactInfo ContactInfo `xml:"contact:info"`
}

type ContactInfo struct {
	XMLName xml.Name `xml:"contact:info"`
	Xmlns   string   `xml:"xmlns:contact,attr"`
	ID      string   `xml:"contact:id"`
}

type InfoContactResponse struct {
	Extension *InfoContactExtension   `xml:"response>extension,omitempty"`
	XMLName   xml.Name                `xml:"epp"`
	Result    Result                  `xml:"response>result"`
	TrID      TrID                    `xml:"response>trID"`
	ResData   InfoContactResponseData `xml:"response>resData"`
}

type InfoContactExtension struct {
	AtExt   *AtContactInfoExtension `xml:"at-ext-contact:infData,omitempty"`
	XMLName xml.Name                `xml:"extension"`
}

type AtContactInfoExtension struct {
	XMLName xml.Name `xml:"at-ext-contact:infData"`
	Xmlns   string   `xml:"xmlns:at-ext-contact,attr"`
	Type    string   `xml:"at-ext-contact:type"`
}

type InfoContactResponseData struct {
	InfData InfoContactData `xml:"infData"`
}

type InfoContactData struct {
	Disclose   *ContactDisclose  `xml:"disclose"`
	PostalInfo ContactPostalInfo `xml:"postalInfo"`
	XMLName    xml.Name          `xml:"infData"`
	UpID       string            `xml:"upID"`
	ROID       string            `xml:"roid"`
	Voice      string            `xml:"voice"`
	Fax        string            `xml:"fax"`
	Email      string            `xml:"email"`
	ClID       string            `xml:"clID"`
	CrID       string            `xml:"crID"`
	CrDate     string            `xml:"crDate"`
	ID         string            `xml:"id"`
	UpDate     string            `xml:"upDate"`
	AuthInfo   string            `xml:"authInfo>pw"`
	Xmlns      string            `xml:"xmlns,attr"`
	Status     []ContactStatus   `xml:"status"`
}

type ContactPostalInfo struct {
	Type string      `xml:"type,attr"`
	Name string      `xml:"contact:name"`
	Org  string      `xml:"contact:org,omitempty"`
	Addr ContactAddr `xml:"contact:addr"`
}

type ContactAddr struct {
	City   string   `xml:"contact:city"`
	SP     string   `xml:"contact:sp,omitempty"`
	PC     string   `xml:"contact:pc"`
	CC     string   `xml:"contact:cc"`
	Street []string `xml:"contact:street"`
}

func normalizeContactPostalInfo(pi *ContactPostalInfo) {
	if pi.Addr.Street == nil {
		return
	}

	const maxLines = 3
	const maxLineLen = 35

	var words []string
	for _, line := range pi.Addr.Street {
		words = append(words, strings.Fields(line)...)
	}

	var lines []string
	var cur strings.Builder

	flush := func() {
		if cur.Len() > 0 {
			lines = append(lines, cur.String())
			cur.Reset()
		}
	}

	for _, w := range words {
		if utf8.RuneCountInString(w) > maxLineLen {
			if cur.Len() > 0 {
				flush()
			}
			rs := []rune(w)
			w = string(rs[:maxLineLen])
		}

		if cur.Len() == 0 {
			cur.WriteString(w)
			continue
		}

		tentative := cur.String() + " " + w
		if utf8.RuneCountInString(tentative) <= maxLineLen {
			cur.WriteString(" ")
			cur.WriteString(w)
			continue
		}

		flush()
		cur.WriteString(w)
		if len(lines) >= maxLines {
			break
		}
	}

	if len(lines) < maxLines {
		flush()
	}

	if len(lines) > maxLines {
		kept := lines[:maxLines]
		overflow := strings.Join(lines[maxLines-1:], " ")
		if utf8.RuneCountInString(overflow) > maxLineLen {
			rs := []rune(overflow)
			overflow = string(rs[:maxLineLen])
		}
		kept[maxLines-1] = overflow
		lines = kept
	}

	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}

	pi.Addr.Street = lines
}

func (c *Client) InfoContact(contactID string) (*InfoContactResponse, error) {
	infoReq := InfoContactRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: InfoContactCommand{
			Info: InfoContact{
				XMLName: xml.Name{Local: "info"},
				ContactInfo: ContactInfo{
					XMLName: xml.Name{Local: "contact:info"},
					Xmlns:   "urn:ietf:params:xml:ns:contact-1.0",
					ID:      contactID,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(infoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal info contact request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send info contact request: %w", err)
	}

	var response InfoContactResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal info contact response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("info contact failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type UpdateContactRequest struct {
	XMLName xml.Name             `xml:"epp"`
	Xmlns   string               `xml:"xmlns,attr"`
	Command UpdateContactCommand `xml:"command"`
}

type UpdateContactCommand struct {
	Update UpdateContact `xml:"update"`
	ClTRID string        `xml:"clTRID"`
}

type UpdateContact struct {
	XMLName       xml.Name      `xml:"update"`
	ContactUpdate ContactUpdate `xml:"contact:update"`
}

type ContactUpdate struct {
	XMLName xml.Name          `xml:"contact:update"`
	Xmlns   string            `xml:"xmlns:contact,attr"`
	ID      string            `xml:"contact:id"`
	Add     *ContactUpdateAdd `xml:"contact:add,omitempty"`
	Rem     *ContactUpdateRem `xml:"contact:rem,omitempty"`
	Chg     *ContactUpdateChg `xml:"contact:chg,omitempty"`
}

type ContactUpdateAdd struct {
	Status []ContactStatus `xml:"contact:status,omitempty"`
}

type ContactUpdateRem struct {
	Status []ContactStatus `xml:"contact:status,omitempty"`
}

type ContactUpdateChg struct {
	PostalInfo *ContactPostalInfo `xml:"contact:postalInfo,omitempty"`
	Disclose   *ContactDisclose   `xml:"contact:disclose,omitempty"`
	Voice      string             `xml:"contact:voice,omitempty"`
	Fax        string             `xml:"contact:fax,omitempty"`
	Email      string             `xml:"contact:email,omitempty"`
	AuthInfo   *ContactAuthInfo   `xml:"contact:authInfo,omitempty"`
}

func (c *Client) UpdateContact(
	contactID string,
	add *ContactUpdateAdd,
	rem *ContactUpdateRem,
	chg *ContactUpdateChg,
) (*Response, error) {
	updateReq := UpdateContactRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: UpdateContactCommand{
			Update: UpdateContact{
				XMLName: xml.Name{Local: "update"},
				ContactUpdate: ContactUpdate{
					XMLName: xml.Name{Local: "contact:update"},
					Xmlns:   "urn:ietf:params:xml:ns:contact-1.0",
					ID:      contactID,
					Add:     add,
					Rem:     rem,
					Chg:     chg,
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	// Normalize postal info if present in chg
	if chg != nil && chg.PostalInfo != nil {
		normalizeContactPostalInfo(chg.PostalInfo)
	}

	var response Response

	requestXML, err := xml.Marshal(updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update contact request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send update contact request: %w", err)
	}

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update contact response: %w", err)
	}

	if !ierr.IsSuccessCode(response.Result.Code) {
		return nil, ierr.NewEPPError(response.Result.Code, response.Result.Msg, "update contact operation failed")
	}

	return &response, nil
}

type DeleteContactRequest struct {
	XMLName xml.Name             `xml:"epp"`
	Xmlns   string               `xml:"xmlns,attr"`
	Command DeleteContactCommand `xml:"command"`
}

type DeleteContactCommand struct {
	Delete DeleteContact `xml:"delete"`
	ClTRID string        `xml:"clTRID"`
}

type DeleteContact struct {
	XMLName xml.Name `xml:"delete"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      string   `xml:"id"`
}

func (c *Client) DeleteContact(contactID string) (*Response, error) {
	deleteReq := DeleteContactRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: DeleteContactCommand{
			Delete: DeleteContact{
				XMLName: xml.Name{Local: "delete"},
				Xmlns:   "urn:ietf:params:xml:ns:contact-1.0",
				ID:      contactID,
			},
			ClTRID: generateTransactionID(),
		},
	}

	var response Response
	requestXML, err := xml.Marshal(deleteReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delete contact request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send delete contact request: %w", err)
	}

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal delete contact response: %w", err)
	}

	if response.Result.Code != "1000" {
		return nil, fmt.Errorf("delete contact failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}
