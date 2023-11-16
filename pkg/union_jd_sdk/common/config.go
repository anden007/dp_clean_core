package common

type Config struct {
	Version     string `url:"v"`
	Method      string `url:"method"`
	AccessToken string `url:"access_token,omitempty"`
	AppKey      string `url:"app_key"`
	Format      string `url:"format"`
	Timestamp   string `url:"timestamp"`
	Sign        string `url:"sign"`
	ParamJson   string `url:"360buy_param_json"`
}
