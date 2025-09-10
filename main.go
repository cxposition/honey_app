package main

type test interface {
	eat()
}

type Man struct {
	test
}

type ImTest struct {
}

func (m *ImTest) eat() {
	println("ImTest.")
}

func (m *Man) eat() {
	println("man eat.")
}

func main() {
	x := &Man{test: &ImTest{}}
	x.test.eat()
}
