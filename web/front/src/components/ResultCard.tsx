import React from "react";
import { MapPin } from "lucide-react";
import TimelineBar from "./TimelineBar";
import { DataItem } from "../types";
import "./ResultCard.css";

interface ResultCardProps {
  data: DataItem;
}

const ResultCard: React.FC<ResultCardProps> = ({ data }) => {
  const startTime = new Date(data.start);
  const endTime = new Date(data.end);

  // --- Duration calculation ---
  const duration = endTime.getTime() - startTime.getTime();
  const hours = Math.floor(duration / (1000 * 60 * 60));
  const minutes = Math.floor((duration % (1000 * 60 * 60)) / (1000 * 60));

  // --- Time strings ---
  const startTimeStr = startTime.toLocaleTimeString("fa-IR", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
  const endTimeStr = endTime.toLocaleTimeString("fa-IR", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });

  // --- Date comparison for "امروز" / "فردا" ---
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const tomorrow = new Date(today);
  tomorrow.setDate(today.getDate() + 1);

  const startDay = new Date(startTime);
  startDay.setHours(0, 0, 0, 0);

  let dateLabel: string;
  if (startDay.getTime() === today.getTime()) {
    dateLabel = "امروز";
  } else if (startDay.getTime() === tomorrow.getTime()) {
    dateLabel = "فردا";
  } else {
    // Fallback: format as Jalali date
    dateLabel = new Intl.DateTimeFormat("fa-IR-u-ca-persian", {
      year: "numeric",
      month: "long",
      day: "numeric",
    }).format(startTime);
  }

  // --- Duration text ---
  let durationText = "";
  if (hours > 0) durationText += `${hours} ساعت`;
  if (minutes > 0)
    durationText += `${durationText ? " و " : ""}${minutes} دقیقه`;

  return (
    <div className="result-card">
      <div className="card-city">
        <MapPin className="city-icon" />
        {data.city}
      </div>
      <div className="card-address">{data.address}</div>

      {/* Display date label */}
      <div className="card-date">{dateLabel}</div>

      <div className="timeline-container">
        <div className="timeline-header">
          <span className="timeline-time">
            {endTimeStr}-{startTimeStr}
          </span>
          <span className="timeline-duration">{durationText}</span>
        </div>
        <TimelineBar startTime={startTime} endTime={endTime} />
        <div className="timeline-labels">
          <span>23:59</span>
          <span>18:00</span>
          <span>12:00</span>
          <span>06:00</span>
          <span>00:00</span>
        </div>
      </div>
    </div>
  );
};

export default ResultCard;
