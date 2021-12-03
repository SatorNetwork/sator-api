package sumsub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

const (
	BasicKYCLevel = "basic-kyc-level"
)

// Service is a sumsub facade
type Service struct {
	appToken  string
	appSecret string
	baseURL   string
	ttl       int
}

// New creates new sumsub facade
func New(appToken, appSecret, baseURL string, ttl int) *Service {
	return &Service{
		appToken:  appToken,
		appSecret: appSecret,
		baseURL:   baseURL,
		ttl:       ttl,
	}
}

// GetSDKAccessToken returns access token for web or mobile SDKs
func (s *Service) GetSDKAccessToken(applicantID, externalUserID, externalAction, levelName string) (string, error) {
	path := "/resources/accessTokens"

	params := url.Values{}
	if applicantID != "" {
		params.Add("userId", applicantID)
	} else if externalUserID != "" {
		params.Add("userId", externalUserID)
	}

	if externalAction != "" {
		params.Add("externalActionId", externalAction)
	}

	if levelName != "" {
		params.Add("levelName", levelName)
	}

	params.Add("ttlInSecs", fmt.Sprint(s.ttl))

	type responseSDKAccessToken struct {
		Token  string `json:"token"`
		UserID string `json:"userId"`
	}

	path = fmt.Sprintf("%s?%s", path, params.Encode())
	resp := responseSDKAccessToken{}

	if err := s.jsonRequest(http.MethodPost, path, nil, &resp); err != nil {
		return "", fmt.Errorf("could not get access token by params=%s: %w", params.Encode(), err)
	}

	if resp.Token == "" || resp.UserID == "" {
		return "", ErrNotFound
	}

	return resp.Token, nil
}

// GetSDKAccessTokenByApplicantID returns access token for web or mobile SDKs by applicant id
func (s *Service) GetSDKAccessTokenByApplicantID(ctx context.Context, applicantID string) (string, error) {
	return s.GetSDKAccessToken(applicantID, "", "", BasicKYCLevel)
}

// GetSDKAccessTokenByUserID returns access token for web or mobile SDKs by user id
func (s *Service) GetSDKAccessTokenByUserID(ctx context.Context, userID, levelName string) (string, error) {
	return s.GetSDKAccessToken("", userID, "", levelName)
}

// Get returns sumsub response for applicant
func (s *Service) Get(applicantID string) (*Response, error) {
	path := "/resources/applicants/" + applicantID

	type responseStruct struct {
		List struct {
			Items      []Response `json:"items"`
			TotalItems uint       `json:"totalItems"`
		} `json:"list"`
	}

	resp := responseStruct{}
	if err := s.jsonRequest(http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("error: Get for applicant with id %s: %w", applicantID, err)
	}

	if resp.List.Items == nil || len(resp.List.Items) < 1 {
		return nil, ErrNotFound
	}

	return &resp.List.Items[0], nil
}

// GetByExternalUserID returns sumsub response for applicant by externalUserID
func (s *Service) GetByExternalUserID(externalUserID uuid.UUID) (*Response, error) {
	path := "/resources/applicants/-;externalUserId=" + externalUserID.String()

	type responseStruct struct {
		List struct {
			Items      []Response `json:"items"`
			TotalItems uint       `json:"totalItems"`
		} `json:"list"`
	}

	resp := responseStruct{}
	if err := s.jsonRequest(http.MethodGet, path, nil, &resp); err != nil {
		return nil, fmt.Errorf("error: GetByExternalUserID for external user id %s: %w", externalUserID.String(), err)
	}

	if resp.List.Items == nil || len(resp.List.Items) < 1 {
		return nil, ErrNotFound
	}

	return &resp.List.Items[0], nil
}

// GetReport returns applicantReport by applicantID and reviewID
func (s *Service) GetReport(applicantID, reviewID string) (map[string][]byte, error) {
	pathID := "/resources/applicants/" + applicantID + "/summary/report?lang=en&report=applicantReport" //"/requiredIdDocsStatus"
	methodID := "GET"
	timestampID := fmt.Sprint(time.Now().UTC().Unix())
	response, err := s.sendRequest(timestampID, methodID, pathID, "application/json", nil)
	if err != nil {
		return nil, fmt.Errorf("can't get report for client: %s error %w", applicantID, err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		errCode, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("can't get report for client: %s error code %s", applicantID, string(errCode))
	}
	documentList := map[string][]byte{}
	documentList["report.pdf"], err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read report for client: %s error code %w", applicantID, err)
	}
	return documentList, nil
}

