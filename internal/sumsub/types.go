package sumsub

import "strings"

type (
	// Address struct
	Address struct {
		Country        string `json:"country"`
		PostCode       string `json:"postCode"`
		Town           string `json:"town"`
		Street         string `json:"street"`
		SubStreet      string `json:"subStreet"`
		State          string `json:"state"`
		BuildingName   string `json:"buildingName"`
		FlatNumber     string `json:"flatNumber"`
		BuildingNumber string `json:"buildingNumber"`
		StartDate      string `json:"startDate"`
		EndDate        string `json:"endDate"`
	}

	// RequiredIDDocs struct
	RequiredIDDocs struct {
		DocSets []DocSets `json:"docSets"`
	}

	// DocSets struct
	DocSets struct {
		IDDocSetType string   `json:"idDocSetType"`
		Types        []string `json:"types"`
		Fields       []string `json:"fields"`
	}

	// Applicant struct
	Applicant struct {
		ExternalUserID string              `json:"externalUserId,omitempty"`
		SourceKey      string              `json:"sourceKey,omitempty"`
		Email          string              `json:"email,omitempty"`
		Lang           string              `json:"lang,omitempty"`
		Metadata       []map[string]string `json:"metadata,omitempty"`
		RequiredIDDocs *RequiredIDDocs     `json:"requiredIdDocs,omitempty"`
		Info           *ApplicantInfo      `json:"info,omitempty"`
	}

	// Response struct
	Response struct {
		ID             string         `json:"id"`
		CreatedAt      string         `json:"createdAt"`
		ClientID       string         `json:"clientId"`
		InspectionID   string         `json:"inspectionId"`
		ExternalUserID string         `json:"externalUserId"`
		Info           ApplicantInfo  `json:"info"`
		Env            string         `json:"env"`
		Email          string         `json:"email"`
		RequiredIDDocs RequiredIDDocs `json:"requiredIdDocs"`
		Review         Review         `json:"review"`
	}

	// RequiredIDDocStatus type
	RequiredIDDocStatus map[string]struct {
		ReviewResult       *ReviewResult           `json:"reviewResult,omitempty"`
		Country            string                  `json:"country"`
		IDDocType          string                  `json:"idDocType"`
		ImageIDs           []int                   `json:"imageIds"`
		ImageReviewResults map[string]ReviewResult `json:"imageReviewResults"`
	}

	// ReviewResult struct
	ReviewResult struct {
		ModerationComment string   `json:"moderationComment"`
		ClientComment     string   `json:"clientComment"`
		ReviewAnswer      string   `json:"reviewAnswer"`
		RejectLabels      []string `json:"rejectLabels"`
		ReviewRejectType  string   `json:"reviewRejectType"`
	}

	// DocumentStatus struct
	DocumentStatus struct {
		IDDocType    string       `json:"idDocType"`
		IDDocSubType string       `json:"idDocSubType"`
		Country      string       `json:"country"`
		ImageID      int          `json:"imageId"`
		ReviewResult ReviewResult `json:"reviewResult"`
		AddedDate    string       `json:"addedDate"`
	}

	// ApplicantInfo struct
	ApplicantInfo struct {
		FirstName      string    `json:"firstName,omitempty"`
		LastName       string    `json:"lastName,omitempty"`
		MiddleName     string    `json:"middleName,omitempty"`
		LegalName      string    `json:"legalName,omitempty"`
		Gender         string    `json:"gender,omitempty"`
		DOB            string    `json:"dob,omitempty"`
		PlaceOfBirth   string    `json:"placeOfBirth,omitempty"`
		CountryOfBirth string    `json:"countryOfBirth,omitempty"`
		StateOfBirth   string    `json:"stateOfBirth,omitempty"`
		Country        string    `json:"country,omitempty"`
		Nationality    string    `json:"nationality,omitempty"`
		Phone          string    `json:"phone,omitempty"`
		Addresses      []Address `json:"addresses,omitempty"`
		Language       string    `json:"lang"`
	}

	// Review struct
	Review struct {
		ID                     string       `json:"id"`
		InspectionID           string       `json:"inspectionId"`
		CreatedDate            string       `json:"createDate"`
		ReviewDate             string       `json:"reviewDate"`
		StartDate              string       `json:"startDate"`
		ReviewResult           ReviewResult `json:"reviewResult"`
		ReviewStatus           string       `json:"reviewStatus"`
		NotificationFailureCnt int          `json:"notificationFailureCnt"`
		ApplicantID            string       `json:"applicantId"`
	}

	WebhookPayload struct {
		ApplicantID               string       `json:"applicantId"`
		ApplicantActionID         string       `json:"applicantActionId"`
		ExternalApplicantActionID string       `json:"externalApplicantActionId"`
		InspectionID              string       `json:"inspectionId"`
		CorrelationID             string       `json:"correlationId"`
		ExternalUserID            string       `json:"externalUserId"`
		EventType                 string       `json:"type"`
		ReviewResult              ReviewResult `json:"reviewResult,omitempty"`
		ReviewStatus              string       `json:"reviewStatus"`
		CreatedAt                 string       `json:"createdAt"`
	}
)

func KYCRejectLabelsMap() func([]string) string {
	innerMap := map[string]string{
		"FORGERY":                         "Forgery attempt has been made.",
		"DOCUMENT_TEMPLATE":               "Documents supplied are templates, downloaded from internet.",
		"SPAM":                            "An applicant has been created by mistake or is just a spam user (irrelevant images were supplied).",
		"NOT_DOCUMENT":                    "Documents supplied are not relevant for the verification procedure.",
		"SELFIE_MISMATCH":                 "A user photo (profile image) does not match a photo on the provided documents.",
		"FOREIGNER":                       "When a client does not accept applicants from a different country or e.g. without a residence permit.",
		"DUPLICATE":                       "This applicant was already created for this client, and duplicates are not allowed by the regulations.",
		"WRONG_USER_REGION":               "When applicants from certain regions/countries are not allowed to be registered.",
		"BLACKLIST":                       "User is blacklisted.",
		"REGULATIONS_VIOLATIONS":          "Regulations violations.",
		"INCONSISTENT_PROFILE":            "Data or documents of different persons were uploaded to one applicant.",
		"AGE_REQUIREMENT_MISMATCH":        "Age requirement is not met (e.g. cannot rent a car to a person below 25yo).",
		"EXPERIENCE_REQUIREMENT_MISMATCH": "Not enough experience (e.g. driving experience is not enough).",
		"CRIMINAL":                        "The user is involved in illegal actions.",
		"COMPROMISED_PERSONS":             "The user does not correspond to Compromised Person Politics.",
		"PEP":                             "The user belongs to the PEP category.",
		"ADVERSE_MEDIA":                   "The user was found in the adverse media.",
		"FRAUDULENT_PATTERNS":             "Fraudulent behavior was detected.",
		"SANCTIONS":                       "The user was found on sanction lists.",
		"FRAUDULENT_LIVENESS":             "There was an attempt to bypass liveness check.",
	}

	return func(key []string) string {
		var resp []string
		for _, k := range key {
			resp = append(resp, innerMap[k])
		}

		return strings.Join(resp, "")
	}
}
