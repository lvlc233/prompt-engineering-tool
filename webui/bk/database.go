package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDatabase() (*sql.DB, error) {
	// 创建或打开 SQLite 数据库文件
	db, err := sql.Open("sqlite", "./prompts.db")
	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		return nil, err
	}

	//数据集表
	datasetMapCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS dataset (
		dataset_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		data_count INTEGER DEFAULT 0,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(datasetMapCreateTableSQL); err != nil {
		return nil, err
	}

	//数据集详情表
	datasetDetailCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS dataset_detail (
		dataset_detail_id TEXT PRIMARY KEY,
		dataset_id TEXT NOT NULL,
		input TEXT,
		target TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(datasetDetailCreateTableSQL); err != nil {
		return nil, err
	}

	//评测集表
	evaluationMapCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS evaluation_map (
		evaluation_map_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		sorce_cap REAL DEFAULT 0,
		EvaluationCriteria TEXT,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(evaluationMapCreateTableSQL); err != nil {
		return nil, err
	}
	//评测集数据集合映射表
	evaluationDatasetMappingCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS evaluation_dataset_mapping (
		evaluation_dataset_mapping_id TEXT PRIMARY KEY,
		evaluation_map_id TEXT NOT NULL,
		dataset_map_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(evaluationDatasetMappingCreateTableSQL); err != nil {
		return nil, err
	}

	jobCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS job (
		job_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(jobCreateTableSQL); err != nil {
		return nil, err
	}
	jobVersionCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS job_version_mapping(
		job_id TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		father_version TEXT,
		description TEXT,
		is_excute BOOL,
		excute_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(jobVersionCreateTableSQL); err != nil {
		return nil, err
	}

	jobDetailCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS job_detail(
		job_id TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		input_prompt TEXT,
		output_promt TEXT,
		optimize_orientation TEXT,
		optimized_prompt TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err := db.Exec(jobDetailCreateTableSQL); err != nil {
		return nil, err
	}

	jobDetailEvaluationMappingCreateTableSQL := `
	CREATE TABLE IF NOT EXISTS job_detail_evaluation_mapping(
		job_id TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		evaluation_map_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(jobDetailEvaluationMappingCreateTableSQL); err != nil {
		return nil, err
	}

	return db, nil
}
