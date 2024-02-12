package temporalcli

import (
	"fmt"
)

func (c *TemporalWorkflowTraceCommand) run(cctx *CommandContext, args []string) error {
	// *** Call describe
	cl, err := c.Parent.ClientOptions.dialClient(cctx)
	if err != nil {
		return err
	}
	defer cl.Close()
	// // ***
	// resp, err := cl.DescribeWorkflowExecution(cctx, c.WorkflowId, c.RunId)
	// if err != nil {
	// 	return fmt.Errorf("failed describing workflow: %w", err)
	// }
	return fmt.Errorf("TODO")
	// return nil
}
