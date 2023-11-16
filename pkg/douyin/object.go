package douyin

type ApiResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type AuthData struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiresIn  int64  `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	UserId                int64  `json:"user_id"`
}

type AuthResult struct {
	ApiResult
	Data AuthData `json:"data"`
}

type AdvertiserInfoList struct {
	List []AdvertiserInfo `json:"list"`
}

type CompanyInfo struct {
	CustomerCompanyId   int64  `json:"customer_company_id"`   // 客户公司id
	CustomerCompanyName string `json:"customer_company_name"` // 客户公司名
}

type AdvertiserInfo struct {
	AdvertiserId   int64         `json:"advertiser_id"`   // 账号id
	AdvertiserName string        `json:"advertiser_name"` // 账号名称
	AdvertiserRole int64         `json:"advertiser_role"` // 旧版账号角色，1-普通广告主，2-纵横组织账户，3-一级代理商，4-二级代理商，6-星图账号
	IsValid        bool          `json:"is_valid"`        // 授权有效性，允许值：true/false；false表示对应的user在客户中心/一站式平台代理商平台变更了对此账号的权限,需要到对应平台进行调整过来；
	AccountRole    string        `json:"account_role"`    // 新版账号角色，详见：https://open.oceanengine.com/labels/7/docs/1696710760171535
	CompanyList    []CompanyInfo `json:"company_list"`
}

type AdvertiserInfoResult struct {
	ApiResult
	Data AdvertiserInfoList `json:"data"`
}

type ADPlanInfo struct {
	ApiResult
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ADPlanInfoResult struct {
	ApiResult
	Data struct {
		List []ADPlanInfo `json:"list"`
	} `json:"data"`
}

type ADGroupInfo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ADGroupInfoResult struct {
	ApiResult
	Data struct {
		List []ADGroupInfo `json:"list"`
	} `json:"data"`
}

type CreativeInfoDetail struct {
	CreativeId int64  `json:"creative_id"`
	Title      string `json:"title"`
	Status     string `json:"status"`
}
type CreativeInfoResult struct {
	ApiResult
	Data struct {
		List []CreativeInfoDetail `json:"list"`
	} `json:"data"`
}
