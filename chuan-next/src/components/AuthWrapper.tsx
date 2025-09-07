'use client';

import { useState, useEffect } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import { clientAPI } from '@/lib/client-api';

export function AuthWrapper({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        // 首先检查是否需要认证
        const configRes = await clientAPI.getConfig();
        const authEnabled = (configRes as any).auth_enabled;

        if (!authEnabled) {
          setIsAuthenticated(true);
          return;
        }

        // 如果需要认证，则检查 cookie
        await clientAPI.checkAuth();
        setIsAuthenticated(true);
      } catch (error) {
        console.log('Authentication check failed, redirecting to login.');
        setIsAuthenticated(false);
      }
    };

    checkAuthStatus();
  }, []);

  useEffect(() => {
    if (isAuthenticated === false && pathname !== '/login') {
      router.push('/login');
    }
  }, [isAuthenticated, pathname, router]);

  if (isAuthenticated === null) {
    // 加载状态，可以显示一个加载动画
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p>正在加载...</p>
      </div>
    );
  }

  if (isAuthenticated === false && pathname !== '/login') {
    // 等待重定向
    return null;
  }

  if (isAuthenticated === true && pathname === '/login') {
    // 如果已经登录，但仍在登录页，则跳转到首页
    router.push('/');
    return null;
  }

  return <>{children}</>;
}
