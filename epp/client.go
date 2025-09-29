package epp

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"net"
	"time"
)

type Client struct {
	conn     net.Conn
	hostname string
	port     int
	username string
	password string
	timeout  time.Duration
}

type Config struct {
	Hostname string        // EPP server hostname
	Port     int           // EPP server port (typically 700)
	Username string        // EPP account username
	Password string        // EPP account password
	Timeout  time.Duration // Connection timeout duration
}

func NewClient(config Config) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		hostname: config.Hostname,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		timeout:  config.Timeout,
	}
}

func (c *Client) Connect() error {
	address := fmt.Sprintf("%s:%d", c.hostname, c.port)

	tlsConfig := &tls.Config{
		ServerName: c.hostname,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: c.timeout}, "tcp", address, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection to EPP server: %w", err)
	}

	c.conn = conn

	_, err = c.readResponse()
	if err != nil {
		c.conn.Close()
		return fmt.Errorf("failed to read server greeting: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) sendRequest(request []byte) ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("client not connected to EPP server")
	}

	length := uint32(len(request) + 4)
	header := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}

	if _, err := c.conn.Write(append(header, request...)); err != nil {
		return nil, fmt.Errorf("failed to send EPP request: %w", err)
	}

	response, err := c.readResponse()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) readResponse() ([]byte, error) {

	header := make([]byte, 4)
	if _, err := io.ReadFull(c.conn, header); err != nil {
		return nil, fmt.Errorf("failed to read EPP response header: %w", err)
	}

	length := uint32(header[0])<<24 | uint32(header[1])<<16 | uint32(header[2])<<8 | uint32(header[3])
	if length < 4 {
		return nil, fmt.Errorf("invalid EPP response length: %d", length)
	}

	body := make([]byte, length-4)
	if _, err := io.ReadFull(c.conn, body); err != nil {
		return nil, fmt.Errorf("failed to read EPP response body: %w", err)
	}

	return body, nil
}

func (c *Client) Login() error {
	loginReq := LoginRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: LoginCommand{
			Login: Login{
				ClID: c.username,
				Pw:   c.password,
				Options: LoginOptions{
					Version: "1.0",
					Lang:    "en",
				},
				Svcs: LoginServices{
					ObjURI: []string{
						"urn:ietf:params:xml:ns:domain-1.0",
						"urn:ietf:params:xml:ns:contact-1.0",
					},
					SvcExtension: &LoginServiceExtension{
						ExtURI: []string{
							"http://www.nic.at/xsd/at-ext-epp-1.0",
							"http://www.nic.at/xsd/at-ext-contact-1.0",
							"http://www.nic.at/xsd/at-ext-domain-1.0",
						},
					},
				},
			},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(loginReq)
	if err != nil {
		return fmt.Errorf("failed to marshal login request: %w", err)
	}

	responseXML, err := c.sendRequest(requestXML)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}

	var response Response
	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	if response.Result.Code != "1000" {
		return fmt.Errorf("EPP login failed: %s - %s", response.Result.Code, response.Result.Msg)
	}

	return nil
}

func (c *Client) Logout() error {
	logoutReq := LogoutRequest{
		XMLName: xml.Name{Local: "epp"},
		Xmlns:   "urn:ietf:params:xml:ns:epp-1.0",
		Command: LogoutCommand{
			Logout: struct{}{},
			ClTRID: generateTransactionID(),
		},
	}

	requestXML, err := xml.Marshal(logoutReq)
	if err != nil {
		return fmt.Errorf("failed to marshal logout request: %w", err)
	}

	_, err = c.sendRequest(requestXML)
	if err != nil {
		return fmt.Errorf("failed to send logout request: %w", err)
	}

	return nil
}
