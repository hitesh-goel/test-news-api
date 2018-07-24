package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/hitesh-goel/test-news-api/apireq"
	"github.com/hitesh-goel/test-news-api/rediscon"
)

const newsAPI = "https://newsapi.org/v2/top-headlines?"

//Response struct for API
type Response struct {
	Country        string `json:"country"`
	Category       string `json:"category"`
	FilterKeyworld string `json:"filter_keyworld"`
	NewsTitle      string `json:"news_title"`
	Description    string `json:"description"`
	SourceNewsURL  string `json:"source_news_url"`
}

//NewsResponse response for this API
type NewsResponse struct {
	Status   string     `json:"status"`
	Articles []Response `json:"articles,omitempty"`
	Error    string     `json:"error,omitempty"`
}

//NewsAPIResponse response from news-api
type NewsAPIResponse struct {
	Status   string    `json:"status"`
	Articles []Article `json:"articles"`
}

//Article single article from news-api
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

//Filter applied filters to string
type Filter struct {
	Country  string
	Category string
	Keyworld string
}

func (f Filter) getRedisKey() string {
	return strings.ToLower(f.Category) + "-" + strings.ToLower(f.Country)
}

func (f Filter) queryFromRedis() (string, error) {
	redisKey := f.getRedisKey()
	r := rediscon.GetRedisService()
	res, err := r.Get(redisKey).Result()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (f Filter) updateDataInRedis(apiResp *NewsAPIResponse) error {

	resp, err := json.Marshal(apiResp)
	if err != nil {
		return err
	}

	redisKey := f.getRedisKey()
	r := rediscon.GetRedisService()

	//Added 10 min caching
	err = r.Set(redisKey, resp, time.Minute*time.Duration(10)).Err()

	if err == redis.Nil {
		return errors.New("error while updating the records")
	}
	return nil
}

func (f Filter) fetchDataFromNewsAPI(resp *NewsAPIResponse) error {
	apiKey := os.Getenv("APT_KEY")
	if apiKey == "" {
		return errors.New("not a valid auth")
	}

	//Create url from filters
	api := fmt.Sprintf("%scountry=%s&category=%s&apiKey=%s", newsAPI, f.Country, f.Category, apiKey)
	err := apireq.GetAPIRequests(api, resp)
	if err != nil {
		return err
	}
	err = f.updateDataInRedis(resp)
	if err != nil {
		return err
	}

	return nil
}

func (f Filter) filterResults(resp *NewsAPIResponse) (*NewsResponse, error) {
	var articles []Response
	for _, val := range resp.Articles {
		articles = append(articles, Response{
			SourceNewsURL:  val.URL,
			Category:       f.Category,
			Country:        f.Country,
			FilterKeyworld: f.Keyworld,
			NewsTitle:      val.Title,
			Description:    val.Description,
		})
	}
	response := &NewsResponse{
		Status:   resp.Status,
		Articles: articles,
	}

	if f.Keyworld != "" {
		searchKeywordInResult(response)
	}

	return response, nil
}

func searchKeywordInResult(resp *NewsResponse) NewsResponse {
	response := *resp
	return response
}

func (f Filter) getNewsResults() (*NewsResponse, error) {
	res, err := f.queryFromRedis()
	resp := &NewsAPIResponse{}
	if err != nil {
		err = f.fetchDataFromNewsAPI(resp)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(res), resp)
	}
	return f.filterResults(resp)
}

//GetNews returns the filtered News
func GetNews(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()

	f := Filter{
		Country:  q.Get("country"),
		Category: q.Get("category"),
		Keyworld: q.Get("keyworld"),
	}

	if f.Category == "" || f.Country == "" {
		returnResponse(w, r, nil, errors.New("pass valid filter parameters"))
		return nil
	}

	resp, err := f.getNewsResults()

	returnResponse(w, r, resp, err)

	return nil
}

func returnResponse(w http.ResponseWriter, r *http.Request, resp *NewsResponse, err error) {
	if err != nil {
		resp = &NewsResponse{}
		resp.Error = err.Error()
		resp.Status = "200"
	}
	b, _ := json.MarshalIndent(resp, "", "\t")
	w.Write(b)
}
