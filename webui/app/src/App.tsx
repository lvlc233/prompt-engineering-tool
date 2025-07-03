import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import DatasetsPage from './pages/DatasetsPage';
import CriteriaPage from './pages/CriteriaPage';
import CriteriaDetailPage from './pages/CriteriaDetailPage';
import TasksPage from './pages/TasksPage';
import TaskEditorPage from './pages/TaskEditorPage';
import DatasetEditorPage from './pages/DatasetEditorPage';
import DatasetDetailPage from './pages/DatasetDetailPage';
import SettingsPage from './pages/SettingsPage';
import './App.css'; // 保留 App.css 以便未来可能的全局样式调整

const App: React.FC = () => {
  return (
    <Router>
      <Layout>
        <Routes>
            <Route path="/" element={<DatasetsPage />} />
            <Route path="/datasets" element={<DatasetsPage />} />
            <Route path="/datasets/:id" element={<DatasetDetailPage />} />
            <Route path="/datasets/editor" element={<DatasetEditorPage />} />
            <Route path="/criteria" element={<CriteriaPage />} />
            <Route path="/criteria/:id" element={<CriteriaDetailPage />} />
            <Route path="/tasks" element={<TasksPage />} />
            <Route path="/tasks/editor" element={<TaskEditorPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Routes>
      </Layout>
    </Router>
  );
};

export default App;
