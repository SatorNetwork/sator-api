package firebase

import "net/http"

//// Error type for firebase client.
//var Error = errs.Class("firebase error")
//
//// ResponseError firebase response type of error.
//type ResponseError struct {
//	Status string `json:"Status"`
//	Detail string `json:"Detail"`
//}
//
//// ResponseErrors firebase response type in case of code 400 (error).
//type ResponseErrors struct {
//	Errors []ResponseError `json:"Errors"`
//}
//
//func (response *ResponseErrors) Error() (err error) {
//	for _, responseError := range response.Errors {
//		err = errs.Combine(err, errs.New(responseError.Detail))
//	}
//
//	return err
//}

// Config is configuration for firebase client.
type (
	Config struct {
		BaseFirebaseURL    string
		WebAPIKey          string
		MainSiteLink       string
		AndroidPackageName string
		IosBundleId        string
		SuffixOption       string
	}
	// Interactor struct
	Interactor struct {
		fbClient FBClient
		http     http.Client
		config   Config
	}
)

// New is a factory function,
// returns a new instance of the firebase interactor.
func New(client FBClient, http http.Client, config Config) *Interactor {
	return &Interactor{
		fbClient: client,
		http:     http,
		config:   config,
	}
}
