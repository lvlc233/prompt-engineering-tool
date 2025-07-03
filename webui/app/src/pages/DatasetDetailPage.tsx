import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import HelpModal from '../components/HelpModal';

interface DataPair {
  dataset_detail_id: string;
  dataset_id: string;
  input: string | null;
  target: string | null;
  created_at: string;
}

interface Dataset {
  dataset_id: string;
  name: string;
  description: string;
  data_count: number;
  created_at: string;
}

const DatasetDetailPage: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const [isHelpModalOpen, setIsHelpModalOpen] = useState(false);
  const [dataset, setDataset] = useState<Dataset | null>(null);
  const [dataPairs, setDataPairs] = useState<DataPair[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // 检查路由状态中是否有数据集信息
  const datasetFromState = location.state?.dataset as Dataset | undefined;

  // 获取数据集基本信息
  const fetchDatasetInfo = async () => {
    if (!id) return;
    
    try {
      const response = await fetch(`http://localhost:8593/dataset_manager/info/${id}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.code === 200) {
        setDataset(result.data);
      } else {
        throw new Error(result.message || '获取数据集信息失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取数据集信息失败');
      console.error('获取数据集信息失败:', err);
    }
  };

  // 获取数据集详细数据
  const fetchDatasetDetails = async () => {
    if (!id) return;
    
    try {
      const response = await fetch(`http://localhost:8593/dataset_manager/${id}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const result = await response.json();
      if (result.code === 200) {
        setDataPairs(result.data || []);
      } else {
        throw new Error(result.message || '获取数据集详情失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取数据集详情失败');
      console.error('获取数据集详情失败:', err);
    }
  };

  // 加载数据
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);
      
      // 如果路由状态中有数据集信息，直接使用
      if (datasetFromState) {
        setDataset(datasetFromState);
        // 只需要获取详细数据
        await fetchDatasetDetails();
      } else {
        // 否则获取完整数据
        await Promise.all([fetchDatasetInfo(), fetchDatasetDetails()]);
      }
      
      setLoading(false);
    };
    
    loadData();
  }, [id, datasetFromState]);

  if (loading) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <p className="text-gray-500">加载中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <p className="text-red-500 mb-4">{error}</p>
          <button 
            onClick={() => navigate('/datasets')}
            className="modern-button"
          >
            返回数据集列表
          </button>
        </div>
      </div>
    );
  }

  if (!dataset) {
    return (
      <div className="p-6">
        <div className="text-center py-12">
          <p className="text-gray-500">数据集不存在</p>
          <button 
            onClick={() => navigate('/datasets')}
            className="mt-4 modern-button"
          >
            返回数据集列表
          </button>
        </div>
      </div>
    );
  }

  const handleOpenHelp = () => {
    setIsHelpModalOpen(true);
  };

  const handleCloseHelp = () => {
    setIsHelpModalOpen(false);
  };

  const handleEdit = () => {
    navigate('/datasets/editor', {
      state: {
        datasetName: dataset.name,
        datasetDescription: dataset.description,
        isEditing: true,
        datasetId: dataset.dataset_id
      }
    });
  };

  const handleCancel = () => {
    navigate('/datasets');
  };

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <h1 className="text-2xl font-semibold text-gray-800">数据集详情</h1>
          <button
            onClick={handleOpenHelp}
            className="w-6 h-6 rounded-full bg-blue-100 hover:bg-blue-200 flex items-center justify-center text-blue-600 hover:text-blue-700 transition-colors text-sm font-bold"
            aria-label="查看使用说明"
            title="查看使用说明"
          >
            !
          </button>
        </div>
        <p className="text-gray-600">查看数据集的详细信息和训练数据</p>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              数据集名称
            </label>
            <div className="modern-input bg-gray-50 cursor-not-allowed">
              {dataset.name}
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              数据集大小
            </label>
            <div className="modern-input bg-gray-50 cursor-not-allowed">
              {dataset.data_count} 条
            </div>
          </div>
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            数据集描述
          </label>
          <div className="modern-input bg-gray-50 cursor-not-allowed min-h-20">
            {dataset.description}
          </div>
        </div>

        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            创建时间
          </label>
          <div className="modern-input bg-gray-50 cursor-not-allowed">
            {new Date(dataset.created_at).toLocaleDateString('zh-CN')}
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-gray-800">训练数据预览</h2>
          <span className="text-sm text-gray-500">共 {dataPairs.length} 条数据</span>
        </div>
        
        {dataPairs.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500">暂无训练数据</p>
          </div>
        ) : (
          <div className="space-y-4">
            {dataPairs.map((pair) => (
              <div key={pair.dataset_detail_id} className="border border-gray-200 rounded-lg p-4 bg-gray-50">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      输入
                    </label>
                    <div className="modern-input bg-white cursor-not-allowed">
                      {pair.input || '无输入内容'}
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      目标
                    </label>
                    <div className="modern-input bg-white cursor-not-allowed">
                      {pair.target || '无目标内容'}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="flex gap-4 mt-6">
        <button 
          onClick={handleCancel}
          className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          返回列表
        </button>
        <button 
          onClick={handleEdit}
          className="modern-button"
        >
          修改数据集
        </button>
      </div>

      <HelpModal
        isOpen={isHelpModalOpen}
        onClose={handleCloseHelp}
      />
    </div>
  );
};

export default DatasetDetailPage;