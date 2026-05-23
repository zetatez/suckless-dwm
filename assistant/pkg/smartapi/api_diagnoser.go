package smartapi

import (
	"context"
	"fmt"

	"assistant/pkg/llm"
	"assistant/pkg/smartapi/prompts"
)

type Diagnoser struct {
	engine *Engine
}

func NewDiagnoser(client llm.Client) *Diagnoser {
	return &Diagnoser{engine: NewEngine(client)}
}

func (d *Diagnoser) Diagnose(
	ctx context.Context,
	input string,
) (*DiagnoseResult, error) {
	prompt := fmt.Sprintf(prompts.DiagnosePromptTpl, input)
	return CompleteJSON[DiagnoseResult](
		ctx,
		d.engine,
		prompt,
		prompts.DiagnoserSystem,
		0.2,
		3000,
	)
}

type DiagnoseIssue struct {
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Location  string `json:"location"`
	ErrorCode string `json:"error_code"`
	Timestamp string `json:"timestamp"`
}

type DiagnoseRootCause struct {
	Primary             string   `json:"primary"`
	Category            string   `json:"category"`
	ContributingFactors []string `json:"contributing_factors"`
	Confidence          string   `json:"confidence"`
}

type DiagnoseSolution struct {
	Description     string   `json:"description"`
	Priority        string   `json:"priority"`
	Category        string   `json:"category"`
	Actionable      bool     `json:"actionable"`
	EstimatedEffort string   `json:"estimated_effort"`
	SideEffects     []string `json:"side_effects"`
}

type DiagnoseResult struct {
	ProblemDomain      string             `json:"problem_domain"`
	ProblemType        string             `json:"problem_type"`
	Severity           string             `json:"severity"`
	ImpactScope        string             `json:"impact_scope"`
	Summary            string             `json:"summary"`
	Issues             []DiagnoseIssue    `json:"issues"`
	RootCause          DiagnoseRootCause  `json:"root_cause"`
	DiagnosisSteps     []string           `json:"diagnosis_steps"`
	Solutions          []DiagnoseSolution `json:"solutions"`
	AffectedComponents []string           `json:"affected_components"`
	Dependencies       []string           `json:"dependencies"`
	PreventionMeasures []string           `json:"prevention_measures"`
	Confidence         float64            `json:"confidence"`
}
