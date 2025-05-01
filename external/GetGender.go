package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetGender(name string) (*string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
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
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Error().Msg("Invalid JSON received")
		return nil, err
	}

	return &result.Gender, nil
}
