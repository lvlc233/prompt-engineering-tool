import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

// 添加CSS动画样式
const slideInAnimation = `
  @keyframes slideInFromRight {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }
`;

// 将样式注入到页面
if (typeof document !== 'undefined') {
  const styleElement = document.createElement('style');
  styleElement.textContent = slideInAnimation;
  if (!document.head.querySelector('style[data-slide-animation]')) {
    styleElement.setAttribute('data-slide-animation', 'true');
    document.head.appendChild(styleElement);
  }
}

interface TaskEditorState {
  taskName?: string;
  taskDescription?: string;
  isEditing?: boolean;
  taskId?: number;
}

interface DataRow {
  id: number;
  input: string;
  output: string;
}

interface EvaluationSet {
  id: string;
  name: string;
  description?: string;
  score: number;
  totalScore: number;
  evaluationMethod: 'LLM评估' | '人类评估' | '';
  createdAt?: string;
}

const TaskEditorPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const state = location.state as TaskEditorState;
  
  const [taskName, setTaskName] = useState(state?.taskName || '');
  const [taskDescription, setTaskDescription] = useState(state?.taskDescription || '');
  const [inputPrompt, setInputPrompt] = useState('');
  const [outputContent, setOutputContent] = useState('');
  const [dataRows, setDataRows] = useState<DataRow[]>([
    { id: 1, input: '', output: '' },
    { id: 2, input: '', output: '' }
  ]);
  const [nextId, setNextId] = useState(3);
  const [evaluationSets, setEvaluationSets] = useState<EvaluationSet[]>([]);
  const [isEvaluationModalOpen, setIsEvaluationModalOpen] = useState(false);
  const [availableEvaluationSets, setAvailableEvaluationSets] = useState<EvaluationSet[]>([]);
  const [selectedEvaluationSets, setSelectedEvaluationSets] = useState<Set<string>>(new Set());
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'createdAt' | 'totalScore'>('name');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

  const handleBack = () => {
    navigate('/tasks');
  };

  const handleSave = () => {
    // 这里应该调用保存API
    console.log('保存任务:', {
      name: taskName,
      description: taskDescription,
      inputPrompt,
      outputContent,
      dataRows
    });
    navigate('/tasks');
  };

  const handleAddRow = () => {
    setDataRows([...dataRows, { id: nextId, input: '', output: '' }]);
    setNextId(nextId + 1);
  };

  const handleDeleteRow = (id: number) => {
    if (dataRows.length > 1) {
      setDataRows(dataRows.filter(row => row.id !== id));
    }
  };

  const handleRowChange = (id: number, field: 'input' | 'output', value: string) => {
    setDataRows(dataRows.map(row => 
      row.id === id ? { ...row, [field]: value } : row
    ));
  };

  const handleRunPrompt = () => {
    // 模拟运行提示词
    if (inputPrompt.trim()) {
      setOutputContent(`处理结果：\n\n基于输入提示词"${inputPrompt}"的处理结果。\n\n这里会显示AI模型的实际输出内容。`);
    }
  };

  // 获取可用的评测集
  const fetchEvaluationSets = async () => {
    try {
      const response = await fetch('http://localhost:3000/evaluation-sets');
      const data = await response.json();
      setAvailableEvaluationSets(data);
    } catch (error) {
      console.error('获取评测集失败:', error);
      // 模拟数据作为备选
      setAvailableEvaluationSets([
        { id: '1', name: '评测集1', description: '这是第一个评测集', score: 0, totalScore: 100, evaluationMethod: '', createdAt: '2024-01-01' },
        { id: '2', name: '评测集2', description: '这是第二个评测集', score: 0, totalScore: 80, evaluationMethod: '', createdAt: '2024-01-02' },
        { id: '3', name: '评测集3', description: '这是第三个评测集', score: 0, totalScore: 120, evaluationMethod: '', createdAt: '2024-01-03' }
      ]);
    }
  };

  // 过滤和排序评测集
  const filteredAndSortedEvaluationSets = availableEvaluationSets
    .filter(evaluationSet => 
      evaluationSet.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      (evaluationSet.description && evaluationSet.description.toLowerCase().includes(searchTerm.toLowerCase()))
    )
    .sort((a, b) => {
      const aValue = a[sortBy];
      const bValue = b[sortBy];
      
      // 处理 undefined 值
      if (aValue === undefined && bValue === undefined) return 0;
      if (aValue === undefined) return sortOrder === 'asc' ? 1 : -1;
      if (bValue === undefined) return sortOrder === 'asc' ? -1 : 1;
      
      const comparison = aValue < bValue ? -1 : aValue > bValue ? 1 : 0;
      return sortOrder === 'asc' ? comparison : -comparison;
    });

  // 选择评测集
  const handleEvaluationSetSelect = (setId: string) => {
    const newSelected = new Set(selectedEvaluationSets);
    if (newSelected.has(setId)) {
      newSelected.delete(setId);
    } else {
      newSelected.add(setId);
    }
    setSelectedEvaluationSets(newSelected);
  };

  // 添加选中的评测集
  const handleAddSelectedEvaluationSets = () => {
    const selectedSets = availableEvaluationSets.filter(set => selectedEvaluationSets.has(set.id));
    selectedSets.forEach(set => {
      const newEvaluationSet: EvaluationSet = {
        id: set.id,
        name: set.name,
        score: 0,
        totalScore: 100,
        evaluationMethod: ''
      };
      setEvaluationSets(prev => [...prev, newEvaluationSet]);
    });
    setSelectedEvaluationSets(new Set());
    setIsEvaluationModalOpen(false);
  };

  // 更新评测集
  const handleEvaluationSetChange = (id: string, field: keyof EvaluationSet, value: any) => {
    setEvaluationSets(prev => prev.map(set => 
      set.id === id ? { ...set, [field]: value } : set
    ));
  };

  // 删除评测集
  const handleDeleteEvaluationSet = (id: string) => {
    setEvaluationSets(prev => prev.filter(set => set.id !== id));
  };

  useEffect(() => {
    if (isEvaluationModalOpen) {
      fetchEvaluationSets();
    }
  }, [isEvaluationModalOpen]);

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      {/* 页面头部 */}
      <div className="mb-6">
        <div className="flex items-center gap-4 mb-4">
          <button
            onClick={handleBack}
            className="px-4 py-2 text-sm border border-gray-300 text-gray-600 rounded hover:bg-gray-50 transition-colors"
          >
            ← 返回任务列表
          </button>
          <h1 className="text-2xl font-semibold text-gray-800">
            {state?.isEditing ? '编辑任务' : '创建新任务'}
          </h1>
        </div>
        

      </div>

      {/* 主要内容区域 - 四个区域布局 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 h-[calc(100vh-200px)]">
        {/* 上半部分：提示词输入输出 */}
        <div className="lg:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6 h-1/2">
          {/* 左上：输入提示词 */}
          <div className="modern-card p-6 flex flex-col">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800">输入提示词</h3>
              <button
                onClick={handleRunPrompt}
                className="modern-button text-sm"
                disabled={!inputPrompt.trim()}
              >
                运行
              </button>
            </div>
            <textarea
              value={inputPrompt}
              onChange={(e) => setInputPrompt(e.target.value)}
              className="modern-input flex-1 resize-none"
              placeholder="请输入您的提示词..."
            />
          </div>

          {/* 右上：输出内容 */}
          <div className="modern-card p-6 flex flex-col">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">输出内容</h3>
            <div className="flex-1 bg-gray-50 rounded border p-4 overflow-auto">
              <pre className="whitespace-pre-wrap text-sm text-gray-700">
                {outputContent || '运行提示词后，输出内容将显示在这里...'}
              </pre>
            </div>
          </div>
        </div>

        {/* 下半部分：数据相关 - 左右布局 */}
        <div className="lg:col-span-2 flex gap-6 h-1/2">
          {/* 左侧：评测集 */}
          <div className="flex-1 modern-card p-6 flex flex-col">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800">评测集</h3>
              <button
                onClick={() => setIsEvaluationModalOpen(true)}
                className="px-3 py-1 text-sm border border-gray-300 text-gray-600 rounded hover:bg-gray-50 transition-colors"
              >
                + 选择评测集
              </button>
            </div>
            
            <div className="flex-1 overflow-auto">
              {evaluationSets.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <p>暂无评测集</p>
                  <p className="text-sm mt-2">点击"选择评测集"按钮添加</p>
                </div>
              ) : (
                <div className="border border-gray-200 rounded-lg overflow-hidden">
                  {/* 表格 */}
                  <table className="w-full">
                     <thead>
                       <tr className="bg-gray-50 border-b border-gray-200">
                         <th className="px-4 py-3 text-left text-sm font-medium text-gray-700 w-1/2">评测集名</th>
                          <th className="px-1 py-3 text-left text-sm font-medium text-gray-700 w-12">得分</th>
                          <th className="px-1 py-3 text-left text-sm font-medium text-gray-700 w-12">总分</th>
                          <th className="px-4 py-3 text-left text-sm font-medium text-gray-700">评测方法</th>
                          <th className="px-4 py-3 text-left text-sm font-medium text-gray-700 w-20">操作</th>
                       </tr>
                     </thead>
                     <tbody className="divide-y divide-gray-200">
                       {evaluationSets.map((set) => (
                         <tr key={set.id} className="hover:bg-gray-50">
                           <td className="px-4 py-3">
                             <span className="text-sm text-gray-900">{set.name}</span>
                           </td>
                           <td className="px-1 py-3">
                              <input
                                type="number"
                                value={set.score}
                                onChange={(e) => handleEvaluationSetChange(set.id, 'score', Number(e.target.value))}
                                className="w-10 px-1 py-1 text-xs border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                disabled={set.evaluationMethod !== '人类评估'}
                                min="0"
                                max={set.totalScore}
                              />
                            </td>
                            <td className="px-1 py-3">
                              <input
                                type="number"
                                value={set.totalScore}
                                onChange={(e) => handleEvaluationSetChange(set.id, 'totalScore', Number(e.target.value))}
                                className="w-10 px-1 py-1 text-xs border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                                min="1"
                              />
                            </td>
                           <td className="px-4 py-3">
                             <select
                               value={set.evaluationMethod}
                               onChange={(e) => handleEvaluationSetChange(set.id, 'evaluationMethod', e.target.value as 'LLM评估' | '人类评估')}
                               className="px-2 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                             >
                               <option value="">请选择</option>
                               <option value="LLM评估">LLM评估</option>
                               <option value="人类评估">人类评估</option>
                             </select>
                           </td>
                           <td className="px-4 py-3">
                             <button
                               onClick={() => handleDeleteEvaluationSet(set.id)}
                               className="px-3 py-1 text-xs border border-red-300 text-red-600 rounded hover:bg-red-50 transition-colors"
                             >
                               删除
                             </button>
                           </td>
                         </tr>
                       ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          </div>

          {/* 右侧：数据可视化 */}
          <div className="flex-1 modern-card p-6 flex flex-col">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">数据可视化</h3>
            <div className="flex-1 bg-gray-50 rounded border p-4">
              {evaluationSets.length > 0 ? (
                <ResponsiveContainer width="100%" height={Math.max(300, evaluationSets.length * 60 + 150)}>
                  <BarChart
                    data={evaluationSets.map(set => ({
                      name: set.name,
                      得分: set.score,
                      满分: set.totalScore
                    }))}
                    margin={{
                      top: 20,
                      right: 30,
                      left: 20,
                      bottom: 80,
                    }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="name" 
                      tick={{ fontSize: 12 }}
                      interval={0}
                      angle={-45}
                      textAnchor="end"
                      height={80}
                    />
                    <YAxis 
                      label={{ value: '分数', angle: -90, position: 'insideLeft' }}
                      tick={{ fontSize: 12 }}
                    />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="得分" fill="#3b82f6" name="得分" />
                    <Bar dataKey="满分" fill="#e5e7eb" name="满分" />
                  </BarChart>
                </ResponsiveContainer>
              ) : (
                <div className="flex items-center justify-center h-full">
                  <div className="text-center text-gray-500">
                    <div className="w-16 h-16 mx-auto mb-4 bg-gray-200 rounded-lg flex items-center justify-center">
                      <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                      </svg>
                    </div>
                    <p className="text-sm">数据可视化图表</p>
                    <p className="text-xs mt-1">添加评测集后将显示相关图表</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* 优化方向和优化后提示词 */}
      <div className="mt-6 modern-card p-6">
        <div className="flex flex-col lg:flex-row gap-6 w-full">
          {/* 左侧：优化方向 */}
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">优化方向</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  输入优化方向
                </label>
                <textarea
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                  rows={6}
                  placeholder="请输入您希望优化的方向和目标..."
                />
              </div>
              <button className="modern-button w-full">
                运行优化
              </button>
            </div>
          </div>

          {/* 右侧：优化后提示词 */}
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">优化后提示词</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  优化后的提示词
                </label>
                <div className="bg-gray-50 border border-gray-200 rounded-md p-3 min-h-[200px]">
                  <p className="text-gray-500 text-sm">优化完成后将显示改进的提示词...</p>
                </div>
              </div>
              <div className="flex gap-3">
                <button className="flex-1 px-4 py-2 border border-gray-300 text-gray-600 rounded hover:bg-gray-50 transition-colors">
                  复制提示词
                </button>
                <button className="flex-1 modern-button">
                  迭代
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 评测集选择模态框 */}
      {isEvaluationModalOpen && (
        <div 
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            display: 'flex',
            alignItems: 'stretch',
            justifyContent: 'flex-end',
            zIndex: 50
          }}
          onClick={(e) => {
             if (e.target === e.currentTarget) {
               setIsEvaluationModalOpen(false);
               setSelectedEvaluationSets(new Set());
               setSearchTerm('');
             }
           }}
        >
          <div style={{
            backgroundColor: 'white',
            width: '60%',
            maxWidth: '800px',
            minWidth: '600px',
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            overflow: 'hidden',
            boxShadow: '-4px 0 15px rgba(0, 0, 0, 0.1)',
            animation: 'slideInFromRight 0.3s ease-out'
          }}>
            {/* 弹窗头部 */}
            <div style={{
              padding: '24px 24px 16px 24px',
              borderBottom: '1px solid #e5e7eb',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              <h3 style={{
                fontSize: '18px',
                fontWeight: '600',
                color: '#111827',
                margin: 0
              }}>选择评测集</h3>
              <button
                 onClick={() => {
                   setIsEvaluationModalOpen(false);
                   setSelectedEvaluationSets(new Set());
                   setSearchTerm('');
                 }}
                style={{
                  background: 'none',
                  border: 'none',
                  fontSize: '24px',
                  color: '#6b7280',
                  cursor: 'pointer',
                  padding: '4px',
                  borderRadius: '4px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  width: '32px',
                  height: '32px'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = '#f3f4f6';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'transparent';
                }}
              >
                ×
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
                placeholder="搜索评测集..."
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
                <option value="totalScore-asc">总分 ↑</option>
                <option value="totalScore-desc">总分 ↓</option>
              </select>
            </div>

            {/* 弹窗内容 */}
            <div style={{
              flex: 1,
              padding: '24px',
              overflow: 'auto'
            }}>
              {filteredAndSortedEvaluationSets.length === 0 ? (
                <div style={{
                  textAlign: 'center',
                  padding: '48px 0',
                  color: '#6b7280'
                }}>
                  {searchTerm ? (
                    <p>未找到匹配的评测集</p>
                  ) : (
                    <p>暂无可用的评测集</p>
                  )}
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {filteredAndSortedEvaluationSets.map((evaluationSet) => (
                    <div 
                      key={evaluationSet.id} 
                      className={`modern-card p-6 flex flex-col relative cursor-pointer hover:bg-gray-50 transition-colors ${
                        selectedEvaluationSets.has(evaluationSet.id) ? 'ring-2 ring-blue-500' : ''
                      }`}
                      onClick={() => handleEvaluationSetSelect(evaluationSet.id)}
                    >
                      {selectedEvaluationSets.has(evaluationSet.id) && (
                        <div className="absolute top-2 right-2 w-5 h-5 bg-blue-500 rounded-full flex items-center justify-center text-white text-xs">
                          ✓
                        </div>
                      )}
                      <div className="mb-3">
                        <h3 className="text-lg font-semibold text-gray-800 mb-2">{evaluationSet.name}</h3>
                        <p className="text-gray-600 text-sm mb-3 flex-1">{evaluationSet.description}</p>
                      </div>
                      
                      <div className="mb-4 space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-500">总分:</span>
                          <span className="text-gray-700">{evaluationSet.totalScore}</span>
                        </div>
                        {evaluationSet.createdAt && (
                          <div className="flex justify-between text-sm">
                            <span className="text-gray-500">创建时间:</span>
                            <span className="text-gray-700">{evaluationSet.createdAt}</span>
                          </div>
                        )}
                      </div>
                      
                      <div className="flex gap-2 mt-auto">
                        <button 
                          onClick={(e) => {
                            e.stopPropagation();
                            // 查看详情逻辑
                          }}
                          className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
                        >
                          查看详情
                        </button>
                        <button 
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEvaluationSetSelect(evaluationSet.id);
                          }}
                          className={`flex-1 px-3 py-2 text-sm rounded transition-colors border-0 ${
                            selectedEvaluationSets.has(evaluationSet.id)
                              ? 'bg-blue-500 text-white hover:bg-blue-600'
                              : 'bg-green-500 text-white hover:bg-green-600'
                          }`}
                        >
                          {selectedEvaluationSets.has(evaluationSet.id) ? '已选择' : '选择'}
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
                已选择 {selectedEvaluationSets.size} 个评测集
              </div>
              <div style={{
                display: 'flex',
                gap: '12px'
              }}>
                <button
                  onClick={() => {
                    setIsEvaluationModalOpen(false);
                    setSelectedEvaluationSets(new Set());
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
                  onClick={handleAddSelectedEvaluationSets}
                  disabled={selectedEvaluationSets.size === 0}
                  style={{
                    padding: '8px 16px',
                    backgroundColor: selectedEvaluationSets.size > 0 ? '#3b82f6' : '#9ca3af',
                    color: 'white',
                    border: 'none',
                    borderRadius: '6px',
                    cursor: selectedEvaluationSets.size > 0 ? 'pointer' : 'not-allowed',
                    fontSize: '14px',
                    fontWeight: '500'
                  }}
                >
                  添加选中项 ({selectedEvaluationSets.size})
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default TaskEditorPage;