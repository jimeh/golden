package golden

var instance = NewStore()

func Update() bool {
	return instance.Update()
}

func Filename(tb TestingTB) string {
	return instance.Filename(tb)
}

func Get(tb TestingTB) []byte {
	return instance.Get(tb)
}

func Set(tb TestingTB, input []byte) {
	instance.Set(tb, input)
}
