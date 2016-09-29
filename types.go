package poeditor

import "errors"

//
// Types
//

type Language string

func (t Language) String() string {
	return string(t)
}

//

type Action string

func (t Action) String() string {
	return string(t)
}

//

type FileType string

func (t FileType) String() string {
	return string(t)
}

//

type UpdateType string

func (t UpdateType) String() string {
	return string(t)
}

//
// Response types
//

type PoTerm struct {
	Term    string `json:"term"`
	Context string `json:"context"`
}

type response struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type listElem struct {
	PoTerm
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	Percentage float32 `json:"percentage"`
}

type details struct {
	Parsed  int32 `json:"parsed"`
	Added   int32 `json:"added"`
	Deleted int32 `json:"deleted"`
	Updated int32 `json:"updated"`
}

type result struct {
	Response response   `json:"response"`
	List     []listElem `json:"list"`
	Details  details    `json:"details"`
	Item string `json:"item"`
}

//
// Constants
//

const (
	apiUrl string = "https://poeditor.com/api/"

	actionAvailableLanguages Action = "available_languages"
	actionListLanguages Action = "list_languages"
	actionListTerms Action = "view_terms"
	actionAddTerms Action = "add_terms"
	actionDeleteTerms Action = "delete_terms"
	actionSyncTerms Action = "sync_terms"
	actionExportTerms Action = "export"
	actionUploadTerms Action = "upload"

	UploadTypeTerms = "terms"
	UploadTypeDefinitions UpdateType = "definitions"
	UploadTypeTermsDefinitions UpdateType = "terms_definitions"

	LanguageNULL Language = ""
	LanguageRU Language = "ru"
	LanguageEN Language = "en"

	FileTypePo FileType = "po"
	FileTypePot FileType = "pot"
	FileTypeMo FileType = "mo"
	FileTypeXls FileType = "xls"
	FileTypeCsv FileType = "csv"
	FileTypeRews FileType = "resw"
	FileTypeResx FileType = "resx"
	FileTypeAndroidStrings FileType = "android_strings"
	FileTypeAccpeStrings FileType = "apple_strings"
	FileTypeXliff FileType = "xliff"
	FileTypeProperties FileType = "properties"
	FileTypeKeyValueJson FileType = "key_value_json"
	FileTypeJson FileType = "json"
)

//
// Errors
//

var (
	ErrPoeErrorResponse = errors.New("Ошибка запроса к https://poeditor.com/api/")
)
