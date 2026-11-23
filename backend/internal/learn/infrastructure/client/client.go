package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"backend/internal/learn/domain"
)

type Client struct {
	baseURL         string
	clientID        string
	clientSecret    string
	accessToken     string
	tokenExpiry     time.Time
	tokenMu         sync.Mutex
	httpClient      *http.Client
}

func NewClient(baseUrl, clientID, clientSecret string) *Client {
	return &Client{
		baseURL:      baseUrl,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) getToken() (string, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/api/v3/token", c.baseURL), data)
	if err != nil {
		return "", fmt.Errorf("запрос токена: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("чтение ответа токена: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка получения токена (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("разбор ответа токена: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	expiresIn := time.Duration(tokenResp.ExpiresIn) * time.Second
	if expiresIn == 0 {
		expiresIn = 30 * time.Minute
	}
	c.tokenExpiry = time.Now().Add(expiresIn - 5*time.Minute)

	log.Printf("Learn API: токен получен, истекает через %d сек", tokenResp.ExpiresIn)
	return c.accessToken, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		c.tokenMu.Lock()
		c.accessToken = ""
		c.tokenMu.Unlock()

		token, err = c.getToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", token)
		return c.httpClient.Do(req)
	}

	return resp, err
}

func (c *Client) FetchAllLearnAccounts(ctx context.Context) ([]domain.MotherLearnAccount, error) {
	var allLearnAccounts []domain.MotherLearnAccount
	pageToken := ""

	for {
		params := url.Values{}
		params.Set("pageSize", "1000")
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		apiURL := fmt.Sprintf("%s/users?%s", c.baseURL, params.Encode())

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiResp struct {
			UserProfiles []struct {
				UserID        string `json:"userId"`
				DepartmentID  string `json:"departmentId"`
				Status        int    `json:"status"`
				AddedDate     string `json:"addedDate"`
				LastLoginDate string `json:"lastLoginDate"`
				Fields        []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"fields"`
			} `json:"userProfiles"`
			NextPageToken string `json:"nextPageToken"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("разбор списка пользователей: %w (ответ: %s)", err, string(body)[:min(len(body), 500)])
		}

		for _, p := range apiResp.UserProfiles {
			getField := func(name string) string {
				for _, f := range p.Fields {
					if f.Name == name {
						return f.Value
					}
				}
				return ""
			}

			addedDate, _ := time.Parse("2006-01-02", p.AddedDate)
			lastLoginDate, _ := time.Parse("2006-01-02", p.LastLoginDate)

			allLearnAccounts = append(allLearnAccounts, domain.MotherLearnAccount{
				UserId:        p.UserID,
				DepartmentId:  p.DepartmentID,
				FirstName:     getField("FIRST_NAME"),
				LastName:      getField("LAST_NAME"),
				Email:         getField("EMAIL"),
				Login:         getField("LOGIN"),
				Phone:         getField("PHONE"),
				Status:        fmt.Sprintf("%d", p.Status),
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

func (c *Client) FetchAllCourseResults(ctx context.Context) ([]domain.CourseResult, error) {
	var allResults []domain.CourseResult
	pageToken := ""

	for {
		params := url.Values{}
		params.Set("pageSize", "1000")
		if pageToken != "" {
			params.Set("pageToken", pageToken)
		}

		params.Add("userStatuses[]", "1")

		apiURL := fmt.Sprintf("%s/learners/results?%s", c.baseURL, params.Encode())

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.doRequest(req)
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

			allResults = append(allResults, domain.CourseResult{
				UserId:           r.UserID,
				CourseId:         r.CourseID,
				NameCourse:       r.CourseTitle,
				AccesDate:        accessDate,
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

func (c *Client) FetchAllIPR(ctx context.Context) ([]domain.IPRPlan, error) {
	var allPlans []domain.IPRPlan
	pageToken := ""

	for {
		bodyBytes := []byte(`{"pageSize": 100}`)
		if pageToken != "" {
			bodyBytes = []byte(fmt.Sprintf(`{"pageSize": 100, "pageToken": "%s"}`, pageToken))
		}

		apiURL := fmt.Sprintf("%s/api/v3/developmentplan/v1/enrolled-plans", c.baseURL)

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		var apiResp struct {
			Items         []struct{ PlanID string `json:"id"` } `json:"items"`
			NextPageToken string                                `json:"nextPageToken"`
		}

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("разбор списка ИПР: %w (ответ: %s)", err, string(body)[:min(len(body), 500)])
		}

		var planIDs []string
		for _, item := range apiResp.Items {
			planIDs = append(planIDs, item.PlanID)
		}

		if apiResp.NextPageToken == "" && len(planIDs) > 0 {
			plans, err := c.fetchIPRDetails(ctx, planIDs)
			if err != nil {
				return nil, err
			}
			allPlans = append(allPlans, plans...)
			break
		}

		if len(planIDs) > 0 {
			plans, err := c.fetchIPRDetails(ctx, planIDs)
			if err != nil {
				return nil, err
			}
			allPlans = append(allPlans, plans...)
		}

		if apiResp.NextPageToken == "" {
			break
		}
		pageToken = apiResp.NextPageToken
	}

	return allPlans, nil
}

func (c *Client) fetchIPRDetails(ctx context.Context, planIDs []string) ([]domain.IPRPlan, error) {
	if len(planIDs) == 0 {
		return nil, nil
	}

	idsJSON, _ := json.Marshal(planIDs)
	bodyStr := fmt.Sprintf(`{"ids": %s}`, string(idsJSON))

	apiURL := fmt.Sprintf("%s/api/v3/developmentplan/v1/enrolled-plans/detailed", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader([]byte(bodyStr)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

 var apiResp struct {
		Items []struct {
			PlanID       string   `json:"id"`
			Title        string   `json:"title"`
			EmployeeID   string   `json:"employeeId"`
			MentorID     string   `json:"mentorId"`
			MentorName   string   `json:"mentorName"`
			StartDate    string   `json:"startDate"`
			EndDate      *string  `json:"endDate"`
			Status       string   `json:"status"`
			Progress     struct {
				CompletedTasksCount int `json:"completedTasksCount"`
				TotalTasksCount     int `json:"totalTasksCount"`
			} `json:"progress"`
			Goals        []string `json:"goals"`
			Stages       []struct {
				StageID     string   `json:"id"`
				Title       string   `json:"title"`
				Description string   `json:"description"`
				StartDate   *string  `json:"startDate"`
				EndDate     *string  `json:"endDate"`
				Tasks       []struct {
					TaskData struct {
						TaskID string `json:"id"`
						Type   string `json:"type"`
						TextTask struct {
							Title       string `json:"title"`
							Description string `json:"description"`
						} `json:"textTask"`
					} `json:"taskData"`
					Status          string  `json:"status"`
					IsOverdue       bool    `json:"isOverdue"`
					AbsoluteDueDate *string `json:"absoluteDueDate"`
					CompletionDate  *string `json:"completionDate"`
				} `json:"tasks"`
			} `json:"stages"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("разбор деталей ИПР: %w (ответ: %s)", err, string(body)[:min(len(body), 500)])
	}

	var plans []domain.IPRPlan
	for _, item := range apiResp.Items {
		startDate, _ := time.Parse("2006-01-02", item.StartDate)

		var endDate *time.Time
		if item.EndDate != nil && *item.EndDate != "" {
			t, _ := time.Parse("2006-01-02", *item.EndDate)
			endDate = &t
		}

		var stages []domain.IPRStage
		for _, s := range item.Stages {
			var sStart, sEnd *time.Time
			if s.StartDate != nil && *s.StartDate != "" {
				t, _ := time.Parse("2006-01-02", *s.StartDate)
				sStart = &t
			}
			if s.EndDate != nil && *s.EndDate != "" {
				t, _ := time.Parse("2006-01-02", *s.EndDate)
				sEnd = &t
			}

			var tasks []domain.IPRTask
			for _, t := range s.Tasks {
				var dueDate, completionDate *time.Time
				if t.AbsoluteDueDate != nil && *t.AbsoluteDueDate != "" {
					d, _ := time.Parse("2006-01-02", *t.AbsoluteDueDate)
					dueDate = &d
				}
				if t.CompletionDate != nil && *t.CompletionDate != "" {
					d, _ := time.Parse("2006-01-02", *t.CompletionDate)
					completionDate = &d
				}

				title := t.TaskData.TextTask.Title
				if title == "" {
					title = t.TaskData.Type
				}

				tasks = append(tasks, domain.IPRTask{
					TaskID:          t.TaskData.TaskID,
					TaskType:        t.TaskData.Type,
					Title:           title,
					Description:     t.TaskData.TextTask.Description,
					Status:          t.Status,
					IsOverdue:       t.IsOverdue,
					AbsoluteDueDate: dueDate,
					CompletionDate:  completionDate,
				})
			}

			stages = append(stages, domain.IPRStage{
				StageID:     s.StageID,
				Title:       s.Title,
				Description: s.Description,
				StartDate:   sStart,
				EndDate:     sEnd,
				Tasks:       tasks,
			})
		}

		var progress float64
		if item.Progress.TotalTasksCount > 0 {
			progress = float64(item.Progress.CompletedTasksCount) / float64(item.Progress.TotalTasksCount) * 100
		}
		plans = append(plans, domain.IPRPlan{
			PlanID:         item.PlanID,
			Title:          item.Title,
			EmployeeID:     item.EmployeeID,
			MentorID:       item.MentorID,
			MentorName:     item.MentorName,
			StartDate:      startDate,
			EndDate:        endDate,
			Status:         item.Status,
			Progress:       progress,
			CompletedTasks: item.Progress.CompletedTasksCount,
			TotalTasks:     item.Progress.TotalTasksCount,
			Goals:          item.Goals,
			Stages:         stages,
		})
	}

	return plans, nil
}
