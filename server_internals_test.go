package raft

import (
	"bytes"
	"testing"
)

func TestFollowerAllegiance(t *testing.T) {
	// a follower with allegiance to leader=2
	s := Server{
		id:     1,
		term:   5,
		state:  &serverState{value: Follower},
		leader: 2,
		log:    NewLog(&bytes.Buffer{}, noop),
	}

	// receives an AppendEntries from a future term and different leader
	_, stepDown := s.handleAppendEntries(AppendEntries{
		Term:     6,
		LeaderId: 3,
	})

	// should now step down and have a new term
	if !stepDown {
		t.Errorf("wasn't told to step down (i.e. abandon leader)")
	}
	if s.term != 6 {
		t.Errorf("no term change")
	}
}

func TestStrongLeader(t *testing.T) {
	// a leader in term=2
	s := Server{
		id:     1,
		term:   2,
		state:  &serverState{value: Leader},
		leader: 1,
		log:    NewLog(&bytes.Buffer{}, noop),
	}

	// receives a RequestVote from someone also in term=2
	resp, stepDown := s.handleRequestVote(RequestVote{
		Term:         2,
		CandidateId:  3,
		LastLogIndex: 0,
		LastLogTerm:  0,
	})

	// and should retain his leadership
	if resp.VoteGranted {
		t.Errorf("shouldn't have granted vote")
	}
	if stepDown {
		t.Errorf("shouldn't have stepped down")
	}
}