package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v55/github"
)

type PR struct {
	number    int
	title     string
	createdBy string
	createdAt *time.Time
	openedAt  *time.Time
	mergedAt  *time.Time
}

func (pr *PR) durationBetweenCreateMerge() time.Duration {
	return pr.mergedAt.Sub(*pr.createdAt)
}

func (pr *PR) durationBwtweenOpenMerge() time.Duration {
	if pr.openedAt != nil {
		return pr.mergedAt.Sub(*pr.openedAt)
	} else {
		return pr.mergedAt.Sub(*pr.createdAt)
	}
}

type PRStat struct {
	pullRequest []PR
}

func (stat *PRStat) json() (string, error) {
	s := struct {
		Count                         int            `json:"count"`
		AverageTimeBetweenCreateMerge string         `json:"averageTimeBetweenCreateMerge"`
		AverageTimeBetweenOpenMerge   string         `json:"averageTimeBetweenOpenMerge"`
		PRCountPerUser                map[string]int `json:"prCountPerUser"`
	}{
		Count:                         stat.prCount(),
		AverageTimeBetweenCreateMerge: stat.calcAverageTimeBetweenCreateMerge().String(),
		AverageTimeBetweenOpenMerge:   stat.calcAverageTimeBetweenOpenMerge().String(),
		PRCountPerUser:                stat.getPRCountPerUser(),
	}
	j, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("Marchal Failed: %v", err)
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, j, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Indent failed: %v", err)
	}
	return buf.String(), nil
}

func (stat *PRStat) prCount() int {
	return len(stat.pullRequest)
}

func (stat *PRStat) calcAverageTimeBetweenCreateMerge() time.Duration {
	var total time.Duration
	for _, pr := range stat.pullRequest {
		total += pr.durationBetweenCreateMerge()
	}
	return time.Duration(total.Nanoseconds() / int64(stat.prCount()))
}

func (stat *PRStat) calcAverageTimeBetweenOpenMerge() time.Duration {
	var total time.Duration
	for _, pr := range stat.pullRequest {
		total += pr.durationBwtweenOpenMerge()
	}
	return time.Duration(total.Nanoseconds() / int64(stat.prCount()))
}

func (stat *PRStat) getPRCountPerUser() map[string]int {
	prCount := map[string]int{}
	for _, pr := range stat.pullRequest {
		var count = 0
		if c, ok := prCount[pr.createdBy]; ok {
			count = c
		}
		prCount[pr.createdBy] = (count + 1)
	}
	return prCount
}

func showStatAsJson(
	owner string,
	repo string,
	includeStart time.Time,
	excludeEnd time.Time,
	accessToken string,
) error {
	client := github.NewClient(nil).WithAuthToken(accessToken)
	pullRequests, err := getPullRequests(client, owner, repo, includeStart, excludeEnd)
	if err != nil {
		return err
	}
	prs := make([]PR, len(*pullRequests))
	for i, pr := range *pullRequests {
		readyForReviewedAt := findReadyForReviewDateTime(client, owner, repo, *pr.Number)
		prs[i] = PR{
			number:    *pr.Number,
			title:     *pr.Title,
			createdBy: *pr.User.Login,
			createdAt: &pr.CreatedAt.Time,
			openedAt:  readyForReviewedAt,
			mergedAt:  &pr.MergedAt.Time,
		}
	}

	stat := &PRStat{pullRequest: prs}
	j, err := stat.json()
	if err != nil {
		return err
	}
	fmt.Println(j)
	return nil
}

func getPullRequests(
	client *github.Client,
	owner string,
	repo string,
	includeStart time.Time,
	excludeEnd time.Time,
) (*[]github.PullRequest, error) {
	page := 1
	var pullRequests [](github.PullRequest)
	for {
		prs, _, err := client.PullRequests.List(
			context.Background(),
			owner,
			repo,
			&github.PullRequestListOptions{
				State:     "closed",
				Sort:      "created",
				Direction: "desc",
				ListOptions: github.ListOptions{
					Page:    page,
					PerPage: 100,
				},
			},
		)

		if err != nil {
			return nil, fmt.Errorf("Get pullrequest failed: %v", err)
		}

		for _, pr := range prs {
			if pr.MergedAt == nil {
				continue
			}
			startOk := pr.MergedAt.After(includeStart) || pr.MergedAt.Time.Equal(includeStart)
			endOk := pr.MergedAt.Before(excludeEnd)
			if startOk && endOk {
				pullRequests = append(pullRequests, *pr)
			}
		}

		if len(prs) == 0 {
			return &pullRequests, nil
		}
		lastCreatedAt := prs[len(prs)-1].CreatedAt.Time
		if lastCreatedAt.Before(includeStart) {
			return &pullRequests, nil
		}
		page++
	}
}

func findReadyForReviewDateTime(
	client *github.Client,
	owner string,
	repo string,
	prNumber int,
) *time.Time {
	timeline, _, err := client.Issues.ListIssueTimeline(
		context.Background(),
		owner,
		repo,
		prNumber,
		&github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	)
	if err != nil {
		log.Printf("findReadyForReviewDateTime failed: %v", err)
		return nil
	}
	for _, event := range timeline {
		if *event.Event == "ready_for_review" {
			return event.CreatedAt.GetTime()
		}
	}
	return nil
}
