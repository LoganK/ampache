package ampache

import (
	"encoding/xml"
	"time"
)

type xmlTime struct {
	time.Time
}

func (c *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse("2006-01-02T03:04:05-07:00", v)
	if err != nil {
		return err
	}
	*c = xmlTime{parse}
	return nil
}
