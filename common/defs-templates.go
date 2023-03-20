package common

const CustomTemplateHelpCLI = `NAME:
{{template "helpNameTemplate" .}}
USAGE:
{{if .VisibleCommands}}command [command options] {{end}}{{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{if .Description}}
DESCRIPTION:
{{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}
COMMANDS:{{template "visibleCommandTemplate" .}}{{end}}{{if .VisibleFlagCategories}}
OPTIONS:{{template "visibleFlagCategoryTemplate" .}}{{else if .VisibleFlags}}
OPTIONS:{{template "visibleFlagTemplate" .}}{{end}}
`