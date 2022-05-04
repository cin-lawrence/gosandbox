package models

import (
        "gorm.io/datatypes"
)

type JobStatus string

const (
        JobStatusPending JobStatus = "PENDING"
        JobStatusRunning JobStatus = "RUNNING"
        JobStatusFinished JobStatus = "FINISHED"
        JobStatusFailed JobStatus = "FAILED"

)

var CREATE_ENUM_JOB_STATUS string = `DO $$
BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'job_status') THEN
                CREATE TYPE job_status AS ENUM('PENDING', 'RUNNING', 'FINISHED', 'FAILED');
        END IF;
END$$;`

func (s *JobStatus) Scan(value interface{}) error {
        *s = JobStatus(value.(string))
        return nil
}

func (s JobStatus) Value() (string, error) {
        return string(s), nil
}

type Job struct {
        BaseModel
        Status JobStatus `json:"status" sql:"type:job_status"`
        Result datatypes.JSON `json:"result"`
        UserID uint `json:"user_id" binding:"required"`
        User *User `json:"user"`
}
