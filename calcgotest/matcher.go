package calcgotest

// ShouldEqualErrors checks if actual and expected errors are equal.
func ShouldEqualErrors(actual interface{}, expected ...interface{}) string {
	e1 := actual.([]error)
	e2 := expected[0].([]error)

	if eqErrors(e1, e2) {
		return ""
	}

	return errorsError(e1, e2) + "(Should be Equal)"
}

func errorsToString(errors []error) string {
	var str string

	str += "(\n"
	for _, err := range errors {
		if err == nil {
			continue
		}
		str += err.Error() + "\n"
	}
	str += ")\n"

	return str
}

func errorsError(actual []error, expected []error) string {
	return "Expected: \n" +
		errorsToString(expected) +
		"Actual: \n" +
		errorsToString(actual)
}

func eqErrors(e1 []error, e2 []error) bool {
	if len(e1) != len(e2) {
		return false
	}

	for i := 0; i < len(e1); i++ {
		if e1[i] != e2[i] {
			return false
		}
	}

	return true
}
