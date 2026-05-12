// Package envchain provides a pipeline for chaining multiple env map
// transformations in sequence, applying each step in order.
package envchain

// Step is a function that transforms an env map.
type Step func(env map[string]string) map[string]string

// Chain holds an ordered list of transformation steps.
type Chain struct {
	steps []Step
}

// New creates a new Chain with the given steps.
func New(steps ...Step) *Chain {
	return &Chain{steps: steps}
}

// Add appends one or more steps to the chain.
func (c *Chain) Add(steps ...Step) *Chain {
	c.steps = append(c.steps, steps...)
	return c
}

// Run executes all steps in order, passing the output of each step as the
// input to the next. The original env map is never mutated; each step is
// responsible for returning a new map (or the same map if no change is needed).
func (c *Chain) Run(env map[string]string) map[string]string {
	current := copyEnv(env)
	for _, step := range c.steps {
		current = step(current)
	}
	return current
}

// Len returns the number of steps in the chain.
func (c *Chain) Len() int {
	return len(c.steps)
}

// copyEnv returns a shallow copy of env.
func copyEnv(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	return out
}
