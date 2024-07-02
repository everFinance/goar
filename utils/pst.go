package utils

import "github.com/everVision/goar/schema"

func PstTransferTags(contractId string, target string, qty int64, warp bool) ([]schema.Tag, error) {
	input := schema.Input{
		"function": "transfer",
		"target":   target,
		"qty":      qty,
	}

	inputStr, err := input.ToString()
	if err != nil {
		return nil, err
	}

	pstTags := []schema.Tag{
		{Name: "App-Name", Value: "SmartWeaveAction"},
		{Name: "App-Version", Value: "0.3.0"},
		{Name: "Contract", Value: contractId},
		{Name: "Input", Value: inputStr},
	}

	if warp {
		pstTags = append(pstTags, schema.Tag{Name: "SDK", Value: "Warp"})
	}

	return pstTags, nil
}
