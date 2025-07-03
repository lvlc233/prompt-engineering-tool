package sql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ===== 评测集相关结构体 =====

// EvaluationMap 评测集表结构
type EvaluationMap struct {
	EvaluationMapID    string    `json:"evaluation_map_id" db:"evaluation_map_id"`
	Name               string    `json:"name" db:"name"`
	SorceCap           *float64  `json:"sorce_cap" db:"sorce_cap"`
	EvaluationCriteria *string   `json:"evaluation_criteria" db:"EvaluationCriteria"`
	Description        *string   `json:"description" db:"description"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

// EvaluationDatasetMapping 评测集数据集合映射表结构
type EvaluationDatasetMapping struct {
	EvaluationDatasetMappingID string    `json:"evaluation_dataset_mapping_id" db:"evaluation_dataset_mapping_id"`
	EvaluationMapID            string    `json:"evaluation_map_id" db:"evaluation_map_id"`
	DatasetMapID               string    `json:"dataset_map_id" db:"dataset_map_id"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
}

// ===== 请求结构体 =====

// CreateEvaluationMapRequest 创建评测集请求
type CreateEvaluationMapRequest struct {
	Name               string   `json:"name" binding:"required"`
	SorceCap           *float64 `json:"sorce_cap,omitempty"`
	EvaluationCriteria *string  `json:"evaluation_criteria,omitempty"`
	Description        *string  `json:"description,omitempty"`
}

// UpdateEvaluationMapRequest 更新评测集请求
type UpdateEvaluationMapRequest struct {
	Name               *string  `json:"name,omitempty"`
	SorceCap           *float64 `json:"sorce_cap,omitempty"`
	EvaluationCriteria *string  `json:"evaluation_criteria,omitempty"`
	Description        *string  `json:"description,omitempty"`
}

// CreateEvaluationDatasetMappingRequest 创建评测集数据集映射请求
type CreateEvaluationDatasetMappingRequest struct {
	EvaluationMapID string `json:"evaluation_map_id" binding:"required"`
	DatasetMapID    string `json:"dataset_map_id" binding:"required"`
}

// ===== EvaluationMap CRUD 操作 =====

// GetAllEvaluationMaps 获取所有评测集
func GetAllEvaluationMaps() ([]EvaluationMap, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluation_map_id, name, sorce_cap, EvaluationCriteria, description, created_at 
			  FROM evaluation_map 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []EvaluationMap
	for rows.Next() {
		var evaluation EvaluationMap
		err := rows.Scan(&evaluation.EvaluationMapID, &evaluation.Name, &evaluation.SorceCap,
			&evaluation.EvaluationCriteria, &evaluation.Description, &evaluation.CreatedAt)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}

	return evaluations, nil
}

