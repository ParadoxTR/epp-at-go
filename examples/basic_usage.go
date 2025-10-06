package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ParadoxTR/epp-at-go/epp"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	host := os.Getenv("EPP_HOST")
	username := os.Getenv("EPP_USERNAME")
	password := os.Getenv("EPP_PASSWORD")

	if host == "" || username == "" || password == "" {
		log.Fatal("Please set EPP_HOST, EPP_USERNAME, and EPP_PASSWORD environment variables")
	}

	config := epp.Config{
		Hostname: host,
		Port:     700,
		Username: username,
		Password: password,
		Timeout:  10 * time.Second,
	}

	client := epp.NewClient(config)
	defer client.Close()

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer func() {
		if err := client.Logout(); err != nil {
			log.Printf("Failed to logout: %v", err)
		}
	}()

	fmt.Println("Successfully connected and logged in to EPP server")

	fmt.Println("\n=== Hello Example ===")
	helloResp, err := client.Hello()
	if err != nil {
		log.Printf("Hello failed: %v", err)
	} else {
		fmt.Printf("Server ID: %s\n", helloResp.Greeting.SvID)
		fmt.Printf("Server Date: %s\n", helloResp.Greeting.SvDate)
	}

	fmt.Println("\n=== Domain Check Example ===")
	domainName := "teswt.at"
	checkResp, err := client.CheckDomain([]string{domainName})
	if err != nil {
		log.Printf("Domain check failed: %v", err)
	} else {
		if len(checkResp.ResData.ChkData.Names) > 0 {
			domain := checkResp.ResData.ChkData.Names[0]
			fmt.Printf("Domain: %s\n", domain.Name.Name)
			fmt.Printf("Available: %s\n", domain.Name.Available)
			if domain.Reason != "" {
				fmt.Printf("Reason: %s\n", domain.Reason)
			}
		} else {
			fmt.Println("No domain check results returned")
		}
	}

	fmt.Println("\n=== Contact Create Example ===")

	var contactID string

	contact := epp.Contact{
		ID: "AUTO", // Austrian EPP uses AUTO for auto-generation
		PostalInfo: epp.ContactPostalInfo{
			Type: "int",
			Name: "John Doe",

			Addr: epp.ContactAddr{
				Street: []string{"123 Main Street"},
				City:   "Vienna",
				PC:     "1010",
				CC:     "AT", // Use Austrian country code
			},
		},
		Voice: "+43.15551234567",
		Email: "john.doe@example.com",
		AuthInfo: epp.ContactAuthInfo{
			Pw: "TestPass123", // Simpler password
		},
		Type: "privateperson", // Austrian EPP extension: privateperson, organisation, role
	}

	createContactResp, err := client.CreateContact(&contact)
	if err != nil {
		log.Printf("Contact creation failed: %v", err)
	} else {
		fmt.Printf("Contact created successfully: %s\n", createContactResp.ResData.CreData.ID)
		fmt.Printf("Creation date: %s\n", createContactResp.ResData.CreData.CrDate)

		contactID = createContactResp.ResData.CreData.ID
	}

	fmt.Println("\n=== Contact Info Example ===")
	infoContactResp, err := client.InfoContact(contactID)
	if err != nil {
		log.Printf("Contact info failed: %v", err)
	} else {
		fmt.Printf("Contact ID: %s\n", infoContactResp.ResData.InfData.ID)
		fmt.Printf("Name: %s\n", infoContactResp.ResData.InfData.PostalInfo.Name)
		fmt.Printf("Email: %s\n", infoContactResp.ResData.InfData.Email)
		if len(infoContactResp.ResData.InfData.Status) > 0 {
			fmt.Printf("Status: %s\n", infoContactResp.ResData.InfData.Status[0].Status)
		}
	}

	if checkResp != nil && len(checkResp.ResData.ChkData.Names) > 0 && checkResp.ResData.ChkData.Names[0].Name.Available == "1" {
		fmt.Println("\n=== Domain Create Example ===")
		domain := epp.Domain{
			Name:        domainName,
			Nameservers: []string{"ns1.example.com", "ns2.example.com"},
			Registrant:  contactID,
			Contacts: []epp.DomainContact{
				{Type: "admin", ID: contactID},
				{Type: "tech", ID: contactID},
			},
			AuthInfo: "domain-auth-password-123",
		}

		createDomainResp, err := client.CreateDomain(domain)
		if err != nil {
			log.Printf("Domain creation failed: %v", err)
		} else {
			fmt.Printf("Domain created successfully: %s\n", createDomainResp.ResData.CreData.Name)
			fmt.Printf("Creation date: %s\n", createDomainResp.ResData.CreData.CrDate)
			fmt.Printf("Expiration date: %s\n", createDomainResp.ResData.CreData.ExDate)
		}

		fmt.Println("\n=== Domain Info Example ===")
		infoDomainResp, err := client.InfoDomain(domainName)
		if err != nil {
			log.Printf("Domain info failed: %v", err)
		} else {
			fmt.Printf("Domain: %s\n", infoDomainResp.ResData.InfData.Name)
			fmt.Printf("Registrant: %s\n", infoDomainResp.ResData.InfData.Registrant)
			if len(infoDomainResp.ResData.InfData.Status) > 0 {
				fmt.Printf("Status: %s\n", infoDomainResp.ResData.InfData.Status[0].Status)
			}
			fmt.Printf("Creation date: %s\n", infoDomainResp.ResData.InfData.CrDate)
			fmt.Printf("Expiration date: %s\n", infoDomainResp.ResData.InfData.ExDate)
		}
	}

	fmt.Println("\n=== Poll Messages Example ===")
	pollResp, err := client.PollMessage()
	if err != nil {
		log.Printf("Poll failed: %v", err)
	} else {
		if pollResp.MsgQ != nil && pollResp.MsgQ.Count > 0 {
			fmt.Printf("Messages in queue: %d\n", pollResp.MsgQ.Count)
			fmt.Printf("Message ID: %s\n", pollResp.MsgQ.ID)
			fmt.Printf("Queue date: %s\n", pollResp.MsgQ.QDate)
			fmt.Printf("Message: %s\n", pollResp.MsgQ.Msg)

			ackResp, err := client.AckPollMessage(pollResp.MsgQ.ID)
			if err != nil {
				log.Printf("Message acknowledgment failed: %v", err)
			} else {
				fmt.Printf("Message acknowledged successfully\n")
				if ackResp.MsgQ != nil {
					fmt.Printf("Remaining messages: %d\n", ackResp.MsgQ.Count)
				}
			}
		} else {
			fmt.Println("No messages in queue")
		}
	}

	fmt.Println("\n=== Example completed successfully ===")
}
