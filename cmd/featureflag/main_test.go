package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SomeConfig struct {
	Param *string `spex_config:"param"`
}

type AppConfig struct {
	SomeConfig *SomeConfig `spex_config:"some_config"`
}

type Config struct {
	AppConfig *AppConfig `spex_config:"app_config"`
}

type ConfigUser struct {
	Param *string
}

// passing param as AppConfig.SomeConfig.Param
func NewConfigUser(param *string) *ConfigUser {
	return &ConfigUser{Param: param}
}

func (cu ConfigUser) PrintParam() {
	fmt.Println(cu.Param)
}



type CustomConfig struct {
	MyFeatureA bool `json:"my_feature_a"`
	MyFeatureB bool `json:"my_feature_b"`
	*FeatureToggle
}

type CustomConfigB struct {
	MyFeatureA bool `json:"my_feature_a"`
	MyFeatureB bool `json:"my_feature_b"`
	FeatureToggle
}

func (c CustomConfigB) GetFeatureToggle() *FeatureToggle {
	return &c.FeatureToggle
}

func (c CustomConfig) GetFeatureToggle() *FeatureToggle {
	return c.FeatureToggle
}

func TestFeatureToggle_IsFeatureEnabled(t *testing.T) {
	cf := &CustomConfig{
		MyFeatureA:    false,
		MyFeatureB:    true,
		FeatureToggle: &FeatureToggle{},
	}
	err := MakeFeatureToggle(cf)
	assert.Nil(t, err)
	assert.False(t, cf.IsFeatureEnabled("my_feature_a"))
	assert.True(t, cf.IsFeatureEnabled("my_feature_b"))
	assert.False(t, cf.IsFeatureEnabled("my_feature_d"))
}

func TestFeatureToggle_IsFeatureEnabled_When_CustomConfig_Is_Not_Ptr(t *testing.T) {
	cf := CustomConfig{
		MyFeatureA:    false,
		MyFeatureB:    true,
		FeatureToggle: &FeatureToggle{},
	}
	err := MakeFeatureToggle(cf)
	assert.Nil(t, err)
	assert.False(t, cf.IsFeatureEnabled("my_feature_a"))
	assert.True(t, cf.IsFeatureEnabled("my_feature_b"))
	assert.False(t, cf.IsFeatureEnabled("my_feature_d"))
}

func TestFeatureToggle_IsFeatureEnabled_When_FeatureToggle_Is_Not_Ptr(t *testing.T) {
	cf := CustomConfigB{
		MyFeatureA:    false,
		MyFeatureB:    true,
		FeatureToggle: FeatureToggle{},
	}
	err := MakeFeatureToggle(cf)
	assert.EqualError(t, err, ErrPtrExpected.Error())
}
