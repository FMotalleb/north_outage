import React from 'react';
import { Search } from 'lucide-react';
import './EmptyState.css';

const EmptyState: React.FC = () => {
  return (
    <div className="empty-state">
      <Search className="empty-state-icon" />
      <h3>هیچ نتیجه‌ای یافت نشد</h3>
      <p>لطفاً فیلترهای خود را تغییر دهید</p>
    </div>
  );
};

export default EmptyState;