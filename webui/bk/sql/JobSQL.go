package sql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ===== Job相关结构体 =====

// Job 任务表结构
type Job struct {
	JobID     string    `json:"job_id" db:"job_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// JobVersionMapping 任务版本映射表结构
type JobVersionMapping struct {
	JobID         string     `json:"job_id" db:"job_id"`
	Version       string     `json:"version" db:"version"`
	FatherVersion *string    `json:"father_version" db:"father_version"`
	Description   *string    `json:"description" db:"description"`
	IsExecute     *bool      `json:"is_excute" db:"is_excute"`
	ExecuteDate   *time.Time `json:"excute_date" db:"excute_date"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

// JobDetail 任务详情表结构
type JobDetail struct {
	JobID               string    `json:"job_id" db:"job_id"`
	Version             string    `json:"version" db:"version"`
	InputPrompt         *string   `json:"input_prompt" db:"input_prompt"`
	OutputPrompt        *string   `json:"output_promt" db:"output_promt"`
	OptimizeOrientation *string   `json:"optimize_orientation" db:"optimize_orientation"`
	OptimizedPrompt     *string   `json:"optimized_prompt" db:"optimized_prompt"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// JobDetailEvaluationMapping 任务详情评测映射表结构
type JobDetailEvaluationMapping struct {
	JobID           string    `json:"job_id" db:"job_id"`
	Version         string    `json:"version" db:"version"`
	EvaluationMapID string    `json:"evaluation_map_id" db:"evaluation_map_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// ===== 请求结构体 =====

// CreateJobRequest 创建任务请求
type CreateJobRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateJobRequest 更新任务请求
type UpdateJobRequest struct {
	Name *string `json:"name,omitempty"`
}

// CreateJobVersionRequest 创建任务版本请求
type CreateJobVersionRequest struct {
	JobID         string  `json:"job_id" binding:"required"`
	Version       string  `json:"version" binding:"required"`
	FatherVersion *string `json:"father_version,omitempty"`
	Description   *string `json:"description,omitempty"`
}

// UpdateJobVersionRequest 更新任务版本请求
type UpdateJobVersionRequest struct {
	FatherVersion *string    `json:"father_version,omitempty"`
	Description   *string    `json:"description,omitempty"`
	IsExecute     *bool      `json:"is_excute,omitempty"`
	ExecuteDate   *time.Time `json:"excute_date,omitempty"`
}

// CreateJobDetailRequest 创建任务详情请求
type CreateJobDetailRequest struct {
	JobID               string  `json:"job_id" binding:"required"`
	Version             string  `json:"version" binding:"required"`
	InputPrompt         *string `json:"input_prompt,omitempty"`
	OutputPrompt        *string `json:"output_promt,omitempty"`
	OptimizeOrientation *string `json:"optimize_orientation,omitempty"`
	OptimizedPrompt     *string `json:"optimized_prompt,omitempty"`
}

// UpdateJobDetailRequest 更新任务详情请求
type UpdateJobDetailRequest struct {
	InputPrompt         *string `json:"input_prompt,omitempty"`
	OutputPrompt        *string `json:"output_promt,omitempty"`
	OptimizeOrientation *string `json:"optimize_orientation,omitempty"`
	OptimizedPrompt     *string `json:"optimized_prompt,omitempty"`
}

// CreateJobDetailEvaluationMappingRequest 创建任务详情评测映射请求
type CreateJobDetailEvaluationMappingRequest struct {
	JobID           string `json:"job_id" binding:"required"`
	Version         string `json:"version" binding:"required"`
	EvaluationMapID string `json:"evaluation_map_id" binding:"required"`
}

// ===== Job CRUD 操作 =====

// GetAllJobs 获取所有任务
func GetAllJobs() ([]Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, name, created_at FROM job ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.Name, &job.CreatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetJobByID 根据ID获取任务
func GetJobByID(jobID string) (*Job, error) {
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
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &job, nil
}

// CreateJob 创建任务
func CreateJob(req CreateJobRequest) (*Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 生成UUID
	jobID := uuid.New().String()

	query := `INSERT INTO job (job_id, name) VALUES (?, ?)`

	_, err = db.Exec(query, jobID, req.Name)
	if err != nil {
		return nil, err
	}

	// 返回创建的任务
	return GetJobByID(jobID)
}

// UpdateJob 更新任务
func UpdateJob(jobID string, req UpdateJobRequest) (*Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查任务是否存在
	existingJob, err := GetJobByID(jobID)
	if err != nil {
		return nil, err
	}
	if existingJob == nil {
		return nil, sql.ErrNoRows
	}

	// 构建动态更新语句
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}

	if len(setParts) == 0 {
		return existingJob, nil // 没有需要更新的字段
	}

	query := "UPDATE job SET " + strings.Join(setParts, ", ") + " WHERE job_id = ?"
	args = append(args, jobID)

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// 返回更新后的任务
	return GetJobByID(jobID)
}

// DeleteJob 删除任务
func DeleteJob(jobID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查任务是否存在
	existingJob, err := GetJobByID(jobID)
	if err != nil {
		return err
	}
	if existingJob == nil {
		return sql.ErrNoRows
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除相关的映射关系
	_, err = tx.Exec("DELETE FROM job_detail_evaluation_mapping WHERE job_id = ?", jobID)
	if err != nil {
		return err
	}

	// 删除任务详情
	_, err = tx.Exec("DELETE FROM job_detail WHERE job_id = ?", jobID)
	if err != nil {
		return err
	}

	// 删除任务版本映射
	_, err = tx.Exec("DELETE FROM job_version_mapping WHERE job_id = ?", jobID)
	if err != nil {
		return err
	}

	// 删除任务
	_, err = tx.Exec("DELETE FROM job WHERE job_id = ?", jobID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// SearchJobs 搜索任务
func SearchJobs(keyword string) ([]Job, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, name, created_at FROM job WHERE name LIKE ? ORDER BY created_at DESC`

	searchPattern := "%" + keyword + "%"
	rows, err := db.Query(query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.JobID, &job.Name, &job.CreatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// ===== JobVersionMapping CRUD 操作 =====

// GetJobVersionsByJobID 根据任务ID获取所有版本
func GetJobVersionsByJobID(jobID string) ([]JobVersionMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, version, father_version, description, is_excute, excute_date, created_at 
			  FROM job_version_mapping 
			  WHERE job_id = ? 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []JobVersionMapping
	for rows.Next() {
		var version JobVersionMapping
		err := rows.Scan(&version.JobID, &version.Version, &version.FatherVersion,
			&version.Description, &version.IsExecute, &version.ExecuteDate, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

// GetJobVersionByJobIDAndVersion 根据任务ID和版本获取版本信息
func GetJobVersionByJobIDAndVersion(jobID, version string) (*JobVersionMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, version, father_version, description, is_excute, excute_date, created_at 
			  FROM job_version_mapping 
			  WHERE job_id = ? AND version = ?`

	row := db.QueryRow(query, jobID, version)

	var jobVersion JobVersionMapping
	err = row.Scan(&jobVersion.JobID, &jobVersion.Version, &jobVersion.FatherVersion,
		&jobVersion.Description, &jobVersion.IsExecute, &jobVersion.ExecuteDate, &jobVersion.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &jobVersion, nil
}

// CreateJobVersion 创建任务版本
func CreateJobVersion(req CreateJobVersionRequest) (*JobVersionMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查任务是否存在
	job, err := GetJobByID(req.JobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, sql.ErrNoRows
	}

	query := `INSERT INTO job_version_mapping (job_id, version, father_version, description) 
			  VALUES (?, ?, ?, ?)`

	_, err = db.Exec(query, req.JobID, req.Version, req.FatherVersion, req.Description)
	if err != nil {
		return nil, err
	}

	// 返回创建的版本
	return GetJobVersionByJobIDAndVersion(req.JobID, req.Version)
}

// UpdateJobVersion 更新任务版本
func UpdateJobVersion(jobID, version string, req UpdateJobVersionRequest) (*JobVersionMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查版本是否存在
	existingVersion, err := GetJobVersionByJobIDAndVersion(jobID, version)
	if err != nil {
		return nil, err
	}
	if existingVersion == nil {
		return nil, sql.ErrNoRows
	}

	// 构建动态更新语句
	var setParts []string
	var args []interface{}

	if req.FatherVersion != nil {
		setParts = append(setParts, "father_version = ?")
		args = append(args, *req.FatherVersion)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}
	if req.IsExecute != nil {
		setParts = append(setParts, "is_excute = ?")
		args = append(args, *req.IsExecute)
	}
	if req.ExecuteDate != nil {
		setParts = append(setParts, "excute_date = ?")
		args = append(args, *req.ExecuteDate)
	}

	if len(setParts) == 0 {
		return existingVersion, nil // 没有需要更新的字段
	}

	query := "UPDATE job_version_mapping SET " + strings.Join(setParts, ", ") + " WHERE job_id = ? AND version = ?"
	args = append(args, jobID, version)

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// 返回更新后的版本
	return GetJobVersionByJobIDAndVersion(jobID, version)
}

// DeleteJobVersion 删除任务版本
func DeleteJobVersion(jobID, version string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查版本是否存在
	existingVersion, err := GetJobVersionByJobIDAndVersion(jobID, version)
	if err != nil {
		return err
	}
	if existingVersion == nil {
		return sql.ErrNoRows
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除相关的映射关系
	_, err = tx.Exec("DELETE FROM job_detail_evaluation_mapping WHERE job_id = ? AND version = ?", jobID, version)
	if err != nil {
		return err
	}

	// 删除任务详情
	_, err = tx.Exec("DELETE FROM job_detail WHERE job_id = ? AND version = ?", jobID, version)
	if err != nil {
		return err
	}

	// 删除任务版本
	_, err = tx.Exec("DELETE FROM job_version_mapping WHERE job_id = ? AND version = ?", jobID, version)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// ===== JobDetail CRUD 操作 =====

// GetJobDetailByJobIDAndVersion 根据任务ID和版本获取任务详情
func GetJobDetailByJobIDAndVersion(jobID, version string) (*JobDetail, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, version, input_prompt, output_promt, optimize_orientation, optimized_prompt, created_at 
			  FROM job_detail 
			  WHERE job_id = ? AND version = ?`

	row := db.QueryRow(query, jobID, version)

	var detail JobDetail
	err = row.Scan(&detail.JobID, &detail.Version, &detail.InputPrompt, &detail.OutputPrompt,
		&detail.OptimizeOrientation, &detail.OptimizedPrompt, &detail.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &detail, nil
}

// CreateJobDetail 创建任务详情
func CreateJobDetail(req CreateJobDetailRequest) (*JobDetail, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查任务版本是否存在
	version, err := GetJobVersionByJobIDAndVersion(req.JobID, req.Version)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, sql.ErrNoRows
	}

	query := `INSERT INTO job_detail (job_id, version, input_prompt, output_promt, optimize_orientation, optimized_prompt) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, req.JobID, req.Version, req.InputPrompt, req.OutputPrompt,
		req.OptimizeOrientation, req.OptimizedPrompt)
	if err != nil {
		return nil, err
	}

	// 返回创建的详情
	return GetJobDetailByJobIDAndVersion(req.JobID, req.Version)
}

// UpdateJobDetail 更新任务详情
func UpdateJobDetail(jobID, version string, req UpdateJobDetailRequest) (*JobDetail, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查详情是否存在
	existingDetail, err := GetJobDetailByJobIDAndVersion(jobID, version)
	if err != nil {
		return nil, err
	}
	if existingDetail == nil {
		return nil, sql.ErrNoRows
	}

	// 构建动态更新语句
	var setParts []string
	var args []interface{}

	if req.InputPrompt != nil {
		setParts = append(setParts, "input_prompt = ?")
		args = append(args, *req.InputPrompt)
	}
	if req.OutputPrompt != nil {
		setParts = append(setParts, "output_promt = ?")
		args = append(args, *req.OutputPrompt)
	}
	if req.OptimizeOrientation != nil {
		setParts = append(setParts, "optimize_orientation = ?")
		args = append(args, *req.OptimizeOrientation)
	}
	if req.OptimizedPrompt != nil {
		setParts = append(setParts, "optimized_prompt = ?")
		args = append(args, *req.OptimizedPrompt)
	}

	if len(setParts) == 0 {
		return existingDetail, nil // 没有需要更新的字段
	}

	query := "UPDATE job_detail SET " + strings.Join(setParts, ", ") + " WHERE job_id = ? AND version = ?"
	args = append(args, jobID, version)

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// 返回更新后的详情
	return GetJobDetailByJobIDAndVersion(jobID, version)
}

// DeleteJobDetail 删除任务详情
func DeleteJobDetail(jobID, version string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查详情是否存在
	existingDetail, err := GetJobDetailByJobIDAndVersion(jobID, version)
	if err != nil {
		return err
	}
	if existingDetail == nil {
		return sql.ErrNoRows
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除相关的映射关系
	_, err = tx.Exec("DELETE FROM job_detail_evaluation_mapping WHERE job_id = ? AND version = ?", jobID, version)
	if err != nil {
		return err
	}

	// 删除任务详情
	_, err = tx.Exec("DELETE FROM job_detail WHERE job_id = ? AND version = ?", jobID, version)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// ===== JobDetailEvaluationMapping CRUD 操作 =====

// GetJobDetailEvaluationMappingsByJobIDAndVersion 根据任务ID和版本获取评测映射
func GetJobDetailEvaluationMappingsByJobIDAndVersion(jobID, version string) ([]JobDetailEvaluationMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, version, evaluation_map_id, created_at 
			  FROM job_detail_evaluation_mapping 
			  WHERE job_id = ? AND version = ? 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query, jobID, version)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []JobDetailEvaluationMapping
	for rows.Next() {
		var mapping JobDetailEvaluationMapping
		err := rows.Scan(&mapping.JobID, &mapping.Version, &mapping.EvaluationMapID, &mapping.CreatedAt)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

// CreateJobDetailEvaluationMapping 创建任务详情评测映射
// func CreateJobDetailEvaluationMapping(req CreateJobDetailEvaluationMappingRequest) (*JobDetailEvaluationMapping, error) {
// 	db, err := getDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	// 检查任务详情和评测集是否存在
// 	detail, err := GetJobDetailByJobIDAndVersion(req.JobID, req.Version)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if detail == nil {
// 		return nil, sql.ErrNoRows
// 	}

// 	evaluation, err := GetEvaluationMapByID(req.EvaluationMapID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if evaluation == nil {
// 		return nil, sql.ErrNoRows
// 	}

// 	query := `INSERT INTO job_detail_evaluation_mapping (job_id, version, evaluation_map_id)
// 			  VALUES (?, ?, ?)`

// 	_, err = db.Exec(query, req.JobID, req.Version, req.EvaluationMapID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 返回创建的映射
// 	return GetJobDetailEvaluationMappingByJobIDVersionAndEvaluationID(req.JobID, req.Version, req.EvaluationMapID)
// }

// GetJobDetailEvaluationMappingByJobIDVersionAndEvaluationID 根据任务ID、版本和评测ID获取映射
func GetJobDetailEvaluationMappingByJobIDVersionAndEvaluationID(jobID, version, evaluationMapID string) (*JobDetailEvaluationMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT job_id, version, evaluation_map_id, created_at 
			  FROM job_detail_evaluation_mapping 
			  WHERE job_id = ? AND version = ? AND evaluation_map_id = ?`

	row := db.QueryRow(query, jobID, version, evaluationMapID)

	var mapping JobDetailEvaluationMapping
	err = row.Scan(&mapping.JobID, &mapping.Version, &mapping.EvaluationMapID, &mapping.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &mapping, nil
}

// DeleteJobDetailEvaluationMapping 删除任务详情评测映射
func DeleteJobDetailEvaluationMapping(jobID, version, evaluationMapID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查映射是否存在
	existingMapping, err := GetJobDetailEvaluationMappingByJobIDVersionAndEvaluationID(jobID, version, evaluationMapID)
	if err != nil {
		return err
	}
	if existingMapping == nil {
		return sql.ErrNoRows
	}

	_, err = db.Exec("DELETE FROM job_detail_evaluation_mapping WHERE job_id = ? AND version = ? AND evaluation_map_id = ?",
		jobID, version, evaluationMapID)
	return err
}

// ===== 复合查询操作 =====

// GetJobWithAllDetails 获取任务及其所有相关信息
func GetJobWithAllDetails(jobID string) (*Job, []JobVersionMapping, []JobDetail, error) {
	job, err := GetJobByID(jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	if job == nil {
		return nil, nil, nil, nil
	}

	versions, err := GetJobVersionsByJobID(jobID)
	if err != nil {
		return nil, nil, nil, err
	}

	var details []JobDetail
	for _, version := range versions {
		detail, err := GetJobDetailByJobIDAndVersion(jobID, version.Version)
		if err != nil {
			return nil, nil, nil, err
		}
		if detail != nil {
			details = append(details, *detail)
		}
	}

	return job, versions, details, nil
}
