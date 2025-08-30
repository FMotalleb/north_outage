import React, { useState, useEffect, useMemo } from "react";
// import { MapPin, Search, Calendar, X } from "lucide-react";
import FilterPanel from "./FilterPanel";
import ResultCard from "./ResultCard";
import EmptyState from "./EmptyState";
import { DataItem, FilterState } from "../types";
import "./DataFilterSystem.css";

const DataFilterSystem: React.FC = () => {
  const [data, setData] = useState<DataItem[]>([]);
  const [filters, setFilters] = useState<FilterState>({
    city: "",
    date: "",
    address: "",
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("/events");
        const jsonData = await response.json();
        setData(jsonData);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };

    fetchData();
  }, []);

  const cities = useMemo(() => {
    return [...new Set(data.map((item) => item.city))];
  }, [data]);

  const filteredData = useMemo(() => {
    return data.filter((item) => {
      // City filter
      if (filters.city && item.city !== filters.city) {
        return false;
      }

      // Date filter
      if (filters.date) {
        const itemDate = new Date(item.start).toISOString().split("T")[0];
        if (itemDate !== filters.date) {
          return false;
        }
      }

      // Address filter
      if (
        filters.address &&
        !item.address.toLowerCase().includes(filters.address.toLowerCase())
      ) {
        return false;
      }

      return true;
    });
  }, [filters, data]);

  const handleFilterChange = (key: keyof FilterState, value: string) => {
    setFilters((prev) => ({
      ...prev,
      [key]: value,
    }));
  };

  const clearAllFilters = () => {
    setFilters({
      city: "ساری",
      date: "",
      address: "",
    });
  };

  return (
    <div className="data-filter-container">
      <div className="header">
        <h1>سیستم مشاهده قطعی برق</h1>
        <p>فیلتر و جستجوی قطعی‌های برق شهری</p>
      </div>

      <FilterPanel
        filters={filters}
        cities={cities}
        onFilterChange={handleFilterChange}
        onClearFilters={clearAllFilters}
      />
      {/* 
      <OutageChart data={filteredData} /> */}

      <div className="results-count">{filteredData.length} نتیجه یافت شد</div>

      <div className="results">
        {data.length === 0 ? (
          <p>Loading data...</p>
        ) : filteredData.length === 0 ? (
          <EmptyState />
        ) : (
          filteredData
            .slice(0, 20)
            .map((item) => <ResultCard key={item.id} data={item} />)
        )}
      </div>
    </div>
  );
};

export default DataFilterSystem;
