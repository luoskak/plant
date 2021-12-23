package abstract

func FreeSubs(subs []Subscription) (noClosedSubs []Subscription) {
	for _, sub := range subs {
		if !sub.IsClosed() {
			noClosedSubs = append(noClosedSubs, sub)
		}
	}
	return
}
