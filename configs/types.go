package configs

var (
	// ref: https://confluence.hive-intel.com/pages/viewpage.action?pageId=12550516#git%E8%A7%84%E8%8C%83-CommitLog
	Types = []string{
		"feat",
		"fix",
		"refactor",
		"style",
		"impr",
		"perf",
		"chore",
		"dep",
		"docs",
		"test",
		"typo",
		"revert",
		"merge",
		"wip",
	}
	BaseUrl              = ""
	Project              = ""
	ChangelogHeaderLines = 2
)
