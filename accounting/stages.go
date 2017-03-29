package accounting

import (
	"context"
	"time"

	"k8s.io/kubernetes/pkg/api/unversioned"
)

type StageInfo struct {
	Name                 string
	StartTime            unversioned.Time
	DurationMilliseconds int64
	Steps                []*StepInfo
}

type StepInfo struct {
	Name                 string
	StartTime            unversioned.Time
	DurationMilliseconds int64
}

func ToStageInfos(ctx context.Context) []*StageInfo {
	var stages []*StageInfo
	var currentStage *StageInfo
	var currentStep *StepInfo

	closeStep := func(endTime time.Time) {
		if currentStep != nil {
			currentStep.DurationMilliseconds = int64(endTime.Sub(currentStep.StartTime.Time) / time.Millisecond)
			currentStep = nil
		}
	}

	closeStage := func(endTime time.Time) {
		closeStep(endTime)
		if currentStage != nil {
			currentStage.DurationMilliseconds = int64(endTime.Sub(currentStage.StartTime.Time) / time.Millisecond)
			currentStage = nil
		}
	}

	for _, rec := range *fromContext(ctx) {
		switch rec.Type {
		case newStage:
			closeStage(rec.Time)
			currentStage = &StageInfo{Name: rec.Name, StartTime: unversioned.Time{Time: rec.Time}}
			stages = append(stages, currentStage)

		case newStep:
			closeStep(rec.Time)
			currentStep = &StepInfo{Name: rec.Name, StartTime: unversioned.Time{Time: rec.Time}}
			if currentStage != nil { // ensure caller called NewStage before NewStep
				currentStage.Steps = append(currentStage.Steps, currentStep)
			}

		case stageDone:
			closeStage(rec.Time)

		case stepDone:
			closeStep(rec.Time)
		}
	}

	closeStage(time.Now()) // only has effect if caller forgot to call StageDone()

	return stages
}
