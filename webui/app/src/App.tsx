import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import DatasetsPage from './pages/DatasetsPage';
import EvaluationPage from './pages/EvaluationPage';
import EvaluationDetailPage from './pages/EvaluationDetailPage';
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
            <Route path="/evaluation" element={<EvaluationPage />} />
            <Route path="/evaluation/:id" element={<EvaluationDetailPage />} />
            <Route path="/jobs" element={<TasksPage />} />
            <Route path="/jobs/editor" element={<TaskEditorPage />} />
            <Route path="/settings" element={<SettingsPage />} />
          </Routes>
      </Layout>
    </Router>
  );
};

export default App;
