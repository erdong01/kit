package db

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

// SetDB uses reflection to assign a *gorm.DB value to the specified field of Dbs
func SetDB(d any, fieldName string, db *gorm.DB) error {
	// Check if the Dbs struct is nil
	if d == nil {
		return errors.New("Dbs struct is nil")
	}

	// Check if the gorm.DB pointer is nil
	if db == nil {
		return errors.New("gorm.DB pointer is nil")
	}

	// Get the reflect.Value of d
	v := reflect.ValueOf(d).Elem()

	// Ensure the specified field exists
	fieldVal := v.FieldByName(fieldName)
	if !fieldVal.IsValid() {
		return errors.New("field does not exist")
	}

	// Ensure the field can be set
	if !fieldVal.CanSet() {
		return errors.New("cannot set field")
	}

	// Ensure the field is of the correct types
	if fieldVal.Type() != reflect.TypeOf(db) {
		return errors.New("provided value types does not match field types")
	}

	// Set the field
	fieldVal.Set(reflect.ValueOf(db))
	return nil
}
