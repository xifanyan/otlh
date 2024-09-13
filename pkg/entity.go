package otlh

import (
	"encoding/json"
)

type DefaultEntityListInfo struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Page struct {
		HasMore    bool `json:"has-more"`
		TotalCount int  `json:"total-count"`
	} `json:"page"`
}

func getAllEntities[T any](c *Client, req Requestor, opts Options, unmarshal func([]byte) ([]T, bool, error)) ([]T, error) {
	var entities []T
	var page int = 1

	for {
		opts.(*ListOptions).WithPageNumber(page)

		resp, err := c.Send(req, opts)
		if err != nil {
			return nil, err
		}

		pageEntities, hasMore, err := unmarshal(resp)
		if err != nil {
			return nil, err
		}

		entities = append(entities, pageEntities...)

		if !hasMore {
			break
		}

		page++
	}

	return entities, nil
}

func unmarshalCustodians(data []byte) ([]Custodian, bool, error) {
	var resp CustodiansResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Custodians, resp.Page.HasMore, err
}

func unmarshalCustodianGroups(data []byte) ([]CustodianGroup, bool, error) {
	var resp CustodianGroupsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.CustodianGroups, resp.Page.HasMore, err
}

func unmarshalGroups(data []byte) ([]Group, bool, error) {
	var resp GroupsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Groups, resp.Page.HasMore, err
}

func unmarshalFolders(data []byte) ([]Folder, bool, error) {
	var resp FoldersResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Folders, resp.Page.HasMore, err
}

func unmarshalMatters(data []byte) ([]Matter, bool, error) {
	var resp MattersResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Matters, resp.Page.HasMore, err
}

func unmarshalLegalholds(data []byte) ([]Legalhold, bool, error) {
	var resp LegalholdsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Legalholds, resp.Page.HasMore, err
}

func unmarshalSilentholds(data []byte) ([]Silenthold, bool, error) {
	var resp SilentholdsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Silentholds, resp.Page.HasMore, err
}
