package dictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Define structure for parsing API response
type Definition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type DictionaryResponse struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

type Dictionary struct {
	DICTIONARY_API_URL string
}

func NewDictionary(apiUrl string) *Dictionary {
	return &Dictionary{
		DICTIONARY_API_URL: apiUrl,
	}
}

func (d *Dictionary) Search(word string) ([]DictionaryResponse, error) {
	urlSearchWord := fmt.Sprintf("%s/%s", d.DICTIONARY_API_URL, word)
	resp, err := http.Get(urlSearchWord)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var results []DictionaryResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (t *Dictionary) FormatDefinition(response []DictionaryResponse) string {
	var sb strings.Builder

	for _, entry := range response {
		sb.WriteString(fmt.Sprintf("📖 *Word:* _%s_\n\n", entry.Word))
		// Group definitions by part of speech
		posMap := make(map[string][]string)

		for _, meaning := range entry.Meanings {
			for _, def := range meaning.Definitions {
				if def.Example != "" {
					posMap[meaning.PartOfSpeech] = append(posMap[meaning.PartOfSpeech],
						fmt.Sprintf("%s\n      _Example: %s_", def.Definition, def.Example))
				} else {
					posMap[meaning.PartOfSpeech] = append(posMap[meaning.PartOfSpeech], def.Definition)
				}
			}
		}
		// Write formatted output
		for pos, defs := range posMap {
			sb.WriteString(fmt.Sprintf("📌 *%s*\n", strings.Title(pos)))
			for i, def := range defs {
				sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, def))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("\n────────────────────\n\n")
	}

	return sb.String()
}
