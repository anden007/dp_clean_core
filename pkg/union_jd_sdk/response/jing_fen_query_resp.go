package response

type UnionOpenGoodsJingfenQueryResponse struct {
	Code       int32          `json:"code"`
	Message    string         `json:"message"`
	RequestId  string         `json:"requestId"`
	TotalCount int64          `json:"totalCount"`
	Data       []*JFGoodsResp `json:"data"`
}

// 数据明细
type JFGoodsResp struct {
	CategoryInfo          *CategoryInfo   `json:"categoryInfo"`          // 类目信息
	Comments              int64           `json:"comments"`              // 评论数
	CommissionInfo        *CommissionInfo `json:"commissionInfo"`        // 佣金信息
	CouponInfo            *CouponInfo     `json:"couponInfo"`            // 优惠券信息，返回内容为空说明该SKU无可用优惠券
	GoodCommentsShare     float64         `json:"goodCommentsShare"`     // 商品好评率
	ImageInfo             *ImageInfo      `json:"imageInfo"`             // 图片信息
	InOrderCount30Days    int64           `json:"inOrderCount30Days"`    // 30天内引单数量
	MaterialUrl           string          `json:"materialUrl"`           // 商品落地页
	PriceInfo             *PriceInfo      `json:"priceInfo"`             // 价格信息
	ShopInfo              *ShopInfo       `json:"shopInfo"`              // 店铺信息
	SkuId                 int64           `json:"skuId"`                 // 商品ID
	SkuName               string          `json:"skuName"`               // 商品名称
	IsHot                 uint8           `json:"isHot"`                 // 是否爆款，1：是，0：否
	Spuid                 float64         `json:"spuid"`                 // spuid，其值为同款商品的主skuid
	BrandCode             string          `json:"brandCode"`             // 品牌code
	BrandName             string          `json:"brandName"`             // 品牌名
	Owner                 string          `json:"owner"`                 // g=自营，p=pop
	PinGouInfo            *PinGouInfo     `json:"pinGouInfo"`            // 拼购信息
	ResourceInfo          *ResourceInfo   `json:"resourceInfo"`          // 资源信息
	InOrderCount30DaysSku int64           `json:"inOrderCount30DaysSku"` // 30天引单数量(sku维度)
	SeckillInfo           *SeckillInfo    `json:"seckillInfo"`           // 秒杀信息(可选）
	JxFlags               []int32         `json:"jxFlags"`               // 京喜商品类型，1京喜、2京喜工厂直供、3京喜优选（包含3时可在京东APP购买）(可选）
	VideoInfo             *VideoInfo      `json:"videoInfo"`             // 视频信息(可选）
	DocumentInfo          *DocumentInfo   `json:"documentInfo"`          // 段子信息(可选）
}
