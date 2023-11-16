package response

type OrderRowResp struct {
	Id                  string        `json:"id"`                  // 标记唯一订单行：订单+sku维度的唯一标识
	OrderId             int64         `json:"orderId"`             // 订单号
	ParentId            int64         `json:"parentId"`            // 父单的订单号：如一个订单拆成多个子订单时，原订单号会作为父单号，拆分的订单号为子单号存储在orderid中。若未发生拆单，该字段为0
	OrderTime           string        `json:"orderTime"`           // 下单时间,格式yyyy-MM-dd HH:mm:ss
	FinishTime          string        `json:"finishTime"`          // 完成时间（购买用户确认收货时间）,格式yyyy-MM-dd HH:mm:ss
	ModifyTime          string        `json:"modifyTime"`          // 更新时间,格式yyyy-MM-dd HH:mm:ss
	OrderEmt            int           `json:"orderEmt"`            // 下单设备 1.pc 2.无线
	Plus                int           `json:"plus"`                // 下单用户是否为PLUS会员 0：否，1：是
	UnionId             int64         `json:"unionId"`             // 推客ID
	SkuId               int64         `json:"skuId"`               // 商品ID
	SkuName             string        `json:"skuName"`             // 商品名称
	SkuNum              int           `json:"skuNum"`              // 商品数量
	SkuReturnNum        int           `json:"skuReturnNum"`        // 商品已退货数量
	SkuFrozenNum        int           `json:"skuFrozenNum"`        // 商品售后中数量
	Price               float64       `json:"price"`               // 商品单价
	CommissionRate      float64       `json:"commissionRate"`      // 佣金比例(投放的广告主计划比例)
	SubSideRate         float64       `json:"subSideRate"`         // 分成比例（单位：%）
	SubsidyRate         float64       `json:"subsidyRate"`         // 补贴比例（单位：%）
	FinalRate           float64       `json:"finalRate"`           // 最终分佣比例（单位：%）=分成比例+补贴比例
	EstimateCosPrice    float64       `json:"estimateCosPrice"`    // 预估计佣金额：由订单的实付金额拆分至每个商品的预估计佣金额，不包括运费，以及京券、东券、E卡、余额等虚拟资产支付的金额。该字段仅为预估值，实际佣金以actualCosPrice为准进行计算
	EstimateFee         float64       `json:"estimateFee"`         // 推客的预估佣金（预估计佣金额*佣金比例*最终比例），如订单完成前发生退款，此金额也会更新。
	ActualCosPrice      float64       `json:"actualCosPrice"`      // 实际计算佣金的金额。订单完成后，会将误扣除的运费券金额更正。如订单完成后发生退款，此金额会更新。
	ActualFee           float64       `json:"actualFee"`           // 推客分得的实际佣金（实际计佣金额*佣金比例*最终比例）。如订单完成后发生退款，此金额会更新。
	ValidCode           int           `json:"validCode"`           // sku维度的有效码（-1：未知,2.无效-拆单,3.无效-取消,4.无效-京东帮帮主订单,5.无效-账号异常,6.无效-赠品类目不返佣,7.无效-校园订单,8.无效-企业订单,9.无效-团购订单,11.无效-乡村推广员下单,13. 违规订单-其他,14.无效-来源与备案网址不符,15.待付款,16.已付款,17.已完成（购买用户确认收货）,19.无效-佣金比例为0,20.无效-此复购订单对应的首购订单无效,21.无效-云店订单，22.无效-PLUS会员佣金比例为0，23.无效-支付有礼，24.已付定金，25. 违规订单-流量劫持，26. 违规订单-流量异常，27. 违规订单-违反京东平台规则，28. 违规订单-多笔交易异常
	TraceType           int           `json:"traceType"`           // 同跨店：2同店 3跨店
	PositionId          int           `json:"positionId"`          // 推广位ID
	SiteId              int           `json:"siteId"`              // 应用id（网站id、appid、社交媒体id）
	UnionAlias          string        `json:"unionAlias"`          // PID所属母账号平台名称（原第三方服务商来源），两方分佣会有该值
	Pid                 string        `json:"pid"`                 // 格式:子推客ID_子站长应用ID_子推客推广位ID
	Cid1                int           `json:"cid1"`                // 一级类目id
	Cid2                int           `json:"cid2"`                // 二级类目id
	Cid3                int           `json:"cid3"`                // 三级类目id
	SubUnionId          string        `json:"subUnionId"`          // 子渠道标识，在转链时可自定义传入，格式要求：字母、数字或下划线，最多支持80个字符(需要联系运营开放白名单才能拿到数据)
	UnionTag            string        `json:"unionTag"`            // 联盟标签数据（32位整型二进制字符串：00000000000000000000000000000001。数据从右向左进行，每一位为1表示符合特征，第1位：红包，第2位：组合推广，第3位：拼购，第5位：有效首次购（0000000000011XXX表示有效首购，最终奖励活动结算金额会结合订单状态判断，以联盟后台对应活动效果数据报表https://union.jd.com/active为准）,第8位：复购订单，第9位：礼金，第10位：联盟礼金，第11位：推客礼金，第12位：京喜APP首购，第13位：京喜首购，第14位：京喜复购，第15位：京喜订单，第16位：京东极速版APP首购，第17位白条首购，第18位校园订单，第19位是0或1时，均代表普通订单，第20位：预售订单，第21位：学生订单，第22位：全球购订单 ，第23位：京喜拼拼首购订单，第24位：京喜拼拼复购订单 例如：00000000000000000000000000000001:红包订单，00000000000000000000000000000010:组合推广订单，00000000000000000000000000000100:拼购订单，00000000000000000000000000011000:有效首购，00000000000000000000000000000111：红包+组合推广+拼购等） 注：一个订单同时使用礼金和红包，仅礼金位数为1，红包位数为0
	PopId               int           `json:"popId"`               // 商家ID
	Ext1                string        `json:"ext1"`                // 推客生成推广链接时传入的扩展字段（需要联系运营开放白名单才能拿到数据）。
	PayMonth            int           `json:"payMonth"`            // 预估结算时间，订单完成后才会返回，格式：yyyyMMdd，默认：0。表示最新的预估结算日期。当payMonth为当前的未来时间时，表示该订单可结算；当payMonth为当前的过去时间时，表示该订单已结算
	CpActId             int           `json:"cpActId"`             // 招商团活动id：当商品参加了招商团会有该值，为0时表示无活动
	UnionRole           int           `json:"unionRole"`           // 站长角色：1 推客 2 团长 3内容服务商
	GiftCouponOcsAmount float64       `json:"giftCouponOcsAmount"` // 礼金分摊金额：使用礼金的订单会有该值
	GiftCouponKey       string        `json:"giftCouponKey"`       // 礼金批次ID：使用礼金的订单会有该值
	BalanceExt          string        `json:"balanceExt"`          // 计佣扩展信息，表示结算月:每月实际佣金变化情况，格式：{20191020:10,20191120:-2}，订单完成后会有该值
	Sign                string        `json:"sign"`                // 数据签名，用来核对出参数据是否被修改，入参fields中写入sign时返回
	ProPriceAmount      float64       `json:"proPriceAmount"`      // 价保赔付金额：订单申请价保或赔付的金额，实际计佣金额已经减去此金额，您无需处理
	Rid                 int           `json:"rid"`                 // 团长渠道ID，仅限招商团长管理渠道使用，团长开通权限后才可使用。
	GoodsInfo           *GoodsInfo    `json:"goodsInfo"`           // 商品信息，入参传入fields，goodsInfo获取
	CategoryInfo        *CategoryInfo `json:"categoryInfo"`        // 类目信息,入参传入fields，categoryInfo获取
	ExpressStatus       int           `json:"expressStatus"`       // 发货状态（10：待发货，20：已发货）
	ChannelId           int           `json:"channelId"`           // 渠道关系ID
}

type OrderRowQueryResult struct {
	UnionJDApiResult
	Data []*OrderRowResp `json:"data"`
}
