package schedule

type ScheduleStorage interface {
	DeleteSchedule(groupId string, date string) error
}
