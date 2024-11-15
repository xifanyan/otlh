package otlh

import (
	"encoding/json"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
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

func getAllEntities[T any](c *Client, req Requestor, opts Options, unmarshal func([]byte) ([]T, bool, int, error)) ([]T, error) {
	var entities []T
	var page int = 1

	log.Debug().Msg("Getting all entities")

	bar := progressbar.Default(100)
	defer bar.Finish()

	for {
		opts.(*ListOptions).WithPageNumber(page)

		resp, err := c.Send(req, opts)
		if err != nil {
			return nil, err
		}

		pageEntities, hasMore, totalCount, err := unmarshal(resp)
		if err != nil {
			return nil, err
		}

		bar.ChangeMax(totalCount)
		bar.Add(opts.(*ListOptions).pageSize)
		time.Sleep(5 * time.Millisecond)

		entities = append(entities, pageEntities...)
		if !hasMore {
			break
		}

		page++
	}

	return entities, nil
}

func unmarshalCustodians(data []byte) ([]Custodian, bool, int, error) {
	var resp CustodiansResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Custodians, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalCustodianGroups(data []byte) ([]CustodianGroup, bool, int, error) {
	var resp CustodianGroupsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.CustodianGroups, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalGroups(data []byte) ([]Group, bool, int, error) {
	var resp GroupsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Groups, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalFolders(data []byte) ([]Folder, bool, int, error) {
	var resp FoldersResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Folders, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalMatters(data []byte) ([]Matter, bool, int, error) {
	var resp MattersResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Matters, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalLegalholds(data []byte) ([]Legalhold, bool, int, error) {
	var resp LegalholdsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Legalholds, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalSilentholds(data []byte) ([]Silenthold, bool, int, error) {
	var resp SilentholdsResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Silentholds, resp.Page.HasMore, resp.Page.TotalCount, err
}

func unmarshalQuestionnaires(data []byte) ([]Questionnaire, bool, int, error) {
	var resp QuestionnairesResponse
	err := json.Unmarshal(data, &resp)
	return resp.Embedded.Questionnaires, resp.Page.HasMore, resp.Page.TotalCount, err
}
