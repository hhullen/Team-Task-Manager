package gracefulterminator

var vault = []func(){}

func Add(f func()) {
	vault = append(vault, f)
}

func Stop() {
	for i := len(vault) - 1; i >= 0; i-- {
		vault[i]()
	}
}
