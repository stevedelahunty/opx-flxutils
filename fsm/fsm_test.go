// fsm_test
package fsm_test

import (
	"fmt"
	"testing"
	"utils/fsm"
)

type MyFSM struct {
	FSM *fsm.Machine
}

type MyStateEvent struct {
	State fsm.State
	Event fsm.Event
}

func (se *MyStateEvent) CurrentState() fsm.State { return se.State }
func (se *MyStateEvent) CurrentEvent() fsm.Event { return se.Event }
func (se *MyStateEvent) SetState(s fsm.State)    { se.State = s }
func (se *MyStateEvent) SetEvent(e fsm.Event)    { se.Event = e }

const (
	exampleState1 = iota
	exampleState2
	exampleState3
)
const (
	exampleEvent1 = iota + 1
	exampleEvent2
)

func TestProcessEventNoStartCalled(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules}}

	rv := myFsm.FSM.ProcessEvent(exampleEvent1, nil)
	if rv != fsm.ErrorMachineNotStarted {
		t.Error("Expected Error", fsm.ErrorMachineNotStarted)
	}

	if 0 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", nil, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}

	if 0 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected no valid event stored\nActual", myFsm.FSM.Curr.CurrentEvent())
	}

}

func TestAddRuleDuplicateAdd(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rv := rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rv2 := rules.AddRule(exampleState1, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })
	rv3 := rules.AddRule(exampleState1, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	if rv != nil {
		t.Error("Expected no error")
	}
	if rv2 != nil {
		t.Error("Expected no error")
	}
	if rv3 != fsm.ErrorMachineStateEventExists {
		t.Error("Expected Error", fsm.ErrorMachineStateEventExists)
	}
}

func TestProcessEventBadEventForGivenState(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules,
		Begin: false}}

	// start state
	begin := myFsm.FSM.Start(exampleState1)
	fmt.Println("Begin", begin)

	rv := myFsm.FSM.ProcessEvent(exampleEvent2, nil)
	if rv != fsm.InvalidStateEvent {
		t.Error("Expected Error", fsm.InvalidStateEvent, "\nActual", rv)
	}
	if exampleState1 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState1, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}

	if 0 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected no valid event stored\nActual", myFsm.FSM.Curr.CurrentEvent())
	}

}

func TestProcessEventGoodStateTransition(t *testing.T) {

	rules := fsm.Ruleset{}

	// example rules
	rules.AddRule(exampleState1, exampleEvent1, func(m fsm.Machine, data interface{}) fsm.State { return exampleState2 })
	rules.AddRule(exampleState2, exampleEvent2, func(m fsm.Machine, data interface{}) fsm.State { return exampleState3 })

	myFsm := &MyFSM{FSM: &fsm.Machine{Curr: &MyStateEvent{},
		Rules: &rules}}

	// start state
	myFsm.FSM.Start(exampleState1)

	// First transition
	rv := myFsm.FSM.ProcessEvent(exampleEvent1, nil)
	if rv != nil {
		t.Error("Expected no error")
	}
	if exampleState2 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState2, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}
	if exampleEvent1 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected event", exampleEvent1, "\nActual event", myFsm.FSM.Curr.CurrentEvent())
	}

	// Second transition
	rv2 := myFsm.FSM.ProcessEvent(exampleEvent2, nil)
	if rv2 != nil {
		t.Error("Expected no error")
	}
	if exampleState3 != myFsm.FSM.Curr.CurrentState() {
		t.Error("Expected state", exampleState2, "\nActual state", myFsm.FSM.Curr.CurrentState())
	}
	if exampleEvent2 != myFsm.FSM.Curr.CurrentEvent() {
		t.Error("Expected event", exampleEvent2, "\nActual event", myFsm.FSM.Curr.CurrentEvent())
	}

}
