package schedule

type ScheduleStorageService struct {
	scheduleStorage ScheduleStorage
}

func NewScheduleStorageService(scheduleStorage ScheduleStorage) *ScheduleStorageService {
	return &ScheduleStorageService{scheduleStorage: scheduleStorage}
}

func (s *ScheduleStorageService) DeleteSchedule(groupId string, date string) error {
	return s.scheduleStorage.DeleteSchedule(groupId, date)
}
