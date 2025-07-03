import React from 'react';

interface HelpModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const HelpModal: React.FC<HelpModalProps> = ({ isOpen, onClose }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 w-full max-w-md mx-auto my-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-gray-800">使用说明</h2>
          <button
            onClick={onClose}
            className="w-8 h-8 rounded-full bg-red-100 hover:bg-red-200 flex items-center justify-center text-red-500 hover:text-red-600 transition-all duration-200 hover:scale-110 focus:outline-none focus:ring-2 focus:ring-red-300"
            aria-label="关闭说明"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div className="space-y-4 text-sm text-gray-700">
          <div>
            <h3 className="font-medium text-gray-800 mb-2">如何添加数据对：</h3>
            <ul className="list-disc list-inside space-y-1 ml-2">
              <li>在"输入"字段中填写训练数据的输入部分</li>
              <li>在"目标"字段中填写对应的期望输出</li>
              <li>当您在最后一个数据对中输入内容时，系统会自动创建新的空白数据对</li>
            </ul>
          </div>
          
          <div>
            <h3 className="font-medium text-gray-800 mb-2">数据对管理：</h3>
            <ul className="list-disc list-inside space-y-1 ml-2">
              <li>点击"删除"按钮可以移除不需要的数据对</li>
              <li>至少需要保留一个数据对</li>
              <li>只有包含内容的数据对才会被保存</li>
            </ul>
          </div>
          
          <div>
            <h3 className="font-medium text-gray-800 mb-2">保存数据集：</h3>
            <ul className="list-disc list-inside space-y-1 ml-2">
              <li>确保至少有一组有效的输入-目标对</li>
              <li>点击"保存数据集"完成创建</li>
              <li>点击"取消"将返回数据集列表页面</li>
            </ul>
          </div>
        </div>
        
        <div className="flex justify-end mt-6">
          <button
            onClick={onClose}
            className="modern-button px-4 py-2"
          >
            我知道了
          </button>
        </div>
      </div>
    </div>
  );
};

export default HelpModal;