package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hunderaweke/codative-codeforces/utils"
	persistent "github.com/juju/persistent-cookiejar"
	"github.com/mgutz/ansi"
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

func Login(handleOrEmail, password string) error {
	jar, _ := persistent.New(&persistent.Options{PublicSuffixList: publicsuffix.List})
	S.Client.Jar = jar
	u, _ := url.Parse(S.Host + "enter")
	b, _ := utils.GetBody(S.Client, u)
	csrfToken, _ := findCsrf(b)
	ftaa := genFtaa()
	bfaa := genBtaa()
	resp, err := utils.PostBody(S.Client, u, url.Values{
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
	S.Handle = string(handle[1])
	S.HandleOrEmail = handleOrEmail
	S.Password = ePass
	S.Cookies = jar.Cookies(u)
	S.Bfaa = bfaa
	S.Ftaa = ftaa
	if err = S.Save(); err != nil {
		return err
	}
	return nil
}

func LoginPrompt() {
	var handle, password string
	handlePrompt := &survey.Input{
		Message: "Enter Your Handle or Email:",
	}
	survey.AskOne(handlePrompt, &handle)
	passwordPrompt := &survey.Password{
		Message: "Enter Password:",
	}
	survey.AskOne(passwordPrompt, &password)
	err := Login(handle, password)
	if err != nil {
		fmt.Println(err)
	}
	msg := ansi.Color("Logged in Successfully", "green+bh")
	fmt.Println(msg)
}
