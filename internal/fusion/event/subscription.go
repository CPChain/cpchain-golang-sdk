package event

import (
	"sync"
)

type Subscription interface {
	Err() <-chan error // returns the error channel
	Unsubscribe()      // cancels sending of events, closing the error channel
}


type SubscriptionScope struct {
	mu     sync.Mutex
	subs   map[*scopeSub]struct{}
	closed bool
}

type scopeSub struct {
	sc *SubscriptionScope
	s  Subscription
}

// Track starts tracking a subscription. If the scope is closed, Track returns nil. The
// returned subscription is a wrapper. Unsubscribing the wrapper removes it from the
// scope.
func (sc *SubscriptionScope) Track(s Subscription) Subscription {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return nil
	}
	if sc.subs == nil {
		sc.subs = make(map[*scopeSub]struct{})
	}
	ss := &scopeSub{sc, s}
	sc.subs[ss] = struct{}{}
	return ss
}

// Close calls Unsubscribe on all tracked subscriptions and prevents further additions to
// the tracked set. Calls to Track after Close return nil.
func (sc *SubscriptionScope) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return
	}
	sc.closed = true
	for s := range sc.subs {
		s.s.Unsubscribe()
	}
	sc.subs = nil
}

// Count returns the number of tracked subscriptions.
// It is meant to be used for debugging.
func (sc *SubscriptionScope) Count() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return len(sc.subs)
}

func (s *scopeSub) Unsubscribe() {
	s.s.Unsubscribe()
	s.sc.mu.Lock()
	defer s.sc.mu.Unlock()
	delete(s.sc.subs, s)
}

func (s *scopeSub) Err() <-chan error {
	return s.s.Err()
}