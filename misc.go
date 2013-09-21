package main

func CloseLater(closable *io.Closer) {
	err := closable.Close()
	if err != nil {
		panic(err)
	}
}
