package vexillum

import (
	"context"

	caesar "github.com/caesar-rocks/core"
)

type VexillumContextKey string

const (
	VEXILLUM_CONTEXT_KEY VexillumContextKey = "vexillum"
)

// Vexillum is a (very, very, very) lightweight feature flagging feature for Caesar.
type Vexillum struct {
	Flags map[string]bool
}

// New returns a new Vexillum instance.
func New() *Vexillum {
	return &Vexillum{
		Flags: make(map[string]bool),
	}
}

// IsActive returns true if the feature is active.
func (v *Vexillum) IsActive(feature string) bool {
	return v.Flags[feature]
}

// Activate activates a feature.
func (v *Vexillum) Activate(feature string) {
	v.Flags[feature] = true
}

// Deactivate deactivates a feature.
func (v *Vexillum) Deactivate(feature string) {
	v.Flags[feature] = false
}

// VexillumMiddleware is a middleware that injects feature flags into the context
// (so that it can be used in the templ views).
func VexillumMiddleware(vexillum *Vexillum) caesar.Handler {
	return func(ctx *caesar.Context) error {
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(ctx.Request.Context(), VEXILLUM_CONTEXT_KEY, vexillum.Flags),
		)

		ctx.Next()

		return nil
	}
}

// EnsureFeatureEnabledMiddleware returns a middleware that checks if a feature is active.
func (v *Vexillum) EnsureFeatureEnabledMiddleware(feature string) caesar.Handler {
	return func(ctx *caesar.Context) error {
		if v.IsActive(feature) {
			ctx.Next()
			return nil
		}
		return caesar.NewError(400)
	}
}

// IsFeatureActive returns true if the feature is active.
// This is a helper function that is aimed to be used in a templ file.
func IsFeatureActive(ctx context.Context, feature string) bool {
	flags := ctx.Value(VEXILLUM_CONTEXT_KEY).(map[string]bool)
	return flags[feature]
}