// CreateSumSubRecord creates an applicant record in sumsub
func (s *Service) CreateSumSubRecord(externalUserID uuid.UUID, sourceKey string, create map[string]interface{}) (string, error) {
	path := "/resources/applicants"

	applicant := sumSubApplicantFromMap(create)
	applicant.ExternalUserID = externalUserID.String()
	applicant.SourceKey = sourceKey

	body, err := json.Marshal(applicant)
	if err != nil {
		return "", fmt.Errorf("CreateSumSubRecord: could not encode applicant data request: %w", err)
	}

	response := Response{}
	if err := s.jsonRequest(http.MethodPost, path, body, &response); err != nil {
		return "", fmt.Errorf("error: CreateSumSubRecord for external user id %s: %w", externalUserID.String(), err)
	}

	return response.ID, nil
}

// SetRequiredIDDocs sets required id docs for applicant
func (s *Service) SetRequiredIDDocs(applicantID string, docSet RequiredIDDocs) error {
	path := "/resources/applicants/" + applicantID + "/requiredIdDocs"

	docSetRaw, err := json.Marshal(docSet)
	if err != nil {
		return fmt.Errorf("could not encode docSet request body: %w", err)
	}

	if err := s.jsonRequest(http.MethodPost, path, docSetRaw, nil); err != nil {
		return fmt.Errorf("error: SetRequiredIDDocs for applicant with id %s: %w", applicantID, err)
	}

	return nil
}

// PostIDDoc sends ID document's bytes of applicant
func (s Service) PostIDDoc(applicantID string, doc []byte) error {
	path := "/resources/applicants/" + applicantID + "/info/idDoc"
	method := "POST"
	timestamp := fmt.Sprint(time.Now().UTC().Unix())

	boundary := "---documents---"

	response, err := s.sendRequest(timestamp, method, path, "multipart/form-data; boundary="+boundary, doc)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errCode, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("wrong response status code=%v, errorCode=%v", response.StatusCode, errCode)
	}

	return nil
}

// GetDocumentStatus returns document status for applicant
func (s *Service) GetDocumentStatus(applicantID string) ([]DocumentStatus, error) {
	type responseStruct struct {
		ID             string           `json:"id"`
		Status         Review           `json:"status"`
		DocumentStatus []DocumentStatus `json:"documentStatus"`
	}

	result := responseStruct{}
	path := "/resources/applicants/" + applicantID + "/state"

	if err := s.jsonRequest(http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("could not make sumsub request: %w", err)
	}

	return result.DocumentStatus, nil
}

// GetRequiredDocumentStatus returns status of required docs
func (s *Service) GetRequiredDocumentStatus(applicantID string) (RequiredIDDocStatus, error) {
	path := "/resources/applicants/" + applicantID + "/requiredIdDocsStatus"
	result := RequiredIDDocStatus{}
	if err := s.jsonRequest(http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("could not make sumsub request: %w", err)
	}
	return result, nil
}

// PatchApplicantInfo updates applicant info
func (s *Service) PatchApplicantInfo(applicantID string, update map[string]interface{}) error {
	info := sumSubApplicantInfoFromMap(update)
	path := "/resources/applicants/" + applicantID + "/info"

	doc, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("could not marshal applicantInfo: %v", err)
	}

	if err := s.jsonRequest(http.MethodPatch, path, doc, nil); err != nil {
		return fmt.Errorf("could not make sumsub request: %w", err)
	}

	return nil
}

//TestStatus tests status of applicant
func (s *Service) TestStatus(applicantID string, body map[string]interface{}) error {
	path := "/resources/applicants/" + applicantID + "/status/testCompleted"
	doc, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not encode request payload: %w", err)
	}
	return s.jsonRequest(http.MethodPost, path, doc, nil)
}

// func (s *Service) getApprovedDocuments(applicantID, reviewID string) (map[string][]byte, error) {
// 	data, err := s.GetRequiredDocumentStatus(applicantID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	documentList := map[string][]byte{}
// 	for key, value := range data {
// 		if value.ReviewResult.ReviewAnswer != "GREEN" {
// 			continue
// 		}
// 		documents := value.ImageReviewResults

// 		for docID, imageReview := range documents {
// 			if imageReview.ReviewAnswer != "GREEN" {
// 				continue
// 			}

// 			review, format, err := s.getReviewResult(reviewID, docID)
// 			if err != nil {
// 				return nil, err
// 			}
// 			documentList[key+format] = review
// 		}
// 	}
// 	return documentList, nil
// }

// func (s Service) getReviewResult(reviewID, docID string) ([]byte, string, error) {
// 	pathID := "/resources/inspections/" + reviewID + "/resources/" + docID
// 	methodID := "GET"
// 	timestampID := fmt.Sprint(time.Now().UTC().Unix())

