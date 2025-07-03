import React, { useState, useEffect } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';

interface CriteriaItem {
  id: string;
  name: string;
  description: string;
  score: number;
  createdAt: string;
}

interface Dataset {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  itemCount: number;
}

const CriteriaDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const isViewMode = searchParams.get('mode') === 'view';
  const [items, setItems] = useState<CriteriaItem[]>([]);
  const [isDeleteMode, setIsDeleteMode] = useState(false);
  const [selectedItems, setSelectedItems] = useState<Set<string>>(new Set());
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [criteriaText, setCriteriaText] = useState('');
  const [isDatasetModalOpen, setIsDatasetModalOpen] = useState(false);
  const [datasets, setDatasets] = useState<Dataset[]>([]);
  const [selectedDatasets, setSelectedDatasets] = useState<Set<string>>(new Set());
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'createdAt' | 'itemCount'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [isScoreLimitModalOpen, setIsScoreLimitModalOpen] = useState(false);
  const [scoreLimit, setScoreLimit] = useState(100);

  // 模拟测评集信息
  const criteriaSet = {
    id: id || '1',
    name: '包含数据集',
    description: '绑定数据集和评测集合关系'
  };

  const handleAddItem = () => {
    const newItem: CriteriaItem = {
      id: Date.now().toString(),
      name: `评测项 ${items.length + 1}`,
      description: '新的评测项描述',
      score: 10,
      createdAt: new Date().toLocaleDateString()
    };
    setItems([...items, newItem]);
  };

  const handleBatchDelete = () => {
    if (isDeleteMode) {
      // 确认删除选中的项目
      setItems(items.filter(item => !selectedItems.has(item.id)));
      setSelectedItems(new Set());
      setIsDeleteMode(false);
    } else {
      // 进入批量删除模式
      setIsDeleteMode(true);
      setSelectedItems(new Set());
    }
  };

  const handleCancelDelete = () => {
    setIsDeleteMode(false);
    setSelectedItems(new Set());
  };

  const handleItemSelect = (itemId: string) => {
    const newSelected = new Set(selectedItems);
    if (newSelected.has(itemId)) {
      newSelected.delete(itemId);
    } else {
      newSelected.add(itemId);
    }
    setSelectedItems(newSelected);
  };

  const handleItemClick = (itemId: string) => {
    if (isDeleteMode) {
      handleItemSelect(itemId);
    } else {
      // 查看详情逻辑
      console.log('查看详情:', itemId);
      // 这里可以添加跳转到详情页面或打开详情弹窗的逻辑
      alert(`查看评测项详情: ${itemId}`);
    }
  };

  const handleSaveCriteria = () => {
    // 这里可以添加保存逻辑
    console.log('保存评测标准:', criteriaText);
    setIsModalOpen(false);
  };

  const handleSaveScoreLimit = () => {
    // 这里可以添加保存分数上限的逻辑
    console.log('保存分数上限:', scoreLimit);
    setIsScoreLimitModalOpen(false);
  };

  // 获取数据集列表
  const fetchDatasets = async () => {
    try {
      const response = await fetch('http://localhost:3000/datasets');
      const data = await response.json();
      setDatasets(data);
    } catch (error) {
      console.error('获取数据集失败:', error);
      // 模拟数据作为备选
      setDatasets([
        { id: '1', name: '数据集1', description: '这是第一个数据集', createdAt: '2024-01-01', itemCount: 100 },
        { id: '2', name: '数据集2', description: '这是第二个数据集', createdAt: '2024-01-02', itemCount: 200 },
        { id: '3', name: '数据集3', description: '这是第三个数据集', createdAt: '2024-01-03', itemCount: 150 }
      ]);
    }
  };

  useEffect(() => {
    if (isDatasetModalOpen) {
      fetchDatasets();
    }
  }, [isDatasetModalOpen]);

  // 数据集选择处理
  const handleDatasetSelect = (datasetId: string) => {
    const newSelected = new Set(selectedDatasets);
    if (newSelected.has(datasetId)) {
      newSelected.delete(datasetId);
    } else {
      newSelected.add(datasetId);
    }
    setSelectedDatasets(newSelected);
  };

  // 添加选中的数据集
  const handleAddSelectedDatasets = () => {
    const selectedDatasetList = datasets.filter(dataset => selectedDatasets.has(dataset.id));
    selectedDatasetList.forEach(dataset => {
      const newItem: CriteriaItem = {
        id: Date.now().toString() + Math.random(),
        name: dataset.name,
        description: dataset.description,
        score: 10,
        createdAt: dataset.createdAt
      };
      setItems(prev => [...prev, newItem]);
    });
    setSelectedDatasets(new Set());
    setIsDatasetModalOpen(false);
  };

  // 过滤和排序数据集
  const filteredAndSortedDatasets = datasets
    .filter(dataset => 
      dataset.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      dataset.description.toLowerCase().includes(searchTerm.toLowerCase())
    )
    .sort((a, b) => {
      const aValue = a[sortBy];
      const bValue = b[sortBy];
      const comparison = aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
      return sortOrder === 'asc' ? comparison : -comparison;
    });

  // 排序处理
  const handleSort = (field: 'name' | 'createdAt' | 'itemCount') => {
    if (sortBy === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };

  const canDelete = items.length > 0;
  const hasSelectedItems = selectedItems.size > 0;

  return (
    <div className="p-6">
        {/* 头部区域 - 标题和右上角按钮 */}
        <div className="mb-6 flex justify-between items-start">
          <div>
            <h1 className="text-2xl font-bold text-gray-800 mb-2">{criteriaSet.name}</h1>
            <p className="text-gray-600">{criteriaSet.description}</p>
          </div>
          
          {/* 右上角按钮组 */}
          {!isViewMode && (
            <div className="flex gap-2">
              <button
                onClick={handleBatchDelete}
                disabled={!canDelete}
                className={`px-4 py-2 text-sm rounded transition-colors border-0 ${
                  canDelete 
                    ? 'bg-red-500 text-white hover:bg-red-600' 
                    : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                }`}
              >
                批量删除
              </button>
              
              {isDeleteMode && (
                <>
                  <button
                    onClick={handleCancelDelete}
                    className="px-4 py-2 text-sm border-0 bg-gray-100 rounded hover:bg-gray-200 transition-colors"
                  >
                    取消
                  </button>
                  
                  <button
                    onClick={handleBatchDelete}
                    disabled={!hasSelectedItems}
                    className={`px-4 py-2 text-sm rounded transition-colors border-0 ${
                      hasSelectedItems 
                        ? 'bg-red-600 text-white hover:bg-red-700' 
                        : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                    }`}
                  >
                    确认删除 ({selectedItems.size})
                  </button>
                </>
              )}
            </div>
          )}
        </div>

        {/* 功能按钮区域 */}
        <div className="mb-6 flex gap-2">
          <button
            onClick={() => setIsModalOpen(true)}
            className="modern-button"
          >
            {isViewMode ? '查看评测标准' : '设置评测标准'}
          </button>
          
          <button
            onClick={() => setIsScoreLimitModalOpen(true)}
            className="modern-button"
          >
            {isViewMode ? '查看分数上限' : '设置分数上限'}
          </button>
          
          {!isViewMode && (
             <button
               onClick={() => setIsDatasetModalOpen(true)}
               className="modern-button"
             >
               + 新增数据集
             </button>
           )}
        </div>

        {/* 评测项统计 */}
        <div className="mb-6">
          <div className="text-sm text-gray-500">
            {items.length === 0 ? '暂无评测项' : `共 ${items.length} 个评测项`}
          </div>
        </div>

      {items.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-4">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p className="text-gray-500 mb-4">暂无评测项</p>
          <p className="text-sm text-gray-400">点击右上角的"新增"按钮添加评测项</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {items.map((item) => (
            <div 
              key={item.id} 
              className={`modern-card p-6 flex flex-col relative ${
                isViewMode 
                  ? 'cursor-default' 
                  : `cursor-pointer hover:bg-gray-50 transition-colors ${
                      isDeleteMode ? 'ring-2 ring-blue-200' : 'hover:shadow-lg'
                    }`
              }`}
              onClick={isViewMode ? undefined : () => handleItemClick(item.id)}
            >
              {isDeleteMode && !isViewMode && (
                <div className="absolute top-4 left-4">
                  <input
                    type="checkbox"
                    checked={selectedItems.has(item.id)}
                    onChange={() => handleItemSelect(item.id)}
                    className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 pointer-events-none"
                  />
                </div>
              )}
              
              <div className={`mb-3 ${isDeleteMode && !isViewMode ? 'ml-8' : ''}`}>
                <div className="flex justify-between items-start mb-2">
                  <h3 className="text-lg font-semibold text-gray-800">{item.name}</h3>
                  {!isDeleteMode && !isViewMode && (
                    <button 
                      onClick={(e) => {
                        e.stopPropagation();
                        // 删除单个项目的逻辑
                        setItems(items.filter(i => i.id !== item.id));
                      }}
                      className="px-2 py-1 text-xs border-0 bg-red-50 text-red-600 rounded hover:bg-red-100 transition-colors flex-shrink-0"
                    >
                      删除
                    </button>
                  )}
                </div>
                <p className="text-gray-600 text-sm mb-3 flex-1">{item.description}</p>
              </div>
              
              <div className="mb-4 space-y-2">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">分数:</span>
                  <span className="text-gray-700">{item.score}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-500">创建时间:</span>
                  <span className="text-gray-700">{item.createdAt}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
      
      {/* 评测标准弹窗 */}
      {isModalOpen && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 50
        }}>
          <div style={{
            backgroundColor: 'white',
            borderRadius: '8px',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
            width: '33.333333%',
            height: '50%',
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden'
          }}>
            {/* 弹窗头部 */}
            <div style={{
              padding: '24px 24px 16px 24px',
              borderBottom: '1px solid #dcfce7'
            }}>
              <h3 style={{
                fontSize: '18px',
                fontWeight: '600',
                color: '#111827',
                margin: 0
              }}>{isViewMode ? '查看评测标准' : '设置评测标准'}</h3>
            </div>

            {/* 弹窗内容 - 可拖动输入框 */}
            <div style={{
              flex: 1,
              padding: '24px',
              overflow: 'hidden'
            }}>
              <textarea
                value={criteriaText}
                onChange={isViewMode ? undefined : (e) => setCriteriaText(e.target.value)}
                placeholder={isViewMode ? "暂无评测标准" : "请输入评测标准..."}
                readOnly={isViewMode}
                style={{
                  width: '100%',
                  height: '100%',
                  padding: '16px',
                  border: '1px solid #bbf7d0',
                  borderRadius: '8px',
                  fontSize: '14px',
                  resize: 'none',
                  outline: 'none',
                  transition: 'all 0.2s',
                  fontFamily: 'inherit',
                  boxSizing: 'border-box',
                  backgroundColor: isViewMode ? '#f9fafb' : 'white',
                  cursor: isViewMode ? 'default' : 'text'
                }}
                onFocus={isViewMode ? undefined : (e) => {
                  e.target.style.borderColor = '#4ade80';
                  e.target.style.boxShadow = '0 0 0 2px rgba(74, 222, 128, 0.3)';
                }}
                onBlur={isViewMode ? undefined : (e) => {
                  e.target.style.borderColor = '#bbf7d0';
                  e.target.style.boxShadow = 'none';
                }}
              />
            </div>

            {/* 弹窗底部 */}
            <div style={{
              padding: '16px 24px 24px 24px',
              borderTop: '1px solid #dcfce7',
              display: 'flex',
              justifyContent: 'center',
              gap: '16px'
            }}>
              {!isViewMode && (
                <button
                  onClick={handleSaveCriteria}
                  disabled={!criteriaText.trim()}
                  style={{
                    padding: '8px 24px',
                    backgroundColor: criteriaText.trim() ? '#10b981' : '#9ca3af',
                    color: 'white',
                    borderRadius: '8px',
                    border: 'none',
                    cursor: criteriaText.trim() ? 'pointer' : 'not-allowed',
                    transition: 'background-color 0.2s',
                    fontSize: '14px',
                    fontWeight: '500'
                  }}
                  onMouseEnter={(e) => {
                    if (criteriaText.trim()) {
                      (e.target as HTMLButtonElement).style.backgroundColor = '#059669';
                    }
                  }}
                  onMouseLeave={(e) => {
                    if (criteriaText.trim()) {
                      (e.target as HTMLButtonElement).style.backgroundColor = '#10b981';
                    }
                  }}
                >
                  保存
                </button>
              )}
              <button
                onClick={() => setIsModalOpen(false)}
                style={{
                  padding: '8px 24px',
                  color: '#374151',
                  backgroundColor: 'white',
                  border: '1px solid #d1d5db',
                  borderRadius: '8px',
                  cursor: 'pointer',
                  transition: 'background-color 0.2s',
                  fontSize: '14px',
                  fontWeight: '500'
                }}
                onMouseEnter={(e) => {
                  (e.target as HTMLButtonElement).style.backgroundColor = '#f9fafb';
                }}
                onMouseLeave={(e) => {
                  (e.target as HTMLButtonElement).style.backgroundColor = 'white';
                }}
              >
                关闭
              </button>
            </div>
          </div>
        </div>
      )}

      {/* 数据集选择弹窗 */}
      {isDatasetModalOpen && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'flex-end',
          zIndex: 50
        }}>
          <div style={{
            backgroundColor: 'white',
            width: '50%',
            height: '100vh',
            display: 'flex',
            flexDirection: 'column',
            boxShadow: '-4px 0 15px rgba(0, 0, 0, 0.1)'
          }}>
            {/* 弹窗头部 */}
            <div style={{
              padding: '24px',
              borderBottom: '1px solid #e5e7eb',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              <h3 style={{
                fontSize: '20px',
                fontWeight: '600',
                color: '#111827',
                margin: 0
              }}>选择数据集</h3>
              <button
                onClick={() => {
                  setIsDatasetModalOpen(false);
                  setSelectedDatasets(new Set());
                  setSearchTerm('');
                }}
                style={{
                  padding: '8px',
                  backgroundColor: 'transparent',
                  border: 'none',
                  cursor: 'pointer',
                  borderRadius: '4px',
                  color: '#6b7280'
                }}
              >
                ✕
              </button>
            </div>

            {/* 搜索和排序区域 */}
            <div style={{
              padding: '16px 24px',
              borderBottom: '1px solid #e5e7eb',
              display: 'flex',
              gap: '16px',
              alignItems: 'center'
            }}>
              <input
                type="text"
                placeholder="搜索数据集..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                style={{
                  flex: 1,
                  padding: '8px 12px',
                  border: '1px solid #d1d5db',
                  borderRadius: '6px',
                  fontSize: '14px',
                  outline: 'none'
                }}
              />
              <select
                value={`${sortBy}-${sortOrder}`}
                onChange={(e) => {
                  const [field, order] = e.target.value.split('-') as [typeof sortBy, typeof sortOrder];
                  setSortBy(field);
                  setSortOrder(order);
                }}
                style={{
                  padding: '8px 12px',
                  border: '1px solid #d1d5db',
                  borderRadius: '6px',
                  fontSize: '14px',
                  backgroundColor: 'white'
                }}
              >
                <option value="name-asc">名称 ↑</option>
                <option value="name-desc">名称 ↓</option>
                <option value="createdAt-asc">创建时间 ↑</option>
                <option value="createdAt-desc">创建时间 ↓</option>
                <option value="itemCount-asc">数据量 ↑</option>
                <option value="itemCount-desc">数据量 ↓</option>
              </select>
            </div>

            {/* 数据集列表 */}
            <div style={{
              flex: 1,
              padding: '24px',
              overflow: 'auto'
            }}>
              {filteredAndSortedDatasets.length === 0 ? (
                <div style={{
                  textAlign: 'center',
                  padding: '48px 0',
                  color: '#6b7280'
                }}>
                  <p>暂无数据集</p>
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {filteredAndSortedDatasets.map((dataset) => (
                    <div
                      key={dataset.id}
                      onClick={() => handleDatasetSelect(dataset.id)}
                      className={`modern-card p-6 flex flex-col cursor-pointer transition-all relative ${
                        selectedDatasets.has(dataset.id) 
                          ? 'border-2 border-blue-500 bg-blue-50' 
                          : 'hover:shadow-md'
                      }`}
                    >
                      {selectedDatasets.has(dataset.id) && (
                        <div className="absolute top-2 right-2 w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center text-white text-xs">
                          ✓
                        </div>
                      )}
                      <div className="mb-3">
                        <h3 className="text-lg font-semibold text-gray-800 mb-2">{dataset.name}</h3>
                        <p className="text-gray-600 text-sm mb-3 flex-1">{dataset.description}</p>
                      </div>
                      
                      <div className="mb-4 space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-500">数据量:</span>
                          <span className="text-gray-700">{dataset.itemCount}</span>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-500">创建时间:</span>
                          <span className="text-gray-700">{dataset.createdAt}</span>
                        </div>
                      </div>
                      
                      <div className="flex gap-2 mt-auto">
                        <button 
                          onClick={(e) => {
                            e.stopPropagation();
                            // 这里可以添加查看详情的逻辑
                          }}
                          className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
                        >
                          查看详情
                        </button>
                        <button 
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDatasetSelect(dataset.id);
                          }}
                          className={`flex-1 px-3 py-2 text-sm rounded transition-colors border-0 ${
                            selectedDatasets.has(dataset.id)
                              ? 'bg-blue-500 text-white hover:bg-blue-600'
                              : 'accent-green hover:bg-accent-light'
                          }`}
                        >
                          {selectedDatasets.has(dataset.id) ? '已选择' : '选择'}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* 底部操作区域 */}
            <div style={{
              padding: '16px 24px',
              borderTop: '1px solid #e5e7eb',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              backgroundColor: '#f9fafb'
            }}>
              <div style={{
                fontSize: '14px',
                color: '#6b7280'
              }}>
                已选择 {selectedDatasets.size} 个数据集
              </div>
              <div style={{
                display: 'flex',
                gap: '12px'
              }}>
                <button
                  onClick={() => {
                    setIsDatasetModalOpen(false);
                    setSelectedDatasets(new Set());
                    setSearchTerm('');
                  }}
                  style={{
                    padding: '8px 16px',
                    backgroundColor: 'white',
                    border: '1px solid #d1d5db',
                    borderRadius: '6px',
                    cursor: 'pointer',
                    fontSize: '14px',
                    color: '#374151'
                  }}
                >
                  取消
                </button>
                <button
                  onClick={handleAddSelectedDatasets}
                  disabled={selectedDatasets.size === 0}
                  style={{
                    padding: '8px 16px',
                    backgroundColor: selectedDatasets.size > 0 ? '#3b82f6' : '#9ca3af',
                    color: 'white',
                    border: 'none',
                    borderRadius: '6px',
                    cursor: selectedDatasets.size > 0 ? 'pointer' : 'not-allowed',
                    fontSize: '14px',
                    fontWeight: '500'
                  }}
                >
                  添加选中项 ({selectedDatasets.size})
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
      
      {/* 分数上限设置弹窗 */}
      {isScoreLimitModalOpen && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 50
        }}>
          <div style={{
            backgroundColor: 'white',
            borderRadius: '8px',
            boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)',
            width: '400px',
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden'
          }}>
            {/* 弹窗头部 */}
            <div style={{
              padding: '24px 24px 16px 24px',
              borderBottom: '1px solid #dcfce7'
            }}>
              <h3 style={{
                fontSize: '18px',
                fontWeight: '600',
                color: '#111827',
                margin: 0
              }}>{isViewMode ? '查看分数上限' : '设置分数上限'}</h3>
            </div>

            {/* 弹窗内容 */}
            <div style={{
              padding: '24px'
            }}>
              <div style={{
                marginBottom: '16px'
              }}>
                <label style={{
                  display: 'block',
                  fontSize: '14px',
                  fontWeight: '500',
                  color: '#374151',
                  marginBottom: '8px'
                }}>分数上限</label>
                <input
                  type="number"
                  value={scoreLimit}
                  onChange={isViewMode ? undefined : (e) => setScoreLimit(Number(e.target.value))}
                  min="1"
                  max="1000"
                  readOnly={isViewMode}
                  style={{
                    width: '100%',
                    padding: '12px',
                    border: '1px solid #d1d5db',
                    borderRadius: '6px',
                    fontSize: '14px',
                    outline: 'none',
                    transition: 'border-color 0.2s',
                    boxSizing: 'border-box',
                    backgroundColor: isViewMode ? '#f9fafb' : 'white',
                    cursor: isViewMode ? 'default' : 'text'
                  }}
                  placeholder={isViewMode ? "当前分数上限" : "请输入分数上限"}
                />
              </div>
              

            </div>

            {/* 弹窗底部 */}
            <div style={{
              padding: '16px 24px 24px 24px',
              display: 'flex',
              justifyContent: 'flex-end',
              gap: '12px',
              borderTop: '1px solid #f3f4f6'
            }}>
              <button
                onClick={() => setIsScoreLimitModalOpen(false)}
                style={{
                  padding: '8px 16px',
                  backgroundColor: 'white',
                  border: '1px solid #d1d5db',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontSize: '14px',
                  color: '#374151'
                }}
              >
                {isViewMode ? '关闭' : '取消'}
              </button>
              {!isViewMode && (
                <button
                   onClick={handleSaveScoreLimit}
                   style={{
                     padding: '8px 16px',
                     backgroundColor: '#10b981',
                     color: 'white',
                     border: 'none',
                     borderRadius: '6px',
                     cursor: 'pointer',
                     fontSize: '14px',
                     fontWeight: '500'
                   }}
                 >
                   保存
                 </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default CriteriaDetailPage;