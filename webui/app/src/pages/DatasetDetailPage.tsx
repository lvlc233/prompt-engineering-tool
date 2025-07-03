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
  description?: string;
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
  
  // 移除对路由状态的依赖，只使用API获取数据

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
      
      // 始终从API获取完整数据
      await Promise.all([fetchDatasetInfo(), fetchDatasetDetails()]);
      
      setLoading(false);
    };
    
    loadData();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

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

  const handleGoBack = () => {
    // 如果有来源页面信息，返回到来源页面，否则返回上一页
    if (location.state?.from) {
      navigate(location.state.from);
    } else {
      // 使用浏览器历史记录返回
      window.history.back();
    }
  };

  return (
    <div className="p-6 relative">
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
           <div className="flex items-center gap-2">
             <h1 className="text-2xl font-semibold text-gray-800">数据集详情</h1>
             <span className="text-sm text-gray-500">({dataset?.name})</span>
             <button
               onClick={handleOpenHelp}
               className="w-6 h-6 rounded-full bg-blue-100 hover:bg-blue-200 flex items-center justify-center text-blue-600 hover:text-blue-700 transition-colors text-sm font-bold"
               aria-label="查看使用说明"
               title="查看使用说明"
             >
               !
             </button>
           </div>
          <button 
             onClick={handleGoBack}
             className="inline-flex items-center gap-2 px-6 py-2 border-0 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
           >
             <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
               <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
             </svg>
             返回
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
            {dataset.description || '暂无描述'}
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

      <div className="bg-white rounded-lg border border-gray-200 p-6 mb-6">
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

      {/* 修改数据集按钮 */}
      <div className="flex justify-end">
        <button 
          onClick={handleEdit}
          className="modern-button px-6 py-2"
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