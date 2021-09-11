// Unless explicitly stated otherwise all files in this repository are licensed under the MIT License.
//
// This product includes software developed at Datadog (https://www.datadoghq.com/). Copyright 2021 Datadog, Inc.

package temporaltest_test

import (
	"context"
	"testing"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/DataDog/temporalite/internal/examples/helloworld"
	"github.com/DataDog/temporalite/server/temporaltest"
)

func TestNewServer(t *testing.T) {
	// Create test Temporal server and client
	ts := temporaltest.NewServer()
	c := ts.Client()
	defer ts.Stop()

	w := worker.New(c, "example", worker.Options{})
	helloworld.RegisterWorkflowsAndActivities(w)

	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wfr, err := c.ExecuteWorkflow(ctx, client.StartWorkflowOptions{TaskQueue: "example"}, helloworld.Greet, "world")
	if err != nil {
		t.Fatal(err)
	}

	var resp string
	if err := wfr.Get(ctx, &resp); err != nil {
		t.Fatal(err)
	}

	if resp != "Hello world" {
		t.Fatalf("unexpected result: %q", resp)
	}
}
