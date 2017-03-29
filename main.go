package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/jim-minter/accounting-example/accounting"
)

func f(ctx context.Context) {
	accounting.NewStep(ctx, "f")
	f1(ctx)
}

func f1(ctx context.Context) {
	accounting.NewStep(ctx, "f1")
	time.Sleep(100 * time.Millisecond)
}

func g(ctx context.Context) {
	g1(ctx)
	time.Sleep(50 * time.Millisecond)
}

func g1(ctx context.Context) {
}

func main() {
	ctx := accounting.NewContext(context.Background())

	accounting.NewStage(ctx, "stage1")
	f(ctx)

	accounting.NewStage(ctx, "stage2")
	g(ctx)

	accounting.StageDone(ctx)

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	e.Encode(accounting.ToStageInfos(ctx))
}
