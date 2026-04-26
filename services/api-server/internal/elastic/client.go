package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	es8 "github.com/elastic/go-elasticsearch/v8"

	"api-server/internal/model"
)

const logsIndex = "logs"

type Client struct {
	es *es8.Client
}

func NewClient(url string) (*Client, error) {
	c, err := es8.NewClient(es8.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}
	return &Client{es: c}, nil
}

// Search queries Elasticsearch with optional full-text (q), level, and service filters.
func (c *Client) Search(ctx context.Context, q, level, service string) ([]model.Log, error) {
	query := buildQuery(q, level, service)

	body, err := json.Marshal(map[string]interface{}{"query": query})
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(logsIndex),
		c.es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		raw, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch error response: %s", raw)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source model.Log `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	logs := make([]model.Log, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		logs = append(logs, h.Source)
	}
	return logs, nil
}

func buildQuery(q, level, service string) map[string]interface{} {
	var mustClauses []interface{}
	var filterClauses []interface{}

	if q != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  q,
				"fields": []string{"message", "service_name"},
			},
		})
	} else {
		mustClauses = append(mustClauses, map[string]interface{}{"match_all": map[string]interface{}{}})
	}

	if level != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{"level": level},
		})
	}
	if service != "" {
		filterClauses = append(filterClauses, map[string]interface{}{
			"term": map[string]interface{}{"service_name": service},
		})
	}

	return map[string]interface{}{
		"bool": map[string]interface{}{
			"must":   mustClauses,
			"filter": filterClauses,
		},
	}
}
