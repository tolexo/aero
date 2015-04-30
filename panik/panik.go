package panik

func On(err error) {
	if err != nil {
		panic(err)
	}
}

func If(condition bool, message string) {
	if condition {
		panic(message)
	}
}
