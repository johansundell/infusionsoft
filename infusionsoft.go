// infusionsoft project infusionsoft.go
package infusionsoft

import (
	"reflect"

	//"github.com/kolo/xmlrpc"
	"github.com/johansundell/xmlrpc"
)

type DupCheck string

const (
	NoDupCheck                     DupCheck = ""
	EmailDupCheck                           = "Email"
	EmailAndNameDupCheck                    = "EmailAndName"
	EmailAndNameAndCompanyDupCheck          = "EmailAndNameAndCompany"
)

type Connection struct {
	apiKey string
	url    string
	client *xmlrpc.Client
}

type Contact struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Company   string
}

type CategoryTag struct {
	Id                  int
	CategoryName        string
	CategoryDescription string
}

type Tag struct {
	Id             int
	CategoryTagId  int    `xmlrpc:"GroupCategoryId"`
	TagName        string `xmlrpc:"GroupName"`
	TagDescription string `xmlrpc:"GroupDescription"`
}

func NewConnection(apiKey, url string) (Connection, error) {
	conn := Connection{apiKey: apiKey, url: url}
	var err error
	conn.client, err = xmlrpc.NewClient(conn.url, nil)
	return conn, err
}

func (conn *Connection) CreateContact(c *Contact, dc DupCheck) error {
	var id int
	var err error
	if dc != NoDupCheck {
		err = conn.client.Call("ContactService.addWithDupCheck", []interface{}{conn.apiKey, c, dc}, &id)
	} else {
		err = conn.client.Call("ContactService.add", []interface{}{conn.apiKey, c}, &id)
	}
	c.Id = id
	return err
}

func (conn *Connection) SearchContacts(limit, page int, queryData Contact) (result []Contact, err error) {
	selectedFields := getNamesFromStruct(queryData)
	params := []interface{}{conn.apiKey, "Contact", limit, page, queryData, selectedFields}
	err = conn.client.Call("DataService.query", params, &result)
	return
}

func (conn *Connection) CreateCategoryTag(c *CategoryTag) error {
	var id int
	err := conn.client.Call("DataService.add", []interface{}{conn.apiKey, "ContactGroupCategory", c}, &id)
	c.Id = id
	return err
}

func (conn *Connection) SearchCategoryTags(limit, page int, queryData CategoryTag) (result []CategoryTag, err error) {
	selectedFields := getNamesFromStruct(queryData)
	params := []interface{}{conn.apiKey, "ContactGroupCategory", limit, page, queryData, selectedFields}
	err = conn.client.Call("DataService.query", params, &result)
	return
}

func (conn *Connection) CreateTag(c *Tag) error {
	var id int
	err := conn.client.Call("DataService.add", []interface{}{conn.apiKey, "ContactGroup", c}, &id)
	c.Id = id
	return err
}

func (conn *Connection) SearchTags(limit, page int, queryData Tag) (result []Tag, err error) {
	selectedFields := getNamesFromStruct(queryData)
	params := []interface{}{conn.apiKey, "ContactGroup", limit, page, queryData, selectedFields}
	err = conn.client.Call("DataService.query", params, &result)
	return
}

func (conn *Connection) OptInEmail(email, reason string) (result bool, err error) {
	params := []interface{}{conn.apiKey, email, reason}
	err = conn.client.Call("APIEmailService.optIn", params, &result)
	return
}

func getNamesFromStruct(i interface{}) []string {
	val := reflect.ValueOf(i)
	t := val.Type()
	selectedFields := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("xmlrpc")
		if name == "" {
			name = f.Name
		}
		selectedFields = append(selectedFields, name)
	}
	return selectedFields
}
