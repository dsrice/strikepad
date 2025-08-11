import React from 'react';
import { useAuth } from '../contexts/AuthContext';

const Dashboard: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-blue-600 text-white p-6">
        <div className="max-w-6xl mx-auto flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold">StrikePad Dashboard</h1>
            <p className="text-blue-100 mt-2">ようこそ、{user?.display_name}さん</p>
          </div>
          <button
              onClick={() => logout()}
            className="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded transition-colors"
          >
            ログアウト
          </button>
        </div>
      </header>
      
      <main className="max-w-6xl mx-auto p-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-4">ユーザー情報</h2>
            <div className="space-y-2">
              <p><span className="font-medium">ID:</span> {user?.id}</p>
              <p><span className="font-medium">メール:</span> {user?.email}</p>
              <p><span className="font-medium">表示名:</span> {user?.display_name}</p>
              <p>
                <span className="font-medium">メール認証:</span>{' '}
                <span className={user?.email_verified ? 'text-green-600' : 'text-red-600'}>
                  {user?.email_verified ? '認証済み' : '未認証'}
                </span>
              </p>
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-4">D3.js Chart Demo</h2>
          </div>
        </div>
      </main>
    </div>
  );
};

export default Dashboard;