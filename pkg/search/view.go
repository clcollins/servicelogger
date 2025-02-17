package search

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	verticalPadding   = 1
	horizontalPadding = 2
	paddingStyle      = lipgloss.NewStyle().Padding(verticalPadding, horizontalPadding)
)

func (m *Model) getPaneWidth() int {
	x, _ := paddingStyle.GetFrameSize()
	if m.windowWidth <= x {
		return 0
	}
	return (m.windowWidth - x) / 2
}

func (m *Model) getPaneHeight() int {
	_, y := paddingStyle.GetFrameSize()
	if m.windowHeight <= y {
		return 0
	}
	return m.windowHeight - y
}

func (m *Model) View() string {
	m.list.SetSize(m.getPaneWidth()-horizontalPadding*2, m.getPaneHeight())
	md := m.templateSelection.String()
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("notty"),
		glamour.WithWordWrap(m.getPaneWidth()-1-horizontalPadding*4),
	)
	renderedMd, err := renderer.Render(md)
	if err != nil {
		renderedMd = md
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().
			Width(m.getPaneWidth()).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("227")).
			BorderLeft(false).BorderTop(false).BorderRight(true).BorderBottom(false).
			Render(
				paddingStyle.Render(m.list.View()),
			),
		paddingStyle.Width(m.getPaneWidth()).Render(renderedMd),
	)
}
