package plugins

import (
	"testing"
)

func TestEmpty(t *testing.T) {
    r := NewRegistry()
    r.Setup()
}

type okPlugin struct {
    ok bool
}

func (o *okPlugin) Init(*plugins) {
    o.ok = true
}

func TestAdd(t *testing.T) {
    r := NewRegistry()

    o := &okPlugin{}
    r.Add(o)

    r.Setup()

    if !o.ok {
        t.Error("Expected ok to be true")
    }
}

type firstGetPlugin struct {
    t *testing.T
    ok bool
}

func (f *firstGetPlugin) Init(p *plugins) {
    p2, ok := Get[*secondGetPlugin](p)

    if !ok {
        f.t.Fatal("Could not retrieve *secondGetPlugin")
    }

    p2.ok = true
}

type secondGetPlugin struct {
    t *testing.T
    ok bool
}

func (s *secondGetPlugin) Init(p *plugins) {
    p1, ok := Get[*firstGetPlugin](p)

    if !ok {
        s.t.Fatal("Could not retrieve *firstGetPlugin")
    }

    p1.ok = true
}

func TestGet(t *testing.T) {
    r := NewRegistry()

    p1 := &firstGetPlugin{t: t}
    r.Add(p1)

    p2 := &secondGetPlugin{t: t}
    r.Add(p2)

    r.Setup()

    if !p1.ok {
        t.Error("Expected p1.ok to be true")
    }

    if !p2.ok {
        t.Error("Expected p2.ok to be true")
    }
}

type selfRetrievingPlugin struct {
    t *testing.T
    ok bool
}

func (s *selfRetrievingPlugin) Init(p *plugins) {
    me, ok := Get[*selfRetrievingPlugin](p)

    if !ok {
        s.t.Fatal("Could not retrieve *selfRetrievingPlugin")
    }

    if s != me {
        s.t.Fatal("s does not reference the same memory as me")
    }

    me.ok = true
}

func TestGetItself(t *testing.T) {
    r := NewRegistry()

    s := &selfRetrievingPlugin{t: t}
    r.Add(s)

    r.Setup()

    if !s.ok {
        t.Error("Expected ok to be true")
    }
}

type missingPlugin struct {
    t *testing.T
}

func (m *missingPlugin) Init(p *plugins) {
    _, ok := Get[*okPlugin](p)

    if ok {
        m.t.Fatal("Expected ok to be false")
    }
}

func TestMissing(t *testing.T) {
    r := NewRegistry()

    mp := &missingPlugin{t: t}
    r.Add(mp)

    r.Setup()
}

func TestDuplicate(t *testing.T) {
    r := NewRegistry()

    o1 := &okPlugin{}

    if !r.Add(o1) {
        t.Error("Expected adding o1 to succeed")
    }

    o2 := &okPlugin{}

    if r.Add(o2) {
        t.Error("Expected adding o2 to fail")
    }

    r.Setup()

    if !o1.ok {
        t.Error("Expected o1.ok to be true")
    }

    if o2.ok {
        t.Error("Expected o2.ok to be false")
    }
}

type valuePlugin struct {}

func (valuePlugin) Init(*plugins) {}

func TestNonReference(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Fail()
        } else if r != "plugin is not a pointer" {
            t.Errorf("Incorrect panic message: %s", r)
        }
    }()

    r := NewRegistry()

    v := valuePlugin{}
    r.Add(v)

    r.Setup()
}

type getNonReferencePlugin struct {}

func (*getNonReferencePlugin) Init(p *plugins) {
    Get[getNonReferencePlugin](p)
}

func TestGetNonReference(t *testing.T) {
    defer func() {
        if r := recover(); r == nil {
            t.Fail()
        } else if r != "T should be a pointer" {
            t.Errorf("Incorrect panic message: %s", r)
        }
    }()

    r := NewRegistry()

    v := getNonReferencePlugin{}
    r.Add(&v)

    r.Setup()
}
