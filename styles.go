package cli

import "github.com/charmbracelet/lipgloss"

var WarningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
var PurpleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
var HeaderStyle = lipgloss.NewStyle().Bold(true)
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

var QuestionTextStyle = HeaderStyle
