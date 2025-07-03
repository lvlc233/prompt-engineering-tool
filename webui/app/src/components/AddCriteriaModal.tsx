import React, { useState } from 'react';

interface AddCriteriaModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (name: string, description: string, scoreLimit: number) => void;
}

const AddCriteriaModal: React.FC<AddCriteriaModalProps> = ({ isOpen, onClose, onConfirm }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [scoreLimit, setScoreLimit] = useState<number>(100);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim() && description.trim() && scoreLimit > 0) {
      onConfirm(name.trim(), description.trim(), scoreLimit);
      setName('');
      setDescription('');
      setScoreLimit(100);
    }
  };

  const handleClose = () => {
    setName('');
    setDescription('');
    setScoreLimit(100);
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 w-full max-w-md mx-auto my-auto">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">添加新测评集</h2>
        
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="criteria-name" className="block text-sm font-medium text-gray-700 mb-2">
              测评集名称
            </label>
            <input
              id="criteria-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="modern-input"
              placeholder="请输入测评集名称"
              required
            />
          </div>
          
          <div className="mb-4">
            <label htmlFor="criteria-description" className="block text-sm font-medium text-gray-700 mb-2">
              测评集描述
            </label>
            <textarea
              id="criteria-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="modern-input h-24 resize-none"
              placeholder="请输入测评集描述"
              required
            />
          </div>
          
          <div className="mb-6">
            <label htmlFor="score-limit" className="block text-sm font-medium text-gray-700 mb-2">
              分数上限
            </label>
            <input
              id="score-limit"
              type="number"
              value={scoreLimit}
              onChange={(e) => setScoreLimit(Number(e.target.value))}
              className="modern-input"
              placeholder="请输入分数上限"
              min="1"
              max="1000"
              required
            />
          </div>
          
          <div className="flex gap-3 justify-end">
            <button
              type="button"
              onClick={handleClose}
              className="px-4 py-2 text-sm border-0 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
            >
              取消
            </button>
            <button
              type="submit"
              className="modern-button"
            >
              确认
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default AddCriteriaModal;