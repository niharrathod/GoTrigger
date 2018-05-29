package scheduler

import (
	"errors"
	"log"
	"sync"
	"time"
)

var container scheduledContainer

// ISchedule : scheduler understands the type which satisfies IScheduler interface,
type ISchedule interface {
	// GetID : returns an Identifier for the Scheduler
	GetID() string

	// GetDuration : return the Duration at which an activity will be schedule by the scheduler
	GetDuration() time.Duration

	// Execute : scheduler call this method on schedule is wake up, i.e. schedule activity triggered.
	Execute()
}

// Schedule : schedules the given activity.
func Schedule(schedule ISchedule) error {
	var sd scheduled
	sd.s = schedule
	container.put(schedule.GetID(), &sd)
	sd.scheduleTimer = time.AfterFunc(schedule.GetDuration(), func() {
		schedule.Execute()
		log.Print("ID:", schedule.GetID(), " is Executed")
		container.remove(schedule.GetID())
	})
	log.Print("ID:", schedule.GetID(), " is Scheduled")
	return nil
}

// Abort : abort already scheduler activity.
func Abort(ID string) error {
	s, err := container.get(ID)
	if err != nil {
		return err
	}

	isSuccess := s.scheduleTimer.Stop()
	if isSuccess {
		container.remove(ID)
		log.Print("Schedule with ID:", ID, " is Aborted")
	} else {
		return errors.New("Could not abort schedule")
	}
	return nil
}

type scheduled struct {
	s             ISchedule
	scheduleTimer *time.Timer
}

type scheduledContainer struct {
	concurrentMap sync.Map
}

func (sc *scheduledContainer) get(id string) (*scheduled, error) {
	ret, ok := sc.concurrentMap.Load(id)
	if ok == false {
		return nil, errors.New("could not find schedule in container by id:" + id)

	}

	s, ok := ret.(*scheduled)
	if ok == false {
		return nil, errors.New("for Id:" + id + " is not type Scheduled")
	}
	return s, nil
}

func (sc *scheduledContainer) put(id string, s *scheduled) {
	sc.concurrentMap.Store(id, s)
}

func (sc *scheduledContainer) remove(id string) {
	sc.concurrentMap.Delete(id)
}
