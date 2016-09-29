package main

import (
	poeditor ".."
	"fmt"
	"os/exec"
	"net/http"
	"io/ioutil"
	"time"
	"os"
)

const (
	apiToken = ""
	serverProjectId = ""
	webProjectId = ""
	testServerProjectId = ""
	testServerWebId = ""

	pathToWebLocaleMakeFile = "./web/thequestion/locales"
	pathToServerLocaleMakeFile = "./server/locale"

	pathToRUServerPo = "./server/locale/ru/LC_MESSAGES/thequestion.po"
	pathToENServerPo = "./server/locale/en/LC_MESSAGES/thequestion.po"

	pathToRUWebJson = "./web/thequestion/locales/messages/ru.json"
	pathToENWebJson = "./web/thequestion/locales/messages/en.json"
)

func main() {
	if err := updateWebLocale(); err != nil {
		Error(err)
	}

	if err := updateServerLocale(); err != nil {
		Error(err)
	}
}

func updateServerLocale() error {
	cmd := exec.Command("make", "merge", "-C", pathToServerLocaleMakeFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	api := poeditor.New(apiToken, serverProjectId)
	if err := updateFile(api, pathToRUServerPo, poeditor.LanguageRU, poeditor.FileTypePo); err != nil {
		return err
	}

	time.Sleep(30*time.Second)
	
	api = poeditor.New(apiToken, serverProjectId)
	if err := updateFile(api, pathToENServerPo, poeditor.LanguageEN, poeditor.FileTypePo); err != nil {
		return err
	}

	return nil
}

func updateWebLocale() error {
	cmd := exec.Command("make", "-C", pathToWebLocaleMakeFile)
	if err := cmd.Run(); err != nil {
		return err
	}

	api := poeditor.New(apiToken, webProjectId)
	if err := updateFile(api, pathToRUWebJson, poeditor.LanguageRU, poeditor.FileTypeKeyValueJson); err != nil {
		return err
	}

	time.Sleep(30*time.Second)

	api = poeditor.New(apiToken, webProjectId)
	if err := updateFile(api, pathToENWebJson, poeditor.LanguageEN, poeditor.FileTypeKeyValueJson); err != nil {
		return err
	}

	return nil
}


func updateFile(api *poeditor.Poeditor, path string, language poeditor.Language, filetype poeditor.FileType) error {
	if err := api.Upload(path, language, poeditor.UploadTypeTermsDefinitions); err != nil {
		return err
	}

	item, err := api.Export(language, filetype)
	if err != nil {
		return err
	}

	resp, err := http.Get(item)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	return nil
}

func Error(err error) {
	fmt.Printf("exit with error: %v", err.Error())
	os.Exit(2)
}
