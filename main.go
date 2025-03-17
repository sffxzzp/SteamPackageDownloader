package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type (
	steamPackage struct {
		urlBase string
		path    string
	}
)

func download(filename string, url string) error {
	fmt.Printf("Downloading: %s\n", url)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	bar := pb.Full.Start64(res.ContentLength)
	defer bar.Finish()
	reader := bar.NewProxyReader(res.Body)
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	return err
}

func newSteamPackage(path string) *steamPackage {
	if _, err := os.Stat(path + "package"); os.IsNotExist(err) {
		os.Mkdir(path+"package", 0777)
	}
	s := &steamPackage{
		path: path + "package/",
	}
	s.testUrlBase()
	return s
}

func (s *steamPackage) tcpping(url string) int64 {
	start := time.Now()
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return -1
	}
	defer conn.Close()
	return int64(time.Since(start))
}

func (s *steamPackage) download(filename string) error {
	return download(s.path+filename, s.urlBase+filename)
}

func (s *steamPackage) getLink(obj map[string]interface{}) string {
	if _, exists := obj["zipvz"]; exists {
		return obj["zipvz"].(string)
	} else if _, exists := obj["file"]; exists {
		return obj["file"].(string)
	} else {
		return ""
	}
}

func (s *steamPackage) downManifest() {
	manifest := s.urlBase + "steam_client_win32"
	download(s.path+"steam_client_win32.manifest", manifest)
	vdf := newVDF()
	vdf.loadVDF(s.path + "steam_client_win32.manifest")
	for _, v := range vdf.data["win32"].(map[string]interface{}) {
		switch v.(type) {
		case string:
			continue
		case map[string]interface{}:
			if _, exists := v.(map[string]interface{})["win7-64"]; exists {
				s.download(s.getLink(v.(map[string]interface{})["win7-64"].(map[string]interface{})))
			} else if _, exists := v.(map[string]interface{})["steamrow"]; exists {
				s.download(s.getLink(v.(map[string]interface{})["steamrow"].(map[string]interface{})))
			} else {
				s.download(s.getLink(v.(map[string]interface{})))
			}
		}
	}
}

func (s *steamPackage) testUrlBase() {
	var urls = []string{
		"https://media.st.dl.eccdnx.com/client/",
		"https://client-update.queniuqe.com/",
		"https://client-update.akamai.steamstatic.com/",
		"https://client-update.fastly.steamstatic.com/",
	}

	var addresses = []string{
		"media.st.dl.eccdnx.com:443",
		"client-update.queniuqe.com:443",
		"client-update.akamai.steamstatic.com:443",
		"client-update.fastly.steamstatic.com:443",
	}

	minTime := int64(999999999)

	for i, address := range addresses {
		time := s.tcpping(address)
		if time > 0 && time < minTime {
			minTime = time
			s.urlBase = urls[i]
		}
	}
}

func main() {
	path := getSteamPath()
	s := newSteamPackage(path)
	s.downManifest()
}
