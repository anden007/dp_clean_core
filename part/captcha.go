package part

import (
	"github.com/anden007/dp_clean_core/pkg"

	"github.com/mojocn/base64Captcha"
)

type Captcha struct {
	captcha *base64Captcha.Captcha
}

type ICaptcha interface {
	New() (result *pkg.CaptchaData, err error)
	Reload(id string) (result *pkg.CaptchaData, err error)
	Draw(id string) (result *pkg.CaptchaData, err error)
	Verify(id string, digits string) (result bool, err error)
}

func NewCaptcha(cache ICache) ICaptcha {
	instance := new(Captcha)
	store := NewRedisStore(cache)
	driver := base64Captcha.NewDriverDigit(32, 110, 4, 0.1, 10)
	instance.captcha = base64Captcha.NewCaptcha(driver, store)
	return instance
}

func (m *Captcha) New() (result *pkg.CaptchaData, err error) {
	err = nil
	id, digits, answer := m.captcha.Driver.GenerateIdQuestionAnswer()
	item, err := m.captcha.Driver.DrawCaptcha(digits)
	if err == nil {
		result = &pkg.CaptchaData{
			ID:          id,
			Digits:      answer,
			Base64Image: item.EncodeB64string(),
		}
		m.captcha.Store.Set(result.ID, result.Digits)
	}
	return result, err
}

func (m *Captcha) Reload(id string) (result *pkg.CaptchaData, err error) {
	err = nil
	_, digits, answer := m.captcha.Driver.GenerateIdQuestionAnswer()
	item, err := m.captcha.Driver.DrawCaptcha(digits)
	if err == nil {
		result = &pkg.CaptchaData{
			ID:          id,
			Digits:      answer,
			Base64Image: item.EncodeB64string(),
		}
		m.captcha.Store.Set(result.ID, result.Digits)
	}
	return result, err
}

func (m *Captcha) Verify(id string, digits string) (result bool, err error) {
	result = m.captcha.Store.Verify(id, digits, true)
	return result, err
}

func (m *Captcha) Draw(id string) (result *pkg.CaptchaData, err error) {
	digits := m.captcha.Store.Get(id, false)
	if item, drawErr := m.captcha.Driver.DrawCaptcha(digits); drawErr == nil {
		result = &pkg.CaptchaData{
			ID:          id,
			Digits:      digits,
			Base64Image: item.EncodeB64string(),
		}
	} else {
		err = drawErr
	}
	return
}
