package internal

func logIfFatal(err error) {
	if err != nil {
		panic(err)
	}
}
