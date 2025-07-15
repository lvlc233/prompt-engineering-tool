import React, { useState, useMemo, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import AddTaskModal from '../components/AddTaskModal';
import { getTimestampForSorting } from '../utils/timeUtils';

// Git风格的版本节点数据结构
interface GitNode {
  id: string;
  message: string;
  author?: string;
  timestamp: string;
  parentIds: string[];
  branch: string;
  status: 'selected' | 'executed' | 'unexecuted';
  isExecute: boolean;
  executeDate: string | null;
}

interface Task {
  id: string;
  name: string;
  description: string;
  createdAt: string; // 完整的时间信息，用于排序
  displayDate: string; // 格式化的日期，用于显示
  versions: GitNode[];
  currentVersionId: string;
}

// Git风格的版本树组件
const GitVersionTree: React.FC<{ 
  nodes: GitNode[]; 
  currentVersionId: string;
  onVersionSelect?: (versionId: string) => void;
  taskId: string;
  onExpandVersions?: (taskId: string) => void;
}> = ({ nodes, currentVersionId, onVersionSelect, taskId, onExpandVersions }) => {
  const [selectedNode, setSelectedNode] = React.useState<string | null>(null);
  const [hoveredNode, setHoveredNode] = React.useState<string | null>(null);
  const [lastHoveredNode, setLastHoveredNode] = React.useState<string | null>(null);
  const [isCollapsed, setIsCollapsed] = React.useState<boolean>(true);

  // 自动选中当前版本
  React.useEffect(() => {
    if (currentVersionId && !selectedNode) {
      setSelectedNode(currentVersionId);
    }
  }, [currentVersionId, selectedNode]);

  // 计算节点布局
  const calculateLayout = () => {
    const nodeMap = new Map<string, GitNode>();
    const branches = new Map<string, number>();
    const positions = new Map<string, { x: number; y: number; branchIndex: number }>();
    
    // 建立节点映射
    nodes.forEach(node => nodeMap.set(node.id, node));
    
    // 按时间排序节点
    const sortedNodes = [...nodes].sort((a, b) => 
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
    
    let branchCounter = 0;
    const nodeSpacing = 80;
    const branchSpacing = 40;
    
    sortedNodes.forEach((node, index) => {
      // 为新分支分配索引
      if (!branches.has(node.branch)) {
        branches.set(node.branch, branchCounter++);
      }
      
      const branchIndex = branches.get(node.branch)!;
      const x = branchIndex * branchSpacing + 30;
      const y = index * nodeSpacing + 40;
      
      positions.set(node.id, { x, y, branchIndex });
    });
    
    return { positions, branches, nodeMap };
  };

  const { positions, branches, nodeMap } = calculateLayout();
  
  // 获取节点颜色
  const getNodeColor = (node: GitNode) => {
    if (node.id === currentVersionId) {
      return { fill: '#93c5fd', stroke: '#60a5fa', text: '#1f2937' }; // 浅蓝色 - 当前版本
    }
    
    switch (node.status) {
      case 'executed':
        return { fill: '#86efac', stroke: '#4ade80', text: '#1f2937' }; // 浅绿色 - 已执行
      case 'unexecuted':
        return { fill: '#d1d5db', stroke: '#9ca3af', text: '#1f2937' }; // 浅灰色 - 未执行
      case 'selected':
        return { fill: '#93c5fd', stroke: '#60a5fa', text: '#1f2937' }; // 浅蓝色 - 选中
      default:
        return { fill: '#d1d5db', stroke: '#9ca3af', text: '#1f2937' }; // 浅灰色 - 默认
    }
  };

  // 获取分支颜色
  const getBranchColor = (branchName: string) => {
    const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4'];
    const branchIndex = Array.from(branches.keys()).indexOf(branchName);
    return colors[branchIndex % colors.length];
  };

  // 处理节点点击
  const handleNodeClick = (nodeId: string) => {
    setSelectedNode(nodeId);
    setLastHoveredNode(null); // 点击时清除悬停状态
    if (onVersionSelect) {
      onVersionSelect(nodeId);
    }
  };

  // 处理鼠标进入节点
  const handleNodeMouseEnter = (nodeId: string) => {
    setHoveredNode(nodeId);
    setLastHoveredNode(nodeId);
  };

  // 处理鼠标离开节点
  const handleNodeMouseLeave = () => {
    setHoveredNode(null);
    // 不立即清除lastHoveredNode，保持详情面板显示
  };

  // 计算SVG尺寸
  const positionValues = Array.from(positions.values());
  const maxX = positionValues.length > 0 ? Math.max(...positionValues.map(p => p.x)) + 60 : 100;
  const maxY = positionValues.length > 0 ? Math.max(...positionValues.map(p => p.y)) + 60 : 100;

  // 处理点击外部区域
  const handleContainerClick = (e: React.MouseEvent) => {
    // 如果点击的不是节点，清除悬停状态
    if (e.target === e.currentTarget) {
      setLastHoveredNode(null);
    }
  };

  return (
    <div className="git-version-tree bg-white border border-gray-200 rounded-lg" onClick={handleContainerClick}>
      {/* 版本标题区域 - 可点击展开收起 */}
      <div 
        className="flex items-center gap-2 pb-2 cursor-pointer hover:bg-gray-50 transition-colors duration-200"
        onClick={() => {
          const newCollapsed = !isCollapsed;
          setIsCollapsed(newCollapsed);
          // 如果展开且没有版本数据，则获取版本数据
          if (!newCollapsed && nodes.length === 0 && onExpandVersions && taskId) {
            onExpandVersions(taskId);
          }
        }}
      >
        <h4 className="text-sm font-medium text-gray-900">版本</h4>
        {(selectedNode || currentVersionId) && (
          <span className="text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded">
            当前: {selectedNode || currentVersionId}
          </span>
        )}
      </div>
      
      {/* Git图形区域 - 不可点击展开收起 */}
      {!isCollapsed && (
        <div className="relative overflow-x-auto px-4 pb-4">
          <svg width={maxX} height={maxY} className="border border-gray-100 rounded">
          {/* 绘制连接线 */}
          {nodes.map(node => {
            const nodePos = positions.get(node.id);
            if (!nodePos) return null;
            
            return node.parentIds.map((parentId, index) => {
              const parentPos = positions.get(parentId);
              if (!parentPos) return null;
              
              const isSameBranch = nodeMap.get(parentId)?.branch === node.branch;
              const strokeColor = isSameBranch ? getBranchColor(node.branch) : '#d1d5db';
              
              return (
                <line
                  key={`${node.id}-${parentId}-${index}`}
                  x1={parentPos.x}
                  y1={parentPos.y}
                  x2={nodePos.x}
                  y2={nodePos.y}
                  stroke={strokeColor}
                  strokeWidth="2"
                  className="transition-all duration-200"
                />
              );
            });
          }).flat()}
          
          {/* 绘制节点 */}
          {nodes.map(node => {
            const pos = positions.get(node.id);
            if (!pos) return null;
            
            const colors = getNodeColor(node);
            const isSelected = selectedNode === node.id;
            const isHovered = hoveredNode === node.id;
            const radius = isSelected || isHovered ? 12 : 10;
            
            return (
              <g key={node.id}>
                {/* 节点圆圈 */}
                <circle
                  cx={pos.x}
                  cy={pos.y}
                  r={radius}
                  fill={colors.fill}
                  stroke={colors.stroke}
                  strokeWidth={isSelected ? "3" : "2"}
                  className="cursor-pointer transition-all duration-200"
                  onClick={() => handleNodeClick(node.id)}
                  onMouseEnter={() => handleNodeMouseEnter(node.id)}
                  onMouseLeave={handleNodeMouseLeave}
                />
                
                {/* 当前版本标识 */}
                {node.id === currentVersionId && (
                  <circle
                    cx={pos.x}
                    cy={pos.y}
                    r={6}
                    fill="white"
                    className="pointer-events-none"
                  />
                )}
                
                {/* 节点标签 */}
                <text
                  x={pos.x + 20}
                  y={pos.y - 5}
                  className="text-xs font-medium fill-gray-900 pointer-events-none"
                >
                  {node.id.substring(0, 8)}
                </text>
                <text
                  x={pos.x + 20}
                  y={pos.y + 8}
                  className="text-xs fill-gray-600 pointer-events-none"
                >
                  {node.message.length > 30 ? node.message.substring(0, 30) + '...' : node.message}
                </text>
                <text
                  x={pos.x + 20}
                  y={pos.y + 20}
                  className="text-xs fill-gray-400 pointer-events-none"
                >
                  {node.timestamp} • {node.branch}
                </text>
              </g>
            );
          })}
          </svg>
        </div>
      )}
      
      {/* 节点详情 - 右侧显示 */}
      {!isCollapsed && (hoveredNode || lastHoveredNode || selectedNode) && (
        <div className="absolute right-4 top-4 w-72 p-4 bg-white border border-gray-300 rounded-lg shadow-lg z-10">
          {(() => {
            // 优先显示悬停节点，然后是最后悬停的节点，最后是选中节点
            const displayNodeId = hoveredNode || lastHoveredNode || selectedNode;
            const node = nodeMap.get(displayNodeId!);
            if (!node) return null;
            
            const parentNodes = node.parentIds.map(id => nodeMap.get(id)).filter(Boolean);
            const isHovering = hoveredNode === displayNodeId;
            
            return (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-gray-900">版本号: {node.id.substring(0, 8)}</span>
                  {isHovering && (
                    <span className="text-xs bg-yellow-100 text-yellow-800 px-2 py-1 rounded">
                      悬停中
                    </span>
                  )}
                  {!isHovering && selectedNode === displayNodeId && (
                    <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                      已选中
                    </span>
                  )}
                </div>
                <div>
                  <span className="text-sm text-gray-700">描述: {node.message}</span>
                </div>
                <div className="text-xs text-gray-500 space-y-1">
                    <div>创建时间: {new Date(node.timestamp).toLocaleString('zh-CN')}</div>
                    <div>执行时间: {node.executeDate ? new Date(node.executeDate).toLocaleString('zh-CN') : '尚未执行'}</div>
                    <div>父版本: {parentNodes.length > 0 ? parentNodes.map(p => p!.id.substring(0, 8)).join(', ') : '无'}</div>
                  </div>
              </div>
            );
          })()}
        </div>
      )}
      
      {/* 版本选择确认 */}
      {selectedNode && selectedNode !== currentVersionId && (
        <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-lg">
          <div className="flex items-center justify-between">
            <span className="text-sm text-blue-800">
              已选择版本: {selectedNode}
            </span>
            <div className="flex gap-2">
              <button 
                onClick={() => {
                  if (onVersionSelect) {
                    onVersionSelect(selectedNode);
                  }
                  setSelectedNode(null);
                }}
                className="px-3 py-1 bg-blue-500 text-white text-xs rounded hover:bg-blue-600 transition-colors"
              >
                切换到此版本
              </button>
              <button 
                onClick={() => setSelectedNode(null)}
                className="px-3 py-1 bg-gray-500 text-white text-xs rounded hover:bg-gray-600 transition-colors"
              >
                取消
              </button>
            </div>
          </div>
        </div>
      )}
      
      {/* 颜色说明 */}
      {!isCollapsed && (
        <div className="mt-2 p-3 bg-gray-50 border border-gray-300 rounded-lg">
          <div className="flex flex-wrap gap-4 text-sm">
            <div className="flex items-center gap-2">
              <svg width="16" height="16" className="flex-shrink-0">
                <circle
                  cx="8"
                  cy="8"
                  r="6"
                  fill="#93c5fd"
                  stroke="#60a5fa"
                  strokeWidth="2"
                  className="pointer-events-none"
                />
              </svg>
              <span className="text-gray-700">选中</span>
            </div>
            <div className="flex items-center gap-2">
              <svg width="16" height="16" className="flex-shrink-0">
                <circle
                  cx="8"
                  cy="8"
                  r="6"
                  fill="#86efac"
                  stroke="#4ade80"
                  strokeWidth="2"
                  className="pointer-events-none"
                />
              </svg>
              <span className="text-gray-700">已执行</span>
            </div>
            <div className="flex items-center gap-2">
              <svg width="16" height="16" className="flex-shrink-0">
                <circle
                  cx="8"
                  cy="8"
                  r="6"
                  fill="#d1d5db"
                  stroke="#9ca3af"
                  strokeWidth="2"
                  className="pointer-events-none"
                />
              </svg>
              <span className="text-gray-700">未执行</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

const TasksPage: React.FC = () => {
  const navigate = useNavigate();
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'name' | 'time'>('time');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [isAddTaskModalOpen, setIsAddTaskModalOpen] = useState(false);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 获取所有任务
  const fetchJobs = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8593/job_manager/');
      if (!response.ok) {
        throw new Error('获取任务列表失败');
      }
      const data = await response.json();
      if (data.code === 200) {
        // 处理后端返回的数据格式
        if (data.data && Array.isArray(data.data)) {
          // 转换后端数据格式为前端格式
          const transformedTasks = data.data.map((job: any) => ({
            id: job.id || job.job_id || '', // 尝试多个可能的id字段
            name: job.name || '',
            description: job.description || '',
            createdAt: job.created_at || job.createdAt || '', // 保留完整的时间信息用于排序
            displayDate: job.created_at ? new Date(job.created_at).toISOString().split('T')[0] : '', // 用于显示的日期格式
            currentVersionId: job.selected_version || 'main', // 使用后端返回的selected_version
            versions: [] // 初始为空，通过fetchJobVersions获取真实数据
          })).filter((task: Task) => task.id); // 过滤掉没有id的任务
          setTasks(transformedTasks);
        } else {
          // 如果data为null或空数组，设置为空数组
          setTasks([]);
        }
      } else {
        throw new Error(data.message || '获取任务列表失败');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取任务列表失败');
      console.error('获取任务列表失败:', err);
    } finally {
      setLoading(false);
    }
  };

  // 获取任务版本信息
  const fetchJobVersions = async (jobId: string) => {
    try {
      const response = await fetch(`http://localhost:8593/job_manager/versions/${jobId}`);
      if (!response.ok) {
        throw new Error('获取任务版本失败');
      }
      const data = await response.json();
      
      if (data.code === 200 && data.data) {
        // 将后端数据转换为GitNode格式
        const gitNodes: GitNode[] = data.data.map((version: any) => ({
          id: version.version,
          message: version.description || '无描述',
          author: '系统',
          timestamp: version.created_at,
          parentIds: version.father_version ? [version.father_version] : [],
          branch: 'main',
          status: version.is_execute ? 'executed' : 'unexecuted',
          isExecute: version.is_execute,
          executeDate: version.exceute_date
        }));
        
        // 更新对应任务的版本数据
        setTasks(prevTasks => 
          prevTasks.map(task => 
            task.id === jobId 
              ? { ...task, versions: gitNodes }
              : task
          )
        );
      }
    } catch (err) {
      console.error('获取任务版本失败:', err);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, []);

  const handleVersionChange = (taskId: string, versionId: string) => {
    // 查找当前任务
    const currentTask = tasks.find(task => task.id === taskId);
    
    // 如果当前任务没有版本数据，则发送请求获取版本列表
    if (!currentTask || !currentTask.versions || currentTask.versions.length === 0) {
      fetchJobVersions(taskId);
    }
    
    // 更新任务的当前版本ID
    setTasks(prevTasks => 
      prevTasks.map(task => 
        task.id === taskId 
          ? { ...task, currentVersionId: versionId }
          : task
      )
    );
  };

  const handleDeleteTask = async (id: string) => {
    if (window.confirm('确定要删除这个任务吗？此操作不可撤销。')) {
      try {
        const response = await fetch(`http://localhost:8593/job_manager/${id}`, {
          method: 'DELETE'
        });
        if (!response.ok) {
          throw new Error('删除任务失败');
        }
        const data = await response.json();
        if (data.code === 200) {
          // 删除成功，从本地状态中移除
          setTasks(prevTasks => prevTasks.filter(task => task.id !== id));
        } else {
          throw new Error(data.message || '删除任务失败');
        }
      } catch (err) {
        alert(err instanceof Error ? err.message : '删除任务失败');
        console.error('删除任务失败:', err);
      }
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

  const filteredAndSortedTasks = useMemo(() => {
    let filtered = tasks.filter(task => 
      task.name.toLowerCase().includes(searchTerm.toLowerCase())
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
   }, [tasks, searchTerm, sortBy, sortOrder]);

   const handleAddTask = () => {
    setIsAddTaskModalOpen(true);
  };

  const handleEditTask = (task: Task) => {
    navigate('/jobs/editor', {
      state: {
        taskName: task.name,
        taskDescription: task.description,
        isEditing: true,
        taskId: task.id
      }
    });
  };

  const handleViewTask = (task: Task) => {
    navigate('/jobs/editor', {
      state: {
        taskName: task.name,
        taskDescription: task.description,
        isEditing: false,
        taskId: task.id,
        isViewOnly: true
      }
    });
  };

   const handleModalClose = () => {
     setIsAddTaskModalOpen(false);
   };

   const handleModalConfirm = async (name: string, description: string) => {
     try {
       const response = await fetch('http://localhost:8593/job_manager/add', {
         method: 'POST',
         headers: {
           'Content-Type': 'application/json'
         },
         body: JSON.stringify({
           name: name,
           description: description
         })
       });
       if (!response.ok) {
         throw new Error('创建任务失败');
       }
       const data = await response.json();
       if (data.code === 200) {
         setIsAddTaskModalOpen(false);
         // 重新获取任务列表
         await fetchJobs();
         navigate('/jobs/editor', {
           state: {
             taskName: name,
             taskDescription: description,
             taskId: data.data.id
           }
         });
       } else {
         throw new Error(data.message || '创建任务失败');
       }
     } catch (err) {
       alert(err instanceof Error ? err.message : '创建任务失败');
       console.error('创建任务失败:', err);
     }
   };
 
   // 加载状态
   if (loading) {
     return (
       <div className="p-6">
         <div className="mb-6">
           <h1 className="text-2xl font-semibold text-gray-800 mb-2">任务管理</h1>
           <p className="text-gray-600">管理和监控所有AI模型训练任务</p>
         </div>
         <div className="text-center py-12">
           <p className="text-gray-500">加载中...</p>
         </div>
       </div>
     );
   }

   // 错误状态
   if (error) {
     return (
       <div className="p-6">
         <div className="mb-6">
           <h1 className="text-2xl font-semibold text-gray-800 mb-2">任务管理</h1>
           <p className="text-gray-600">管理和监控所有AI模型训练任务</p>
         </div>
         <div className="text-center py-12">
           <p className="text-red-500">错误: {error}</p>
           <button 
             onClick={fetchJobs}
             className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
           >
             重试
           </button>
         </div>
       </div>
     );
   }

   return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-gray-800 mb-2">任务管理</h1>
        <p className="text-gray-600">管理和监控所有AI模型训练任务</p>
      </div>
      
      <div className="mb-6 flex flex-wrap gap-4 items-center">
        <button onClick={handleAddTask} className="modern-button">
          + 创建新任务
        </button>
        
        <div className="flex gap-2">
          <input
            type="text"
            placeholder="搜索任务名称..."
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
        {filteredAndSortedTasks.map((task) => (
          <div key={task.id} className="modern-card p-6 flex flex-col">
            <div className="mb-3">
              <h3 className="text-lg font-semibold text-gray-800 mb-2">{task.name}</h3>
              <p className="text-gray-600 text-sm mb-3 flex-1">{task.description}</p>
            </div>
            
            <div className="mb-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">创建时间:</span>
                <span className="text-gray-700">{task.displayDate}</span>
              </div>
            </div>
            
            <div className="mb-4">
              <GitVersionTree 
                  nodes={task.versions} 
                  currentVersionId={task.currentVersionId}
                  onVersionSelect={(versionId) => handleVersionChange(task.id, versionId)}
                  taskId={task.id}
                  onExpandVersions={fetchJobVersions}
                />
            </div>
            
            <div className="flex gap-2 mt-auto">
              <button 
                onClick={() => handleViewTask(task)}
                className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors"
              >
                查看详情
              </button>
              <button 
                onClick={() => handleEditTask(task)}
                className="flex-1 px-3 py-2 text-sm accent-green hover:bg-accent-light rounded transition-colors border-0"
              >
                编辑
              </button>
              <button 
                onClick={() => handleDeleteTask(task.id)}
                className="px-3 py-2 text-sm border border-red-300 text-red-600 rounded hover:bg-red-50 transition-colors"
              >
                删除
              </button>
            </div>
          </div>
        ))}
      </div>
      
      {filteredAndSortedTasks.length === 0 && !loading && (
        <div className="text-center py-12">
           <p className="text-gray-500 mb-4">
             {tasks.length === 0 ? '暂无任务，点击上方按钮创建新任务' : '没有找到匹配的任务'}
           </p>
           {tasks.length === 0 && (
             <button 
               onClick={fetchJobs}
               className="modern-button"
             >
               刷新列表
             </button>
           )}
         </div>
       )}
       
       <AddTaskModal
         isOpen={isAddTaskModalOpen}
         onClose={handleModalClose}
         onConfirm={handleModalConfirm}
       />
     </div>
   );
 };
 
 export default TasksPage;