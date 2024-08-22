package importer

type SilentholdInfo struct {
	MatterName            string
	MatterID              string `json:"Matter id"`
	HoldName              string `json:"Hold Name"`
	AdvisoryNoticeSubject string `json:"Advisory notice subject"`
	AdvisoryNoticeTitle   string `json:"Advisory notice title"`
	AdvisoryNoticeBody    string `json:"Advisory notice body"`
}

type SilentholdDetail struct {
	FolderName       string
	SilentholdInfo   SilentholdInfo
	CustodianDetails []CustodianDetail
}

type SilentholdDetails []SilentholdDetail
