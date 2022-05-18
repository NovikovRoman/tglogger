package tglogger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

const DefaultApiURL = "https://api.telegram.org"

type file struct {
	name     string
	caption  string
	reader   io.Reader
	fileSize int
}

type bot struct {
	client *http.Client
	token  string
	url    string
}

func newBot(pref Settings) (b *bot, err error) {
	if pref.Client == nil {
		pref.Client = http.DefaultClient
	}

	if pref.URL == "" {
		pref.URL = DefaultApiURL
	}

	b = &bot{
		client: pref.Client,
		url:    pref.URL,
		token:  pref.Token,
	}

	_, err = b.getMe()
	return
}

func (bot *bot) getMe() ([]byte, error) {
	return bot.raw("getMe", nil)
}

func (bot *bot) sendDocument(chatID int64, f *file) (err error) {
	payload := map[string]string{
		"chat_id":   strconv.Itoa(int(chatID)),
		"caption":   f.caption,
		"file_name": f.name,
	}

	if f.fileSize != 0 {
		payload["file_size"] = strconv.Itoa(f.fileSize)
	}

	_, err = bot.sendFiles("sendDocument", map[string]*file{"document": f}, payload)
	return
}

func (bot *bot) sendFiles(method string, files map[string]*file, payload map[string]string) (b []byte, err error) {
	var (
		resp *http.Response
		part io.Writer
	)

	if len(files) == 0 {
		err = errors.New("No files. ")
		return
	}

	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)

	go func() {
		defer func() {
			_ = pipeWriter.Close()
		}()

		for field, f := range files {
			part, err = writer.CreateFormFile(field, f.name)
			if err != nil {
				_ = pipeWriter.CloseWithError(err)
				return
			}

			if _, err = io.Copy(part, f.reader); err != nil {
				_ = pipeWriter.CloseWithError(err)
				return
			}
		}

		for field, value := range payload {
			if err = writer.WriteField(field, value); err != nil {
				_ = pipeWriter.CloseWithError(err)
				return
			}
		}

		if err = writer.Close(); err != nil {
			_ = pipeWriter.CloseWithError(err)
			return
		}
	}()

	url := bot.url + "/bot" + bot.token + "/" + method

	resp, err = bot.client.Post(url, writer.FormDataContentType(), pipeReader)
	if err != nil {
		_ = pipeReader.CloseWithError(err)
		return
	}

	if resp == nil {
		err = errors.New("Response is nil. ")
		return
	}

	resp.Close = true
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusInternalServerError {
		err = errors.New("StatusCode 500")
		return
	}

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	err = extractTelegramResponse(b)
	return
}

func (bot *bot) sendMessage(chatID int64, msg string) (err error) {
	payload := map[string]string{
		"chat_id":                  strconv.Itoa(int(chatID)),
		"text":                     msg,
		"parse_mode":               "Markdown",
		"disable_web_page_preview": "true",
	}

	_, err = bot.raw("sendMessage", payload)
	return
}

func (bot *bot) raw(method string, payload interface{}) (b []byte, err error) {
	var (
		buf  bytes.Buffer
		resp *http.Response
	)

	if err = json.NewEncoder(&buf).Encode(payload); err != nil {
		return
	}

	resp, err = bot.client.Post(bot.url+"/bot"+bot.token+"/"+method, "application/json", &buf)
	if err != nil {
		return
	}

	if resp == nil {
		err = errors.New("Response is nil. ")
		return
	}

	resp.Close = true
	defer func() {
		_ = resp.Body.Close()
	}()

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	err = extractTelegramResponse(b)
	return
}

func extractTelegramResponse(b []byte) (err error) {
	var e struct {
		Ok          bool   `json:"ok"`
		Code        int    `json:"error_code"`
		Description string `json:"description"`
		// Parameters  map[string]interface{} `json:"parameters"`
	}

	// fixme
	if json.NewDecoder(bytes.NewReader(b)).Decode(&e) != nil || e.Ok {
		return
	}

	err = fmt.Errorf("telegram: %s (%d)", e.Description, e.Code)
	return
}
