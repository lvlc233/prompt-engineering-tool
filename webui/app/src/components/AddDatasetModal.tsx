import React, { useState } from 'react';

interface AddDatasetModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (name: string, description: string) => void;
}

const AddDatasetModal: React.FC<AddDatasetModalProps> = ({ isOpen, onClose, onConfirm }) => {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (name.trim() && description.trim()) {
      onConfirm(name.trim(), description.trim());
      setName('');
      setDescription('');
    }
  };

  const handleClose = () => {
    setName('');
    setDescription('');
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 w-full max-w-md mx-auto my-auto">
        <h2 className="text-xl font-semibold text-gray-800 mb-4">添加新数据集</h2>
        
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="dataset-name" className="block text-sm font-medium text-gray-700 mb-2">
              数据集名称
            </label>
            <input
              id="dataset-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="modern-input"
              placeholder="请输入数据集名称"
              required
            />
          </div>
          
          <div className="mb-6">
            <label htmlFor="dataset-description" className="block text-sm font-medium text-gray-700 mb-2">
              数据集描述
            </label>
            <textarea
              id="dataset-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="modern-input h-24 resize-none"
              placeholder="请输入数据集描述"
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

export default AddDatasetModal;