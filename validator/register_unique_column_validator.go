package validator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"template-go-api/config"
)

func RegisterUniqueColumnValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("unique_column", func(fl validator.FieldLevel) bool {
			param := fl.Param()
			parts := strings.Fields(param)
			if len(parts) != 2 && len(parts) != 3 {
				return false
			}

			table := parts[0]
			column := parts[1]
			value := fl.Field().Interface()

			query := config.DB.Table(table).Where(fmt.Sprintf("%s = ?", column), value)

			if len(parts) == 3 {
				exceptField := fl.Parent().FieldByName(parts[2])
				if exceptField.IsValid() && exceptField.String() != "" {
					query = query.Where("id != ?", exceptField.Interface())
				}
			}

			var count int64
			if err := query.Count(&count).Error; err != nil {
				return false
			}

			return count == 0
		})

		if err != nil {
			return
		}
	}
}
