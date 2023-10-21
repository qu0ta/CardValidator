package errorValidator

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type JSONRespond struct {
	Ok      bool      `json:"ok"`
	Error   string    `json:"error"`
	Respond BinResult `json:"respond"`
}
type BinResult struct {
	Bin         string `json:"bin"`
	Bank        string `json:"bank"`
	Card        string `json:"card"`
	Type        string `json:"type"`
	Level       string `json:"level"`
	County      string `json:"county"`
	CountryCode string `json:"countrycode"`
	Website     string `json:"website"`
	Phone       string `json:"phone"`
	Valid       string `json:"valid"`
	Message     string `json:"message"`
	Error       string `json:"error"`
}

type ErrorValidator struct {
	Err     error
	Message string
}
type OSENV struct {
	APIKEY string
}
type Card struct {
	CardNumber int `json:"card-number"`
}

func (e *ErrorValidator) Decode(r *http.Request) int {
	if e.Err != nil {
		return 0
	}

	var c Card
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		e.Err = err
		e.Message = "Error while decoding request: " + err.Error()
	}
	return c.CardNumber
}
func (e *ErrorValidator) GetApiKey() string {
	if e.Err != nil {
		return ""
	}

	if err := godotenv.Load(); err != nil {
		e.Err = err
		e.Message = "Error while loading .env values: " + err.Error()
	}
	apiKey := os.Getenv("API_TOKEN")
	return apiKey
}

func (e *ErrorValidator) SendRequest(baseUrl, apiKey string, cardNumber int) *http.Response {
	if e.Err != nil {
		return nil
	}
	requestUrl := baseUrl + "api_key=" + apiKey + "&" + "bin=" + strconv.Itoa(cardNumber)[:6]
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		e.Err = err
		e.Message = "Error while making a GET request: " + err.Error()
	}

	res, err := client.Do(req)
	if err != nil {
		e.Err = err
		e.Message = "Error while doing a GET request: " + err.Error()
	}
	defer res.Body.Close()
	return res
}

func (e *ErrorValidator) UnmarshalBin(res *http.Response) *BinResult {
	if e.Err != nil {
		return nil
	}
	var b BinResult
	body, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &b)
	if err != nil {
		e.Err = err
		if strings.Contains(err.Error(), "unexpected end of JSON input") {
			e.Message = "Daily free api usage count has expired"
		} else {
			e.Message = "Error while unmarshalling json response: " + err.Error()
		}
	}
	return &b
}
func (e *ErrorValidator) MakeJson(b *BinResult) JSONRespond {
	if e.Err != nil {
		return JSONRespond{Ok: false, Respond: BinResult{}}
	}

	var ok bool
	if b.Valid == "true" {
		ok = true
	}
	respond := JSONRespond{Ok: ok, Respond: BinResult{
		Bin:         b.Bin,
		Bank:        b.Bank,
		Card:        b.Card,
		Type:        b.Type,
		Level:       b.Level,
		County:      b.County,
		CountryCode: b.CountryCode,
		Website:     b.Website,
		Phone:       b.Phone,
		Valid:       b.Valid,
		Message:     b.Message,
		Error:       b.Error,
	}, Error: ""}
	return respond
}
