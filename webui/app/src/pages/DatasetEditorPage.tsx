import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import HelpModal from '../components/HelpModal';

interface DataPair {
  id: string;
  input: string;
  target: string;
}

interface LocationState {
  datasetName: string;
  datasetDescription: string;
  datasetId: string;
}

const DatasetEditorPage: React.FC = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const { datasetName, datasetDescription, datasetId } = (location.state as LocationState) || {
    datasetName: '数据集',
    datasetDescription: '数据集描述',
    datasetId: ''
  };

  // 如果没有数据集ID，返回数据集列表页面
  useEffect(() => {
    if (!datasetId) {
      navigate('/datasets');
    }
  }, [datasetId, navigate]);

  const [dataPairs, setDataPairs] = useState<DataPair[]>([
    { id: '', input: '', target: '' }
  ]);
  const [deletedPairIds, setDeletedPairIds] = useState<string[]>([]);
  const [nextId, setNextId] = useState(2);
  const [keyCounter, setKeyCounter] = useState(1);
  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editingName, setEditingName] = useState(datasetName);
  const [editingDescription, setEditingDescription] = useState(datasetDescription);
  const [currentDatasetName, setCurrentDatasetName] = useState(datasetName);
  const [currentDatasetDescription, setCurrentDatasetDescription] = useState(datasetDescription);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  // 获取数据集的详细数据
  const fetchDatasetDetails = async () => {
    if (!datasetId) return;
    
    try {
      setLoading(true);
      const response = await fetch(`http://localhost:8593/dataset_manager/${datasetId}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.code === 200) {
        const existingPairs = result.data || [];
        const editPairs: DataPair[] = existingPairs.map((pair: any) => ({
          id: pair.dataset_detail_id || '',
          input: pair.input || '',
          target: pair.target || ''
        }));
        
        // 如果没有数据，添加一个空的输入对
        if (editPairs.length === 0) {
          editPairs.push({ id: '', input: '', target: '' });
        }
        
        setDataPairs(editPairs);
        setDeletedPairIds([]);
        setNextId(editPairs.length + 1);
        setKeyCounter(editPairs.length + 1);
      }
    } catch (err) {
      console.error('获取数据集详情失败:', err);
      alert('获取数据集详情失败: ' + (err instanceof Error ? err.message : '未知错误'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDatasetDetails();
  }, [datasetId]);

  const handleOpenHelp = () => {
    setIsHelpModalOpen(true);
  };

  const handleCloseHelp = () => {
    setIsHelpModalOpen(false);
  };

  const handleOpenEditModal = () => {
    setEditingName(currentDatasetName);
    setEditingDescription(currentDatasetDescription);
    setIsEditModalOpen(true);
  };

  const handleCloseEditModal = () => {
    setIsEditModalOpen(false);
  };



  const handleInputChange = (index: number, field: 'input' | 'target', value: string) => {
    setDataPairs(prev => prev.map((pair, i) => 
      i === index ? { ...pair, [field]: value } : pair
    ));
  };

  const handleAddPair = () => {
    const newPair: DataPair = {
      id: '', // 新增数据的ID为空，由后端生成
      input: '',
      target: ''
    };
    setDataPairs(prev => [...prev, newPair]);
    setNextId(prev => prev + 1);
    setKeyCounter(prev => prev + 1);
  };

  const handleDeletePair = (index: number) => {
    if (dataPairs.length <= 1) {
      alert('至少需要保留一个数据对');
      return;
    }
    
    const pairToDelete = dataPairs[index];
    console.log('删除数据对:', pairToDelete, '索引:', index);
    
    // 如果是已存在的数据对（有ID），记录到删除列表中
    if (pairToDelete.id && pairToDelete.id !== '') {
      setDeletedPairIds(prev => [...prev, pairToDelete.id]);
    }
    
    setDataPairs(prev => prev.filter((_, i) => i !== index));
  };

  const handleSaveDataset = async () => {
    try {
      setSaving(true);
      
      // 根据后端API逻辑处理数据
      const editDataTuples: { id: string; input: string; output: string }[] = [];
      
      dataPairs.forEach(pair => {
        const hasData = pair.input.trim() || pair.target.trim();
        
        if (pair.id === '') {
          // ID为空，表示新增（只有有数据的才新增）
          if (hasData) {
            editDataTuples.push({
              id: '',
              input: pair.input,
              output: pair.target
            });
          }
        } else {
          // ID不为空的现有数据
          if (hasData) {
            // 若id不为空,且input/output有数据为修改
            editDataTuples.push({
              id: pair.id,
              input: pair.input,
              output: pair.target
            });
          } else {
            // 若id不为空,而input/output皆无数据为删除
            editDataTuples.push({
              id: pair.id,
              input: '',
              output: ''
            });
          }
        }
      });
      
      // 添加被删除的数据对
      deletedPairIds.forEach(deletedId => {
        editDataTuples.push({
          id: deletedId,
          input: '',
          output: ''
        });
      });
      
      if (editDataTuples.length === 0) {
        alert('没有需要保存的数据变更');
        return;
      }
      
      const response = await fetch('http://localhost:8593/dataset_manager/editor', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: datasetId,
          name: currentDatasetName,
          description: currentDatasetDescription,
          edit_data_tuples: editDataTuples
        })
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const result = await response.json();
      if (result.code === 200) {
        alert('数据集保存成功！');
        setDeletedPairIds([]); // 清空删除列表
        navigate('/datasets');
      } else {
        throw new Error(result.message || '保存失败');
       }
    } catch (err) {
      console.error('保存数据集失败:', err);
      alert('保存数据集失败: ' + (err instanceof Error ? err.message : '未知错误'));
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = () => {
    if (window.confirm('确定要取消编辑吗？未保存的数据将丢失。')) {
      navigate('/datasets');
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <p className="text-gray-500">加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h1 className="text-2xl font-semibold text-gray-800">
            {currentDatasetName} (编辑模式)
          </h1>
          <button
            onClick={handleOpenEditModal}
            className="px-3 py-1 text-sm bg-blue-50 text-blue-600 rounded hover:bg-blue-100 transition-colors border-0"
            aria-label="编辑数据集信息"
          >
            编辑信息
          </button>
        </div>
        <p className="text-gray-600">{currentDatasetDescription}</p>
      </div>

      <div className="mb-6">
        <div className="flex items-center gap-2 mb-4">
          <h2 className="text-lg font-medium text-gray-800">数据编辑</h2>
          <button
            onClick={handleOpenHelp}
            className="w-5 h-5 rounded-full bg-blue-500 text-white text-xs flex items-center justify-center hover:bg-blue-600 transition-colors"
            aria-label="查看使用说明"
            title="查看使用说明"
          >
            !
          </button>
        </div>
      </div>

      <div className="space-y-4 mb-6">
        {dataPairs.map((pair, index) => (
          <div key={pair.id || `new-${index}-${keyCounter}`} className="modern-card p-4">
            <div className="flex items-start gap-4">
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  输入 #{index + 1}
                </label>
                <textarea
                  value={pair.input}
                  onChange={(e) => handleInputChange(index, 'input', e.target.value)}
                  className="modern-input h-20 resize-none"
                  placeholder="输入部分"
                />
              </div>
              
              <div className="flex-1">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  目标 #{index + 1}
                </label>
                <textarea
                  value={pair.target}
                  onChange={(e) => handleInputChange(index, 'target', e.target.value)}
                  className="modern-input h-20 resize-none"
                  placeholder="目标输出"
                />
              </div>
              
              <div className="flex flex-col justify-start pt-6 gap-2">
                <button
                  onClick={() => handleDeletePair(index)}
                  className="px-2 py-1 text-xs border-0 bg-red-50 text-red-600 rounded hover:bg-red-100 transition-colors"
                  aria-label={`删除第${index + 1}个数据对`}
                  disabled={dataPairs.length <= 1}
                >
                  删除
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="mb-6">
        <button
          onClick={handleAddPair}
          className="px-4 py-2 border-2 border-dashed border-gray-300 text-gray-600 rounded hover:border-blue-400 hover:text-blue-600 transition-colors w-full"
        >
          + 添加新的数据对
        </button>
      </div>

      <div className="flex gap-3 justify-end">
        <button
          onClick={handleCancel}
          className="px-6 py-2 border-0 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
          disabled={saving}
        >
          取消
        </button>
        <button
          onClick={handleSaveDataset}
          className="modern-button px-6 py-2"
          disabled={saving}
        >
          {saving ? '保存中...' : '保存数据集'}
        </button>
      </div>
      
      <HelpModal 
        isOpen={isHelpModalOpen} 
        onClose={handleCloseHelp} 
      />
      
      {/* 编辑数据集信息弹窗 */}
      {isEditModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-auto my-auto">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">编辑数据集信息</h2>
            
            <div className="mb-4">
              <label htmlFor="edit-dataset-name" className="block text-sm font-medium text-gray-700 mb-2">
                数据集名称
              </label>
              <input
                id="edit-dataset-name"
                type="text"
                value={editingName}
                onChange={(e) => setEditingName(e.target.value)}
                className="modern-input"
                placeholder="请输入数据集名称"
                required
              />
            </div>
            
            <div className="mb-6">
              <label htmlFor="edit-dataset-description" className="block text-sm font-medium text-gray-700 mb-2">
                数据集描述
              </label>
              <textarea
                id="edit-dataset-description"
                value={editingDescription}
                onChange={(e) => setEditingDescription(e.target.value)}
                className="modern-input h-24 resize-none"
                placeholder="请输入数据集描述"
                required
              />
            </div>
            
            <div className="flex gap-3 justify-end">
              <button
                type="button"
                onClick={handleCloseEditModal}
                className="px-4 py-2 text-sm border-0 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
              >
                取消
              </button>
              <button
                type="button"
                onClick={() => {
                  setCurrentDatasetName(editingName);
                  setCurrentDatasetDescription(editingDescription);
                  handleCloseEditModal();
                }}
                className="modern-button"
                disabled={!editingName.trim()}
              >
                确认
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default DatasetEditorPage;