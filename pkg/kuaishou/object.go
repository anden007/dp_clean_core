package kuaishou

type ApiResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type AuthData struct {
	AccessToken           string  `json:"access_token"`
	AccessTokenExpiresIn  int64   `json:"access_token_expires_in"`
	RefreshToken          string  `json:"refresh_token"`
	RefreshTokenExpiresIn int64   `json:"refresh_token_expires_in"`
	AdvertiserId          int64   `json:"advertiser_id"`
	AdvertiserIds         []int64 `json:"advertiser_ids"`
	UserId                int64   `json:"user_id"`
}

type AuthResult struct {
	ApiResult
	Data AuthData `json:"data"`
}

type AdvertiserInfo struct {
	UserId              int64  `json:"user_id"`
	CorporationName     string `json:"corporation_name"`
	UserName            string `json:"user_name"`
	PrimaryIndustryId   int64  `json:"primary_industry_id"`
	PrimaryIndustryName string `json:"primary_industry_name"`
	IndustryId          int64  `json:"industry_id"`
	IndustryName        string `json:"industry_name"`
}

type AdvertiserInfoResult struct {
	ApiResult
	Data AdvertiserInfo `json:"data"`
}

type ADPlanInfo struct {
	ApiResult
	CampaignId   int64  `json:"campaign_id"`
	CampaignName string `json:"campaign_name"`
	Status       int    `json:"status"`
}

type ADPlanInfoResult struct {
	ApiResult
	Data struct {
		TotalCount int64        `json:"total_count"`
		Details    []ADPlanInfo `json:"details"`
	} `json:"data"`
}

type ADGroupInfo struct {
	UnitId   int64  `json:"unit_id"`
	UnitName string `json:"unit_name"`
	Status   int    `json:"status"`
}
type ADGroupInfoResult struct {
	ApiResult
	Data struct {
		TotalCount int64         `json:"total_count"`
		Details    []ADGroupInfo `json:"details"`
	} `json:"data"`
}

type CreativeInfoDetail struct {
	CreativeId   int64  `json:"creative_id"`
	CreativeName string `json:"creative_name"`
	Status       int    `json:"status"`
}
type CreativeInfoResult struct {
	ApiResult
	Data struct {
		TotalCount int64                `json:"total_count"`
		Details    []CreativeInfoDetail `json:"details"`
	} `json:"data"`
}
