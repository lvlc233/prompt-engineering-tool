import React, { useState } from 'react';

const SettingsPage: React.FC = () => {
  const [baseUrl, setBaseUrl] = useState('');
  const [model, setModel] = useState('');
  const [apiKey, setApiKey] = useState('');

  const handleSave = () => {
    // 这里可以添加保存设置的逻辑
    console.log('保存设置:', { baseUrl, model, apiKey });
    alert('设置已保存！');
  };

  const handleReset = () => {
    setBaseUrl('');
    setModel('');
    setApiKey('');
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800 mb-2">设置</h1>
        <p className="text-gray-600">配置API连接和模型参数</p>
      </div>

      <div className="modern-card p-6">
        <div className="space-y-6">
          {/* Base URL */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Base URL
            </label>
            <input
              type="url"
              value={baseUrl}
              onChange={(e) => setBaseUrl(e.target.value)}
              className="modern-input"
              placeholder="https://api.openai.com/v1"
              aria-label="API Base URL"
            />
            <p className="text-xs text-gray-500 mt-1">
              API服务的基础URL地址
            </p>
          </div>

          {/* Model */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Model
            </label>
            <input
              type="text"
              value={model}
              onChange={(e) => setModel(e.target.value)}
              className="modern-input"
              placeholder="gpt-3.5-turbo"
              aria-label="AI模型名称"
            />
            <p className="text-xs text-gray-500 mt-1">
              要使用的AI模型名称
            </p>
          </div>

          {/* API Key */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              API Key
            </label>
            <input
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              className="modern-input"
              placeholder="sk-..."
              aria-label="API密钥"
            />
            <p className="text-xs text-gray-500 mt-1">
              您的API密钥，将安全存储
            </p>
          </div>
        </div>

        {/* 操作按钮 */}
        <div className="flex gap-4 mt-8 pt-6 border-t border-gray-200">
          <button
            onClick={handleSave}
            className="modern-button"
            aria-label="保存设置"
          >
            保存设置
          </button>
          <button
            onClick={handleReset}
            className="px-4 py-2 border border-gray-300 text-gray-600 rounded hover:bg-gray-50 transition-colors"
            aria-label="重置设置"
          >
            重置
          </button>
        </div>
      </div>

      {/* 帮助信息 */}
      <div className="mt-6 modern-card p-6 bg-blue-50">
        <h3 className="text-lg font-semibold text-blue-800 mb-3">配置说明</h3>
        <div className="space-y-2 text-sm text-blue-700">
          <p><strong>Base URL:</strong> API服务的基础地址，通常以 /v1 结尾</p>
          <p><strong>Model:</strong> 要使用的AI模型，如 gpt-3.5-turbo、gpt-4 等</p>
          <p><strong>API Key:</strong> 您的API访问密钥，请妥善保管</p>
        </div>
      </div>
    </div>
  );
};

export default SettingsPage;