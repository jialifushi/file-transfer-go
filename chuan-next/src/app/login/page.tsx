'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { clientAPI } from '@/lib/client-api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export default function LoginPage() {
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await clientAPI.login(code);
      if (response.success) {
        // 登录成功，跳转到首页
        router.push('/');
      } else {
        setError('无效的访问码，请重试。');
      }
    } catch (err) {
      setError('登录失败，请检查网络或联系管理员。');
      console.error(err);
    }

    setLoading(false);
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-lg shadow-md dark:bg-gray-800">
        <h1 className="text-2xl font-bold text-center text-gray-900 dark:text-white">P2P服务 - 访问授权</h1>
        <p className="text-center text-gray-600 dark:text-gray-400">请输入访问码以继续</p>
        <form onSubmit={handleLogin} className="space-y-6">
          <div>
            <Input
              id="access-code"
              type="password"
              placeholder="请输入访问码"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              required
              className="w-full px-4 py-2 border rounded-md"
            />
          </div>
          {error && <p className="text-sm text-red-500 text-center">{error}</p>}
          <div>
            <Button type="submit" disabled={loading} className="w-full">
              {loading ? '验证中...' : '进入'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
