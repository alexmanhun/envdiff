package envchain_test

import (
	"strings"
	"testing"

	"envdiff/internal/envchain"
)

func upperStep(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = strings.ToUpper(v)
	}
	return out
}

func dropEmptyStep(env map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if v != "" {
			out[k] = v
		}
	}
	return out
}

func prefixStep(prefix string) envchain.Step {
	return func(env map[string]string) map[string]string {
		out := make(map[string]string, len(env))
		for k, v := range env {
			out[k] = prefix + v
		}
		return out
	}
}

func TestNew_EmptyChain(t *testing.T) {
	c := envchain.New()
	if c.Len() != 0 {
		t.Fatalf("expected 0 steps, got %d", c.Len())
	}
}

func TestRun_NoSteps_ReturnsCopy(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	c := envchain.New()
	out := c.Run(env)
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("unexpected output: %v", out)
	}
	// must be a copy
	out["A"] = "mutated"
	if env["A"] == "mutated" {
		t.Fatal("Run mutated the original env")
	}
}

func TestRun_SingleStep(t *testing.T) {
	env := map[string]string{"KEY": "hello"}
	c := envchain.New(upperStep)
	out := c.Run(env)
	if out["KEY"] != "HELLO" {
		t.Fatalf("expected HELLO, got %s", out["KEY"])
	}
}

func TestRun_MultipleSteps_OrderMatters(t *testing.T) {
	env := map[string]string{"K": "world"}
	c := envchain.New(upperStep, prefixStep(">>>"))
	out := c.Run(env)
	if out["K"] != ">>>WORLD" {
		t.Fatalf("expected >>>WORLD, got %s", out["K"])
	}
}

func TestRun_DropEmpty_ThenUpper(t *testing.T) {
	env := map[string]string{"A": "value", "B": ""}
	c := envchain.New(dropEmptyStep, upperStep)
	out := c.Run(env)
	if _, ok := out["B"]; ok {
		t.Fatal("expected B to be dropped")
	}
	if out["A"] != "VALUE" {
		t.Fatalf("expected VALUE, got %s", out["A"])
	}
}

func TestAdd_AppendsSteps(t *testing.T) {
	c := envchain.New(upperStep)
	c.Add(dropEmptyStep)
	if c.Len() != 2 {
		t.Fatalf("expected 2 steps, got %d", c.Len())
	}
}
