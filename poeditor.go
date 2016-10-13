package poeditor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type Poeditor struct {
	ApiToken  string
	ProjectId string
}

func New(apiToken string, projectId string) *Poeditor {
	return &Poeditor{
		ApiToken:  apiToken,
		ProjectId: projectId,
	}
}

func (t *Poeditor) ListProjectLanguages() ([]Language, error) {
	response, err := t.request(actionListLanguages, map[string]string{})
	if err != nil {
		return nil, err
	}

	result := []Language{}
	for _, listElem := range response.List {
		result = append(result, Language(listElem.Code))
	}

	return result, nil
}

func (t *Poeditor) ListProjectTerms(language Language) ([]*PoTerm, error) {
	response, err := t.request(actionListTerms, map[string]string{
		"language": language.String(),
	})
	if err != nil {
		return nil, err
	}

	result := []*PoTerm{}
	for _, listElem := range response.List {
		result = append(result, &PoTerm{
			Term:    listElem.Term,
			Context: listElem.Context,
		})
	}

	return result, nil
}

func (t *Poeditor) AddTermToProject(terms []*PoTerm) error {
	bytes, err := json.Marshal(terms)
	if err != nil {
		return err
	}

	_, err = t.request(actionAddTerms, map[string]string{
		"data": string(bytes),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Poeditor) DeleteTermFromProject(terms []*PoTerm) error {
	bytes, err := json.Marshal(terms)
	if err != nil {
		return err
	}

	_, err = t.request(actionDeleteTerms, map[string]string{
		"data": string(bytes),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Poeditor) SyncTerms(terms []*PoTerm) error {
	bytes, err := json.Marshal(terms)
	if err != nil {
		return err
	}

	_, err = t.request(actionSyncTerms, map[string]string{
		"data": string(bytes),
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Poeditor) Export(language Language, filetype FileType, filters string) (string, error) {
	response, err := t.request(actionExportTerms, map[string]string{
		"language": language.String(),
		"type":     filetype.String(),
		"filters":  filters,
	})
	if err != nil {
		return "", err
	}

	return response.Item, nil
}

func (t *Poeditor) Upload(pathToFile string, language Language, updateType UpdateType) error {
	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = t.requestWithFile(pathToFile, map[string]string{
		"language":  language.String(),
		"updating":  updateType.String(),
		"api_token": t.ApiToken,
		"action":    actionUploadTerms.String(),
		"id":        t.ProjectId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Poeditor) request(action Action, params map[string]string) (*result, error) {
	data, err := t.parseParams(t.ProjectId, action, params)
	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm(apiUrl, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return t.parseResponse(resp)
}

func (t *Poeditor) requestWithFile(pathToFile string, params map[string]string) (*result, error) {
	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("file", pathToFile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}

	for key, value := range params {
		writer.WriteField(key, value)
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return t.parseResponse(response)
}

func (t *Poeditor) parseResponse(response *http.Response) (*result, error) {
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, ErrPoeErrorResponse
	}

	result := new(result)
	if err := json.Unmarshal(bytes, result); err != nil {
		return nil, err
	}

	if result.Response.Status != "success" {
		return nil, errors.New(result.Response.Message)
	}

	return result, nil
}

func (t *Poeditor) parseParams(projectId string, action Action, params map[string]string) (url.Values, error) {
	urlParams := url.Values{}
	urlParams.Add("api_token", t.ApiToken)
	urlParams.Add("action", action.String())
	urlParams.Add("id", projectId)

	for param, value := range params {
		urlParams.Add(param, value)
	}

	return urlParams, nil
}
