package schedule

type ScheduleService struct {
	scheduleStorage ScheduleStorage
}

func NewScheduleService(scheduleStorage ScheduleStorage) *ScheduleService {
	return &ScheduleService{scheduleStorage: scheduleStorage}
}

func (s *ScheduleService) DeleteSchedule(groupId string, date string) error {
	return s.scheduleStorage.DeleteSchedule(groupId, date)
}
