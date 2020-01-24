package ddd

import (
	en_locales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func (p *CommandProcessor) ConfigureValidator(customMessages map[string]string) error {
	p.Validate = validator.New()
	en := en_locales.New()
	universalTranslator := ut.New(en, en)
	englishTranslator, _ := universalTranslator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(p.Validate, englishTranslator)
	for tag,value := range customMessages {
		err := registerValidationMessage(p.Validate, tag, englishTranslator, value)
		if err != nil {
			return err
		}
	}
	p.Translator = englishTranslator
	return nil
}
func registerValidationMessage(v *validator.Validate, tag string, translator ut.Translator, text string) error {
	return v.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}