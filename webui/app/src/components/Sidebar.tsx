import React from 'react';
import { Link } from 'react-router-dom';

const Sidebar: React.FC = () => {
  const menuItems = [
    { name: '数据集', path: '/datasets' },
    { name: '评测集', path: '/evaluation' },
    { name: '任务', path: '/jobs' },
    { name: '设置', path: '/settings' },
  ];

  return (
    <aside className="w-64 bg-gray-100 p-4 border-r border-gray-200">
      <nav>
        <ul>
          {menuItems.map((item) => (
            <li key={item.name} className="mb-2">
              <Link
                to={item.path}
                className="block px-4 py-2 text-gray-700 rounded hover:bg-light-green hover:text-white transition-colors duration-200"
                aria-label={`导航到${item.name}`}
                tabIndex={0}
              >
                {item.name}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;