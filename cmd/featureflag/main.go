package main

import (
	"errors"
	"reflect"
	"sync"
)

var ErrPtrExpected = errors.New("expected ptr value")

type FeatureToggleNamesProvider interface {
	GetFeatureToggle() *FeatureToggle
}

type FeatureToggle struct {
	parentStruct           interface{}
	featureToggleFieldsMap map[string]string
	initDone               bool
	initLock               sync.Mutex
}

func MakeFeatureToggle(keysProvider FeatureToggleNamesProvider) error {
	ft := keysProvider.GetFeatureToggle()
	p := getValueForStructOrPtr(reflect.ValueOf(keysProvider))
	if p.FieldByName("FeatureToggle").Type().Kind() != reflect.Ptr {
		return ErrPtrExpected
	}
	ft.parentStruct = keysProvider
	ft.featureToggleFieldsMap = make(map[string]string)
	return nil
}

func (ft *FeatureToggle) IsFeatureEnabled(featureKey string) bool {
	if !ft.initDone {
		ft.initLock.Lock()
		if !ft.initDone {
			ft.InitFieldMap()
			ft.initDone = true
		}
		ft.initLock.Unlock()
	}
	fieldName, ok := ft.featureToggleFieldsMap[featureKey]
	if !ok {
		return false
	}
	parentStruct := getValueForStructOrPtr(reflect.ValueOf(ft.parentStruct))
	field := parentStruct.FieldByName(fieldName)
	if field.Type().Kind() != reflect.Bool {
		return false
	}
	return field.Bool()
}

func (ft *FeatureToggle) InitFieldMap() {
	featureToggleFieldsMap := ft.featureToggleFieldsMap
	eType := reflect.TypeOf(ft.parentStruct)
	elem := getElementForStructOrPtr(eType)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		tag := field.Tag.Get("json")
		if field.Type.Kind() != reflect.Bool {
			continue
		}
		featureToggleFieldsMap[tag] = field.Name
	}
}

func getElementForStructOrPtr(eType reflect.Type) reflect.Type {
	if eType.Kind() == reflect.Ptr {
		return eType.Elem()
	}
	return eType
}

func getValueForStructOrPtr(value reflect.Value) reflect.Value {
	if value.Type().Kind() == reflect.Ptr {
		return value.Elem()
	}
	return value
}
