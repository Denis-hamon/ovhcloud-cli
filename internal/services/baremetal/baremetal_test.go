// SPDX-FileCopyrightText: 2026 OVH SAS <opensource@ovh.net>
//
// SPDX-License-Identifier: Apache-2.0

package baremetal

import (
	"testing"

	"github.com/maxatome/go-testdeep/td"
	"github.com/ovh/ovhcloud-cli/internal/display"
	"github.com/ovh/ovhcloud-cli/internal/flags"
)

// baseBaremetalObject mirrors the keys GetBaremetal always sets on the
// result map ("network", "serviceInfo", "tasks"), regardless of what the
// dedicated server API response itself contains.
func baseBaremetalObject() map[string]any {
	return map[string]any{
		"name":        "sd42-42",
		"region":      "EU",
		"os":          "debian12_64",
		"state":       "ok",
		"network":     map[string]any{},
		"serviceInfo": map[string]any{},
		"tasks":       []any{},
	}
}

// Regression test: some endpoints (e.g. dev/staging environments) don't
// populate "iam" on the dedicated server API response, either omitting the
// key entirely or returning it as null. The baremetal.tmpl template used to
// unconditionally chain-index into it (e.g. `index .Result "iam"
// "displayName"`), which made text/template fail with "index of nil
// pointer" as soon as "iam" was missing.
func TestBaremetalTemplate_RendersWithoutIAM(t *testing.T) {
	prevExitFunc := display.ExitFunc
	display.ExitFunc = func(int) {}
	t.Cleanup(func() { display.ExitFunc = prevExitFunc })

	object := baseBaremetalObject()
	// "iam" intentionally absent, as returned by some custom endpoints

	display.ResultError = nil
	display.OutputObject(object, "sd42-42", baremetalTemplate, &flags.OutputFormatConfig)

	td.CmpNoError(t, display.ResultError)
	td.Cmp(t, display.ResultString, td.Not(td.Contains("<no value>")))
}

func TestBaremetalTemplate_RendersWithIAM(t *testing.T) {
	prevExitFunc := display.ExitFunc
	display.ExitFunc = func(int) {}
	t.Cleanup(func() { display.ExitFunc = prevExitFunc })

	object := baseBaremetalObject()
	object["iam"] = map[string]any{
		"displayName": "my-server",
		"urn":         "urn:v1:eu:resource:dedicatedServer:sd42-42",
	}

	display.ResultError = nil
	display.OutputObject(object, "sd42-42", baremetalTemplate, &flags.OutputFormatConfig)

	td.CmpNoError(t, display.ResultError)
	td.Cmp(t, display.ResultString, td.Contains("my-server"))
}
