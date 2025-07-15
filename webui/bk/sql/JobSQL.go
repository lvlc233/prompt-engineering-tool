package sql

import (
	"database/sql"
	"fmt"
	"time"
	"webui-server/model"
	"webui-server/util"
)

// ===== Job相关结构体 =====

// Job 任务表结构
type Job struct {
	JobID            string    `json:"job_id" db:"job_id"`
	Name             string    `json:"name" db:"name"`
	Selected_version string    `json:"selected_version" db:"selected_version"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// JobVersionMapping 任务版本映射表结构
type JobVersionMapping struct {
	Version       string    `json:"version" db:"version"`
	JobID         string    `json:"job_id" db:"job_id"`
	FatherVersion *string    `json:"father_version" db:"father_version"`
	Description   *string    `json:"description" db:"description"`
	IsExecute     bool      `json:"is_execute" db:"is_execute"`
	ExecuteDate   *time.Time `json:"exceute_date" db:"execute_date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// JobDetail 任务详情表结构
type JobDetail struct {
	JobID               string    `json:"job_id" db:"job_id"`
	Version             string    `json:"version" db:"version"`
	InputPrompt         string    `json:"input_prompt" db:"input_prompt"`
	OutputPrompt        string    `json:"output_promt" db:"output_promt"`
	OptimizeOrientation string    `json:"optimize_orientation" db:"optimize_orientation"`
	OptimizedPrompt     string    `json:"optimized_prompt" db:"optimized_prompt"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// JobDetailEvaluationMapping 任务详情评测映射表结构
type JobDetailEvaluationMapping struct {
	JobID           string    `json:"job_id" db:"job_id"`
	Version         string    `json:"version" db:"version"`
	EvaluationsetID string    `json:"evaluationset_id" db:"evaluation_map_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// ===== CRUD 操作结构体 =====
type JobCRUD struct{}
type JobVersionMappingCRUD struct{}
type JobDetailCRUD struct{}
type JobDetailEvaluationMappingCRUD struct{}

// ===== Job CRUD 操作 =====

// GetAllJobs 获取所有任务
func (jc *JobCRUD) GetAllJobs() ([]Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			job_id, name, selected_version, description, created_at
		FROM job
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.Name, &job.Selected_version, &job.Description, &job.CreatedAt)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJobVersionsByID 获取单个任务的所有版本信息
func (jc *JobCRUD) GetJobVersionsByID(JobID string) ([]JobVersionMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			*
		FROM job_version_mapping
		WHERE job_id = ?
	`
	fmt.Println(JobID)
	rows, err := db.Query(query, JobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []JobVersionMapping
	for rows.Next() {

		var version JobVersionMapping
		err := rows.Scan( &version.Version,&version.JobID, &version.FatherVersion, &version.Description, &version.IsExecute, &version.ExecuteDate, &version.CreatedAt)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}
	return versions, nil

}

// AddJob 添加任务
func (jc *JobCRUD) AddJob(name, description string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	jobID := util.GenerateUUID()
	version := util.GenerateUUID()
	jobInsertSQL := `INSERT INTO job (
		job_id, 
		name, 
		selected_version,
		description,
		created_at) VALUES (?, ?, ?,?,?)`
	_, err = db.Exec(jobInsertSQL, jobID, name, version, description, time.Now())
	if err != nil {
		return err
	}
	versionInsertSQL := `INSERT INTO job_version_mapping (
		version,	
		job_id, 
		is_execute,
		created_at) VALUES (?, ?, ?,?)`
	_, err = db.Exec(versionInsertSQL, version,jobID, false, time.Now())
	if err != nil {
		return err
	}
	detailInsertSQL := `INSERT INTO job_detail (
		job_id, 
		version,
		created_at) VALUES (?, ?, ?)`
	_, err = db.Exec(detailInsertSQL, jobID, version, time.Now())

	return err
}

// GetJobByID 根据ID获取任务
func (jc *JobCRUD) GetJobByID(jobID string) (*Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, name, created_at FROM job WHERE job_id = ?`
	row := db.QueryRow(query, jobID)

	var job Job
	err = row.Scan(&job.JobID, &job.Name, &job.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &job, nil
}

// DeleteJob 删除任务
func (jc *JobCRUD) DeleteJob(jobID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除相关的评测映射
	_, err = tx.Exec(`DELETE FROM job_detail_evaluationset_mapping WHERE job_id = ?`, jobID)
	if err != nil {
		return err
	}

	// 删除任务详情
	_, err = tx.Exec(`DELETE FROM job_detail WHERE job_id = ?`, jobID)
	if err != nil {
		return err
	}

	// 删除版本映射
	_, err = tx.Exec(`DELETE FROM job_version_mapping WHERE job_id = ?`, jobID)
	if err != nil {
		return err
	}

	// 删除任务
	_, err = tx.Exec(`DELETE FROM job WHERE job_id = ?`, jobID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetJobDetail 获取任务详情
func (jc *JobCRUD) GetJobDetail(jobID string) (*model.JobDetailResponse, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT 
			j.job_id, j.name, j.created_at,
			jvm.version, jvm.father_version, jvm.description, jvm.is_execute, jvm.execute_date,
			jd.input_prompt, jd.output_promt, jd.optimize_orientation, jd.optimized_prompt
		FROM job j
		LEFT JOIN job_version_mapping jvm ON j.job_id = jvm.job_id
		LEFT JOIN job_detail jd ON j.job_id = jd.job_id AND jvm.version = jd.version
		WHERE j.job_id = ?
		ORDER BY jvm.version DESC
		LIMIT 1
	`

	row := db.QueryRow(query, jobID)

	var detail model.JobDetailResponse
	var createdAt time.Time
	var executeDate sql.NullTime
	var version, fatherVersion, description, inputPrompt, outputPrompt, optimizeOrientation, optimizedPrompt sql.NullString
	var isExecute sql.NullBool

	err = row.Scan(
		&detail.JobID, &detail.Name, &createdAt,
		&version, &fatherVersion, &description, &isExecute, &executeDate,
		&inputPrompt, &outputPrompt, &optimizeOrientation, &optimizedPrompt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	detail.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	if version.Valid {
		detail.Version = version.String
	}
	if fatherVersion.Valid {
		detail.FatherVersion = fatherVersion.String
	}
	if description.Valid {
		detail.Description = description.String
	}
	if inputPrompt.Valid {
		detail.InputPrompt = inputPrompt.String
	}
	if outputPrompt.Valid {
		detail.OutputPrompt = outputPrompt.String
	}
	if optimizeOrientation.Valid {
		detail.OptimizeOrientation = optimizeOrientation.String
	}
	if optimizedPrompt.Valid {
		detail.OptimizedPrompt = optimizedPrompt.String
	}
	if isExecute.Valid {
		detail.IsExecute = isExecute.Bool
	}
	if executeDate.Valid {
		detail.ExecuteDate = executeDate.Time.Format("2006-01-02 15:04:05")
	}

	return &detail, nil
}
