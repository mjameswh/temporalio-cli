package temporalcli_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	workflowpb "go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

func (s *SharedServerSuite) TestWorkflow_Signal_SingleWorkflowSuccess() {
	// Make workflow wait for signal and then return it
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		var ret any
		workflow.GetSignalChannel(ctx, "my-signal").Receive(ctx, &ret)
		return ret, nil
	})

	// Start the workflow
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{TaskQueue: s.Worker.Options.TaskQueue},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)

	// Send signal
	res := s.Execute(
		"workflow", "signal",
		"--address", s.Address(),
		"-w", run.GetID(),
		"--name", "my-signal",
		"-i", `{"foo": "bar"}`,
	)
	s.NoError(res.Err)

	// Confirm workflow result was as expected
	var actual any
	s.NoError(run.Get(s.Context, &actual))
	s.Equal(map[string]any{"foo": "bar"}, actual)
}

func (s *SharedServerSuite) TestWorkflow_Signal_BatchWorkflowSuccess() {
	res := s.testSignalBatchWorkflow(false)
	s.Contains(res.Stdout.String(), "approximately 5 workflow(s)")
	s.Contains(res.Stdout.String(), "Started batch")
}

func (s *SharedServerSuite) TestWorkflow_Signal_BatchWorkflowSuccessJSON() {
	res := s.testSignalBatchWorkflow(true)
	var jsonRes map[string]any
	s.NoError(json.Unmarshal(res.Stdout.Bytes(), &jsonRes))
	s.NotEmpty(jsonRes["batchJobId"])
}

func (s *SharedServerSuite) testSignalBatchWorkflow(json bool) *CommandResult {
	// Make workflow wait for signal and then return it
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		var ret any
		workflow.GetSignalChannel(ctx, "my-signal").Receive(ctx, &ret)
		return ret, nil
	})

	// Start 5 workflows
	runs := make([]client.WorkflowRun, 5)
	searchAttr := "keyword-" + uuid.NewString()
	for i := range runs {
		run, err := s.Client.ExecuteWorkflow(
			s.Context,
			client.StartWorkflowOptions{
				TaskQueue:        s.Worker.Options.TaskQueue,
				SearchAttributes: map[string]any{"CustomKeywordField": searchAttr},
			},
			DevWorkflow,
			"ignored",
		)
		s.NoError(err)
		runs[i] = run
	}

	// Wait for all to appear in list
	s.Eventually(func() bool {
		resp, err := s.Client.ListWorkflow(s.Context, &workflowservice.ListWorkflowExecutionsRequest{
			Query: "CustomKeywordField = '" + searchAttr + "'",
		})
		s.NoError(err)
		return len(resp.Executions) == len(runs)
	}, 3*time.Second, 100*time.Millisecond)

	// Send batch signal with a "y" for non-json or "--yes" for json
	args := []string{
		"workflow", "signal",
		"--address", s.Address(),
		"--query", "CustomKeywordField = '" + searchAttr + "'",
		"--name", "my-signal",
		"-i", `{"key": "val"}`,
	}
	if json {
		args = append(args, "--yes", "-o", "json")
	} else {
		s.CommandHarness.Stdin.WriteString("y\n")
	}
	res := s.Execute(args...)
	s.NoError(res.Err)

	// Confirm that all workflows complete with the signal value
	for _, run := range runs {
		var ret map[string]string
		s.NoError(run.Get(s.Context, &ret))
		s.Equal(map[string]string{"key": "val"}, ret)
	}
	return res
}