// 	resp, err := s.sendRequest(timestampID, methodID, pathID, "application/json", nil)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode < 200 || resp.StatusCode > 299 {
// 		errCode, _ := ioutil.ReadAll(resp.Body)
// 		return nil, "", fmt.Errorf("can't get kyc document image, error code=%v", errCode)
// 	}

// 	contentType := strings.Split(resp.Header.Get("Content-Type"), "/")
// 	format := ""
// 	if len(contentType) > 0 {
// 		format = "." + contentType[len(contentType)-1]
// 	}
// 	res, err := ioutil.ReadAll(resp.Body)
// 	return res, format, err
// }

func sumSubApplicantFromMap(data map[string]interface{}) Applicant {
	result := Applicant{}
	if val, ok := data["email"].(string); ok {
		result.Email = val
	}
	if val, ok := data["lang"].(string); ok {
		result.Lang = val
	}

	result.Info = sumSubApplicantInfoFromMap(data)
	return result
}

func sumSubApplicantInfoFromMap(data map[string]interface{}) *ApplicantInfo {
	isCreated := false
	result := ApplicantInfo{}
	if val, ok := data["first_name"].(string); ok {
		isCreated = true
		result.FirstName = val
	}
	if val, ok := data["last_name"].(string); ok {
		isCreated = true
		result.LastName = val
	}
	if val, ok := data["middle_name"].(string); ok {
		isCreated = true
		result.MiddleName = val
	}
	if val, ok := data["mobile_number"].(string); ok {
		isCreated = true
		result.Phone = val
	}
	if val, ok := data["gender"].(string); ok {
		isCreated = true
		result.Gender = val
	}
	if val, ok := data["birth_date"].(string); ok {
		isCreated = true
		result.DOB = val
	}
	if val, ok := data["country"].(string); ok {
		isCreated = true
		result.Country = val
	}
	if val, ok := data["nationality"].(string); ok {
		isCreated = true
		result.Nationality = val
	}
	if val, ok := data["place_of_birth"].(string); ok {
		isCreated = true
		result.PlaceOfBirth = val
	}
	if val, ok := data["country_of_birth"].(string); ok {
		isCreated = true
		result.CountryOfBirth = val
	}
	if val, ok := data["state_of_birth"].(string); ok {
		isCreated = true
		result.StateOfBirth = val
	}
	if val, ok := data["lang"].(string); ok {
		isCreated = true
		result.Language = val
	}
	if isCreated {
		return &result
	}
	return nil
}

func (s *Service) jsonRequest(method, path string, body []byte, resp interface{}) error {
	timestampID := fmt.Sprint(time.Now().UTC().Unix())

	data, err := s.sendRequest(timestampID, method, path, "application/json", body)
	if err != nil {
		return err
	}
	defer data.Body.Close()

	if data.StatusCode == 404 {
		errCode, err := ioutil.ReadAll(data.Body)
		if err != nil {
			return fmt.Errorf("could not read response body from %v: %w", path, err)
		}
		return fmt.Errorf("wrong status code from %v, statusCode=%v, errCode=%s: %w", path, data.StatusCode, errCode, ErrNotFound)
	} else if data.StatusCode < 200 || data.StatusCode > 299 {
		errCode, err := ioutil.ReadAll(data.Body)
		if err != nil {
			return fmt.Errorf("could not read response body from %v: %w", path, err)
		}
		return fmt.Errorf("wrong status code from %v, statusCode=%v, errCode=%s", path, data.StatusCode, errCode)
	}

	if resp != nil {
		if err = json.NewDecoder(data.Body).Decode(resp); err != nil {
			return fmt.Errorf("could not parse json response for request: %s:%s: %w", method, path, err)
		}
	}

	return nil
}

func (s *Service) sendRequest(timestamp, method, path, contentType string, body []byte) (*http.Response, error) {
	client := http.DefaultClient
	sign := s.signRequest(timestamp, s.appSecret, method, path, body)
	request, _ := http.NewRequest(method, s.baseURL+path, bytes.NewReader(body))
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-App-Token", s.appToken)
	request.Header.Set("X-App-Access-Ts", timestamp)
	request.Header.Set("X-App-Access-Sig", sign)
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send request to %v: %v", path, err)
	}

	return resp, nil
}

func (s *Service) signRequest(timestamp, secret, method, path string, body []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	dataToSign := timestamp + method + path
	if body != nil {
		dataToSign += string(body)
	}
	h.Write([]byte(dataToSign))
	return hex.EncodeToString(h.Sum(nil))
}