// GetEvaluationMapByID 根据ID获取评测集
func GetEvaluationMapByID(evaluationMapID string) (*EvaluationMap, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluation_map_id, name, sorce_cap, EvaluationCriteria, description, created_at 
			  FROM evaluation_map 
			  WHERE evaluation_map_id = ?`

	row := db.QueryRow(query, evaluationMapID)

	var evaluation EvaluationMap
	err = row.Scan(&evaluation.EvaluationMapID, &evaluation.Name, &evaluation.SorceCap,
		&evaluation.EvaluationCriteria, &evaluation.Description, &evaluation.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &evaluation, nil
}

// CreateEvaluationMap 创建评测集
func CreateEvaluationMap(req CreateEvaluationMapRequest) (*EvaluationMap, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 生成UUID
	evaluationMapID := uuid.New().String()

	query := `INSERT INTO evaluation_map (evaluation_map_id, name, sorce_cap, EvaluationCriteria, description) 
			  VALUES (?, ?, ?, ?, ?)`

	_, err = db.Exec(query, evaluationMapID, req.Name, req.SorceCap, req.EvaluationCriteria, req.Description)
	if err != nil {
		return nil, err
	}

	// 返回创建的评测集
	return GetEvaluationMapByID(evaluationMapID)
}

// UpdateEvaluationMap 更新评测集
func UpdateEvaluationMap(evaluationMapID string, req UpdateEvaluationMapRequest) (*EvaluationMap, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 检查评测集是否存在
	existingEvaluation, err := GetEvaluationMapByID(evaluationMapID)
	if err != nil {
		return nil, err
	}
	if existingEvaluation == nil {
		return nil, sql.ErrNoRows
	}

	// 构建动态更新语句
	var setParts []string
	var args []interface{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}
	if req.SorceCap != nil {
		setParts = append(setParts, "sorce_cap = ?")
		args = append(args, *req.SorceCap)
	}
	if req.EvaluationCriteria != nil {
		setParts = append(setParts, "EvaluationCriteria = ?")
		args = append(args, *req.EvaluationCriteria)
	}
	if req.Description != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *req.Description)
	}

	if len(setParts) == 0 {
		return existingEvaluation, nil // 没有需要更新的字段
	}

	query := "UPDATE evaluation_map SET " + strings.Join(setParts, ", ") + " WHERE evaluation_map_id = ?"
	args = append(args, evaluationMapID)

	_, err = db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// 返回更新后的评测集
	return GetEvaluationMapByID(evaluationMapID)
}

// DeleteEvaluationMap 删除评测集
func DeleteEvaluationMap(evaluationMapID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 检查评测集是否存在
	existingEvaluation, err := GetEvaluationMapByID(evaluationMapID)
	if err != nil {
		return err
	}
	if existingEvaluation == nil {
		return sql.ErrNoRows
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除相关的映射关系
	_, err = tx.Exec("DELETE FROM evaluation_dataset_mapping WHERE evaluation_map_id = ?", evaluationMapID)
	if err != nil {
		return err
	}

	// 删除评测集
	_, err = tx.Exec("DELETE FROM evaluation_map WHERE evaluation_map_id = ?", evaluationMapID)
	if err != nil {
		return err
	}

	// 提交事务
	return tx.Commit()
}

// SearchEvaluationMaps 搜索评测集
func SearchEvaluationMaps(keyword string) ([]EvaluationMap, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluation_map_id, name, sorce_cap, EvaluationCriteria, description, created_at 
			  FROM evaluation_map 
			  WHERE name LIKE ? OR description LIKE ? OR EvaluationCriteria LIKE ?
			  ORDER BY created_at DESC`

	searchPattern := "%" + keyword + "%"
	rows, err := db.Query(query, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []EvaluationMap
	for rows.Next() {
		var evaluation EvaluationMap
		err := rows.Scan(&evaluation.EvaluationMapID, &evaluation.Name, &evaluation.SorceCap,
			&evaluation.EvaluationCriteria, &evaluation.Description, &evaluation.CreatedAt)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}

	return evaluations, nil
}

// ===== EvaluationDatasetMapping CRUD 操作 =====

// GetEvaluationDatasetMappingsByEvaluationID 根据评测集ID获取所有映射
func GetEvaluationDatasetMappingsByEvaluationID(evaluationMapID string) ([]EvaluationDatasetMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluation_dataset_mapping_id, evaluation_map_id, dataset_map_id, created_at 
			  FROM evaluation_dataset_mapping 
			  WHERE evaluation_map_id = ? 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query, evaluationMapID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []EvaluationDatasetMapping
	for rows.Next() {
		var mapping EvaluationDatasetMapping
		err := rows.Scan(&mapping.EvaluationDatasetMappingID, &mapping.EvaluationMapID,
			&mapping.DatasetMapID, &mapping.CreatedAt)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

// // CreateEvaluationDatasetMapping 创建评测集数据集映射
// func CreateEvaluationDatasetMapping(req CreateEvaluationDatasetMappingRequest) (*EvaluationDatasetMapping, error) {
// 	db, err := getDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	// 检查评测集和数据集是否存在
// 	evaluation, err := GetEvaluationMapByID(req.EvaluationMapID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if evaluation == nil {
// 		return nil, sql.ErrNoRows
// 	}

// 	dataset, err := GetDatasetByID(req.DatasetMapID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if dataset == nil {
// 		return nil, sql.ErrNoRows
// 	}

// 	// 生成UUID
// 	mappingID := uuid.New().String()

// 	query := `INSERT INTO evaluation_dataset_mapping (evaluation_dataset_mapping_id, evaluation_map_id, dataset_map_id)
// 			  VALUES (?, ?, ?)`

// 	_, err = db.Exec(query, mappingID, req.EvaluationMapID, req.DatasetMapID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 返回创建的映射
// 	return GetEvaluationDatasetMappingByID(mappingID)
// }

// // GetEvaluationDatasetMappingByID 根据ID获取评测集数据集映射
// func GetEvaluationDatasetMappingByID(mappingID string) (*EvaluationDatasetMapping, error) {
// 	db, err := getDB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	query := `SELECT evaluation_dataset_mapping_id, evaluation_map_id, dataset_map_id, created_at
// 			  FROM evaluation_dataset_mapping
// 			  WHERE evaluation_dataset_mapping_id = ?`

// 	row := db.QueryRow(query, mappingID)

// 	var mapping EvaluationDatasetMapping
// 	err = row.Scan(&mapping.EvaluationDatasetMappingID, &mapping.EvaluationMapID,
// 		&mapping.DatasetMapID, &mapping.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // 未找到记录
// 		}
// 		return nil, err
// 	}

// 	return &mapping, nil
// }

// // DeleteEvaluationDatasetMapping 删除评测集数据集映射
// func DeleteEvaluationDatasetMapping(mappingID string) error {
// 	db, err := getDB()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	// 检查映射是否存在
// 	existingMapping, err := GetEvaluationDatasetMappingByID(mappingID)
// 	if err != nil {
// 		return err
// 	}
// 	if existingMapping == nil {
// 		return sql.ErrNoRows
// 	}

// 	_, err = db.Exec("DELETE FROM evaluation_dataset_mapping WHERE evaluation_dataset_mapping_id = ?", mappingID)
// 	return err
// }