func (s *SharedServerSuite) TestWorkflow_Terminate_SingleWorkflowSuccess_WithoutReason() {
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		ctx.Done().Receive(ctx, nil)
		return nil, ctx.Err()
	})

	// Start the workflow
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{TaskQueue: s.Worker.Options.TaskQueue},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)

	// Send terminate
	res := s.Execute(
		"workflow", "terminate",
		"--address", s.Address(),
		"-w", run.GetID(),
	)
	s.NoError(res.Err)

	// Confirm workflow was terminated
	s.Contains(run.Get(s.Context, nil).Error(), "terminated")
	// Ensure the termination reason was recorded
	iter := s.Client.GetWorkflowHistory(s.Context, run.GetID(), run.GetRunID(), false, enums.HISTORY_EVENT_FILTER_TYPE_CLOSE_EVENT)
	var foundReason bool
	for iter.HasNext() {
		event, err := iter.Next()
		s.NoError(err)
		if term := event.GetWorkflowExecutionTerminatedEventAttributes(); term != nil {
			foundReason = true
			// We're not going to check the value here so we don't pin ourselves to our particular default, but there _should_ be a default reason
			s.NotEmpty(term.Reason)
		}
	}
	s.True(foundReason)
}

func (s *SharedServerSuite) TestWorkflow_Terminate_SingleWorkflowSuccess_WithReason() {
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		ctx.Done().Receive(ctx, nil)
		return nil, ctx.Err()
	})

	// Start the workflow
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{TaskQueue: s.Worker.Options.TaskQueue},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)

	// Send terminate
	res := s.Execute(
		"workflow", "terminate",
		"--address", s.Address(),
		"-w", run.GetID(),
		"--reason", "terminate-test",
	)
	s.NoError(res.Err)

	// Confirm workflow was terminated
	s.Contains(run.Get(s.Context, nil).Error(), "terminated")

	// Ensure the termination reason was recorded
	iter := s.Client.GetWorkflowHistory(s.Context, run.GetID(), run.GetRunID(), false, enums.HISTORY_EVENT_FILTER_TYPE_CLOSE_EVENT)
	var foundReason bool
	for iter.HasNext() {
		event, err := iter.Next()
		s.NoError(err)
		if term := event.GetWorkflowExecutionTerminatedEventAttributes(); term != nil {
			foundReason = true
			s.Equal("terminate-test", term.Reason)
		}
	}
	s.True(foundReason)
}

func (s *SharedServerSuite) TestWorkflow_Terminate_BatchWorkflowSuccess() {
	res := s.testTerminateBatchWorkflow(false)
	s.Contains(res.Stdout.String(), "approximately 5 workflow(s)")
	s.Contains(res.Stdout.String(), "Started batch")
}

func (s *SharedServerSuite) TestWorkflow_Terminate_BatchWorkflowSuccessJSON() {
	res := s.testTerminateBatchWorkflow(true)
	var jsonRes map[string]any
	s.NoError(json.Unmarshal(res.Stdout.Bytes(), &jsonRes))
	s.NotEmpty(jsonRes["batchJobId"])
}

func (s *SharedServerSuite) testTerminateBatchWorkflow(json bool) *CommandResult {
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		ctx.Done().Receive(ctx, nil)
		return nil, ctx.Err()
	})

	// Start 5 workflows
	runs := make([]client.WorkflowRun, 5)
	searchAttr := "keyword-" + uuid.NewString()
	for i := range runs {
		run, err := s.Client.ExecuteWorkflow(
			s.Context,
			client.StartWorkflowOptions{
				TaskQueue:        s.Worker.Options.TaskQueue,
				SearchAttributes: map[string]any{"CustomKeywordField": searchAttr},
			},
			DevWorkflow,
			"ignored",
		)
		s.NoError(err)
		runs[i] = run
	}

	// Wait for all to appear in list
	s.Eventually(func() bool {
		resp, err := s.Client.ListWorkflow(s.Context, &workflowservice.ListWorkflowExecutionsRequest{
			Query: "CustomKeywordField = '" + searchAttr + "'",
		})
		s.NoError(err)
		return len(resp.Executions) == len(runs)
	}, 3*time.Second, 100*time.Millisecond)

	// Send batch terminate with a "y" for non-json or "--yes" for json
	args := []string{
		"workflow", "terminate",
		"--address", s.Address(),
		"--query", "CustomKeywordField = '" + searchAttr + "'",
		"--reason", "terminate-test",
	}
	if json {
		args = append(args, "--yes", "-o", "json")
	} else {
		s.CommandHarness.Stdin.WriteString("y\n")
	}
	res := s.Execute(args...)
	s.NoError(res.Err)

	// Confirm that all workflows are terminated
	for _, run := range runs {
		s.Contains(run.Get(s.Context, nil).Error(), "terminated")
		// Ensure the termination reason was recorded
		iter := s.Client.GetWorkflowHistory(s.Context, run.GetID(), run.GetRunID(), false, enums.HISTORY_EVENT_FILTER_TYPE_CLOSE_EVENT)
		var foundReason bool
		for iter.HasNext() {
			event, err := iter.Next()
			s.NoError(err)
			if term := event.GetWorkflowExecutionTerminatedEventAttributes(); term != nil {
				foundReason = true
				s.Equal("terminate-test", term.Reason)
			}
		}
		s.True(foundReason)
	}
	return res
}

