package smtpserver

type Pipelining struct {
	*Extension
	OldProcessOperation
	OldHandleMore
}

var GROUP_COMMANDS []string

func (p *Pipelining) Init(parent interface{}) interface{} {
	GROUP_COMMANDS = []string{"RSET", "MAIL", "SEND", "SOML", "SAML", "RCPT"}
	p.Parent = parent
	return p
}

func (p *Pipelining) ExtendMode(mode bool) {
	if mode {
		p.OldProcessOperation = p.Parent.ProcessOperation
		p.Parent.ProcessOperation = p.ProcessOperation
		p.OldHandleMore = p.Parent.DataHandleMoreData
		p.Parent.DataHandleMoreData = true
	} else {
		if _, ok := p.OldProcessOperation; ok {
			p.Parent.ProcessOperation = p.OldProcessOperation
		}
		if _, ok := p.OldHandleMore; ok {
			p.Parent.DataHandleMoreData = p.OldHandleMore
		}
	}
}

func (p *Pipelining) ProcessOperation(operation interface{}) {
	commands := []string{}
	for i := 0; i <= len(commands); i++ {
		verb, params := obj.TokenizeCommand(commands[i])

		// Once the client SMTP has confirmed that support exists for
		// the pipelining extension, the client SMTP may then elect to
		// transmit groups of SMTP commands in batches without waiting
		// for a response to each individual command. In particular,
		// the commands RSET, MAIL FROM, SEND FROM, SOML FROM, SAML
		// FROM, and RCPT TO can all appear anywhere in a pipelined
		// command group. The EHLO, DATA, VRFY, EXPN, TURN, QUIT, and
		// NOOP commands can only appear as the last command in a group
		// since their success or failure produces a change of state
		// which the client SMTP must accommodate. (NOOP is included in
		// this group so it can be used as a synchronization point.)
		if i < len(commands) {
			obj.Reply(550, fmt.Sprintf("Protocol error: '%s' not allowed in a group of commands", verb))
			return false
		}

		return obj.ProcessCommand(verb, params)
	}

	return false
}

func (p *Pipelining) Keyword() string {
	return "PIPELINING"
}
