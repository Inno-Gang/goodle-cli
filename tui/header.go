package tui

func (s *stateError) Header() string {
	return "Error"
}

func (s *stateLoading) Header() string {
	return "Loading"
}

func (s *stateLogin) Header() string {
	return "Login"
}

func (s *stateCourseSelection) Header() string {
	return ""
}
