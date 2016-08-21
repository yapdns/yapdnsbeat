package input

import (
	"github.com/elastic/beats/libbeat/common"
)

const (
	A     = "A"
	AAAA  = "AAAA"
	CNAME = "CNAME"
	MX    = "MX"
	TXT   = "TXT"
	NS    = "NS"
)

type DnsRecord struct {
	Domain string
	Rtype  string
	Ttl    int
	Rdata  string
}

func NewDnsRecord() DnsRecord {
	return DnsRecord{
		Domain: "",
		Rtype:  A,
		Ttl:    -1,
		Rdata:  "",
	}
}
func (r *DnsRecord) ToMapStr() common.MapStr {
	event := common.MapStr{
	// common.EventMetadataKey: f.EventMetadata,
	// "@timestamp":            common.Time(f.ReadTime),
	// "source":                f.Source,
	// "offset":                f.Offset, // Offset here is the offset before the starting char.
	// "type":                  f.DocumentType,
	// "input_type":            f.InputType,
	}

	// if f.JSONConfig != nil && len(f.JSONFields) > 0 {
	// 	mergeJSONFields(f, event)
	// } else {
	// 	event["message"] = f.Text
	// }

	return event
}
