package tlwpa4220

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Username string
	Password string
	IP       string
}

func md5hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (c Client) authCookie() string {
	auth := "Basic " + c.Username + ":" + md5hash(c.Password)
	return "Authorization=" + auth
}

func (c Client) request(path string, values url.Values, data interface{}) error {
	body := strings.NewReader(values.Encode())

	url := fmt.Sprintf("http://%s/%s", c.IP, path)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", c.authCookie())
	req.Header.Set("Referer", fmt.Sprintf("http://%s/", c.IP))
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	if data != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(respBody, data)
		if err != nil {
			return err
		}
	}

	return nil
}
