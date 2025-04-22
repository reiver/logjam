package scheduler

type Scheduler struct {
}

func NewSchedulerSrv() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) CreateSchedule(reqModel CreateScheduleRequestModel) error {
	return nil
}

func (s *Scheduler) GetNotificationsToSend() {

}

func (s *Scheduler) Start() {

}
