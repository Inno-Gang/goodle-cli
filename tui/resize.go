package tui

func (m *model) size() (width, height int) {
	return m.width, m.height
}

func (m *model) resize(width, height int) {
	m.width = width
	m.height = height

	m.state.Resize(m.size())
}

func (s *stateError) Resize(_, _ int)   {}
func (s *stateLoading) Resize(_, _ int) {}
func (s *stateLogin) Resize(_, _ int)   {}

func (s *stateCourseSelection) Resize(width, height int) {
	s.list.SetSize(width, height)
}
