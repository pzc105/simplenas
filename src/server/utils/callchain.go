package utils

type thencall interface {
	Then(call func() error) *onecall
	Start()
}

type onecall struct {
	thiscall  func() error
	onsuccess func()
	onfailed  func()
	parent    *onecall
}

func (c *onecall) OnSuccess(onsuccess func()) thencall {
	c.onsuccess = onsuccess
	return c
}

func (c *onecall) OnFailed(onfailed func()) thencall {
	c.onfailed = onfailed
	return c
}

func (c *onecall) Then(call func() error) *onecall {
	then := &onecall{
		thiscall: call,
		parent:   c,
	}
	return then
}

func (c *onecall) Start() {
	t := c
	var stack []*onecall
	stack = append(stack, t)
	for len(stack) > 0 {
		t = stack[len(stack)-1]
		if t.parent != nil {
			t.parent = nil
			stack = append(stack, t.parent)
			t.parent = nil
		} else {
			stack = stack[:len(stack)-1]
			err := t.thiscall()
			if err == nil {
				if t.onsuccess != nil {
					t.onsuccess()
				}
			} else {
				if t.onfailed != nil {
					t.onfailed()
				}
			}
		}
	}
}

func Call(c func() error) *onecall {
	return &onecall{
		thiscall: c,
	}
}
