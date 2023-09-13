package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/osbuild/pulp-client/pulpclient"
)

func check(err error, resp *http.Response) {
	if err != nil {
		b, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "%s: %s\n", err.Error(), string(b))
		os.Exit(1)
	}
}

func main() {
	ctx := context.WithValue(context.Background(), pulpclient.ContextServerIndex, 0)
	transport := &http.Transport{}
	httpClient := http.Client{Transport: transport}

	url := "http://localhost:8080"
	pulpConfig := pulpclient.NewConfiguration()
	pulpConfig.HTTPClient = &httpClient
	pulpConfig.Servers = pulpclient.ServerConfigurations{pulpclient.ServerConfiguration{
		URL: url,
	}}
	client := pulpclient.NewAPIClient(pulpConfig)

	ctx = context.WithValue(ctx, pulpclient.ContextBasicAuth, pulpclient.BasicAuth{
		UserName: "admin",
		Password: "password42",
	})

	fmt.Println("Listing repos")
	repoList, resp, err := client.RepositoriesOstreeAPI.RepositoriesOstreeOstreeList(ctx).Execute()
	check(err, resp)

	fmt.Printf("Deleting %d repositories\n", repoList.GetCount())
	for _, repo := range repoList.GetResults() {
		fmt.Printf("Deleting %s (%s)\n", repo.Name, repo.GetPulpHref())

		oper, resp, err := client.RepositoriesOstreeAPI.RepositoriesOstreeOstreeDelete(ctx, repo.GetPulpHref()).Execute()
		check(err, resp)
		fmt.Printf("Task: %s\n", oper.GetTask())
	}

	fmt.Println("Listing distributions")
	distList, resp, err := client.DistributionsOstreeAPI.DistributionsOstreeOstreeList(ctx).Execute()
	check(err, resp)

	fmt.Printf("Deleting %d distributions\n", distList.GetCount())
	for _, dist := range distList.GetResults() {
		fmt.Printf("Deleting %s (%s)\n", dist.Name, dist.GetPulpHref())

		oper, resp, err := client.DistributionsOstreeAPI.DistributionsOstreeOstreeDelete(ctx, dist.GetPulpHref()).Execute()
		check(err, resp)
		fmt.Printf("Task: %s\n", oper.GetTask())
	}

	fmt.Printf("Deleting orphaned content\n")
	orphansCleanup := *pulpclient.NewOrphansCleanup()
	orphansCleanup.SetOrphanProtectionTime(0)
	task, resp, err := client.OrphansCleanupAPI.OrphansCleanupCleanup(ctx).OrphansCleanup(orphansCleanup).Execute()
	check(err, resp)
	fmt.Printf("Delete task: %s%s\n", url, task.GetTask())

	state := "waiting"
	for state == "running" || state == "waiting" {
		res, resp, err := client.TasksAPI.TasksRead(ctx, task.GetTask()).Execute()
		check(err, resp)
		state = res.GetState()
		fmt.Printf("Task is %s\n", state)
		time.Sleep(2 * time.Second)
	}

}
