package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	name        = flag.String("name", "", "Name of repo to create in authenticated user's GitHub account.")
	description = flag.String("description", "", "Description of created repo.")
	private     = flag.Bool("private", false, "Will created repo be private.")
	autoInit    = flag.Bool("auto-init", false, "Pass true to create an initial commit with empty README.")
)

func creatRepo(ctx context.Context, client *github.Client) {
	r := &github.Repository{Name: name, Private: private, Description: description, AutoInit: autoInit}
	repo, _, err := client.Repositories.Create(ctx, "", r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created new repo: %v\n", repo.GetName())

}

func listRepo(ctx context.Context, client *github.Client) {
	repos, _, err := client.Repositories.List(ctx, "ashtatripathi", nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range repos {
		fmt.Printf("List repo: %s\n", *r.FullName)
	}
}

func listBranch(ctx context.Context, client *github.Client, repoName string) {
	repo, resp, err := client.Repositories.Get(ctx, "ashtatripathi", repoName)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	p := resp.Body

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(p)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	fmt.Printf("List branch: %s\n", repo.GetName())

	fmt.Println(repo.GetBranchesURL())

	fmt.Println(repo.GetGitRefsURL())
}

func getRef(ctx context.Context, client *github.Client) (ref *github.Reference, err error) {
	sourceOwner := "ashtatripathi"
	sourceRepo := "babble"
	commitBranch := "test-1"

	baseBranch := "master"

	if ref, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	// We consider that an error means the branch has not been found and needs to
	// be created.
	if commitBranch == baseBranch {
		return nil, errors.New("the commit branch does not exist but `-base-branch` is the same as `-commit-branch`")
	}

	if baseBranch == "" {
		return nil, errors.New("the `-base-branch` should not be set to an empty string when the branch specified by `-commit-branch` does not exists")
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+baseBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, sourceOwner, sourceRepo, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ctx context.Context, client *github.Client, ref *github.Reference, sourceFiles string) (tree *github.Tree, err error) {

	sourceOwner := "ashtatripathi"
	sourceRepo := "babble"

	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ctx context.Context, client *github.Client,
	ref *github.Reference, tree *github.Tree) (err error) {

	sourceOwner := "ashtatripathi"
	sourceRepo := "babble"
	authorName := "ashtatripathi"
	authorEmail := "ashta.tripathi@gmail.com"
	commitMessage := "Message test commit"

	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	return err
}

func createPR(ctx context.Context, client *github.Client) (err error) {

	prSubject := "test PR"
	prBranch := "master"

	sourceOwner := "ashtatripathi"
	prRepo := "babble"

	prDescription := "Description--"

	commitBranch := fmt.Sprintf("%s:%s", sourceOwner, prBranch)

	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &commitBranch,
		Base:                &prBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, sourceOwner, prRepo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}

func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

func main() {
	flag.Parse()
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if *name == "" {
		log.Fatal("No name: New repos must be given a name")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	creatRepo(ctx, client)

	listRepo(ctx, client)

	listBranch(ctx, client, "babble")

	ref, err := getRef(ctx, client)
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	sourceFiles := "git.go,go.sum"

	fmt.Println("done0")
	tree, err := getTree(ctx, client, ref, sourceFiles)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	fmt.Println("done1")

	if err := pushCommit(ctx, client, ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}
	fmt.Println("done2")
	if err := createPR(ctx, client); err != nil {
		log.Fatalf("Error while creating the pull request: %s", err)
	}
}
