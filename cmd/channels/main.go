package main

import (
	"encoding/json"
	"fmt"
)

type logisticInfoExtraData struct {
	// shipping fee after applying the discounts
	DiscountedShippingFees *struct {
		Shopee *int64 `json:"shopee"`
		Seller *int64 `json:"seller"`
	} `json:"discounted_shipping_fees"`

	// values of discount provided
	DiscountOnShippingFees *struct {
		Shopee *int64 `json:"shopee"`
		Seller *int64 `json:"seller"`
	} `json:"discount_shipping_fees"`
}

func decodeLogisticInfoExtraData(extraData []byte) (logisticInfoExtraData, error) {
	d := logisticInfoExtraData{}
	err := json.Unmarshal(extraData, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}

func main() {
	testJson := `{
                "discounted_shipping_fees": {
                    "shopee": 1000000000, 
                    "seller": 1000000000
                }, 
                "discount_shipping_fees": {
                    "shopee": 0
                }, 
                "extra_flag": 0
            }`
	v, err := decodeLogisticInfoExtraData([]byte(testJson))
	if err != nil {
		panic(err)
	}
	if v.DiscountOnShippingFees.Seller == nil {
		fmt.Println("hello nil")
	}
	fmt.Println(*v.DiscountedShippingFees.Seller)
}
