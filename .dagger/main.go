package main

import (
	"context"
	"dagger/test/internal/dagger"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/google/go-github/v74/github"
)

type Test struct{}

func (m *Test) Handle(ctx context.Context, name string, event *dagger.File, token *dagger.Secret) error {
	if name != "issue_comment" {
		fmt.Println("not an issue comment event")

		return nil
	}

	var e github.IssueCommentEvent

	eventJSON, err := event.Contents(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(eventJSON), &e)
	if err != nil {
		return err
	}

	if e.GetAction() != "created" {
		fmt.Println("ignore event")

		return nil
	}

	authorizedUsers := []string{"sagikazarmark", "spacez320"}

	if !slices.Contains(authorizedUsers, e.GetComment().GetUser().GetLogin()) {
		return errors.New("unauthorized user: this incident will be reported")
	}

	fmt.Print("START", e.GetComment().GetBody(), "END")

	if e.GetComment().GetBody() != "/close" {
		return nil
	}

	tokenString, err := token.Plaintext(ctx)
	if err != nil {
		return err
	}

	client := github.NewClient(nil).WithAuthToken(tokenString)

	_, _, err = client.Issues.Edit(ctx, e.GetRepo().GetOwner().GetLogin(), e.GetRepo().GetName(), e.GetIssue().GetNumber(), &github.IssueRequest{
		State: github.Ptr("closed"),
	})
	if err != nil {
		return err
	}

	return nil
}
