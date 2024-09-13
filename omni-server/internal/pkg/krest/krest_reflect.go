/*
* This file contains reflection helpers for the krest package.
 */
package krest

import (
	"fmt"
	"reflect"
)

/*
* ReflectExpandableFields returns a map of expandable fields in a struct.
* The keys are the JSON names of the fields, and the values are the struct fields.
 */
func ReflectExpandableFields[T any]() ([]reflect.StructField, error) {
	fields := []reflect.StructField{}

	typ := reflect.TypeOf((*T)(nil)).Elem() // Get the type of T without needing a value.

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not a struct", typ.Name())
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("krest") == "expandable" {
			fields = append(fields, field)
		}
	}

	// Reflection stuff
	return fields, nil
}

/*
* ReflectNonExpandableFields returns a map of non-expandable fields in a struct type.
* The keys are the JSON names of the fields, and the values are the struct fields.
 */
func ReflectNonExpandableFields[T any]() ([]reflect.StructField, error) {
	fields := []reflect.StructField{}

	typ := reflect.TypeOf((*T)(nil)).Elem() // Get the type of T without needing a value.
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not a struct", typ.Name())
	}

	// Iterate over the fields of the struct.
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Tag.Get("krest") != "expandable" {
			fields = append(fields, field)
		}
	}

	return fields, nil
}