func (s *SharedServerSuite) awaitResetWorkflow(searchAttr string) {
	var lastExecs []*workflowpb.WorkflowExecutionInfo
	s.Eventually(func() bool {
		resp, err := s.Client.ListWorkflow(s.Context, &workflowservice.ListWorkflowExecutionsRequest{
			Query: "CustomKeywordField = '" + searchAttr + "'",
		})
		s.NoError(err)
		lastExecs = resp.Executions
		return len(resp.Executions) == 2 && resp.Executions[0].Status == enums.WORKFLOW_EXECUTION_STATUS_COMPLETED
	}, 3*time.Second, 100*time.Millisecond, "Reset execution failed to complete", lastExecs)
}

func (s *SharedServerSuite) TestWorkflow_Reset_ToFirstWorkflowTask() {
	var wfExecutions, activityExecutions int
	s.Worker.OnDevActivity(func(ctx context.Context, a any) (any, error) {
		activityExecutions++
		return nil, nil
	})
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		workflow.ExecuteActivity(ctx, DevActivity, 1).Get(ctx, nil)
		wfExecutions++
		return nil, nil
	})

	// Start the workflow
	searchAttr := "keyword-" + uuid.NewString()
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{
			TaskQueue:        s.Worker.Options.TaskQueue,
			SearchAttributes: map[string]any{"CustomKeywordField": searchAttr},
		},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)
	var junk any
	s.NoError(run.Get(s.Context, &junk))
	s.Equal(1, wfExecutions)

	// Reset to the first workflow task
	res := s.Execute(
		"workflow", "reset",
		"--address", s.Address(),
		"-w", run.GetID(),
		"-t", "FirstWorkflowTask",
		"--reason", "test-reset-FirstWorkflowTask",
	)
	require.NoError(s.T(), res.Err)
	s.awaitResetWorkflow(searchAttr)
	s.Equal(2, wfExecutions, "Should have re-executed the workflow from the beginning")
	s.Greater(activityExecutions, 1, "Should have re-executed the workflow from the beginning")
}

func (s *SharedServerSuite) TestWorkflow_Reset_ToLastWorkflowTask() {
	var wfExecutions, activityExecutions int
	s.Worker.OnDevActivity(func(ctx context.Context, a any) (any, error) {
		activityExecutions++
		return nil, nil
	})
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		workflow.ExecuteActivity(ctx, DevActivity, 1).Get(ctx, nil)
		wfExecutions++
		return nil, nil
	})

	// Start the workflow
	searchAttr := "keyword-" + uuid.NewString()
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{
			TaskQueue:        s.Worker.Options.TaskQueue,
			SearchAttributes: map[string]any{"CustomKeywordField": searchAttr},
		},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)
	var junk any
	s.NoError(run.Get(s.Context, &junk))
	s.Equal(1, wfExecutions)

	// Reset to the final workflow task
	res := s.Execute(
		"workflow", "reset",
		"--address", s.Address(),
		"-w", run.GetID(),
		"-t", "LastWorkflowTask",
		"--reason", "test-reset-LastWorkflowTask",
	)
	require.NoError(s.T(), res.Err)
	s.awaitResetWorkflow(searchAttr)
	s.Equal(2, wfExecutions, "Should re-executed the workflow")
	s.Equal(1, activityExecutions, "Should not have re-executed the activity")
}

