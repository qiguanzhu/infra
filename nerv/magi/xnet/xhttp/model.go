package xhttp

const (
	ResTypeJSON = "json"
	ResTypeXML  = "xml"

	TypeJSON              = "json"
	TypeXML               = "xml"
	TypeFormData          = "form-data"
	TypeMultipartFormData = "multipart-form-data"
	TypeSmartFormData     = "smart" // 智能版本，自动判断 payload 类型，选择合适的读取方式
)

var (
	_ReqContentTypeMap = map[string]string{
		TypeJSON:              "application/json",
		TypeXML:               "application/xml",
		TypeFormData:          "application/x-www-form-urlencoded",
		TypeMultipartFormData: "multipart/form-data",
	}

	_ResContentTypeMap = map[string]string{
		ResTypeJSON: "application/json",
		ResTypeXML:  "application/xml",
	}
)
