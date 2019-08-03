package gitcount

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type GitCounter struct {
	client    *gitlab.Client
	NameSpace string
}

func New(username string, password string, NameSpace string) *GitCounter {
	client, _ := gitlab.NewBasicAuthClient(nil, "https://gitlab.com", username, password)
	return &GitCounter{
		client:    client,
		NameSpace: NameSpace,
	}
}

func (counter *GitCounter) currentUser() (*gitlab.User, error) {
	usr, _, err := counter.client.Users.CurrentUser()
	if err != nil {
		fmt.Println("Username or password invalid")
		return nil, errors.New("Invalid username or password")
	}
	return usr, nil
}

func (counter *GitCounter) GetProject() (*gitlab.Project, error) {
	buckets := strings.Split(counter.NameSpace, "/")
	groupname := buckets[0]
	groups, _, gerr := counter.client.Groups.ListGroups(&gitlab.ListGroupsOptions{Search: &groupname})
	if gerr != nil {
		return nil, errors.New("Failed fetching Group Information for " + groupname)
	}
	var project *gitlab.Project
	if len(groups) > 0 {
		group := groups[0]
		projects, _, perr := counter.client.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{
			Search: &buckets[1],
		})
		if perr != nil && len(projects) <= 0 {
			return nil, errors.New("Failed loading reports from " + groupname)
		}
		project = projects[0]
	}
	return project, nil
}

func (counter *GitCounter) GetAllUsers(project *gitlab.Project) ([]*gitlab.ProjectUser, error) {
	users, _, urr := counter.client.Projects.ListProjectsUsers(project.ID, &gitlab.ListProjectUserOptions{})
	if urr != nil {
		return nil, errors.New("Unable to get users")
	}
	return users, nil
}

func (counter *GitCounter) PrintStats(project *gitlab.Project, label string) {
	label_data := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 1000},
		Labels:      []string{label},
	}
	issues, _, _ := counter.client.Issues.ListProjectIssues(project.ID, label_data)
	issue_maps := make(map[string]int)

	for _, issue := range issues {
		_, ok := issue_maps[issue.Author.Name]
		if ok {
			issue_maps[issue.Assignee.Name] += issue.TimeStats.TotalTimeSpent / 3600
		} else {
			issue_maps[issue.Assignee.Name] = issue.TimeStats.TotalTimeSpent / 3600
		}
	}

	for key, value := range issue_maps {
		pstr := fmt.Sprintf("%s has spent about %d hours", key, value)
		fmt.Println(pstr)
	}
}
