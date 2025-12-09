package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Registrar validações customizadas se necessário
	// validate.RegisterValidation("custom_tag", customValidationFunc)
}

// Validate valida uma struct usando tags de validação
func Validate(s interface{}) []string {
	var errors []string
	
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, getErrorMessage(err))
		}
	}
	
	return errors
}

// getErrorMessage retorna uma mensagem de erro amigável baseada no tipo de validação
func getErrorMessage(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()
	
	// Mensagens customizadas por tag
	switch tag {
	case "required":
		return fmt.Sprintf("O campo '%s' é obrigatório", toLowerFirst(field))
	case "min":
		return fmt.Sprintf("O campo '%s' deve ser no mínimo %s", toLowerFirst(field), param)
	case "max":
		return fmt.Sprintf("O campo '%s' deve ser no máximo %s", toLowerFirst(field), param)
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um email válido", toLowerFirst(field))
	case "gte":
		return fmt.Sprintf("O campo '%s' deve ser maior ou igual a %s", toLowerFirst(field), param)
	case "lte":
		return fmt.Sprintf("O campo '%s' deve ser menor ou igual a %s", toLowerFirst(field), param)
	case "gt":
		return fmt.Sprintf("O campo '%s' deve ser maior que %s", toLowerFirst(field), param)
	case "lt":
		return fmt.Sprintf("O campo '%s' deve ser menor que %s", toLowerFirst(field), param)
	default:
		return fmt.Sprintf("O campo '%s' é inválido", toLowerFirst(field))
	}
}

// toLowerFirst converte a primeira letra para minúscula
func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ValidateStruct valida uma struct e retorna erro se houver problemas
func ValidateStruct(s interface{}) error {
	errors := Validate(s)
	if len(errors) > 0 {
		return fmt.Errorf("validação falhou: %s", strings.Join(errors, "; "))
	}
	return nil
}

