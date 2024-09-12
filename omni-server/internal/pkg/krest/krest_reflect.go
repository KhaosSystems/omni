/*
* This file contains reflection helpers for the krest package.
 */
package krest

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
* ReflectExpandableFields returns a map of expandable fields in a struct.
* The keys are the JSON names of the fields, and the values are the struct fields.
 */
func ReflectExpandableFields[T any]() (map[string]reflect.StructField, error) {
	expandableFields := make(map[string]reflect.StructField)

	typ := reflect.TypeOf((*T)(nil)).Elem() // Get the type of T without needing a value.

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not a struct", typ.Name())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("krest") == "expandable" {
			jsonTokens := strings.Split(field.Tag.Get("json"), ",")
			if len(jsonTokens) <= 0 {
				return nil, errors.New("expandable props NEED to have a json name")
			}

			expandableFields[jsonTokens[0]] = field
		}
	}

	// Reflection stuff
	return expandableFields, nil
}

/*
* ReflectNonExpandableFields returns a map of non-expandable fields in a struct type.
* The keys are the JSON names of the fields, and the values are the struct fields.
 */
func ReflectNonExpandableFields[T any]() (map[string]reflect.StructField, error) {
	nonExpandableFields := make(map[string]reflect.StructField)

	typ := reflect.TypeOf((*T)(nil)).Elem() // Get the type of T without needing a value.

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not a struct", typ.Name())
	}

	// Iterate over the fields of the struct.
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("krest") != "expandable" {
			jsonTokens := strings.Split(field.Tag.Get("json"), ",")
			if len(jsonTokens) == 0 || jsonTokens[0] == "" {
				return nil, errors.New("non-expandable props NEED to have a json name")
			}

			nonExpandableFields[jsonTokens[0]] = field
		}
	}

	return nonExpandableFields, nil
}
