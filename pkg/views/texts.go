package views

func OnOff(on bool) string {
	if on {
		return TextOn
	}

	return TextOff
}
