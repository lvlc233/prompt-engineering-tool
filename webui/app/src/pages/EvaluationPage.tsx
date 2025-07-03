import React, { useState, useMemo, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AddEvaluationModal from '../components/AddEvaluationModal';

interface Evaluation {
  evaluationset_id: string;
  name: string;
  evaluation_criteria?: string;
  sorce_cap?: number;
  description?: string;
  created_at: string;
}

const EvaluationPage: React.FC = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'time'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const [evaluations, setEvaluations] = useState<Evaluation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 获取评测集数据
  const fetchEvaluations = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await fetch('http://localhost:8593/evaluationset_manager/');
      if (!response.ok) {
        throw new Error(`HTTP错误: ${response.status}`);
      }
      
      const text = await response.text();
      if (!text.trim()) {
        setEvaluations([]);
        return;
      }
      
      try {
        const response = JSON.parse(text);
        // 检查响应格式，如果有 data 字段则使用 data，否则直接使用响应
        const data = response.data || response;
        setEvaluations(Array.isArray(data) ? data : []);
      } catch (jsonError) {
        console.error('JSON解析错误:', jsonError);
        console.error('响应内容:', text);
        throw new Error('服务器返回的数据格式错误');
      }
    } catch (err) {
      console.error('获取评测集失败:', err);
      setError(err instanceof Error ? err.message : '获取评测集失败');
      setEvaluations([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvaluations();
  }, []);



  const handleSortChange = (newSortBy: 'name' | 'time') => {
    if (sortBy === newSortBy) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(newSortBy);
      setSortOrder('asc');
    }
  };

  const filteredAndSortedEvaluations = useMemo(() => {
    let filtered = evaluations.filter(evaluation => 
      evaluation.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (evaluation.evaluation_criteria && evaluation.evaluation_criteria.toLowerCase().includes(searchTerm.toLowerCase())) ||
      (evaluation.description && evaluation.description.toLowerCase().includes(searchTerm.toLowerCase()))
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
  }, [evaluations, searchTerm, sortBy, sortOrder]);

  const handleModalClose = () => {
    setIsModalOpen(false);
  };

  const handleDeleteEvaluation = async (id: string) => {
    if (window.confirm('确定要删除这个评测集吗？此操作不可撤销。')) {
      try {
        const response = await fetch(`http://localhost:8593/evaluationset_manager/${id}`, {
          method: 'DELETE',
        });
        if (response.ok) {
          fetchEvaluations();
        } else {
          alert('删除失败，请重试');
        }
      } catch (error) {
        console.error('删除评测集失败:', error);
        alert('删除失败，请重试');
      }
    }
  };

  const handleModalConfirm = async (name: string, description: string, scoreLimit: number) => {
    try {
      const response = await fetch('http://localhost:8593/evaluationset_manager/add', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name,
          description,
          sorce_cap: scoreLimit,
        }),
      });
      if (response.ok) {
        setIsModalOpen(false);
        await fetchEvaluations(); // 等待数据刷新完成
        // 不再自动跳转到详情页
      } else {
        alert('创建失败，请重试');
      }
    } catch (error) {
      console.error('创建评测集失败:', error);
      alert('创建失败，请重试');
    }
  };

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-gray-800 mb-2">评测集管理</h1>
        <p className="text-gray-600">配置和管理模型评测的各项标准</p>
      </div>
      
      <div className="mb-6 flex flex-wrap gap-4 items-center">
        <button 
          onClick={() => setIsModalOpen(true)}
          className="modern-button"
        >
          + 创建新测评集
        </button>
        
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="搜索评测名称..."
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

      {loading ? (
        <div className="text-center py-12">
          <p className="text-gray-500">加载中...</p>
        </div>
      ) : error ? (
        <div className="text-center py-12">
          <p className="text-red-500">{error}</p>
          <button 
            onClick={fetchEvaluations}
            className="mt-3 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
          >
            重试
          </button>
        </div>
      ) : null}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredAndSortedEvaluations.map((evaluation) => (
          <div key={evaluation.evaluationset_id} className="modern-card p-6 flex flex-col">
            <div className="mb-3">
              <h3 className="text-lg font-semibold text-gray-800 mb-2">{evaluation.name}</h3>
              <p className="text-gray-600 text-sm mb-3 flex-1">{evaluation.description}</p>
            </div>
            
            <div className="mb-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">分数上限:</span>
                <span className="text-gray-700">{evaluation.sorce_cap || 'N/A'}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">创建时间:</span>
                <span className="text-gray-700">{new Date(evaluation.created_at).toLocaleDateString()}</span>
              </div>
            </div>
            
            <div className="flex gap-2 mt-auto">
              <button 
                onClick={() => navigate(`/evaluation/${evaluation.evaluationset_id}?mode=view`)}
                className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
              >
                查看详情
              </button>
              <button 
                onClick={() => navigate(`/evaluation/${evaluation.evaluationset_id}`)}
                className="flex-1 px-3 py-2 text-sm accent-green hover:bg-accent-light rounded transition-colors border-0"
              >
                编辑
              </button>
              <button 
                onClick={() => handleDeleteEvaluation(evaluation.evaluationset_id)}
                className="px-3 py-2 text-sm border border-red-300 text-red-600 rounded hover:bg-red-50 transition-colors"
              >
                删除
              </button>
            </div>
          </div>
        ))}
      </div>
      
      {filteredAndSortedEvaluations.length === 0 && !loading && (
         <div className="text-center py-12">
           <p className="text-gray-500">暂无评估集，点击上方按钮添加新评估集</p>
         </div>
       )}
       
       <AddEvaluationModal
         isOpen={isModalOpen}
         onClose={handleModalClose}
         onConfirm={handleModalConfirm}
       />
     </div>
   );
 };

export default EvaluationPage;