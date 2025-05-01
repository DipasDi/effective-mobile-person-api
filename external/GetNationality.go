package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetNationality(name string) (*string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		log.Warn().Err(err).Msg("failed to make API request")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg("API returned status")
		return nil, err
	}

	var result struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error().Msg("Invalid JSON received")
		return nil, err
	}

	if len(result.Country) == 0 {
		return nil, nil
	}

	nationality := result.Country[0].CountryID
	return &nationality, nil
}
