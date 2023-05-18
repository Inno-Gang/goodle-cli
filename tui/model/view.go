package model

func (m *Model) View() (view string) {
	keyMapHelp := m.help.View(m)

	if header := m.state.Header(); header != "" {
		view = header + "\n\n" + m.state.View(m)
	} else {
		view = m.state.View(m)
	}

	return view + "\n\n" + keyMapHelp
}
