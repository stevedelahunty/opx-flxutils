package netUtils

import (
	"errors"
	"syscall"
)

const MAXEPOLLEVENTS = 32

type EPoll struct {
	fd     int
	event  syscall.EpollEvent
	events [MAXEPOLLEVENTS]syscall.EpollEvent
}

func NewEPoll(fd int) (*EPoll, error) {
	var err error
	var eFd int
	if eFd, err = syscall.EpollCreate1(0); err != nil {
		return nil, err
	}

	e := EPoll{}
	e.fd = eFd
	e.event.Events = syscall.EPOLLOUT
	e.event.Fd = int32(fd)
	if err = syscall.EpollCtl(eFd, syscall.EPOLL_CTL_ADD, fd, &e.event); err != nil {
		e.Close()
		return nil, err
	}

	return &e, nil
}

func (e *EPoll) Close() error {
	return syscall.Close(e.fd)
}

func (e *EPoll) Wait(msec int) error {
	nevents, err := syscall.EpollWait(e.fd, e.events[:], msec)
	if err != nil {
		return err
	}
	if nevents <= 0 {
		return errors.New("i/o timeout")
	}
	return nil
}
