package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zeebo/errs"
)

type (
	DynamicLinkRequest struct {
		DynamicLinkInfo DynamicLinkInfo `json:"dynamicLinkInfo"`
		Suffix          Suffix          `json:"suffix"`
	}

	DynamicLinkInfo struct {
		DomainUriPrefix   string            `json:"domainUriPrefix"`
		Link              string            `json:"link"`
		AndroidInfo       AndroidInfo       `json:"androidInfo"`
		IosInfo           IosInfo           `json:"iosInfo"`
		NavigationInfo    NavigationInfo    `json:"navigationInfo"`
		AnalyticsInfo     AnalyticsInfo     `json:"analyticsInfo"`
		SocialMetaTagInfo SocialMetaTagInfo `json:"socialMetaTagInfo"`
	}

	AndroidInfo struct {
		AndroidPackageName           string `json:"androidPackageName"`
		AndroidFallbackLink          string `json:"androidFallbackLink"`
		AndroidMinPackageVersionCode string `json:"androidMinPackageVersionCode"`
	}

	IosInfo struct {
		IosBundleId         string `json:"iosBundleId"`
		IosFallbackLink     string `json:"iosFallbackLink"`
		IosCustomScheme     string `json:"iosCustomScheme"`
		IosIpadFallbackLink string `json:"iosIpadFallbackLink"`
		IosIpadBundleId     string `json:"iosIpadBundleId"`
		IosAppStoreId       string `json:"iosAppStoreId"`
	}

	NavigationInfo struct {
		EnableForcedRedirect bool `json:"enableForcedRedirect"`
	}

	AnalyticsInfo struct {
		GooglePlayAnalytics    GooglePlayAnalytics    `json:"googlePlayAnalytics"`
		ItunesConnectAnalytics ItunesConnectAnalytics `json:"itunesConnectAnalytics"`
	}

	GooglePlayAnalytics struct {
		UtmSource   string `json:"utmSource"`
		UtmMedium   string `json:"utmMedium"`
		UtmCampaign string `json:"utmCampaign"`
		UtmTerm     string `json:"utmTerm"`
		UtmContent  string `json:"utmContent"`
		Gclid       string `json:"gclid"`
	}

	ItunesConnectAnalytics struct {
		At string `json:"at"`
		Ct string `json:"ct"`
		Mt string `json:"mt"`
		Pt string `json:"pt"`
	}

	SocialMetaTagInfo struct {
		SocialTitle       interface{} `json:"socialTitle"`
		SocialDescription interface{} `json:"socialDescription"`
		SocialImageLink   interface{} `json:"socialImageLink"`
	}

	Suffix struct {
		Option string `json:"option"`
	}

	DynamicLinkResponse struct {
		ShortLink   string `json:"shortLink"`
		PreviewLink string `json:"previewLink"`
	}

	// FBClient interface
	FBClient interface {
		GenerateDynamicLink(ctx context.Context, request DynamicLinkRequest, code string) (DynamicLinkResponse, error)
	}
)

// GenerateDynamicLink used to create firebase dynamic link.
func (i *Interactor) GenerateDynamicLink(ctx context.Context, request DynamicLinkRequest) (DynamicLinkResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return DynamicLinkResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://firebasedynamiclinks.googleapis.com/v1/shortLinks?key="+i.config.WebAPIKey, bytes.NewBuffer(jsonBody))
	if err != nil {
		return DynamicLinkResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.http.Do(req.WithContext(ctx))
	if err != nil {
		return DynamicLinkResponse{}, err
	}

	fmt.Println("qwe//////////////////////") // TODO:REMOVE IT!!!
	qwe, _ := ioutil.ReadAll(resp.Body)      // TODO:REMOVE IT!!!
	fmt.Println(string(qwe))                 // TODO:REMOVE IT!!!

	defer func() {
		err = errs.Combine(err, resp.Body.Close())
	}()

	if resp.StatusCode != http.StatusOK {
		return DynamicLinkResponse{}, err
	}

	var dynamicLinkResponse DynamicLinkResponse
	if err = json.NewDecoder(resp.Body).Decode(&dynamicLinkResponse); err != nil {
		return DynamicLinkResponse{}, err
	}
	return dynamicLinkResponse, nil
}
