// internal/utils/validation.go
package utils

import (
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

func ValidateStruct(s interface{}) error {
    v := reflect.ValueOf(s)
    t := reflect.TypeOf(s)
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        tag := fieldType.Tag.Get("validate")
        
        if tag == "" {
            continue
        }
        
        rules := strings.Split(tag, ",")
        for _, rule := range rules {
            if err := validateRule(field, rule, fieldType.Name); err != nil {
                return err
            }
        }
    }
    
    return nil
}

func validateRule(field reflect.Value, rule, fieldName string) error {
    switch {
    case rule == "required":
        if field.Kind() == reflect.String && field.String() == "" {
            return fmt.Errorf("%s is required", fieldName)
        }
    case rule == "email":
        if field.Kind() == reflect.String {
            email := field.String()
            emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
            if !emailRegex.MatchString(email) {
                return fmt.Errorf("%s must be a valid email", fieldName)
            }
        }
    case strings.HasPrefix(rule, "min="):
        minLenStr := strings.TrimPrefix(rule, "min=")
        minLen, err := strconv.Atoi(minLenStr)
        if err != nil {
            return fmt.Errorf("invalid min validation rule")
        }
        if field.Kind() == reflect.String && len(field.String()) < minLen {
            return fmt.Errorf("%s must be at least %d characters", fieldName, minLen)
        }
    }
    return nil
}