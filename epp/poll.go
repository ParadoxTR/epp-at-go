package epp

import (
	"encoding/xml"
	"fmt"
)

type PollRequest struct {
	XMLName xml.Name    `xml:"epp"`
	Xmlns   string      `xml:"xmlns,attr"`
	Command PollCommand `xml:"command"`
}

type PollCommand struct {
	Poll   Poll   `xml:"poll"`
	ClTRID string `xml:"clTRID"`
}

type Poll struct {
	Op    string `xml:"op,attr"`
	MsgID string `xml:"msgID,attr,omitempty"`
}

type PollResponse struct {
	XMLName xml.Name          `xml:"epp"`
	Result  Result            `xml:"response>result"`
	MsgQ    *PollMessageQueue `xml:"response>msgQ,omitempty"`
	ResData *PollResponseData `xml:"response>resData,omitempty"`
	TrID    TrID              `xml:"response>trID"`
}

type PollMessageQueue struct {
	Count int    `xml:"count,attr"`
	ID    string `xml:"id,attr"`
	QDate string `xml:"qDate"`
	Msg   string `xml:"msg"`
}

type PollResponseData struct {
	XMLName xml.Name `xml:"resData"`
	Content string   `xml:",innerxml"`
}

func (c *Client) PollMessage() (*PollResponse, error) {
	pollReq := PollRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: PollCommand{
			Poll: Poll{
				Op: "req",
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(pollReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal poll request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send poll request: %w", err)
	}

	var response PollResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal poll response: %w", err)
	}

	if response.Result.Code != "1000" && response.Result.Code != "1300" {
		return nil, fmt.Errorf("poll failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

func (c *Client) AckPollMessage(msgID string) (*PollResponse, error) {
	pollReq := PollRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: PollCommand{
			Poll: Poll{
				Op:    "ack",
				MsgID: msgID,
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(pollReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal poll ack request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return nil, fmt.Errorf("failed to send poll ack request: %w", err)
	}

	var response PollResponse
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal poll ack response: %w", err)
	}

	if response.Result.Code != "1000" && response.Result.Code != "1300" {
		return nil, fmt.Errorf("poll ack failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return &response, nil
}

type ChangePasswordRequest struct {
	XMLName xml.Name              `xml:"epp"`
	Xmlns   string                `xml:"xmlns,attr"`
	Command ChangePasswordCommand `xml:"command"`
}

type ChangePasswordCommand struct {
	Login  ChangePassword `xml:"login"`
	ClTRID string         `xml:"clTRID"`
}

type ChangePassword struct {
	ClID    string        `xml:"clID"`
	Pw      string        `xml:"pw"`
	NewPw   string        `xml:"newPW"`
	Options LoginOptions  `xml:"options"`
	Svcs    LoginServices `xml:"svcs"`
}

func (c *Client) ChangePassword(newPassword string) error {
	changeReq := ChangePasswordRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: ChangePasswordCommand{
			Login: ChangePassword{
				ClID:  c.username,
				Pw:    c.password,
				NewPw: newPassword,
				Options: LoginOptions{
					Version: "1.0",
					Lang:    "en",
				},
				Svcs: LoginServices{
					ObjURI: []string{
						"urn:ietf:params:xml:ns:domain-1.0",
						"urn:ietf:params:xml:ns:contact-1.0",
					},
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(changeReq)
	if err != nil {
		return fmt.Errorf("failed to marshal change password request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return fmt.Errorf("failed to send change password request: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return fmt.Errorf("failed to unmarshal change password response: %w", err)
	}

	if response.Result.Code != "1000" {
		return fmt.Errorf("change password failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	c.password = newPassword

	return nil
}
