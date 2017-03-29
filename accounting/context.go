package accounting

import (
	"context"
	"time"
)

type recType int

const (
	newStage recType = iota
	newStep
	stageDone
	stepDone
)

type rec struct {
	Type recType
	Time time.Time
	Name string
}

// key is an unexported type for keys defined in this package.  This prevents
// collisions with keys defined in other packages.
type key int

// accountingKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var accountingKey key

// NewContext returns a context initialised for accounting
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, accountingKey, &[]rec{})
}

func fromContext(ctx context.Context) *[]rec {
	return ctx.Value(accountingKey).(*[]rec)
}

// NewStage records the start of a new stage.  Closing any previous stage by
// calling StageDone is not required.
func NewStage(ctx context.Context, name string) {
	items := fromContext(ctx)
	*items = append(*items, rec{Type: newStage, Time: time.Now(), Name: name})
}

// NewStep records the start of a new step.  NewStage must have been called at
// some previous time, but closing any previous step by calling StepDone is not
// required.
func NewStep(ctx context.Context, name string) {
	items := fromContext(ctx)
	*items = append(*items, rec{Type: newStep, Time: time.Now(), Name: name})
}

// StageDone records the end of a stage.  No further steps may be recorded
// before NewStage is called.
func StageDone(ctx context.Context) {
	items := fromContext(ctx)
	*items = append(*items, rec{Type: stageDone, Time: time.Now()})
}

// StepDone records the end of a step.
func StepDone(ctx context.Context) {
	items := fromContext(ctx)
	*items = append(*items, rec{Type: stepDone, Time: time.Now()})
}
