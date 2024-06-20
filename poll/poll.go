package poll

type Poll struct {
	question string
	options  []string
	fee      float64
	duration int
}

func CreatePoll(question string, options []string, duration int) Poll {
	p := Poll{question: question, options: options, duration: duration}
	p.calculateFee()
	return p
}

func (p *Poll) SetQuestion(question string) {
	p.question = question
	p.calculateFee()
}

func (p *Poll) SetOptions(options []string) {
	p.options = options
	p.calculateFee()
}

func (p *Poll) SetDuration(duration int) {
	p.duration = duration
	p.calculateFee()
}

func (p *Poll) GetQuestion() string {
	return p.question
}

func (p *Poll) GetOptions() []string {
	return p.options
}

func (p *Poll) GetDuration() int {
	return p.duration
}

func (p *Poll) GetFee() float64 {
	return p.fee
}

func (p *Poll) GetSignature() float64 {
	return p.signature
}

func (p *Poll) calculateFee() {
	p.fee = 0.01
}
