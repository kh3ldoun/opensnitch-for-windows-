package wf

type Options struct{ Name string }

type Session struct{}

func New(*Options) (*Session, error) { return &Session{}, nil }
func (s *Session) Close() error      { return nil }

type Provider struct {
	id          string
	Name        string
	Description string
}

func (p Provider) ID() string { return p.id }

type Sublayer struct {
	id       string
	Name     string
	Provider string
	Weight   uint16
}

func (s Sublayer) ID() string { return s.id }

type Layer string

const (
	LayerALEAuthConnectV4 Layer = "ale_auth_connect_v4"
	LayerALEAuthConnectV6 Layer = "ale_auth_connect_v6"
)

type Action int

const (
	ActionPermit Action = iota
)

type Weight struct{}

func EmptyWeight() Weight { return Weight{} }

type Filter struct {
	Name       string
	Layer      Layer
	Sublayer   string
	Action     Action
	Weight     Weight
	Provider   string
	Persistent bool
}

func (s *Session) AddProvider(p *Provider) (Provider, error) {
	if p == nil {
		return Provider{}, nil
	}
	return Provider{id: "provider", Name: p.Name, Description: p.Description}, nil
}

func (s *Session) AddSublayer(sl *Sublayer) (Sublayer, error) {
	if sl == nil {
		return Sublayer{}, nil
	}
	return Sublayer{id: "sublayer", Name: sl.Name, Provider: sl.Provider, Weight: sl.Weight}, nil
}

func (s *Session) AddFilter(f *Filter) (Filter, error) {
	if f == nil {
		return Filter{}, nil
	}
	return *f, nil
}
