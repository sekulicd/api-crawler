package collegeapplication

import (
	"api-crawler/core/collegescorecard/collegedomain"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/go-resty/resty/v2"
	"strconv"
)

const (
	pageSize = 100
)

type collegeService struct {
	apiKey            string
	url               string
	restyClient       *resty.Client
	apiResponse       chan collegedomain.School
	collegeRepository collegedomain.CollegeRepository
}

func NewCollegeService(apiKey string, url string, collegeRepository collegedomain.CollegeRepository) CollegeService {
	return &collegeService{
		apiKey:            apiKey,
		url:               url,
		restyClient:       resty.New(),
		apiResponse:       make(chan collegedomain.School, 5),
		collegeRepository: collegeRepository,
	}
}

type CollegeService interface {
	GetAllSchools() ([]collegedomain.School, error)
	CrawlApi() error
}

type ApiFetcherService interface {
	CrawlApi() error
}

func (c collegeService) CrawlApi() error {

	go c.loadToDb()

	err := c.fetchAllFromApi()
	if err != nil {
		return err
	}

	return nil
}

func (c collegeService) fetchAllFromApi() error {
	pageNum := 0
	lastPage := 1

	bytes, err := c.callApi(pageNum)
	if err != nil {
		return err
	}

	metadata, schools, err := c.mapToDomain(pageNum, bytes)
	if err != nil {
		fmt.Println(err)
	}
	for _, school := range schools {
		c.apiResponse <- school
	}

	lastPage = setUpLastPage(pageNum, metadata.Total)

	for pageNum < lastPage {
		go c.callApiAndMapToDomain(pageNum)
		pageNum++
	}
	return nil
}

func (c *collegeService) callApi(pageNum int) ([]byte, error) {
	resp, err := c.restyClient.R().
		SetQueryParams(map[string]string{
			"api_key":                            c.apiKey,
			"school.degrees_awarded.predominant": "2,3",
			"fields":                             "id,school.name,school.city",
			"_page":                              strconv.Itoa(pageNum),
			"_per_page":                          strconv.Itoa(pageSize),
		}).
		Get(c.url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func (c *collegeService) callApiAndMapToDomain(pageNum int) {
	bytes, err := c.callApi(pageNum)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("page %v fetched\n", pageNum)

	_, schools, err := c.mapToDomain(pageNum, bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, school := range schools {
		c.apiResponse <- school
	}
}

func setUpLastPage(pageNum int, total int) int {
	var lastPage int
	if pageNum == 0 {
		if total%pageSize == 0 {
			lastPage = total / pageSize
		} else {
			lastPage = total/pageSize + 1
		}
	}
	return lastPage
}

func (c collegeService) mapToDomain(pageNum int, data []byte) (*collegedomain.Metadata, []collegedomain.School, error) {
	var metadata collegedomain.Metadata
	var schools = make([]collegedomain.School, 0)

	if pageNum == 0 {
		total, err := jsonparser.GetInt(data, "metadata", "total")
		if err != nil {
			return nil, nil, err
		}
		page, err := jsonparser.GetInt(data, "metadata", "page")
		if err != nil {
			return nil, nil, err
		}
		perPage, err := jsonparser.GetInt(data, "metadata", "per_page")
		if err != nil {
			return nil, nil, err
		}

		metadata = collegedomain.Metadata{
			Total:      int(total),
			PageNumber: int(page),
			PageSize:   int(perPage),
		}
	}

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		id, err := jsonparser.GetInt(value, "id")
		if err != nil {
			fmt.Println(err)
		}
		name, err := jsonparser.GetString(value, "school.name")
		if err != nil {
			fmt.Println(err)
		}
		city, err := jsonparser.GetString(value, "school.city")
		if err != nil {
			fmt.Println(err)
		}
		school := collegedomain.School{
			SchoolId: int(id),
			Name:     name,
			City:     city,
		}
		schools = append(schools, school)
	}, "results")
	if err != nil {
		return nil, nil, err
	}
	return &metadata, schools, nil
}

func (c collegeService) loadToDb() {
	for school := range c.apiResponse {
		err := c.collegeRepository.Create(school)
		if err != nil {
			fmt.Printf("Id: %v, error: %v", school.ID, err)
		}
	}
}

func (c collegeService) GetAllSchools() ([]collegedomain.School, error) {
	schools, err := c.collegeRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return schools, nil
}
