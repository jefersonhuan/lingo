package cmd

import (
	"encoding/xml"
)

type Server struct {
	XMLName  xml.Name `xml:"Server"`
	Host     string
	Port     int
	User     string
	Password string
}
