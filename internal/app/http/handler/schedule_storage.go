package handler

type ScheduleStorage interface {
	DeleteSchedule(groupId string, date string) error
}
