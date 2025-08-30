import React from 'react';
import './TimelineBar.css';

interface TimelineBarProps {
  startTime: Date;
  endTime: Date;
}

const TimelineBar: React.FC<TimelineBarProps> = ({ startTime, endTime }) => {
  // Calculate position and width for the timeline bar
  const startHour = startTime.getHours();
  const startMinute = startTime.getMinutes();
  const endHour = endTime.getHours();
  const endMinute = endTime.getMinutes();
  
  const startTotalMinutes = startHour * 60 + startMinute;
  const endTotalMinutes = endHour * 60 + endMinute;
  const totalDayMinutes = 24 * 60; // 1440 minutes in a day
  
  const startPercentage = (startTotalMinutes / totalDayMinutes) * 100;
  const durationMinutes = endTotalMinutes - startTotalMinutes;
  const durationPercentage = (durationMinutes / totalDayMinutes) * 100;

  return (
    <div className="timeline-bar">
      <div 
        className="timeline-fill" 
        style={{ 
          left: `${startPercentage}%`, 
          width: `${durationPercentage}%` 
        }}
      />
    </div>
  );
};

export default TimelineBar;