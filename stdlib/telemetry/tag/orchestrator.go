package tag

import tags "go.opencensus.io/tag"

var (
	TagOrchEngine, _              = tags.NewKey(`go.bpm.engine`)
	TagOrchBPMMethod, _           = tags.NewKey(`go.bpm.method`)
	TagOrchWorkflowKey, _         = tags.NewKey(`go.bpm.workflowkey`)
	TagOrchWorkflowInstanceKey, _ = tags.NewKey(`go.bpm.workflowinstancekey`)
	TagOrchWorkflowVersion, _     = tags.NewKey(`go.bpm.workflowversion`)
	TagOrchBPMNProcessID, _       = tags.NewKey(`go.bpm.bpmnprocessid`)
	TagOrchElementID, _           = tags.NewKey(`go.bpm.elementid`)
	TagOrchElementInstanceKey, _  = tags.NewKey(`go.bpm.elementinstancekey`)
	TagOrchWorkerName, _          = tags.NewKey(`go.bpm.workername`)
	TagOrchJobRetries, _          = tags.NewKey(`go.bpm.jobretries`)
	TagOrchJobType, _             = tags.NewKey(`go.bpm.jobtype`)
	TagOrchJobKey, _              = tags.NewKey(`go.bpm.jobkey`)
	TagOrchJobTimeout, _          = tags.NewKey(`go.bpm.jobtimeout`)
	TagOrchStatusCode, _          = tags.NewKey(`go.bpm.status.code`)
	TagOrchMessageID, _           = tags.NewKey(`go.bpm.messageid`)
	TagOrchMessageName, _         = tags.NewKey(`go.bpm.messagename`)
	TagOrchMessageCorrKey, _      = tags.NewKey(`go.bpm.messagecorrkey`)
	TagOrchIncidentKey, _         = tags.NewKey(`go.bpm.incidentkey`)
)
