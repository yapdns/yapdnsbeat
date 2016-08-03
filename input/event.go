package input

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

// FileEvent is sent to the output and must contain all relevant information
type FileEvent struct {
	common.EventMetadata
	ReadTime     time.Time
	Source       string
	InputType    string
	DocumentType string
	Offset       int64
	Bytes        int
	Text         *string
	Fileinfo     os.FileInfo
	JSONFields   common.MapStr
	JSONConfig   *JSONConfig
	FileState    FileState
	DnsRecord    DnsRecord
}

type JSONConfig struct {
	MessageKey    string `config:"message_key"`
	KeysUnderRoot bool   `config:"keys_under_root"`
	OverwriteKeys bool   `config:"overwrite_keys"`
	AddErrorKey   bool   `config:"add_error_key"`
}

type MultilineConfig struct {
	Negate   bool           `config:"negate"`
	Match    string         `config:"match"       validate:"required"`
	MaxLines *int           `config:"max_lines"`
	Pattern  *regexp.Regexp `config:"pattern"`
	Timeout  *time.Duration `config:"timeout"     validate:"positive"`
}

func (c *MultilineConfig) Validate() error {
	if c.Match != "after" && c.Match != "before" {
		return fmt.Errorf("unknown matcher type: %s", c.Match)
	}
	return nil
}

// mergeJSONFields writes the JSON fields in the event map,
// respecting the KeysUnderRoot and OverwriteKeys configuration options.
// If MessageKey is defined, the Text value from the event always
// takes precedence.
func mergeJSONFields(f *FileEvent, event common.MapStr) {

	// The message key might have been modified by multiline
	if len(f.JSONConfig.MessageKey) > 0 && f.Text != nil {
		f.JSONFields[f.JSONConfig.MessageKey] = *f.Text
	}

	if f.JSONConfig.KeysUnderRoot {
		for k, v := range f.JSONFields {
			if f.JSONConfig.OverwriteKeys {
				if k == "@timestamp" {
					vstr, ok := v.(string)
					if !ok {
						logp.Err("JSON: Won't overwrite @timestamp because value is not string")
						event[jsonErrorKey] = "@timestamp not overwritten (not string)"
						continue
					}
					// @timestamp must be of time common.Time
					ts, err := common.ParseTime(vstr)
					if err != nil {
						logp.Err("JSON: Won't overwrite @timestamp because of parsing error: %v", err)
						event[jsonErrorKey] = fmt.Sprintf("@timestamp not overwritten (parse error on %s)", vstr)
						continue
					}
					event[k] = ts
				} else if k == "type" {
					vstr, ok := v.(string)
					if !ok {
						logp.Err("JSON: Won't overwrite type because value is not string")
						event[jsonErrorKey] = "type not overwritten (not string)"
						continue
					}
					if len(vstr) == 0 || vstr[0] == '_' {
						logp.Err("JSON: Won't overwrite type because value is empty or starts with an underscore")
						event[jsonErrorKey] = fmt.Sprintf("type not overwritten (invalid value [%s])", vstr)
						continue
					}
					event[k] = vstr
				} else {
					event[k] = v
				}
			} else if _, exists := event[k]; !exists {
				event[k] = v
			}
		}
	} else {
		event["json"] = f.JSONFields
	}
}

func (f *FileEvent) ToMapStr() common.MapStr {
	hostname, _ := os.Hostname()

	event := common.MapStr{
		"timestamp": common.Time(f.ReadTime),
		"domain":    f.DnsRecord.Domain,
		"rdata":     f.DnsRecord.Rdata,
		"rtype":     f.DnsRecord.Rtype,
		"client": common.MapStr{
			"service_type": hostname,
			"ip":       Ip,
		},
	}

	if f.DnsRecord.Ttl != -1 {
		event["ttl"] = f.DnsRecord.Ttl
	}

	return event
}

func (f *FileEvent) ExtractDnsRecord(regex *regexp.Regexp) bool {
	if f.Text == nil {
		return false
	}
	match := regex.FindStringSubmatch(*f.Text)
	subexpNames := regex.SubexpNames()

	if len(match) < len(subexpNames) {
		logp.Err("Not able to match all subExpNames")
		return false
	}

	result := make(map[string]string)
	for i, name := range subexpNames {
		result[name] = match[i]
	}

	// check for valid domain and rdata
	if _, ok := result["domain"]; ok {
		f.DnsRecord.Domain = result["domain"]
	}

	if _, ok := result["rdata"]; ok {
		f.DnsRecord.Rdata = result["rdata"]
	}

	return true
}
