package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"backend/internal/lms"
)

type Client struct {
	baseUrl string
	accessToken string
	httpClient *http.Client
}

func NewClient(baseUrl, token string) *Client {
	return &Client{
		baseUrl: baseUrl, 
		accessToken: token.
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Clinet) FetchAllLearnAccounts() ([]learn.MotherLearnAccount, error) {
	var allLearnAccounts []learn.MotherLearnAccount
	pageToken := ""

	for {
		params := url.Values{}
		params.Set("pageSize", "1000")
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		apiURL := fmt.Sprintf("%s/users/v2?%s", c.baseURL, params.Encode())

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.accessToken)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiPesponse struct {
			UserProfiles []struct {
				UserID        string `json:"userId"`
				DepartmentID  string `json:"departmentId"`
				Status        int    `json:"status"`
				AddedDate     string `json:"addedDate"`
				LastLoginDate string `json:"lastLoginDate"`
				Fields        struct {
					FirstName string `json:"first_name"`
					LastName  string `json:"last_name"`
					Email     string `json:"email"`
					Login     string `json:"login"`
				} `json:"fields"`
			} `json:"userProfiles"`
			NextPageToken string `json:"nextPageToken"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, err
		}

		for _, p := range apiResp.UserProfiles {
			addedDate, _ := time.Parse(time.RFC3339, p.AddedDate)
			lastLoginDate, _ := time.Parse(time.RFC3339, p.LastLoginDate)

			allLearnAccounts = append(allLearners, lms.Learner{
				UserID:        p.UserID,
				DepartmentID:  p.DepartmentID,
				FirstName:     p.Fields.FirstName,
				LastName:      p.Fields.LastName,
				Email:         p.Fields.Email,
				Login:         p.Fields.Login,
				Status:        p.Status,
				AddedDate:     addedDate,
				LastLoginDate: lastLoginDate,
			})
		}

		if apiResp.NextPageToken == "" {
			break
		}
		pageToken = apiResp.NextPageToken
	}

	return allLearnAccounts, nil
}

func (c *Client) FetchAllCourseResults() ([]lms.CourseResult, error) {
	var allResults []lms.CourseResult
	pageToken := ""

	for {
		params := url.Values{}
		params.Set("pageSize", "1000")
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		params.Add("userStatuses[]", "1")

		apiURL := fmt.Sprintf("%s/learners/results?%s", c.baseURL, params.Encode())

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.accessToken)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiResp struct {
			Results []struct {
				UserID           string      `json:"userId"`
				CourseID         string      `json:"courseId"`
				CourseTitle      string      `json:"courseTitle"`
				EnrollmentID     string      `json:"enrollmentId"`
				AccessDate       string      `json:"accessDate"`
				DueDate          string      `json:"dueDate"`
				CompletionDate   string      `json:"completionDate"`
				TimeSpent        int         `json:"timeSpent"`
				Duration         int         `json:"duration"`
				CompletionStatus string      `json:"completionStatuses"`
				Progress         float64     `json:"progress"`
				AwardedScore     *float64    `json:"awardedScore"`
				IsOverdue        interface{} `json:"isOverdue"` 
				ViewsCount       int         `json:"viewsCount"`
			} `json:"results"`
			NextPageToken string `json:"nextPageToken"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, err
		}

		for _, r := range apiResp.Results {
			accessDate, _ := time.Parse(time.RFC3339, r.AccessDate)
			
			var dueDate *time.Time
			if r.DueDate != "" {
				t, _ := time.Parse(time.RFC3339, r.DueDate)
				dueDate = &t
			}

			var completionDate *time.Time
			if r.CompletionDate != "" {
				t, _ := time.Parse(time.RFC3339, r.CompletionDate)
				completionDate = &t
			}

			isOverdue := false
			switch v := r.IsOverdue.(type) {
			case bool:
				isOverdue = v
			case float64:
				isOverdue = v == 1
			case string:
				isOverdue = v == "1" || v == "true"
			}

			allResults = append(allResults, lms.CourseResult{
				UserID:           r.UserID,
				CourseID:         r.CourseID,
				CourseTitle:      r.CourseTitle,
				EnrollmentID:     r.EnrollmentID,
				AccessDate:       accessDate,
				DueDate:          dueDate,
				CompletionDate:   completionDate,
				TimeSpent:        r.TimeSpent,
				Duration:         r.Duration,
				CompletionStatus: r.CompletionStatus,
				Progress:         r.Progress,
				AwardedScore:     r.AwardedScore,
				IsOverdue:        isOverdue,
				ViewsCount:       r.ViewsCount,
			})
		}

		if apiResp.NextPageToken == "" {
			break
		}
		pageToken = apiResp.NextPageToken
	}

	return allResults, nil
}