package rum

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gebitang.com/maorum/config"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type (
	Payload struct {
		Type   string `json:"type"`
		Target Target `json:"target"`
		Object Object `json:"object"`
	}
	Target struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	}
	Object struct {
		Type    string  `json:"type"`
		Content string  `json:"content"`
		Name    string  `json:"name"`
		Image   []Image `json:"image"`
	}
	Image struct {
		MediaType string `json:"mediaType,omitempty"`
		Name      string `json:"name,omitempty"`
		Content   string `json:"content,omitempty"`
	}

	SenderList struct {
		Senders []string `json:"senders"`
	}

	ContentItem struct {
		TrxId     string        `json:"TrxId"`
		Publisher string        `json:"Publisher"`
		TypeUrl   string        `json:"TypeUrl"`
		TimeStamp int64         `json:"TimeStamp"`
		Content   ContentDetail `json:"Content"`
	}
	ContentDetail struct {
		Type    string  `json:"type"`
		Content string  `json:"content,omitempty"`
		Image   []Image `json:"image,omitempty"`
	}
)

func createPayload(img []byte, content string) Payload {

	m := Image{
		MediaType: "png",
		Name:      uuid.NewV4().String(),
		Content:   base64.StdEncoding.EncodeToString(img),
	}

	obj := Object{
		Type:    "Note",
		Content: content,
		Name:    "",
		Image:   []Image{m},
	}
	return Payload{
		Type: "Add",
		Target: Target{
			Id:   config.RumPostGroup,
			Type: "Group",
		},
		Object: obj,
	}
}

var (
	clientOnce    sync.Once
	client        *http.Client
	APIPostGroup  = "/api/v1/group/content"
	APIReadGroup  = "/app/api/v1/group/%s/content"
	APIGroupUsers = "/api/v1/group/%s/announced/users"
)

func initClient() {
	var cert tls.Certificate
	var err error
	caCert, err := ioutil.ReadFile(config.RumCertPath)
	if err != nil {
		log.Fatalf("Error opening cert file %s, Error: %s", config.RumCertPath, err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}
	client = &http.Client{Transport: t, Timeout: 15 * time.Second}
}

func GetGroupUsers() (bool, string) {
	//path := fmt.Sprintf(APIGroupUsers, config.RumReadGroup)
	req, err := http.NewRequest("GET", config.RumUrl+"/api/v1/group/4e784292-6a65-471e-9f80-e91202e3358c/producers", nil)
	if err != nil {
		log.Printf("unable to create http request due to error %s", err)
		return false, err.Error()
	}
	return doReq(req)
}

func ReadFromGroup(latestNum int) (bool, string) {

	sender := &SenderList{
		Senders: []string{},
	}
	marshal, err := json.Marshal(sender)
	if err != nil {
		log.Printf("marshal error: %s", err.Error())
		return false, err.Error()
	}
	path := fmt.Sprintf(APIReadGroup, config.RumReadGroup)
	req, err := http.NewRequest("POST", config.RumUrl+path, bytes.NewBuffer(marshal))
	if err != nil {
		log.Printf("unable to create http request due to error %s", err)
		return false, err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	q := req.URL.Query()
	q.Add("num", fmt.Sprintf("%d", latestNum))
	q.Add("reverse", "true")
	req.URL.RawQuery = q.Encode()
	return doReq(req)
}

// PostToGroup https://github.com/youngkin/gohttps/blob/master/client/client.go
func PostToGroup(img []byte, content string) (bool, string) {

	p := createPayload(img, content)
	marshal, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return false, err.Error()
	}
	// https://golangbyexample.com/http-pdf-post-go/ https://stackoverflow.com/a/24455606/1087122 https://stackoverflow.com/a/51455855/1087122
	req, err := http.NewRequest("POST", config.RumUrl+APIPostGroup, bytes.NewBuffer(marshal))
	if err != nil {
		log.Printf("unable to create http request due to error %s", err)
		return false, err.Error()
	}

	return doReq(req)
}

func doReq(req *http.Request) (bool, string) {
	clientOnce.Do(initClient)
	resp, err := client.Do(req)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			log.Printf("url.Error received on http request: %s", e)
			return false, e.Error()
		default:
			log.Printf("Unexpected error received: %s", err)
			return false, e.Error()
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("unexpected error reading response body: %s", err)
		return false, err.Error()
	}
	return resp.StatusCode == http.StatusOK, string(body)
}
