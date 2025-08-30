import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { DataItem } from '../types';

interface OutageChartProps {
  data: DataItem[];
}

const OutageChart: React.FC<OutageChartProps> = ({ data }) => {
  const chartData = data.reduce((acc, item) => {
    const city = item.city;
    if (!acc[city]) {
      acc[city] = { name: city, outages: 0 };
    }
    acc[city].outages++;
    return acc;
  }, {} as { [key: string]: { name: string; outages: number } });

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="custom-tooltip">
          <p className="label">{`${label} : ${payload[0].value}`}</p>
        </div>
      );
    }

    return null;
  };

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={Object.values(chartData)}>
        <defs>
          <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8}/>
            <stop offset="95%" stopColor="#8884d8" stopOpacity={0}/>
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="name" />
        <YAxis />
        <Tooltip content={<CustomTooltip />} />
        <Legend />
        <Bar dataKey="outages" fill="url(#colorUv)" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default OutageChart;
