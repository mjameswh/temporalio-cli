package temporalcli_test

import (
	"go.temporal.io/sdk/client"
)

func (s *SharedServerSuite) TestWorkflow_Trace_WIP() {
	// Start the workflow and wait until it has at least reached activity failure
	run, err := s.Client.ExecuteWorkflow(
		s.Context,
		client.StartWorkflowOptions{TaskQueue: s.Worker.Options.TaskQueue},
		DevWorkflow,
		"ignored",
	)
	s.NoError(err)
	// s.Eventually(func() bool {
	// 	resp, err := s.Client.DescribeWorkflowExecution(s.Context, run.GetID(), run.GetRunID())
	// 	s.NoError(err)
	// 	return len(resp.PendingActivities) > 0 && resp.PendingActivities[0].LastFailure != nil
	// }, 5*time.Second, 100*time.Millisecond)

	res := s.Execute(
		"workflow", "trace",
		"--address", s.Address(),
		"-w", run.GetID(),
	)
	s.NoError(res.Err)
	// out := res.Stdout.String()
	// s.ContainsOnSameLine(out, "WorkflowId", run.GetID())
	// s.Contains(out, "Pending Activities: 1")
	// s.ContainsOnSameLine(out, "LastFailure", "intentional error")
	// s.Contains(out, "Pending Child Workflows: 0")

	// // JSON
	// res = s.Execute(
	// 	"workflow", "describe",
	// 	"-o", "json",
	// 	"--address", s.Address(),
	// 	"-w", run.GetID(),
	// )
	// s.NoError(res.Err)
	// var jsonOut workflowservice.DescribeWorkflowExecutionResponse
	// s.NoError(temporalcli.UnmarshalProtoJSONWithOptions(res.Stdout.Bytes(), &jsonOut, true))
	// s.Equal("intentional error", jsonOut.PendingActivities[0].LastFailure.Message)
}
