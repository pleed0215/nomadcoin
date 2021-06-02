package person

type Person struct {
	Name string
	Age int
}

func (p *Person) SayHello() string {
	return p.Name + " Says Hello!"
}

func (p *Person) AddAge() {
	p.Age += 1
}