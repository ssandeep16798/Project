package models

import "encoding/xml"

type Feedback struct {
	XMLName        xml.Name `xml:"feedback"`
	Text           string   `xml:",chardata"`
	ReportMetadata struct {
		Text             string `xml:",chardata"`
		OrgName          string `xml:"org_name"`
		Email            string `xml:"email"`
		ExtraContactInfo string `xml:"extra_contact_info"`
		ReportID         string `xml:"report_id"`
		DateRange        struct {
			Text  string `xml:",chardata"`
			Begin string `xml:"begin"`
			End   string `xml:"end"`
		} `xml:"date_range"`
	} `xml:"report_metadata"`
	PolicyPublished struct {
		Text   string `xml:",chardata"`
		Domain string `xml:"domain"`
		Adkim  string `xml:"adkim"`
		Aspf   string `xml:"aspf"`
		P      string `xml:"p"`
		Sp     string `xml:"sp"`
		Pct    string `xml:"pct"`
	} `xml:"policy_published"`
	Record []struct {
		Text string `xml:",chardata"`
		Row  struct {
			Text            string `xml:",chardata"`
			SourceIp        string `xml:"source_ip"`
			Count           string `xml:"count"`
			PolicyEvaluated struct {
				Text        string `xml:",chardata"`
				Disposition string `xml:"disposition"`
				Dkim        string `xml:"dkim"`
				Spf         string `xml:"spf"`
				Reason      struct {
					Text    string `xml:",chardata"`
					Type    string `xml:"type"`
					Comment string `xml:"comment"`
				} `xml:"reason"`
			} `xml:"policy_evaluated"`
		} `xml:"row"`
		Identifiers struct {
			Text       string `xml:",chardata"`
			HeaderFrom string `xml:"header_from"`
		} `xml:"identifiers"`
		AuthResults struct {
			Text string `xml:",chardata"`
			Dkim struct {
				Text     string `xml:",chardata"`
				Domain   string `xml:"domain"`
				Result   string `xml:"result"`
				Selector string `xml:"selector"`
			} `xml:"dkim"`
			Spf struct {
				Text   string `xml:",chardata"`
				Domain string `xml:"domain"`
				Result string `xml:"result"`
			} `xml:"spf"`
		} `xml:"auth_results"`
	} `xml:"record"`
}