func (s *SharedServerSuite) TestWorkflow_Reset_ToLastContinuedAsNew() {
	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		ctx.Done().Receive(ctx, nil)
		return nil, ctx.Err()
	})
	s.Fail("Not implemented")
}

func (s *SharedServerSuite) TestWorkflow_Reset_ToEventID() {
	// We execute two activities and will resume just before the second one. We use the same activity for both
	// but a unique input so we can check which fake activity is executed
	var oneExecutions, twoExecutions int
	s.Worker.OnDevActivity(func(ctx context.Context, a any) (any, error) {
		n, ok := a.(float64)
		if !ok {
			return nil, fmt.Errorf("expected int, not %T (%v)", a, a)
		}
		switch n {
		case 1:
			oneExecutions++
		case 2:
			twoExecutions++
		default:
			return 0, errors.New("you've broken the test!")
		}
		return n, nil
	})

	s.Worker.OnDevWorkflow(func(ctx workflow.Context, a any) (any, error) {
		var res any
		if err := workflow.ExecuteActivity(ctx, DevActivity, 1).Get(ctx, &res); err != nil {
			return res, err
		}
		err := workflow.ExecuteActivity(ctx, DevActivity, 2).Get(ctx, &res)
		return res, err
	})

	// Start the workflow
	searchAttr := "keyword-" + uuid.NewString()
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{
			TaskQueue:        s.Worker.Options.TaskQueue,
			SearchAttributes: map[string]any{"CustomKeywordField": searchAttr},
		},
		DevWorkflow,
		"ignored",
	)
	require.NoError(s.T(), err)
	var ignored any
	s.NoError(run.Get(s.Context, &ignored))
	s.Equal(1, oneExecutions)
	s.Equal(1, twoExecutions)

	// We want to reset to the last WFTCompleted event before the second activity, so we
	// need to search the history for it. I could just pick the event ID, but I don't want
	// this test to break if new event types are added in the future
	wfsvc := workflowservice.NewWorkflowServiceClient(s.GRPCConn)
	ctx := context.Background()
	req := workflowservice.GetWorkflowExecutionHistoryReverseRequest{
		Namespace: s.Namespace(),
		Execution: &commonpb.WorkflowExecution{
			WorkflowId: run.GetID(),
			RunId:      run.GetRunID(),
		},
		MaximumPageSize: 250,
		NextPageToken:   nil,
	}
	beforeSecondActivity := int64(-1)
	var takeNextWorkflowTaskCompleted bool
	for beforeSecondActivity == -1 {
		resp, err := wfsvc.GetWorkflowExecutionHistoryReverse(ctx, &req)
		s.NoError(err)
		for _, e := range resp.GetHistory().GetEvents() {
			s.T().Logf("Event: %d %s", e.GetEventId(), e.GetEventType())
			if e.GetEventType() == enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED && beforeSecondActivity == -1 {
				takeNextWorkflowTaskCompleted = true
			} else if e.GetEventType() == enums.EVENT_TYPE_WORKFLOW_TASK_COMPLETED && takeNextWorkflowTaskCompleted {
				beforeSecondActivity = e.EventId
				break
			}
		}
		if len(resp.NextPageToken) != 0 {
			req.NextPageToken = resp.NextPageToken
		} else {
			break
		}
	}

	// Reset to the before the second activity execution
	res := s.Execute(
		"workflow", "reset",
		"--address", s.Address(),
		"-w", run.GetID(),
		"-e", fmt.Sprintf("%d", beforeSecondActivity),
		"--reason", "test-reset-event-id",
	)
	require.NoError(s.T(), res.Err)

	s.awaitResetWorkflow(searchAttr)
	s.Equal(1, oneExecutions, "Should not have re-executed the first activity")
	s.Equal(2, twoExecutions, "Should have re-executed the second activity")
}
