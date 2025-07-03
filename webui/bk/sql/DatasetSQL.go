package sql

import (
	"strings"
	"time"
	"webui-server/model"
	"webui-server/util"

	"github.com/google/uuid"
)

// Dataset 数据集结构体 - 对应 dataset_map 表
type Dataset struct {
	DatasetMapID string    `json:"dataset_id" db:"dataset_id"`
	Name         string    `json:"name" db:"name"`
	DataCount    int       `json:"data_count" db:"data_count"`
	Description  *string   `json:"description,omitempty" db:"description"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
type DatasetCRUD struct{}

// DatasetDetail 数据集详情结构体 - 对应 dataset_detail 表
type DatasetDetail struct {
	DatasetDetailID string    `json:"dataset_detail_id" db:"dataset_detail_id"`
	DatasetID       string    `json:"dataset_id" db:"dataset_id"`
	Input           *string   `json:"input,omitempty" db:"input"`
	Target          *string   `json:"target,omitempty" db:"target"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
type DatasetDetailCRUD struct{}

// ===== Dataset CRUD 操作 =====

// GetAllDatasets 获取所有数据集详情
func (dc *DatasetCRUD) GetAllDatasets() ([]Dataset, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT dataset_id, name, data_count, description, created_at FROM dataset ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []Dataset
	for rows.Next() {
		var dataset Dataset
		err := rows.Scan(&dataset.DatasetMapID, &dataset.Name, &dataset.DataCount, &dataset.Description, &dataset.CreatedAt)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}

	return datasets, nil
}

// 添加数据集
func (dc *DatasetCRUD) AddDataset(name string, description string) error {

	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 生成UUID
	datasetMapID := uuid.New().String()
	// 插入数据集
	insertDatasetSQL := `INSERT INTO dataset (dataset_id, name, data_count, description, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(insertDatasetSQL, datasetMapID, name, 0, description, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// 获取数据集根据ID
func (dc *DatasetCRUD) GetDatasetByID(datasetID string) (*Dataset, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT * FROM dataset WHERE dataset_id = ?`
	row := db.QueryRow(query, datasetID)
	var dataset Dataset
	err = row.Scan(
		&dataset.DatasetMapID,
		&dataset.Name,
		&dataset.DataCount,
		&dataset.Description,
		&dataset.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

//获取数据集根据ID列表
func (dc *DatasetCRUD) GetDatasetByBatch(datasetID []string) ([]Dataset, error) {
	if len(datasetID) == 0 {
		return []Dataset{}, nil
	}

	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 构建IN查询的占位符
	placeholders := make([]string, len(datasetID))
	args := make([]interface{}, len(datasetID))
	for i, id := range datasetID {
		placeholders[i] = "?"
		args[i] = id
	}

	// 构建查询语句
	query := `SELECT dataset_id, name, data_count, description, created_at FROM dataset WHERE dataset_id IN (` + 
			strings.Join(placeholders, ",") + `) ORDER BY created_at DESC`

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []Dataset
	for rows.Next() {
		var dataset Dataset
		err := rows.Scan(&dataset.DatasetMapID, &dataset.Name, &dataset.DataCount, &dataset.Description, &dataset.CreatedAt)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return datasets, nil
}

// 删除数据集
func (dc *DatasetCRUD) DeleteDataset(datasetID string) error {
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

	// 先删除数据集的具体数据
	deleteDatasetDetailSQL := `DELETE FROM dataset_detail WHERE dataset_id = ?`
	_, err = tx.Exec(deleteDatasetDetailSQL, datasetID)
	if err != nil {
		return err
	}

	// 再删除数据集本身
	deleteDatasetSQL := `DELETE FROM dataset WHERE dataset_id = ?`
	_, err = tx.Exec(deleteDatasetSQL, datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 根据数据集基本
func (dc *DatasetCRUD) UpdateDataset(datasetID string, name string, description string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// 更新
	updateDatasetSQL := `UPDATE dataset SET name = ?, description = ? WHERE dataset_id = ?`
	_, err = db.Exec(updateDatasetSQL, name, description, datasetID)
	if err != nil {
		return err
	}
	return nil
}

// =====DatasetDetail CRUD 操作 =====
// 获取数据集byDatasetId
func (ddc *DatasetDetailCRUD) GetDatasetDetailByDatasetID(datasetID string) ([]DatasetDetail, error) {

	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT dataset_detail_id, dataset_id, input, target, created_at FROM dataset_detail WHERE dataset_id = ?`
	rows, err := db.Query(query, datasetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasetDetails []DatasetDetail
	for rows.Next() {
		var detail DatasetDetail
		err := rows.Scan(
			&detail.DatasetDetailID,
			&detail.DatasetID,
			&detail.Input,
			&detail.Target,
			&detail.CreatedAt)
		if err != nil {
			return nil, err
		}
		datasetDetails = append(datasetDetails, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return datasetDetails, nil
}

// 获取数据集byDatasetDetailID
func (ddc *DatasetDetailCRUD) GetDatasetDetailByID(datasetDetailID string) (*DatasetDetail, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT * FROM dataset_detail WHERE dataset_detail_id = ?`
	row := db.QueryRow(query, datasetDetailID)
	var datasetDetail DatasetDetail
	err = row.Scan(
		&datasetDetail.DatasetDetailID,
		&datasetDetail.DatasetID,
		&datasetDetail.Input,
		&datasetDetail.Target,
		&datasetDetail.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &datasetDetail, nil
}

// 添加数据集具体数据
func (ddc *DatasetDetailCRUD) AddDatasetDetail(datasetID, input, target string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务处理插入和计数更新
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	datasetDetailID := util.GenerateUUID()
	// 插入数据集
	insertDatasetSQL := `INSERT INTO dataset_detail (dataset_detail_id, dataset_id, input, target, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err = tx.Exec(insertDatasetSQL, datasetDetailID, datasetID, input, target, time.Now())
	if err != nil {
		return err
	}

	// 更新 dataset 表中的 data_count 字段
	updateCountSQL := `UPDATE dataset SET data_count = data_count + 1 WHERE dataset_id = ?`
	_, err = tx.Exec(updateCountSQL, datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 批量添加数据集具体数据
func (ddc *DatasetDetailCRUD) AddDatasetDetailByBatch(datasetID string, tuples []model.EditDataTuple) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if len(tuples) == 0 {
		return nil
	}

	// 使用事务处理批量插入
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 使用预处理语句提高性能
	insertSQL := `INSERT INTO dataset_detail (dataset_detail_id, dataset_id, input, target, created_at) VALUES (?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, detail := range tuples {
		datasetDetailID := util.GenerateUUID()
		_, err = stmt.Exec(datasetDetailID, datasetID, detail.Input, detail.Output, time.Now())
		if err != nil {
			return err
		}
	}

	// 更新 dataset 表中的 data_count 字段
	updateCountSQL := `UPDATE dataset SET data_count = data_count + ? WHERE dataset_id = ?`
	_, err = tx.Exec(updateCountSQL, len(tuples), datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteDatasetDetail 删除数据集详情
func (ddc *DatasetDetailCRUD) DeleteDatasetDetail(datasetID, datasetDetailID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务处理删除和计数更新
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// // 先获取要删除记录的 dataset_id
	// var datasetID string
	// selectSQL := `SELECT dataset_id FROM dataset_detail WHERE dataset_detail_id = ?`
	// err = tx.QueryRow(selectSQL, datasetDetailID).Scan(&datasetID)
	// if err != nil {
	// 	return err
	// }

	// 删除记录
	deleteSQL := `DELETE FROM dataset_detail WHERE dataset_detail_id = ?`
	_, err = tx.Exec(deleteSQL, datasetDetailID)
	if err != nil {
		return err
	}

	// 更新 dataset 表中的 data_count 字段
	updateCountSQL := `UPDATE dataset SET data_count = data_count - 1 WHERE dataset_id = ?`
	_, err = tx.Exec(updateCountSQL, datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (ddc *DatasetDetailCRUD) DeleteDatasetDetailByBatch(datasetID string, datasetDetailIDs []string) error {
	if len(datasetDetailIDs) == 0 {
		return nil
	}

	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务处理批量删除
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 使用预处理语句逐条删除，确保每个操作都能被正确处理
	deleteSQL := `DELETE FROM dataset_detail WHERE dataset_detail_id = ?`
	deleteStmt, err := tx.Prepare(deleteSQL)
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	for _, id := range datasetDetailIDs {
		_, err = deleteStmt.Exec(id)
		if err != nil {
			return err
		}
	}

	// 批量更新各个数据集的 data_count 字段
	updateCountSQL := `UPDATE dataset SET data_count = data_count - ? WHERE dataset_id = ?`
	updateStmt, err := tx.Prepare(updateCountSQL)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	_, err = updateStmt.Exec(len(datasetDetailIDs), datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 清空数据集具体数据根据数据集ID
func (dc *DatasetDetailCRUD) ClearDatasetDetail(datasetID string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 使用事务处理清空数据和更新计数
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 清空数据集数据
	deleteSQL := `DELETE FROM dataset_detail WHERE dataset_id = ?`
	_, err = tx.Exec(deleteSQL, datasetID)
	if err != nil {
		return err
	}

	// 更新 dataset 表中的 data_count 字段
	updateCountSQL := `UPDATE dataset SET data_count = 0 WHERE dataset_id = ?`
	_, err = tx.Exec(updateCountSQL, datasetID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// 更新数据集详情
func (dc *DatasetDetailCRUD) UpdateDatasetDetail(datasetDetailID, datasetID, input, target string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// 更新数据集
	updateDatasetSQL := `UPDATE dataset_detail SET input = ?, target = ? WHERE dataset_detail_id = ? AND dataset_id = ?`
	_, err = db.Exec(updateDatasetSQL, input, target, datasetDetailID, datasetID)
	if err != nil {
		return err
	}
	return nil
}

// 批量更新数据集详情
func (dc *DatasetDetailCRUD) UpdateDatasetDetailByBatch(datasetID string, tuples []model.EditDataTuple) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if len(tuples) == 0 {
		return nil
	}

	// 构建批量更新SQL - 使用事务处理
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 使用简单的逐条更新方式，避免复杂的CASE语句
	updateSQL := `UPDATE dataset_detail SET input = ?, target = ? WHERE dataset_detail_id = ? AND dataset_id = ?`
	stmt, err := tx.Prepare(updateSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, detail := range tuples {
		_, err = stmt.Exec(detail.Input, detail.Output, detail.ID, datasetID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
