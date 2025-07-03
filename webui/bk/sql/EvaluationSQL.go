package sql

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ===== 评测集相关结构体 =====

// Evaluation 评测集表结构
type Evaluationset struct {
	EvaluationsetID    string    `json:"evaluationset_id" db:"evaluationset_id"`
	Name               string    `json:"name" db:"name"`
	SorceCap           *float64  `json:"sorce_cap" db:"sorce_cap"`
	EvaluationCriteria *string   `json:"evaluation_criteria" db:"evaluation_criteria"`
	Description        *string   `json:"description" db:"description"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
type EvaluationsetCRUD struct{}

// EvaluationsetDatasetMapping 评测集数据集合映射表结构
type EvaluationsetDatasetMapping struct {
	EvaluationsetDatasetMappingID string    `json:"evaluationset_dataset_mapping_id" db:"evaluationset_dataset_mapping_id"`
	EvaluationsetID               string    `json:"evaluationset_id" db:"evaluationset_id"`
	DatasetID                     string    `json:"dataset_id" db:"dataset_id"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
}
type EvaluationsetDatasetMappingCRUD struct{}

// ===== Evaluationset CRUD 操作 =====
// GetAllEvaluationset 获取所有评测集
func (ec *EvaluationsetCRUD) GetAllEvaluationsets() ([]Evaluationset, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT *
			  FROM evaluationset 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var evaluations []Evaluationset
	for rows.Next() {
		var evaluation Evaluationset
		err := rows.Scan(&evaluation.EvaluationsetID, &evaluation.Name, &evaluation.SorceCap,
			&evaluation.EvaluationCriteria, &evaluation.Description, &evaluation.CreatedAt)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}

	return evaluations, nil
}

// AddEvaluationset 创建评测集
func (ec *EvaluationsetCRUD) AddEvaluationset(name, description string, scoreCap float64) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 生成UUID
	evaluationsetID := uuid.New().String()

	// 插入评测集
	insertSQL := `INSERT INTO evaluationset (evaluationset_id, name, sorce_cap, description, created_at) 
				  VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, evaluationsetID, name, scoreCap, description, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// GetEvaluationsetByID 根据ID获取评测集
func (ec *EvaluationsetCRUD) GetEvaluationsetByID(evaluationsetID string) (*Evaluationset, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluationset_id, name, sorce_cap, evaluation_criteria, description, created_at 
			  FROM evaluationset 
			  WHERE evaluationset_id = ?`

	row := db.QueryRow(query, evaluationsetID)

	var evaluation Evaluationset
	err = row.Scan(&evaluation.EvaluationsetID, &evaluation.Name, &evaluation.SorceCap,
		&evaluation.EvaluationCriteria, &evaluation.Description, &evaluation.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 未找到记录
		}
		return nil, err
	}

	return &evaluation, nil
}

