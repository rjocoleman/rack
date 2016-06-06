package manifest

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/convox/rack/changes"
)

type Sync struct {
	Local  string
	Remote string

	lock     sync.Mutex
	incoming []changes.Change
	outgoing []changes.Change
}

func NewSync(container, local, remote string) (*Sync, error) {
	l, err := filepath.Abs(local)

	if err != nil {
		return nil, err
	}

	sync := &Sync{
		Local:  l,
		Remote: remote,
	}

	return sync, nil
}

func (s *Sync) Contains(t Sync) bool {
	if !filepath.HasPrefix(t.Local, s.Local) {
		return false
	}

	lr, err := filepath.Rel(s.Local, t.Local)

	if err != nil {
		return false
	}

	rr, err := filepath.Rel(s.Remote, t.Remote)

	if err != nil {
		return false
	}

	return lr == rr
}

func (s *Sync) Start() error {
	go s.watchOutgoing()

	for range time.Tick(1 * time.Second) {
		go s.syncOutgoing()
	}

	return nil
}

func (s *Sync) watchOutgoing() {
	ch := make(chan changes.Change)

	go changes.Watch(s.Local, ch)

	for c := range ch {
		fmt.Printf("s: %+v\n", s)
		fmt.Printf("c: %+v\n", c)
		s.lock.Lock()
		s.outgoing = append(s.outgoing, c)
		s.lock.Unlock()
	}
}

func (s *Sync) syncOutgoing() {
	defer s.lock.Unlock()
	s.lock.Lock()
	fmt.Printf("s.outgoing: %+v\n", s.outgoing)
}

func resolvePath(path string) (string, error) {
	a, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	r, err := filepath.EvalSymlinks(a)

	if err != nil {
		return "", err
	}

	return r, nil
}
