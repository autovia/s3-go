package s3

import (
	"encoding/xml"
	"net/http"
)

type Error struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	Message   string   `xml:"Message"`
	Resource  string   `xml:"Resource"`
	RequestId string   `xml:"RequestId"`
}

func RespondError(w http.ResponseWriter, r *http.Request, code string, message string, resource string) error {

	e := Error{
		Code:     code,
		Message:  message,
		Resource: resource,
	}

	out, _ := xml.MarshalIndent(e, " ", "  ")
	//fmt.Println(string(out))

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(out))

	return nil
}
