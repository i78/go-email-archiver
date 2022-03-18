package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFilename(t *testing.T) {

	type testcases struct {
		description string
		call        func(r Repository) string
		expected    string
	}
	tests := []testcases{
		{
			description: "Should return repo root when no subrepo and path given",
			call: func(r Repository) string {
				return createFilename("")
			},
			expected: "mailArchive/",
		}, {
			description: "Should return subrepo root when no further arguments given",
			call: func(r Repository) string {
				return createFilename(configSubfolderPrefix)
			},
			expected: "mailArchive" + "/" + configSubfolderPrefix,
		}, {
			description: "Should return subrepo 1st level",
			call: func(r Repository) string {
				return createFilename(configSubfolderPrefix, "sub1")
			},
			expected: "mailArchive" + "/" + configSubfolderPrefix + "/sub1",
		},
		{
			description: "Should return subrepo 1st level file",
			call: func(r Repository) string {
				return createFilename(configSubfolderPrefix, "sub1", "otter.json")
			},
			expected: "mailArchive" + "/" + configSubfolderPrefix + "/sub1/otter.json",
		}, {
			description: "Should return subrepo 1st level file with dashes",
			call: func(r Repository) string {
				return createFilename(configSubfolderPrefix, "sub1", "awesome-otter.json")
			},
			expected: "mailArchive" + "/" + configSubfolderPrefix + "/sub1/awesome-otter.json",
		}, {
			description: "Should sanitize backticks away",
			call: func(r Repository) string {
				return createFilename(configSubfolderPrefix, "sub1", "../zz", "otter.json")
			},
			expected: "mailArchive" + "/" + configSubfolderPrefix + "/sub1/zz/otter.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			subject := Repository{}
			result := tc.call(subject)
			assert.Equal(t, tc.expected, result)
		})
	}

}
