import React from 'react';
import { X } from 'lucide-react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { FilterState } from '../types';
import './FilterPanel.css';

interface FilterPanelProps {
  filters: FilterState;
  cities: string[];
  onFilterChange: (key: keyof FilterState, value: string) => void;
  onClearFilters: () => void;
}

const FilterPanel: React.FC<FilterPanelProps> = ({
  filters,
  cities,
  onFilterChange,
  onClearFilters
}) => {
  const handleDateChange = (date: Date | null) => {
    if (date) {
      onFilterChange('date', date.toISOString().split('T')[0]);
    } else {
      onFilterChange('date', '');
    }
  };

  return (
    <div className="filters">
      <div className="filter-group">
        <label className="filter-label" htmlFor="cityFilter">شهر</label>
        <select
          id="cityFilter"
          className="filter-select"
          value={filters.city}
          onChange={(e) => onFilterChange('city', e.target.value)}
        >
          <option value="">همه شهرها</option>
          {cities.map(city => (
            <option key={city} value={city}>{city}</option>
          ))}
        </select>
      </div>

      <div className="filter-group">
        <label className="filter-label" htmlFor="dateFilter">تاریخ</label>
        <DatePicker
          id="dateFilter"
          className="filter-input date-input"
          placeholderText="انتخاب تاریخ"
          selected={filters.date ? new Date(filters.date) : null}
          onChange={handleDateChange}
          dateFormat="yyyy/MM/dd"
        />
      </div>

      <div className="filter-group">
        <label className="filter-label" htmlFor="addressFilter">جستجو در آدرس</label>
        <input
          type="text"
          id="addressFilter"
          className="filter-input"
          placeholder="آدرس مورد نظر را جستجو کنید..."
          value={filters.address}
          onChange={(e) => onFilterChange('address', e.target.value)}
        />
      </div>

      <div className="filter-group">
        <button className="clear-filters" onClick={onClearFilters}>
          <X size={16} />
          پاک کردن فیلترها
        </button>
      </div>
    </div>
  );
};

export default FilterPanel;
