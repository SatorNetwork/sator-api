# KYC

Interface of KYC provider.

> *Documentation is under development*


## Available KYC methods

```
Get(kycID string) (*Response, error)
GetReport(kycID, reviewID string) (map[string][]byte, error)
```


## Implemented adapters

- SumSub (https://developers.sumsub.com/api-reference/)

## Usage

```
import "gitlab.maiv.biz/exprowhite/kyc"

...

var(
    appToken = "sumsub app token"
    appSecret = "sumsub app secret"
    baseURL = "sumsub base url"
)
c := kyc.New(appToken, appSecret, baseURL)

resp, err := c.Get(subSubID) // returns sumsub user data
...
report, err := c.GetClientReport(subSubID) // returns sumsub user report, e.g.: required documents,
...
```