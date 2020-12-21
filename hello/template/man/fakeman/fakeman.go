package fakeman

import (
	"context"
	"hello/template/man"
	"sync/atomic"
)

type FakeMan struct {
	once atomic.Value
	opt  man.Options
	seq  int64
}

func NewFakeMan(opt ...man.Option) man.Man {
	opts := man.NewOptions(opt...)
	fakeMan := &FakeMan{
		opt: opts,
		seq: 0,
	}
	fakeMan.once.Store(false)
	man := man.Man(fakeMan)
	for i := range opts.Wrapper {
		man = opts.Wrapper[i](man)
	}
	return man
}

func (f *FakeMan) Init(opt ...man.Options) {

}
func (f *FakeMan) Eat(food string) {

}
func (f *FakeMan) Motion() {

}
func (f *FakeMan) Think() string {
	return ""
}
func (f *FakeMan) Play(pet man.Pet) {

}

func (f *FakeMan) Work(ctx context.Context, req string, msg string, opt man.WorkOptions) string {
	return ""
}
