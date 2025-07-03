import React, { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import AddTaskModal from '../components/AddTaskModal';

// Git风格的版本节点数据结构
interface GitNode {
  id: string;
  message: string;
  author?: string;
  timestamp: string;
  parentIds: string[];
  branch: string;
  status: 'committed' | 'current' | 'draft';
}

interface Task {
  id: number;
  name: string;
  description: string;
  createdAt: string;
  versions: GitNode[];
  currentVersionId: string;
}

// Git风格的版本树组件
const GitVersionTree: React.FC<{ 
  nodes: GitNode[]; 
  currentVersionId: string;
  onVersionSelect?: (versionId: string) => void;
}> = ({ nodes, currentVersionId, onVersionSelect }) => {
  const [selectedNode, setSelectedNode] = React.useState<string | null>(null);
  const [hoveredNode, setHoveredNode] = React.useState<string | null>(null);
  const [isCollapsed, setIsCollapsed] = React.useState<boolean>(false);

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
    const nodeHeight = 60;
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
      return { fill: '#3b82f6', stroke: '#1d4ed8', text: '#ffffff' }; // 蓝色 - 当前版本
    }
    
    switch (node.status) {
      case 'committed':
        return { fill: '#10b981', stroke: '#059669', text: '#ffffff' }; // 绿色 - 已提交
      case 'draft':
        return { fill: '#f59e0b', stroke: '#d97706', text: '#ffffff' }; // 橙色 - 草稿
      default:
        return { fill: '#6b7280', stroke: '#4b5563', text: '#ffffff' }; // 灰色 - 默认
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
    if (onVersionSelect) {
      onVersionSelect(nodeId);
    }
  };

  // 计算SVG尺寸
  const maxX = Math.max(...Array.from(positions.values()).map(p => p.x)) + 60;
  const maxY = Math.max(...Array.from(positions.values()).map(p => p.y)) + 60;

  return (
    <div className="git-version-tree bg-white border border-gray-200 rounded-lg">
      {/* 版本标题区域 - 可点击展开收起 */}
      <div 
        className="flex items-center gap-2 pb-2 cursor-pointer hover:bg-gray-50 transition-colors duration-200"
        onClick={() => setIsCollapsed(!isCollapsed)}
      >
        <h4 className="text-sm font-medium text-gray-900">版本</h4>
        {selectedNode && (
          <span className="text-xs text-blue-600 bg-blue-50 px-2 py-1 rounded">
            当前: {selectedNode}
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
            
            return node.parentIds.map(parentId => {
              const parentPos = positions.get(parentId);
              if (!parentPos) return null;
              
              const isSameBranch = nodeMap.get(parentId)?.branch === node.branch;
              const strokeColor = isSameBranch ? getBranchColor(node.branch) : '#d1d5db';
              
              return (
                <line
                  key={`${parentId}-${node.id}`}
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
          })}
          
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
                  onMouseEnter={() => setHoveredNode(node.id)}
                  onMouseLeave={() => setHoveredNode(null)}
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
                  {node.id}
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
      {!isCollapsed && (hoveredNode || selectedNode) && (
        <div className="absolute right-4 top-4 w-72 p-4 bg-white border border-gray-300 rounded-lg shadow-lg z-10">
          {(() => {
            // 优先显示悬停节点，否则显示选中节点
            const displayNodeId = hoveredNode || selectedNode;
            const node = nodeMap.get(displayNodeId!);
            if (!node) return null;
            
            const parentNodes = node.parentIds.map(id => nodeMap.get(id)).filter(Boolean);
            const isHovering = hoveredNode === displayNodeId;
            
            return (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <span className="font-semibold text-gray-900">版本号: {node.id}</span>
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
                  <div>执行时间: {node.author ? `由 ${node.author} 执行` : '未知'}</div>
                  <div>父版本: {parentNodes.length > 0 ? parentNodes.map(p => p!.id).join(', ') : '无'}</div>
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
        <div className="mt-2 p-2 bg-white border border-gray-200 rounded-lg">
          <div className="flex flex-wrap gap-4 text-xs">
            <div className="flex items-center gap-1">
              <span>🔵</span>
              <span className="text-gray-600">当前版本</span>
            </div>
            <div className="flex items-center gap-1">
              <span>🟢</span>
              <span className="text-gray-600">已提交</span>
            </div>
            <div className="flex items-center gap-1">
              <span>🟡</span>
              <span className="text-gray-600">草稿</span>
            </div>
            <div className="flex items-center gap-1">
              <span>⚫</span>
              <span className="text-gray-600">默认</span>
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
  const [tasks, setTasks] = React.useState<Task[]>([
    {
      id: 1,
      name: "智能机器人提取训练",
      description: "使用机器学习算法对数据进行训练",
      createdAt: "2024-01-15",
      currentVersionId: "c4f2a1b",
      versions: [
        {
          id: "a1b2c3d",
          message: "初始项目设置和基础架构",
          author: "张三",
          timestamp: "2024-01-15 10:00",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "b2c3d4e",
          message: "添加数据预处理模块",
          author: "李四",
          timestamp: "2024-01-16 14:30",
          parentIds: ["a1b2c3d"],
          branch: "main",
          status: "committed"
        },
        {
          id: "c3d4e5f",
          message: "实现机器学习模型训练",
          author: "王五",
          timestamp: "2024-01-17 09:15",
          parentIds: ["b2c3d4e"],
          branch: "main",
          status: "committed"
        },
        {
          id: "d4e5f6g",
          message: "创建实验分支用于新算法测试",
          author: "赵六",
          timestamp: "2024-01-17 16:45",
          parentIds: ["b2c3d4e"],
          branch: "feature/new-algorithm",
          status: "draft"
        },
        {
          id: "c4f2a1b",
          message: "优化模型性能和添加评估指标",
          author: "张三",
          timestamp: "2024-01-18 11:20",
          parentIds: ["c3d4e5f"],
          branch: "main",
          status: "current"
        }
      ]
    },
    {
      id: 2,
      name: "数据分析报告",
      description: "生成详细的数据分析报告",
      createdAt: "2024-01-16",
      currentVersionId: "f6g7h8i",
      versions: [
        {
          id: "e5f6g7h",
          message: "创建基础报告模板",
          author: "李四",
          timestamp: "2024-01-16 11:20",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "f6g7h8i",
          message: "添加交互式图表和高级分析",
          author: "王五",
          timestamp: "2024-01-19 13:45",
          parentIds: ["e5f6g7h"],
          branch: "main",
          status: "current"
        },
        {
          id: "g7h8i9j",
          message: "实验性数据可视化功能",
          author: "赵六",
          timestamp: "2024-01-20 10:30",
          parentIds: ["e5f6g7h"],
          branch: "feature/visualization",
          status: "draft"
        }
      ]
    },
    {
      id: 3,
      name: "自然语言处理模型",
      description: "开发先进的NLP模型用于文本理解",
      createdAt: "2024-01-17",
      currentVersionId: "j9k0l1m",
      versions: [
        {
          id: "h8i9j0k",
          message: "初始化NLP项目结构",
          author: "陈七",
          timestamp: "2024-01-17 09:00",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "i9j0k1l",
          message: "添加词向量训练模块",
          author: "周八",
          timestamp: "2024-01-18 15:20",
          parentIds: ["h8i9j0k"],
          branch: "main",
          status: "committed"
        },
        {
          id: "j9k0l1m",
          message: "实现Transformer架构",
          author: "吴九",
          timestamp: "2024-01-19 16:30",
          parentIds: ["i9j0k1l"],
          branch: "main",
          status: "current"
        }
      ]
    },
    {
      id: 4,
      name: "图像识别系统",
      description: "构建高精度的图像分类和目标检测系统",
      createdAt: "2024-01-18",
      currentVersionId: "m1n2o3p",
      versions: [
        {
          id: "k0l1m2n",
          message: "搭建CNN基础架构",
          author: "郑十",
          timestamp: "2024-01-18 08:45",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "l1m2n3o",
          message: "集成数据增强技术",
          author: "孙十一",
          timestamp: "2024-01-19 12:15",
          parentIds: ["k0l1m2n"],
          branch: "main",
          status: "committed"
        },
        {
          id: "m1n2o3p",
          message: "优化模型准确率",
          author: "李十二",
          timestamp: "2024-01-20 14:50",
          parentIds: ["l1m2n3o"],
          branch: "main",
          status: "current"
        }
      ]
    },
    {
      id: 5,
      name: "推荐系统算法",
      description: "开发个性化推荐算法提升用户体验",
      createdAt: "2024-01-19",
      currentVersionId: "p3q4r5s",
      versions: [
        {
          id: "n2o3p4q",
          message: "实现协同过滤算法",
          author: "王十三",
          timestamp: "2024-01-19 10:30",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "o3p4q5r",
          message: "添加深度学习推荐模型",
          author: "张十四",
          timestamp: "2024-01-20 09:20",
          parentIds: ["n2o3p4q"],
          branch: "main",
          status: "committed"
        },
        {
          id: "p3q4r5s",
          message: "集成实时推荐引擎",
          author: "李十五",
          timestamp: "2024-01-21 11:40",
          parentIds: ["o3p4q5r"],
          branch: "main",
          status: "current"
        }
      ]
    },
    {
      id: 6,
      name: "语音识别引擎",
      description: "构建多语言语音识别和转换系统",
      createdAt: "2024-01-20",
      currentVersionId: "s5t6u7v",
      versions: [
        {
          id: "q4r5s6t",
          message: "初始化语音处理框架",
          author: "赵十六",
          timestamp: "2024-01-20 13:15",
          parentIds: [],
          branch: "main",
          status: "committed"
        },
        {
          id: "r5s6t7u",
          message: "实现声学模型训练",
          author: "钱十七",
          timestamp: "2024-01-21 08:30",
          parentIds: ["q4r5s6t"],
          branch: "main",
          status: "committed"
        },
        {
          id: "s5t6u7v",
          message: "优化识别准确率和速度",
          author: "孙十八",
          timestamp: "2024-01-22 15:45",
          parentIds: ["r5s6t7u"],
          branch: "main",
          status: "current"
        }
      ]
    }
  ]);

  const handleVersionChange = (taskId: number, versionId: string) => {
    console.log(`切换任务 ${taskId} 到版本 ${versionId}`);
    // 更新任务的当前版本ID
    setTasks(prevTasks => 
      prevTasks.map(task => 
        task.id === taskId 
          ? { ...task, currentVersionId: versionId }
          : task
      )
    );
  };

  const handleDeleteTask = (id: number) => {
    if (window.confirm('确定要删除这个任务吗？此操作不可撤销。')) {
      setTasks(prevTasks => prevTasks.filter(task => task.id !== id));
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
        comparison = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime();
      }
      return sortOrder === 'asc' ? comparison : -comparison;
    });

    return filtered;
   }, [tasks, searchTerm, sortBy, sortOrder]);

   const handleAddTask = () => {
     setIsAddTaskModalOpen(true);
   };

   const handleModalClose = () => {
     setIsAddTaskModalOpen(false);
   };

   const handleModalConfirm = (name: string, description: string) => {
     setIsAddTaskModalOpen(false);
     navigate('/tasks/editor', {
       state: {
         taskName: name,
         taskDescription: description
       }
     });
   };
 
   return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-semibold text-white mb-2">任务管理</h1>
        <p className="text-white">管理和监控所有AI模型训练任务</p>
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
              <h3 className="text-xl font-bold text-white mb-2">{task.name}</h3>
              <p className="text-white text-base mb-3 flex-1">{task.description}</p>
            </div>
            
            <div className="mb-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-300">创建时间:</span>
                <span className="text-white">{task.createdAt}</span>
              </div>
            </div>
            
            <div className="mb-4">
              <GitVersionTree 
                  nodes={task.versions} 
                  currentVersionId={task.currentVersionId}
                  onVersionSelect={(versionId) => handleVersionChange(task.id, versionId)}
                />
            </div>
            
            <div className="flex gap-2 mt-auto">
              <button className="flex-1 px-3 py-2 text-sm border-0 bg-gray-50 rounded hover:bg-gray-100 transition-colors text-gray-800">
                查看详情
              </button>
              <button className="flex-1 px-3 py-2 text-sm accent-green hover:bg-accent-light rounded transition-colors border-0">
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
      
      {filteredAndSortedTasks.length === 0 && (
        <div className="text-center py-12">
           <p className="text-gray-500">没有找到匹配的任务</p>
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