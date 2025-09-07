"use client";

import { Suspense } from 'react';
import HomePage from './HomePage';
import { AuthWrapper } from '@/components/AuthWrapper';

function HomePageWrapper() {
  return (
    <AuthWrapper>
      <Suspense fallback={<div className="min-h-screen flex items-center justify-center">加载中...</div>}>
        <HomePage />
      </Suspense>
    </AuthWrapper>
  );
}

export default HomePageWrapper;
