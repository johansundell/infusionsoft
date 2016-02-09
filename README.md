# infusionsoft
TODO

## Usage
	package main
	
	import (
		"fmt"
		"os"
	
		"github.com/johansundell/infusionsoft"
	)
	
	var key, client string
	
	func init() {
		key = os.Getenv("IS_KEY")
		client = os.Getenv("IS_CLIENT")
	}
	
	func main() {
	
		// Create a new connection to infusionsoft; key is the api key, client is the url (https://test.infusionsoft.com/api/xmlrpc)
		conn, err := infusionsoft.NewConnection(key, client)
		if err != nil {
			fmt.Println("Could not create a connection", err)
			return
		}
	
		// Search for a contact, if none found create the contact
		contact := infusionsoft.Contact{FirstName: "Test", LastName: "Person", Company: "Test Company", Email: "test@test.com"}
		contacts, err := conn.SearchContacts(100, 0, contact)
		if err != nil {
			fmt.Println("Err in search for contacts", err)
			return
		}
		if len(contacts) == 0 {
			if err := conn.CreateContact(&contact, infusionsoft.EmailDupCheck); err != nil {
				fmt.Println("Could not create contact", err)
				return
			}
		} else {
			contact = contacts[0]
		}
	
		emailStatus, err := conn.GetEmailStatus(contact.Email)
		if err != nil {
			fmt.Println("Failed to get email status", contact.Email, err)
		}
		if emailStatus == 0 {
			// Optin the email address of the contact
			if result, err := conn.OptInEmail(contact.Email, "API optin"); !result || err != nil {
				fmt.Println("Could not optin email", contact.Email, err)
			}
		}
	
		// Search for a category tag, if none found create it
		category := infusionsoft.CategoryTag{CategoryName: "My new category"}
		categories, err := conn.SearchCategoryTags(100, 0, category)
		if err != nil {
			fmt.Println("Err in categories search", err)
			return
		}
		if len(categories) == 0 {
			if err := conn.CreateCategoryTag(&category); err != nil {
				fmt.Println("Could not create category", err)
			}
		} else {
			category = categories[0]
		}
	
		// Search for a tag, if none found create it
		tag := infusionsoft.Tag{TagName: "My new tag", CategoryTagId: category.Id}
		tags, err := conn.SearchTags(100, 0, tag)
		if err != nil {
			fmt.Println("Err in search for tags", err)
			return
		}
		if len(tags) == 0 {
			if err := conn.CreateTag(&tag); err != nil {
				fmt.Println("Could not create tag", err)
				return
			}
		} else {
			tag = tags[0]
		}
	
		// Add the tag to the contact
		if result, err := conn.AddTagToContact(contact, tag); !result || err != nil {
			fmt.Println("Could not add tag to contact", contact, tag, err)
		}
	}
