import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import AddCriteriaModal from '../components/AddCriteriaModal';
import { getTimestampForSorting } from '../utils/timeUtils';

interface Criterion {
  id: number;
  name: string;
  evaluationCriteria: string;
  maxScore: number;
  createdAt: string; // 完整的时间信息，用于排序
  displayDate: string; // 格式化的日期，用于显示
}

const CriteriaPage: React.FC = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'time'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const criteria: Criterion[] = useMemo(() => [
    { 
      id: 1, 
      name: '准确性评测', 
      evaluationCriteria: '评估模型回答的准确性和事实正确性，包括信息的正确性、逻辑的合理性以及与问题的匹配度', 
      maxScore: 100,
      createdAt: '2024-01-15T10:30:00Z',
      displayDate: '2024-01-15'
    },
    { 
      id: 2, 
      name: '相关性评测', 
      evaluationCriteria: '评估回答与问题的相关程度，检查回答是否直接回应了用户的问题', 
      maxScore: 80,
      createdAt: '2024-01-10T14:20:00Z',
      displayDate: '2024-01-10'
    },
    { 
      id: 3, 
      name: '流畅性评测', 
      evaluationCriteria: '评估回答的语言流畅性和可读性，包括语法正确性、表达清晰度和语言自然度', 
      maxScore: 60,
      createdAt: '2024-01-20T09:15:00Z',
      displayDate: '2024-01-20'
    },
    { 
      id: 4, 
      name: '安全性评测', 
      evaluationCriteria: '评估回答的安全性和合规性，确保内容不包含有害信息或违规内容', 
      maxScore: 50,
      createdAt: '2024-01-12T16:45:00Z',
      displayDate: '2024-01-12'
    },
  ], []);

  const handleSortChange = (newSortBy: 'name' | 'time') => {
    if (sortBy === newSortBy) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(newSortBy);
      setSortOrder('asc');
    }
  };

  const filteredAndSortedCriteria = useMemo(() => {
    let filtered = criteria.filter(criterion => 
      criterion.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    filtered.sort((a, b) => {
      let comparison = 0;
      if (sortBy === 'name') {
        comparison = a.name.localeCompare(b.name);
      } else {
        comparison = getTimestampForSorting(a.createdAt) - getTimestampForSorting(b.createdAt);
      }
      return sortOrder === 'asc' ? comparison : -comparison;
    });

    return filtered;
  }, [criteria, searchTerm, sortBy, sortOrder]);

  const handleModalClose = () => {
    setIsModalOpen(false);
  };

  const handleDeleteCriterion = (id: number) => {
    if (window.confirm('确定要删除这个评测集吗？此操作不可撤销。')) {
      // 这里应该调用删除API
      console.log('删除评测集:', id);
    }
  };

  const handleModalConfirm = (name: string, description: string, scoreLimit: number) => {
    // 创建新测评集并跳转到详情页面
    const newId = Date.now().toString(); // 模拟生成新ID
    console.log('创建新测评集:', { id: newId, name, description, scoreLimit });
    setIsModalOpen(false);
    navigate(`/criteria/${newId}`);
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

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredAndSortedCriteria.map((criterion) => (
          <div key={criterion.id} className="modern-card p-6 flex flex-col">
            <div className="mb-3">
              <h3 className="text-lg font-semibold text-gray-800 mb-2">{criterion.name}</h3>
              <p className="text-gray-600 text-sm mb-3 flex-1">{criterion.evaluationCriteria}</p>
            </div>
            
            <div className="mb-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">分数上限:</span>
                <span className="text-gray-700">{criterion.maxScore}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">创建时间:</span>
                <span className="text-gray-700">{criterion.displayDate}</span>
              </div>
            </div>
            
            <div className="flex gap-2 mt-auto">
              <button 
                onClick={() => navigate(`/criteria/${criterion.id}?mode=view`)}
                className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
              >
                查看详情
              </button>
              <button 
                onClick={() => navigate(`/criteria/${criterion.id}`)}
                className="flex-1 px-3 py-2 text-sm accent-green hover:bg-accent-light rounded transition-colors border-0"
              >
                编辑
              </button>
              <button 
                onClick={() => handleDeleteCriterion(criterion.id)}
                className="px-3 py-2 text-sm border border-red-300 text-red-600 rounded hover:bg-red-50 transition-colors"
              >
                删除
              </button>
            </div>
          </div>
        ))}
      </div>
      
      {filteredAndSortedCriteria.length === 0 && (
         <div className="text-center py-12">
           <p className="text-gray-500">没有找到匹配的评测标准</p>
         </div>
       )}
       
       <AddCriteriaModal
         isOpen={isModalOpen}
         onClose={handleModalClose}
         onConfirm={handleModalConfirm}
       />
     </div>
   );
 };

export default CriteriaPage;