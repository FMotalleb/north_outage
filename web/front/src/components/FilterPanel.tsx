import React from "react";
import { X } from "lucide-react";
import DatePicker from "react-multi-date-picker";
import persian from "react-date-object/calendars/persian";
import persian_fa from "react-date-object/locales/persian_fa";

import { FilterState } from "../types";
import "./FilterPanel.css";

// Optional: a safe, minimal theme for the popup (doesn't affect your app styles)
import "react-multi-date-picker/styles/layouts/mobile.css";
import DateObject from "react-date-object";

interface FilterPanelProps {
  filters: FilterState;
  cities: string[];
  onFilterChange: (key: keyof FilterState, value: string) => void;
  onClearFilters: () => void;
}

const toLocalISODate = (d: Date): string => {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`; // no timezone shift
};

const normalizeToDate = (val: DateObject | null): Date | null => {
  if (!val) return null;
  return val.toDate();
};

const FilterPanel: React.FC<FilterPanelProps> = ({
  filters,
  cities,
  onFilterChange,
  onClearFilters,
}) => {
  const handleDateChange = (value: DateObject | null) => {
    const d = normalizeToDate(value);
    if (d) {
      onFilterChange("date", toLocalISODate(d));
    } else {
      onFilterChange("date", "");
    }
  };

  return (
    <div>
      <div className="filters">
        <div className="filter-group">
          <label className="filter-label" htmlFor="cityFilter">
            شهر
          </label>
          <select
            id="cityFilter"
            className="filter-select"
            value={filters.city}
            onChange={(e) => onFilterChange("city", e.target.value)}
          >
            <option value="">همه شهرها</option>
            {cities.map((city) => (
              <option key={city} value={city}>
                {city}
              </option>
            ))}
          </select>
        </div>

        <div className="filter-group">
          <label className="filter-label" htmlFor="dateFilter">
            تاریخ
          </label>
          <DatePicker
            id="dateFilter"
            value={filters.date ? new Date(`${filters.date}T00:00:00`) : null}
            onChange={handleDateChange}
            calendar={persian}
            locale={persian_fa}
            format="YYYY/MM/DD"
            calendarPosition="bottom-right"
            inputClass="filter-input date-input" // style only the input, not global .rmdp-*
            portal // render popup in a portal; less layout interference
            editable={false}
            placeholder="انتخاب تاریخ"
          />
        </div>

        <div className="filter-group">
          <label className="filter-label" htmlFor="addressFilter">
            جستجو در آدرس
          </label>
          <input
            type="text"
            id="addressFilter"
            className="filter-input"
            placeholder="آدرس مورد نظر را جستجو کنید..."
            value={filters.address}
            onChange={(e) => onFilterChange("address", e.target.value)}
          />
        </div>
      </div>
      <button className="clear-filters" onClick={onClearFilters}>
        <X size={16} />
        پاک کردن فیلترها
      </button>
    </div>
  );
};

export default FilterPanel;
