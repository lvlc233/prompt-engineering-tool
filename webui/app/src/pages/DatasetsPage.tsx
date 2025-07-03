import React, { useState, useMemo, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AddDatasetModal from '../components/AddDatasetModal';
import HelpModal from '../components/HelpModal';

interface Dataset {
  dataset_id: string;
  name: string;
  description?: string;
  data_count: number;
  created_at: string;
}

const DatasetsPage: React.FC = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'time'>('time');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);
  const [datasets, setDatasets] = useState<Dataset[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 获取数据集列表
  const fetchDatasets = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch('http://localhost:8593/dataset_manager/');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.code === 200) {
        setDatasets(result.data || []);
      } else {
        throw new Error(result.message || '获取数据集失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取数据集失败');
      console.error('获取数据集失败:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDatasets();
  }, []);

  const handleDeleteDataset = async (dataset_id: string) => {
    if (window.confirm('确定要删除这个数据集吗？此操作不可撤销。')) {
      try {
        const response = await fetch(`http://localhost:8593/dataset_manager/${dataset_id}`, {
          method: 'DELETE',
        });
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.code === 200) {
          // 删除成功，重新获取数据集列表
          await fetchDatasets();
        } else {
          throw new Error(result.message || '删除数据集失败');
        }
      } catch (err) {
        console.error('删除数据集失败:', err);
        alert('删除数据集失败: ' + (err instanceof Error ? err.message : '未知错误'));
      }
    }
  };

  const handleViewDetail = (dataset: Dataset) => {
    // 跳转到详情页面时，传递完整的数据集信息
    navigate(`/datasets/${dataset.dataset_id}`, {
      state: {
        dataset: dataset
      }
    });
  };

  const handleEditDataset = (dataset: Dataset) => {
    navigate('/datasets/editor', {
      state: {
        datasetName: dataset.name,
        datasetDescription: dataset.description,
        isEditing: true,
        datasetId: dataset.dataset_id
      }
    });
  };

  const handleAddDataset = () => {
    setIsModalOpen(true);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
  };

  const handleOpenHelp = () => {
    setIsHelpModalOpen(true);
  };

  const handleCloseHelp = () => {
    setIsHelpModalOpen(false);
  };

  const handleModalConfirm = async (name: string, description: string) => {
    try {
      const response = await fetch('http://localhost:8593/dataset_manager/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: name,
          description: description
        })
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.code === 200) {
        setIsModalOpen(false);
        // 创建成功，重新获取数据集列表
        await fetchDatasets();
      } else {
        throw new Error(result.message || '创建数据集失败');
      }
    } catch (err) {
      console.error('创建数据集失败:', err);
      alert('创建数据集失败: ' + (err instanceof Error ? err.message : '未知错误'));
    }
  };

  const handleSortChange = (newSortBy: 'name' | 'time') => {
    if (sortBy === newSortBy) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(newSortBy);
      setSortOrder('asc');
    }
  };

  const filteredAndSortedDatasets = useMemo(() => {
    let filtered = datasets.filter(dataset => 
      dataset.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    filtered.sort((a, b) => {
      let comparison = 0;
      if (sortBy === 'name') {
        comparison = a.name.localeCompare(b.name);
      } else {
        comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
      }
      return sortOrder === 'asc' ? comparison : -comparison;
    });

    return filtered;
  }, [datasets, searchTerm, sortBy, sortOrder]);

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <h1 className="text-2xl font-semibold text-gray-800">数据集管理</h1>
          <button
            onClick={handleOpenHelp}
            className="w-6 h-6 rounded-full bg-blue-100 hover:bg-blue-200 flex items-center justify-center text-blue-600 hover:text-blue-700 transition-colors text-sm font-bold"
            aria-label="查看使用说明"
            title="查看使用说明"
          >
            !
          </button>
        </div>
        <p className="text-gray-600">管理和查看所有训练数据集</p>
      </div>
      
      <div className="mb-6 flex flex-wrap gap-4 items-center">
        <button onClick={handleAddDataset} className="modern-button">
          + 添加新数据集
        </button>
        
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="搜索数据集名称..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="modern-input w-64"
          />
        </div>
        
        <div className="flex gap-2">
          <button
            onClick={() => handleSortChange('name')}
            className={`px-3 py-2 text-sm rounded border transition-colors ${
              sortBy === 'name' 
                ? 'bg-accent-light accent-green border-light-green' 
                : 'border-gray-300 text-gray-600 hover:bg-gray-50'
            }`}
          >
            按名称排序 {sortBy === 'name' && (sortOrder === 'asc' ? '↑' : '↓')}
          </button>
          <button
            onClick={() => handleSortChange('time')}
            className={`px-3 py-2 text-sm rounded border transition-colors ${
              sortBy === 'time' 
                ? 'bg-accent-light accent-green border-light-green' 
                : 'border-gray-300 text-gray-600 hover:bg-gray-50'
            }`}
          >
            按时间排序 {sortBy === 'time' && (sortOrder === 'asc' ? '↑' : '↓')}
          </button>
        </div>
      </div>

      {loading && (
        <div className="text-center py-12">
          <p className="text-gray-500">加载中...</p>
        </div>
      )}
      
      {error && (
        <div className="text-center py-12">
          <p className="text-red-500">错误: {error}</p>
          <button 
            onClick={fetchDatasets}
            className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
          >
            重试
          </button>
        </div>
      )}
      
      {!loading && !error && (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredAndSortedDatasets.map((dataset) => (
              <div key={dataset.dataset_id} className="modern-card p-6 flex flex-col">
                <div className="mb-3">
                  <h3 className="text-lg font-semibold text-gray-800 mb-2">{dataset.name}</h3>
                  <p className="text-gray-600 text-sm mb-3 flex-1">{dataset.description || '暂无描述'}</p>
                </div>
                
                <div className="mb-4 space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">数据量:</span>
                    <span className="text-gray-700">{dataset.data_count} 条</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">创建时间:</span>
                    <span className="text-gray-700">{new Date(dataset.created_at).toLocaleDateString()}</span>
                  </div>
                </div>
                
                <div className="flex gap-2 mt-auto">
                  <button 
                    onClick={() => handleViewDetail(dataset)}
                    className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
                  >
                    查看详情
                  </button>
                  <button 
                    onClick={() => handleEditDataset(dataset)}
                    className="flex-1 px-3 py-2 text-sm accent-green hover:bg-accent-light rounded transition-colors border-0"
                  >
                    编辑
                  </button>
                  <button 
                    onClick={() => handleDeleteDataset(dataset.dataset_id)}
                    className="px-3 py-2 text-sm border border-red-300 text-red-600 rounded hover:bg-red-50 transition-colors"
                  >
                    删除
                  </button>
                </div>
              </div>
            ))}
          </div>
          
          {filteredAndSortedDatasets.length === 0 && datasets.length > 0 && (
            <div className="text-center py-12">
              <p className="text-gray-500">没有找到匹配的数据集</p>
            </div>
          )}
          
          {datasets.length === 0 && (
            <div className="text-center py-12">
              <p className="text-gray-500">暂无数据集，点击上方按钮添加新数据集</p>
            </div>
          )}
        </>
      )}
      
      <AddDatasetModal
        isOpen={isModalOpen}
        onClose={handleModalClose}
        onConfirm={handleModalConfirm}
      />
      
      <HelpModal
        isOpen={isHelpModalOpen}
        onClose={handleCloseHelp}
      />
    </div>
  );
};

export default DatasetsPage;