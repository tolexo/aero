package panik

func If(err error) {
	if err != nil {
		panic(err)
	}
}
