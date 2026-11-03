// Package harness holds what the plugin carries for the apps it makes besides
// its programs: terminu project, the family's repo templates (project/),
// which moved here from tdp at tdp v1.2.0-beta.1.
package harness

import "embed"

// Project is project/: the page, its checklist items and the templates.
// project/version.go is a template, a package of its own when Go builds this
// module; nothing imports it.
//
//go:embed project
var Project embed.FS
