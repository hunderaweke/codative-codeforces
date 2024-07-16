package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/url"
	"regexp"

	"github.com/hunderaweke/codative-codeforces/utils"
	"github.com/juju/persistent-cookiejar"
	"golang.org/x/net/publicsuffix"
)

func createHash(hash string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(hash))
	return hasher.Sum(nil)
}

func encrypt(password, handle string) (string, error) {
	block, err := aes.NewCipher(createHash("abds" + handle + "2"))
	if err != nil {
		return "", nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	text := gcm.Seal(nonce, nonce, []byte(password), nil)
	plaintext := hex.EncodeToString(text)
	return plaintext, nil
}

func decrypt(password, handle string) (string, error) {
	data, err := hex.DecodeString(password)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(createHash("abds" + handle + "2"))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, text := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func findCsrf(body []byte) (string, error) {
	reg := regexp.MustCompile(`csrf='(.+?)'`)
	tmp := reg.FindSubmatch(body)
	if len(tmp) < 2 {
		return "", errors.New("cannot find csrf")
	}
	return string(tmp[1]), nil
}

func genFtaa() string {
	return utils.RandString(18)
}

func genBtaa() string {
	return "44fdcff4e443a6be61d26650b259dd90"
}

func (c *Client) Login(handleOrEmail, password string) error {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	c.client.Jar = jar
	u, _ := url.Parse("https://codeforces.com/enter")
	b, _ := utils.GetBody(c.client, u)
	csrfToken, _ := findCsrf(b)
	ftaa := genFtaa()
	bfaa := genBtaa()
	// TODO: Write the implementation for the login well
	resp, err := utils.PostBody(c.client, u, url.Values{
		"csrf_token":    {csrfToken},
		"action":        {"enter"},
		"ftaa":          {ftaa},
		"bfaa":          {bfaa},
		"handleOrEmail": {handleOrEmail},
		"password":      {password},
		"_tta":          {"176"},
		"remember":      {"on"},
	})
	if err != nil {
		return err
	}
	reg := regexp.MustCompile("Invalid handle/email or password")
	loginError := reg.FindSubmatch(resp)
	if len(loginError) != 0 {
		return errors.New(string(loginError[0]))
	}
	reg = regexp.MustCompile(`handle = "(.+?)"`)
	handle := reg.FindSubmatch(resp)
	ePass, _ := encrypt(password, string(handle[1]))
	c.Jar = jar
	c.Ftaa = ftaa
	c.Bfaa = bfaa
	c.Password = ePass
	c.Handle = string(handle[1])
	c.HandleOrEmail = handleOrEmail
	if err = c.save(); err != nil {
		return err
	}
	return nil
}