// DeleteEvaluationset 删除评测集
func (ec *EvaluationsetCRUD) DeleteEvaluationset(evaluationsetID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务确保数据一致性
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 先删除评测集与数据集的映射关系
	deleteMapping := `DELETE FROM evaluationset_dataset_mapping WHERE evaluationset_id = ?`
	_, err = tx.Exec(deleteMapping, evaluationsetID)
	if err != nil {
		return err
	}

	// 再删除评测集本身
	deleteEvaluationset := `DELETE FROM evaluationset WHERE evaluationset_id = ?`
	_, err = tx.Exec(deleteEvaluationset, evaluationsetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 设置评价标准
func (ec *EvaluationsetCRUD) SetEvaluationCriteria(evaluationsetID, evaluationCriteria string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 更新评测集
	updateSQL := `UPDATE evaluationset SET evaluation_criteria = ? WHERE evaluationset_id = ?`
	_, err = db.Exec(updateSQL, evaluationCriteria, evaluationsetID)
	if err != nil {
		return err
	}
	return nil
}

// 设置分数上限
func (ec *EvaluationsetCRUD) SetScoreCap(evaluationsetID string, scoreCap float64) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 更新评测集
	updateSQL := `UPDATE evaluationset SET sorce_cap = ? WHERE evaluationset_id = ?`
	_, err = db.Exec(updateSQL, scoreCap, evaluationsetID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEvaluationset 更新评测集
func (ec *EvaluationsetCRUD) UpdateEvaluationset(evaluationsetID, name, description string, scoreCap float64, evaluationCriteria string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 更新评测集
	updateSQL := `UPDATE evaluationset SET name = ?, description = ?, sorce_cap = ?, evaluation_criteria = ? WHERE evaluationset_id = ?`
	_, err = db.Exec(updateSQL, name, description, scoreCap, evaluationCriteria, evaluationsetID)
	if err != nil {
		return err
	}
	return nil
}

// ===== EvaluationsetDatasetMapping CRUD 操作 =====

// GetMappingsByEvaluationsetID 根据评测集ID获取所有映射关系
func (edmc *EvaluationsetDatasetMappingCRUD) GetMappingsByEvaluationsetID(evaluationsetID string) ([]EvaluationsetDatasetMapping, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT evaluationset_dataset_mapping_id, evaluationset_id, dataset_id, created_at 
			  FROM evaluationset_dataset_mapping 
			  WHERE evaluationset_id = ? 
			  ORDER BY created_at DESC`

	rows, err := db.Query(query, evaluationsetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []EvaluationsetDatasetMapping
	for rows.Next() {
		var mapping EvaluationsetDatasetMapping
		err := rows.Scan(&mapping.EvaluationsetDatasetMappingID, &mapping.EvaluationsetID,
			&mapping.DatasetID, &mapping.CreatedAt)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}


// AddDatasetMapping 添加评测集与数据集的映射关系
func (edmc *EvaluationsetDatasetMappingCRUD) AddDatasetMapping(evaluationsetID, datasetID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 生成UUID
	mappingID := uuid.New().String()

	// 插入映射关系
	insertSQL := `INSERT INTO evaluationset_dataset_mapping (evaluationset_dataset_mapping_id, evaluationset_id, dataset_id, created_at) 
				  VALUES (?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, mappingID, evaluationsetID, datasetID, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// 批量添加评测集与数据集的映射关系
func (edmc *EvaluationsetDatasetMappingCRUD) AddDatasetMappingByBatch(evaluationsetID string, datasetID []string) error {
	// 如果数据集ID列表为空，直接返回
	if len(datasetID) == 0 {
		return nil
	}

	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务确保数据一致性
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 准备批量插入的SQL语句
	insertSQL := `INSERT INTO evaluationset_dataset_mapping (evaluationset_dataset_mapping_id, evaluationset_id, dataset_id, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 批量插入映射关系
	currentTime := time.Now()
	for _, dsID := range datasetID {
		// 生成UUID
		mappingID := uuid.New().String()

		// 执行插入
		_, err = stmt.Exec(mappingID, evaluationsetID, dsID, currentTime)
		if err != nil {
			return err
		}
	}

	// 提交事务
	return tx.Commit()
}

// DeleteDatasetMapping 删除评测集与数据集的映射关系
func (edmc *EvaluationsetDatasetMappingCRUD) DeleteDatasetMapping(mappingID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM evaluationset_dataset_mapping WHERE evaluationset_dataset_mapping_id = ?`
	_, err = db.Exec(deleteSQL, mappingID)
	if err != nil {
		return err
	}
	return nil
}
//删除评测集与数据集的映射关系批次执行
func (edmc *EvaluationsetDatasetMappingCRUD) DeleteDatasetMappingByBatch(mappingIDs []string) error {
	// 如果映射ID列表为空，直接返回
	if len(mappingIDs) == 0 {
		return nil
	}

	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务确保数据一致性
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 准备批量删除的SQL语句
	deleteSQL := `DELETE FROM evaluationset_dataset_mapping WHERE evaluationset_dataset_mapping_id = ?`
	stmt, err := tx.Prepare(deleteSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// 批量删除映射关系
	for _, mappingID := range mappingIDs {
		_, err = stmt.Exec(mappingID)
		if err != nil {
			return err
		}
	}

	// 提交事务
	return tx.Commit()
}


// DeleteDatasetMappingByDatasetID 根据数据集ID删除映射关系
func (edmc *EvaluationsetDatasetMappingCRUD) DeleteDatasetMappingByDatasetID(evaluationsetID, datasetID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	deleteSQL := `DELETE FROM evaluationset_dataset_mapping WHERE evaluationset_id = ? AND dataset_id = ?`
	_, err = db.Exec(deleteSQL, evaluationsetID, datasetID)
	if err != nil {
		return err
	}
	return nil
}
